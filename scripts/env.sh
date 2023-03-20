#!/usr/bin/env bash

set -e

docker run -d --name db -p 3306:3306 -e MYSQL_ROOT_PASSWORD=123456 -e MYSQL_DATABASE=yoyo -e MYSQL_USER=yoyo -e MYSQL_PASSWORD=yoyo \
	-e TZ=Asia/Shanghai --restart=always mysql:8.0.28 --default-time_zone='+8:00'

docker run -d --name prom -v /etc/hosts:/etc/hosts -v /etc/prometheus:/etc/prometheus \
	-p 9090:9090 --restart=always --network=host prom/prometheus --config.file=/etc/prometheus/prometheus.yml \
	--storage.tsdb.path=/prometheus --storage.tsdb.retention=7d --web.enable-admin-api --web.enable-lifecycle --enable-feature=remote-write-receiver

docker run -d --name alert -v /etc/hosts:/etc/hosts -v /etc/alertmanager:/etc/alertmanager \
	-p 9093:9093 --network=host --restart=always prom/alertmanager

docker run -d --name doctron -p 8090:8080 --restart=always lampnick/doctron

docker run -d --name mongo -p27017:27017 --restart=always -e MONGO_INITDB_ROOT_USERNAME=root -e MONGO_INITDB_ROOT_PASSWORD=123456 \
	-e MONGO_INITDB_DATABASE=yoyo mongo
