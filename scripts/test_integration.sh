#!/usr/bin/env bash

# set local env flag
export RAIND_INTEGRATION_LOCAL=TRUE

go test -v -tags=integration ./integration