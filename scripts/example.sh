#!/bin/bash
set -ex

docker run -it --rm stevemcquaid/grypt:latest grypt -e -f grypt.go