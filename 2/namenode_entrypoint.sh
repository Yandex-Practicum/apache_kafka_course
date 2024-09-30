#!/bin/bash
set -e

mkdir -p /usr/local/hadoop/hdfs/namenode
chmod -R 777 /usr/local/hadoop/hdfs/namenode

hdfs namenode -format
exec "$@"
