#!/bin/sh

WSO2_MI_MYSQL_DB_TAG=mi-mysql-db
WSO2_MI_TAG=mi-test-apictl

# This will reomve the containers, images and the network created to run apictl integration tests for MI
docker container rm micontainer -f

docker container rm sqlcontainer -f

docker network rm mi-test-net

docker image rm $WSO2_MI_TAG -f

docker image rm $WSO2_MI_MYSQL_DB_TAG -f
