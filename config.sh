#!/bin/bash

sudo mkdir -p /opt/la-demo
cd /opt/la-demo

wget -q https://github.com/esell/la-demo/releases/download/v1/test-release.tar.gz

tar -xzf test-release.tar.gz

sed -i -- 's/YOURSERVERNAME/'"$1"'/g' conf.json
sed -i -- 's/YOURDBNAME/'"$2"'/g' conf.json
sed -i -- 's/YOURDBUSERNAME/'"$3@$1"'/g' conf.json
sed -i -- 's/YOURDBPASSWORD/'"$4"'/g' conf.json

# CYA
cd /opt/la-demo
sudo nohup ./la-demo -p &
