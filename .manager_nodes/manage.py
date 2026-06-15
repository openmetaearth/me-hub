#!/usr/bin/env python3
"""Batch run commands over SSH on hosts listed in hostlist.txt.

Usage:
    ./manage.py ps-med
    ./manage.py docker-ps
    ./manage.py list                       # list registered command aliases
    ./manage.py run "uname -a"             # one-off arbitrary command
    ./manage.py ps-med --hosts c1-143 c1-144
    ./manage.py docker-ps --parallel 4 --timeout 30

Hosts in hostlist.txt should be resolvable by ssh (typically via ~/.ssh/config).
For sudo commands, the remote user needs NOPASSWD sudo, otherwise the call hangs
because we use BatchMode=yes (no interactive password prompts).
"""

import argparse
import concurrent.futures
import shlex
import subprocess
import sys
import time
from pathlib import Path

HOSTLIST = Path(__file__).parent / "hostlist.txt"

COMMANDS = {
    "ps-med":      "ps -ef | grep med | grep -v grep",
    "docker-ps":   "sudo docker ps -a",
    "docker-pull": "sudo docker pull ghcr.io/openmetaearth/med:v2.0.14",
}

SSH_OPTS = [
    "-o", "BatchMode=yes",
    "-o", "ConnectTimeout=10",
    "-o", "StrictHostKeyChecking=accept-new",
    "-o", "ServerAliveInterval=15",
]


def load_hosts():
    hosts = []
    for line in HOSTLIST.read_text().splitlines():
        line = line.strip()
        if line and not line.startswith("#"):
            hosts.append(line)
    return hosts


def run_on_host(host, command, timeout):
    start = time.monotonic()
    try:
        result = subprocess.run(
            ["ssh", *SSH_OPTS, host, command],
            capture_output=True, text=True, timeout=timeout,
        )
        return host, result.returncode, result.stdout, result.stderr, time.monotonic() - start
    except subprocess.TimeoutExpired:
        return host, 124, "", f"timeout after {timeout}s", time.monotonic() - start
    except Exception as e:
        return host, -1, "", f"{type(e).__name__}: {e}", time.monotonic() - start


def color(s, code):
    return f"\033[{code}m{s}\033[0m" if sys.stdout.isatty() else s


def print_result(host, rc, stdout, stderr, elapsed):
    status = color("OK", "32") if rc == 0 else color(f"FAIL rc={rc}", "31")
    print(f"\n===== {color(host, '1;36')}  [{status}, {elapsed:.1f}s] =====")
    if stdout.strip():
        print(stdout.rstrip())
    if stderr.strip():
        print(color("[stderr]", "33"), stderr.rstrip())


def main():
    p = argparse.ArgumentParser(description=__doc__, formatter_class=argparse.RawDescriptionHelpFormatter)
    p.add_argument("command", help="alias from COMMANDS, 'list', or 'run' (with positional cmd)")
    p.add_argument("extra", nargs="*", help="for 'run': the command to execute")
    p.add_argument("--hosts", nargs="+", help="limit to these hosts (default: all in hostlist.txt)")
    p.add_argument("--parallel", type=int, default=10, help="max concurrent SSH connections (default 10)")
    p.add_argument("--timeout", type=int, default=20, help="per-host timeout seconds (default 20)")
    args = p.parse_args()

    if args.command == "list":
        print("Available commands:")
        for name, cmd in COMMANDS.items():
            print(f"  {name:<12} {cmd}")
        print("  run <cmd>    arbitrary command")
        return 0

    if args.command == "run":
        if not args.extra:
            p.error("'run' requires a command, e.g. ./manage.py run 'uname -a'")
        cmd = " ".join(shlex.quote(a) for a in args.extra) if len(args.extra) > 1 else args.extra[0]
    elif args.command in COMMANDS:
        cmd = COMMANDS[args.command]
    else:
        p.error(f"unknown command '{args.command}'. Try './manage.py list'.")

    hosts = load_hosts()
    if args.hosts:
        want = set(args.hosts)
        hosts = [h for h in hosts if h in want]
        missing = want - set(hosts)
        if missing:
            print(f"warning: hosts not in hostlist.txt: {', '.join(sorted(missing))}", file=sys.stderr)
    if not hosts:
        print("no hosts to run on", file=sys.stderr)
        return 1

    print(f"Running on {len(hosts)} host(s) [parallel={args.parallel}, timeout={args.timeout}s]: {cmd}")

    failed = []
    with concurrent.futures.ThreadPoolExecutor(max_workers=args.parallel) as ex:
        futures = {ex.submit(run_on_host, h, cmd, args.timeout): h for h in hosts}
        for fut in concurrent.futures.as_completed(futures):
            host, rc, out, err, elapsed = fut.result()
            print_result(host, rc, out, err, elapsed)
            if rc != 0:
                failed.append(host)

    print()
    if failed:
        print(color(f"FAILED on {len(failed)}/{len(hosts)}: {', '.join(failed)}", "31"))
        return 2
    print(color(f"OK on all {len(hosts)} host(s)", "32"))
    return 0


if __name__ == "__main__":
    sys.exit(main())
