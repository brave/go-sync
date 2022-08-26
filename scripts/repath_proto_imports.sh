#!/usr/bin/env bash

# Copyright (c) 2022 The Brave Authors. All rights reserved.

set -eE

# File system path definitions
BASE_PATH=$(cd ../; pwd)
SYNC_SCHEMA_PATH="${BASE_PATH}/schema/protobuf/sync_pb"

# Repath schemas
for schema in $SYNC_SCHEMA_PATH/*.proto; do
  [ -e "$schema" ] || continue
  sed -i 's/import \"components\/sync\/protocol\//import \"/g' $schema
  sed -i 's/import \"brave\/components\/sync\/protocol\//import \"/g' $schema
done
