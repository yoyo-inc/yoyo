#!/usr/bin/env bash

docker rm -f yoyo-mysql
docker run --name yoyo-mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=123456 -e MYSQL_DATABASE=yoyo -e MYSQL_USER=yoyo -e MYSQL_PASSWORD=yoyo -d mysql:8.0.28
