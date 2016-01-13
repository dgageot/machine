#!/usr/bin/env bash

set -e

eval $(docker-machine env default)
TOKEN=$(docker run swarm create)

docker-machine create -d virtualbox --swarm --swarm-master --swarm-discovery="token://$TOKEN" swarm-master
docker-machine create -d virtualbox --swarm --swarm-discovery="token://$TOKEN" swarm-node-01
docker-machine create -d virtualbox --swarm --swarm-discovery="token://$TOKEN" swarm-node-02
docker-machine create -d virtualbox --swarm --swarm-discovery="token://$TOKEN" swarm-node-03

eval $(docker-machine env --swarm swarm-master)

docker info

docker run --name=interlock -p 80:80 -d -v $HOME/.docker/machine/certs:/etc/docker ehazlett/interlock --swarm-url $(docker-machine url swarm-master)  --swarm-url $DOCKER_HOST --swarm-tls-ca-cert=/etc/docker/ca.pem --swarm-tls-cert=/etc/docker/cert.pem --swarm-tls-key=/etc/docker/key.pem --plugin haproxy start
open http://stats:interlock@$(docker port interlock 80)/haproxy?stats
docker run -d -P --hostname test.local ehazlett/docker-demo
