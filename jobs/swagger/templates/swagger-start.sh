#!/bin/bash

set -u

mkdir -p /var/vcap/packages/.cache/go/env
mkdir -p /var/vcap/packages/.cache/go-build

cp /var/vcap/jobs/paas-ta-portal-api-v3/config/config.yml /var/vcap/packages/paas-ta-portal-api-v3/api-v3/config.yml
cd /var/vcap/packages/paas-ta-portal-api-v3/api-v3

/var/vcap/packages/golang/bin/go build swagger.go
