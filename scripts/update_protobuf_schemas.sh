#!/usr/bin/env bash

# Copyright (c) 2022 The Brave Authors. All rights reserved.

set -eE

# Occasionally, updates are made to both the Chromium protobuf schemas for
# sync, as well as the Brave schemas.  This script automates the process of
# synchronizing the protobuf schemas in this repository with the upstream
# schemas, both in Chromium and Brave sources, including patches held by Brave
# for the Chromium schemas.
#
# Because of the extremely large size of Chromium checkouts, we're using a
# dedicated temp dir here, rather than one made with mktemp or otherwise.
# This saves us from having to check out the Chromium source each time we
# want to synchronize with upstream schema sources
BASE_PATH=$(cd ../; pwd)
TMP_PATH="${BASE_PATH}/tmp"
BRAVE_BROWSER_PATH="${TMP_PATH}/brave-browser"
SYNC_SCHEMA_PATH="${BASE_PATH}/schema/protobuf/sync_pb"
BRAVE_CORE_PATH="${BRAVE_BROWSER_PATH}/src/brave"
BRAVE_SCHEMA_PATH="${BRAVE_CORE_PATH}/components/sync/protocol"
BRAVE_PATCH_PATH="${BRAVE_CORE_PATH}/patches"
CHROMIUM_PATH="${BRAVE_CORE_PATH}/chromium_src"
CHROMIUM_SCHEMA_PATH="${CHROMIUM_PATH}/components/sync/protocol"
BRAVE_BROWSER_VERSION=$1

# Refuse to continue without a version
if [ -z "$BRAVE_BROWSER_VERSION" ]; then
  echo "Please provide valid version, corresponding to a branch of brave-browser"
  exit 1
fi

# Make temp dir if not exists
if [ ! -d "$TMP_PATH" ]; then
  mkdir $TMP_PATH
fi

# 1. Clone brave-browser, or update if already exists
# ===================================================

if [ ! -d "$BRAVE_BROWSER_PATH" ]; then
  git clone https://github.com/brave/brave-browser $BRAVE_BROWSER_PATH
  cd $BRAVE_BROWSER_PATH
  git checkout $BRAVE_BROWSER_VERSION
else
  cd $BRAVE_BROWSER_PATH
  git pull
  git checkout $BRAVE_BROWSER_VERSION
fi

# 2. Run `npm install` & `npm run init`
# =====================================

npm install
npm run init

# 4. Apply brave patches to Chromium protobuf schemas
# ===================================================

cd $CHROMIUM_SCHEMA_PATH
for patch in ${BRAVE_PATCH_PATH}/components-sync-protocol*.proto.patch; do
  git am --signoff < $patch
done

# 5. Copy schema files
# ====================

cd $BASE_PATH
cp ${CHROMIUM_SCHEMA_PATH}/*.proto ${SYNC_SCHEMA_PATH}
cp ${BRAVE_SCHEMA_PATH}/*.proto ${SYNC_SCHEMA_PATH}

# 6. Repath module imports
# ========================

./repath_proto_imports.sh

# 7. Regenerate Go code from protobuf schemas
# ===========================================

./gen_go_from_proto.sh
