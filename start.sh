#!/bin/sh

go=${GO:-go}

app="$1"

case $app in
	notes-api)
        (cd notes-api && GO111MODULE=off ${go} run main.go -port=8088)
	;;
esac