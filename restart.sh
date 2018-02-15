#!/bin/sh
git pull
go build
sudo systemctl restart remember-items