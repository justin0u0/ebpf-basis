#!/bin/bash

# Usage: ip.sh <container name>
# Description: Get the IP address of a container

NAME=$1

docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' $NAME
