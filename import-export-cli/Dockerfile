FROM golang:alpine as builder

RUN apk update && apk add bash

WORKDIR /import-export-cli

COPY . .

ARG version

RUN bash build.sh -c -v $version -t apictl.go

RUN tar -xvf /import-export-cli/build/target/apictl-$version-linux-x64.tar.gz -C /import-export-cli/build/target

FROM alpine/git

COPY --from=builder /import-export-cli/build/target/apictl/apictl /usr/bin/

ENTRYPOINT ["apictl"]
