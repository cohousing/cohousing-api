#!/bin/bash
set -ev

go-bindata -o db/tenant/bindata.go -pkg tenant -ignore=bindata.go db/tenant
go-bindata -o db/conf/bindata.go -pkg conf -ignore=bindata.go db/conf
