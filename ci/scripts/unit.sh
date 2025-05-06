#!/bin/bash -eux

pushd dis-redis
  make test
popd
