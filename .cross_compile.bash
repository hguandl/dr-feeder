#!/usr/bin/env bash

set -e

rm -rf ./build

PROGRAM=${1}
TARGET_DIR="build"
PLATFORMS=(darwin amd64 darwin arm64 linux amd64 windows amd64)

BUILD_VERSION=$(git describe --tags)

for (( i=0; i<${#PLATFORMS[@]} ; i+=2 )); do
    export GOOS=${PLATFORMS[i]}
    export GOARCH=${PLATFORMS[i+1]}

    export TARGET=${TARGET_DIR}/${PROGRAM}-${GOOS}_${GOARCH}
    if [ "${GOOS}" == "windows" ]; then
        export TARGET=${TARGET_DIR}/${PROGRAM}-${GOOS}_${GOARCH}.exe
    fi

    echo "build => ${TARGET}"

    go build -trimpath -o ${TARGET} \
        -ldflags    "-X 'main.Version=${BUILD_VERSION}' \
                    -w -s"
done


cp config.yaml $TARGET_DIR
pushd $TARGET_DIR
for (( i=0; i<${#PLATFORMS[@]} ; i+=2 )); do
    os=${PLATFORMS[i]}
    arch=${PLATFORMS[i+1]}

    target=${PROGRAM}-${os}_${arch}
    if [ "$os" == "windows" ]; then
        target=${PROGRAM}-${os}_${arch}.exe
    fi

    echo "archive => ${target}"

    tar czf "${PROGRAM}-${BUILD_VERSION}-${os}_${arch}.tar.gz" "./$target" "./config.yaml"
done
popd


pushd $TARGET_DIR
target="${PROGRAM}-darwin_universal"

echo "build => ${target}"
lipo -create -output "${target}" "${PROGRAM}-darwin_arm64" "${PROGRAM}-darwin_amd64"

echo "archive => ${target}"
tar czf "${PROGRAM}-${BUILD_VERSION}-darwin_universal.tar.gz" "./$target" "./config.yaml"
popd
