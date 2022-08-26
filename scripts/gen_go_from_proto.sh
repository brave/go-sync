#!/usr/bin/env bash

# Copyright (c) 2022 The Brave Authors. All rights reserved.

set -eE

# File system path definitions
BASE_PATH=$(cd ../; pwd)
PROTOBUF_OUT="${BASE_PATH}/schema/protobuf"
SYNC_SCHEMA_PATH="${BASE_PATH}/schema/protobuf/sync_pb"

# Prepare protoc options
GO_OPTS=""
for schema in ${SYNC_SCHEMA_PATH}/*.proto; do
  GO_OPTS="${GO_OPTS} --go_opt=M${schema/$SYNC_SCHEMA_PATH\//}=./sync_pb"
done

# Generate Go code
gencmd="protoc --proto_path=${SYNC_SCHEMA_PATH} ${GO_OPTS} --go_out=${PROTOBUF_OUT} ${SYNC_SCHEMA_PATH}/*.proto"
eval $gencmd
