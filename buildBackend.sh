#!/bin/sh
cd backend
go mod tidy
go build
cp ChatWS ../build