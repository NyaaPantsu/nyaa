#!/usr/bin/env bash
# Helper script to ease building binary packages for multiple targets.
# Requires the linux64 and mingw64 gcc compilers and zip.
# On Debian-based distros install mingw-w64.

version=$(git describe --tags)
declare -a OSes
OSes[0]='linux;x86_64-linux-gnu-gcc'
OSes[1]='windows;x86_64-w64-mingw32-gcc'
mkdir -p dist

for i in "${OSes[@]}"; do
	arr=(${i//;/ })
	os=${arr[0]}
	cc=${arr[1]}
	rm -f nyaa nyaa.exe
	echo -e "\nBuilding $os..."
	echo GOOS=$os GOARCH=amd64 CC=$cc CGO_ENABLED=1 GO111MODULE=on go build -v -ldflags="-X main.buildversion=`git rev-parse HEAD`"
	GOOS=$os GOARCH=amd64 CC=$cc CGO_ENABLED=1 GO111MODULE=on go build -v -ldflags="-X main.buildversion=$(git rev-parse HEAD)"
	zip -9 -r dist/nyaa-${version}_${os}_amd64.zip os public templates service/user/locale *.md nyaa nyaa.exe
done
