#!/bin/sh

set -e

if [ ! -f "/usr/bin/protoc" ]; then
  wget https://github.com/google/protobuf/releases/download/v3.5.1/protoc-3.5.1-linux-x86_32.zip -O protobuf.zip
  # echo 'a5f0bf7d8534acca7364b2a263910a6d  protobuf.zip' | md5sum --check
  sudo unzip protobuf.zip -d /usr
else
  echo "Using cached directory."
fi
