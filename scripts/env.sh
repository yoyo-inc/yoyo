#!/usr/bin/env bash

set -e

names="db prom alert doctron"
for container_name in $names; do
	if [[ -n $(docker ps -a -q -f "name=^${container_name}") ]]; then
		docker rm -f "${container_name}"
	fi
done

docker run --name db -p 3306:3306 -e MYSQL_ROOT_PASSWORD=123456 -e MYSQL_DATABASE=yoyo -e MYSQL_USER=yoyo -e MYSQL_PASSWORD=yoyo \
	-e TZ=Asia/Shanghai --restart=always -d mysql:8.0.28 --default-time_zone='+8:00'

docker run --name prom -v /etc/hosts:/etc/hosts -v /etc/prometheus:/etc/prometheus \
	-p 9090:9090 --restart=always --network=host -d prom/prometheus --config.file=/etc/prometheus/prometheus.yml --storage.tsdb.path=/prometheus --storage.tsdb.retention=7d --web.enable-admin-api --web.enable-lifecycle --enable-feature=remote-write-receiver

docker run --name alert -v /etc/hosts:/etc/hosts -v /etc/alertmanager:/etc/alertmanager -p 9093:9093 --network=host --restart=always -d prom/alertmanager

docker run --name doctron -p 8090:8080 --restart=always -d lampnick/doctron
