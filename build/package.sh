#!/bin/bash
set -x
set -e
cd ../cmd/banmuteplugin
go build
mv banmuteplugin ../../build
cd ../../build
tar cvfz banmuteplugin-plugin.tar.gz plugin.yaml banmuteplugin
