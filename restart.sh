#!/bin/sh
pkill -9 remember-items
git pull
go build
nohup ./remember-items &