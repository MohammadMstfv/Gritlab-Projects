#!/bin/sh
set -eu

docker image build -f Dockerfile -t ascii-art-web-docker .
docker container run --detach --name dockerize --publish 8080:8080 ascii-art-web-docker

printf 'Container dockerize is running. Open http://localhost:8080\n'
