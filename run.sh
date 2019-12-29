#!/usr/bin/bash

script_dir="$( cd "$(dirname "$0")" ; pwd -P )"
export GOPATH=$script_dir

go run cmd/webchan 0.0.0.0:80
