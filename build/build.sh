#!/bin/bash

docker start go
docker exec go go install -ldflags "-s -w" github.com/zouxinjiang/le


