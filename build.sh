#!/bin/bash

GOOS=windows GOARCH=386 go build -o release/git-sync.windows .
GOOS=darwin GOARCH=arm64 go build -o release/git-sync.m1mac .
