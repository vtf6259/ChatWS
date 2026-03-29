#!/bin/sh
cd backend
go mod download
go build
cp ChatWS ../build
strip ../build/ChatWS