#!/bin/bash
cd utils
go test -coverprofile=coverage.out
go tool cover -html=coverage.out
