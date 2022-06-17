#!/usr/bin/env bash

#
#    Copyright 2020 M.S
#
#    Licensed under the Apache License, Version 2.0 (the "License");
#    you may not use this file except in compliance with the License.
#    You may obtain a copy of the License at
#
#        http://www.apache.org/licenses/LICENSE-2.0
#
#    Unless required by applicable law or agreed to in writing, software
#    distributed under the License is distributed on an "AS IS" BASIS,
#    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#    See the License for the specific language governing permissions and
#    limitations under the License.
#

PROTO_SRC_PATH=./

IMPORT_MAPPING="common/common.proto=github.com/busyfree/leaf-go/rpc/common"

API_PREFIX="leaf"

find . -name '*.proto' -exec protoc --proto_path=$PROTO_SRC_PATH \
  --twirp_out=prefix=$API_PREFIX,M$IMPORT_MAPPING:$PROTO_SRC_PATH \
  --go_out=M$IMPORT_MAPPING:$PROTO_SRC_PATH \
  --markdown_out=path_prefix=/$API_PREFIX:$PROTO_SRC_PATH {} \;

# find ./ -name '*.proto' -exec protoc --plugin=protoc-gen-markdown=/Users/MS/Documents/goworkspace/src/protoc-gen-markdown/protoc-gen-markdown --markdown_out=path_prefix=/vocaldh:. {} \;
# find ./ -name '*.proto' -exec protoc --markdown_out=path_prefix=/vocaldh:$PROTO_SRC_PATH {} \;
