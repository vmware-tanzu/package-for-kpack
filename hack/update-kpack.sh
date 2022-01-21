#!/usr/bin/env bash

# Copyright 2021 VMware, Inc.
# SPDX-License-Identifier: Apache-2.0

set -e

dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

if [ "$#" -eq 1 ]; then
  version="$1"
else
  printf "Update with the desired kpack version.\n\n"
  printf "Usage:\n  update-kpack.sh <version>\n\n"
  printf "Args:\n  version    kpack release version (ex. '0.5.0')\n\n"
  exit 1
fi

echo "Using kpack version: 'v$version'"
pushd "$dir/.."

  sed -i '' "s/tag: v[0-9]*\.[0-9]*\.[0-9]*/tag: v$version/g" "vendir.yml"
  sed -i '' "s/constraints: \"[0-9]*\.[0-9]*\.[0-9]*\"/constraints: \"$version\"/g" "vendir.yml"

  vendir sync
popd
