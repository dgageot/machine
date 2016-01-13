#!/usr/bin/env bash

set -e

export GOOGLE_PROJECT=code-story-blog
export GOOGLE_MACHINE_IMAGE=https://www.googleapis.com/compute/v1/projects/code-story-blog/global/images/ubuntu-aufs
export GOOGLE_DISK_SIZE=100

docker-machine create -d google mh-keystore

docker $(docker-machine config mh-keystore) run -d -p "8500:8500" --name="consul" -h "consul" progrium/consul -server -bootstrap

docker-machine create -d google --swarm --swarm-master \
    --swarm-discovery="consul://$(docker-machine ip mh-keystore):8500" \
    --engine-opt="cluster-store=consul://$(docker-machine ip mh-keystore):8500" \
    --engine-opt="cluster-advertise=eth0:2376" \
    master

docker-machine create -d google --swarm \
    --swarm-discovery="consul://$(docker-machine ip mh-keystore):8500" \
    --engine-opt="cluster-store=consul://$(docker-machine ip mh-keystore):8500" \
    --engine-opt="cluster-advertise=eth0:2376" \
    node-01

docker-machine create -d google --swarm \
    --swarm-discovery="consul://$(docker-machine ip mh-keystore):8500" \
    --engine-opt="cluster-store=consul://$(docker-machine ip mh-keystore):8500" \
    --engine-opt="cluster-advertise=eth0:2376" \
    node-02

docker-machine create -d google --swarm \
    --swarm-discovery="consul://$(docker-machine ip mh-keystore):8500" \
    --engine-opt="cluster-store=consul://$(docker-machine ip mh-keystore):8500" \
    --engine-opt="cluster-advertise=eth0:2376" \
    node-03

docker-machine create -d google --swarm \
    --swarm-discovery="consul://$(docker-machine ip mh-keystore):8500" \
    --engine-opt="cluster-store=consul://$(docker-machine ip mh-keystore):8500" \
    --engine-opt="cluster-advertise=eth0:2376" \
    node-04

docker-machine create -d google --swarm \
    --swarm-discovery="consul://$(docker-machine ip mh-keystore):8500" \
    --engine-opt="cluster-store=consul://$(docker-machine ip mh-keystore):8500" \
    --engine-opt="cluster-advertise=eth0:2376" \
    node-05


eval $(docker-machine env --swarm master)
docker network create --driver overlay my-net
docker network ls
