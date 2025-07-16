#!/bin/sh
set -x

sudo wget https://github.com/CosmWasm/wasmvm/releases/download/v1.4.1/libwasmvm.x86_64.so -O /lib/libwasmvm.x86_64.so
export LD_LIBRARY_PATH=/lib:$LD_LIBRARY_PATH
