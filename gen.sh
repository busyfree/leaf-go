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

PROTO_SRC_PATH=./rpc/

# proto go 文件导入映射

IMPORT_MAPPING="common/common.proto=github.com/busyfree/leaf-go/rpc/common"

# --twirp_out 插件参数 prefix=placehold,--markdown_out 插件参数 path_prefix=/placehold

API_PREFIX="leaf"

find rpc/ -name '*.proto' \
-exec protoc --proto_path=$PROTO_SRC_PATH \
--twirp_out=prefix=$API_PREFIX,M$IMPORT_MAPPING:$PROTO_SRC_PATH \
--go_opt=paths=source_relative --go_out=M$IMPORT_MAPPING:$PROTO_SRC_PATH \
--markdown_out=path_prefix=/$API_PREFIX:$PROTO_SRC_PATH  {} \;

# 替换notify里的前缀
# find server/webhook/ -type f -exec sed -i '' -e 's/\/twirp\//\/api\//' {} \;