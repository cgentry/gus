#!/bin/sh

export SOURCE="${BASH_SOURCE[0]}"
startdir=`dirname "$SOURCE"`
find "$startdir" -name '*.go' | xargs gofmt -w=true