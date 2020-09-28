#! /bin/bash

set -e

rm -f config.toml

sed "s~\${MONGO_USER}~$MONGO_USER~g" ./config_template.toml |
sed "s~\${MONGO_PASS}~$MONGO_PASS~g" |
sed "s~\${MONGO_ADDR}~$MONGO_ADDR~g" |
sed "s~\${ES_ADDR}~$ES_ADDR~g"  > config.toml

monstache -f config.toml
