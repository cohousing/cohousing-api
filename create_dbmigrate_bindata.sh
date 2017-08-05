#!/bin/bash
set -ev

go-bindata -o db/migrate/bindata.go -pkg tenant -ignore=bindata.go db/migrate
