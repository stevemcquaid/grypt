#!/bin/bash
set -ex

sudo cat test/foo.txt
docker run -it -v $PWD/test:/test --rm stevemcquaid/grypt:latest grypt -e -f /test/foo.txt

sudo cat test/foo.txt

docker run -it -v $PWD/test:/test --rm stevemcquaid/grypt:latest grypt -D -f /test/foo.txt
sudo cat test/foo.txt

sudo chown $USER:$USER test/foo.txt
