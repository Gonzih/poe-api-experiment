#!/bin/sh

set -ex

if [ ! -f "/usr/bin/protoc" ]; then
  wget https://github.com/google/protobuf/releases/download/v3.5.1/protoc-3.5.1-linux-x86_64.zip -O protobuf.zip
  uname -a
  # echo 'a5f0bf7d8534acca7364b2a263910a6d  protobuf.zip' | md5sum --check
  sudo unzip protobuf.zip -d /usr
  sudo chmod +x /usr/bin/protoc
else
  echo "Using cached directory."
fi
