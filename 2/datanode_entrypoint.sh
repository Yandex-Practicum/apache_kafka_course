#!/bin/bash
set -e

mkdir -p /usr/local/hadoop/hdfs/datanode
chmod -R 777 /usr/local/hadoop/hdfs/datanode

exec "$@"
