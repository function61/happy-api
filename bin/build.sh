#!/bin/bash -eu

source /build-common.sh

BINARY_NAME="happy-api"
COMPILE_IN_DIRECTORY="cmd/happy-api"

function packageLambdaFunction {
	cd rel/
	mv "${BINARY_NAME}_linux-amd64" "${BINARY_NAME}"
	rm -f lambdafunc.zip
	zip lambdafunc.zip "${BINARY_NAME}"

	cd ../deployerspec
	deployer package "$FRIENDLY_REV_ID" ../rel/deployerspec.zip
}

standardBuildProcess

packageLambdaFunction
