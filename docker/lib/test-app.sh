#!/bin/bash -xeu

docker build -t fabrizio2210/webstorage-api -f docker/x86_64/Dockerfile-api .

stack="docker/lib/stack-test-x86_64.yml"

docker_compose="docker-compose"
docker compose version && docker_compose="docker compose"

${docker_compose} -f ${stack} run api
${docker_compose} -f ${stack} down
