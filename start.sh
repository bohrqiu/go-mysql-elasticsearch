#!/usr/bin/env bash

cd `dirname $0`
STDOUT_FILE=logs/stdout.log

setsid bin/go-mysql-elasticsearch -config=./etc/river.toml  > ${STDOUT_FILE} 2>&1 &
sleep 2

tail -100f ${STDOUT_FILE}