### Environment Variables

This package can be configured through various environment variables.

#### `GOLOG_LOG_LEVEL`

Specifies the log-level, both globally and on a per-subsystem basis.

For example, the following will set the global minimum log level to `error`, but reduce the minimum
log level for `subsystem1` to `info` and reduce the minimum log level for `subsystem2` to debug.

```bash
export GOLOG_LOG_LEVEL="error,subsystem1=info,subsystem2=debug"
```

`IPFS_LOGGING` is a deprecated alias for this environment variable.

#### `GOLOG_FILE`

Specifies that logs should be written to the specified file. If this option is _not_ specified, logs are written to standard error.

```bash
export GOLOG_FILE="/path/to/my/file.log"
```

#### `GOLOG_OUTPUT`

Specifies where logging output should be written. Can take one or more of the following values, combined with `+`:

- `stdout` -- write logs to standard out.
- `stderr` -- write logs to standard error.
- `file` -- write logs to the file specified by `GOLOG_FILE`

For example, if you want to log to both a file and standard error:

```bash
export GOLOG_FILE="/path/to/my/file.log"
export GOLOG_OUTPUT="stderr+file"
```

Setting _only_ `GOLOG_FILE` will prevent logs from being written to standard error.

#### `GOLOG_LOG_FMT`

Specifies the log message format. It supports the following values:

- `color` -- human readable, colorized (ANSI) output
- `nocolor` -- human readable, plain-text output.
- `json` -- structured JSON.

For example, to log structured JSON (for easier parsing):

```bash
export GOLOG_LOG_FMT="json"
```

The logging format defaults to `color` when the output is a terminal, and `nocolor` otherwise.

`IPFS_LOGGING_FMT` is a deprecated alias for this environment variable.

#### `GOLOG_LOG_LABELS`

Specifies a set of labels that should be added to all log messages as comma-separated key-value
pairs. For example, the following add `{"app": "example_app", "dc": "sjc-1"}` to every log entry.

```bash
export GOLOG_LOG_LABELS="app=example_app,dc=sjc-1"

```

#### `GOLOG_LOG_CALLER_FORMAT`

Specifies the file line format.

- `full` Serializes a caller in /full/path/to/package/file:line.
- `short` serializes a caller in package/file:line format, trimming all but the final directory from the full path.
```bash
export GOLOG_LOG_CALLER_FORMAT="full"
```
