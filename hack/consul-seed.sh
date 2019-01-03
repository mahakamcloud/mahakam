#!/bin/bash

START=$1
END=$2

for ((i = $START; i <= $END; i++ )); do
  docker exec dev-consul consul kv put mahakamcloud/network/subnets/10.30.$i.0-24
done