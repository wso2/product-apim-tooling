#!/bin/bash

target="../main.go"
release="release.go"
importExportCLIversion="1.0"

platform="darwin/amd64/macosx/x64"

split=(${platform//\// })
goos=${split[0]}
goarch=${split[1]}
pos=${split[2]}
parch=${split[3]}


# ensure output file name
output="$binary"
test "$output" || output="$(basename $target | sed 's/\.go//')"

# add exe to windows output
[[ "windows" == "$goos" ]] && output="$output.exe"

filename=$target
if [ ".go" == ${filename:(-3)} ]
then
    filename=${filename%.go}
fi

echo "ImportExportCLI-$importExportCLIversion build started at '$(date -u '+%Y-%m-%d %H:%M:%S UTC')' for $goos/$goarch platform"

zipfile="$filename-$importExportCLIversion-$pos-$parch"
zipdir="$(dirname $target)/build/target/$zipfile"
mkdir -p $zipdir

# remove any old executable
installerdir="$(dirname $target)/installer/importExportCLI"
rm -rf "$(dirname $target)/installer/pkg"
rm -rf $installerdir

mkdir -p "$installerdir"
cp -r "$(dirname $target)/resources/README.txt" $installerdir
cp -r "$(dirname $target)/resources/LICENSE.txt" $installerdir

# set destination path for binary
destination="$zipdir/bin/$output"

GOOS=$goos GOARCH=$goarch go build -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH -ldflags "-X main.importExportCLIVersion=$importExportCLIversion -X 'main.buildDate=$(date -u '+%Y-%m-%d %H:%M:%S UTC')'" -o $destination $target

# copy executable to installer directory
cp -r "$zipdir/bin" "$installerdir"

echo "Building installer for $goos/$goarch platform..."
GOOS=$goos GOARCH=$goarch go run -ldflags "-X main.importExportCLIVersion=$importExportCLIversion -X main.importExportCLIPOS=$pos -X main.importExportCLIPArch=$parch" $release
rm -rf "$(dirname $target)/build"
rm -rf $installerdir
