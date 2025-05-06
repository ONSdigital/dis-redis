#!/bin/bash -eux

pushd dis-redis
  make audit
popd
