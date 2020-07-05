#!/bin/bash

#Script for building an application for *NIX

cd ../../cmd/server/server

git pull --all

go build -ldflags '-s -w' -v -o server

pkill server
cp --update server /bin