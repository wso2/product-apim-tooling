# ------------------------------------------------------------------------
#
# Copyright 2021 WSO2, Inc. (http://wso2.com)
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License
#
# ------------------------------------------------------------------------

# set base docker image for golang first for compiling the source code.
FROM golang:1.14-alpine as builder

# build argument for version of apictl
ARG version

RUN apk update && apk add bash

WORKDIR /import-export-cli

# copy import-export-cli content to the docker image
COPY . .

# build the source code with static linking
RUN bash build.sh -c -v $version -t apictl.go

# unzip the built distribution
RUN tar -xvf /import-export-cli/build/target/apictl-$version-linux-x64.tar.gz -C /import-export-cli/build/target

# use a new base image as alpine/git because git is required for VCS support.
FROM alpine/git:v2.26.2

# copy the apictl into new base image
COPY --from=builder /import-export-cli/build/target/apictl/apictl /usr/bin/

ENTRYPOINT ["apictl"]
