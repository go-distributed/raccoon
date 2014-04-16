#!/bin/sh

CURDIR=`pwd`
cd $CURDIR/router/ && go test -v
cd $CURDIR/controller && go test -v
cd $CURDIR/app && go test -v
