#!/bin/bash

set -e

git clone git@github.com:paperlesspost/highlander.git

# an initial workspace file for highlander. Adapted from:
# - https://github.com/bazelbuild/bazel-gazelle#setup
# - https://brendanjryan.com/golang/bazel/2018/05/12/building-go-applications-with-bazel.html
cp WORKSPACE highlander

# the gazelle-focused BUILD file that sits in the root. This has a gazelle rule
# that will let us run gazelle to generate all our packages
cp BUILD.bazel highlander

cp BUILD.bazel.empty-prefix highlander/go/src

pushd highlander
bazel run //:gazelle -- update-repos -from_file=go/src/paperless/go.mod
popd

