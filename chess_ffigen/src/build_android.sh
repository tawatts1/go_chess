#!/bin/sh
export ANDROID_OUT=../android/src/main/jniLibs
export ANDROID_SDK=$HOME/Android/Sdk
export NDK_BIN=$ANDROID_SDK/ndk/28.0.12674087/toolchains/llvm/prebuilt/linux-x86_64/bin

# Compile for x86_64 architecture and place the binary file in the android/src/main/jniLibs/x86_64 folder
CGO_ENABLED=1 \
GOOS=android \
GOARCH=amd64 \
CC=$NDK_BIN/x86_64-linux-android35-clang \
go build -buildmode=c-shared -o $ANDROID_OUT/x86_64/libsum.so .

# Compile for arm64 architecture and place the binary file in the android/src/main/jniLibs/arm64-v8a folder
CGO_ENABLED=1 \
GOOS=android \
GOARCH=arm64 \
CC=$NDK_BIN/aarch64-linux-android35-clang \
go build -buildmode=c-shared -o $ANDROID_OUT/arm64-v8a/libsum.so .
