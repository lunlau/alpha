#!/bin/sh
protoc --go_out=plugins=grpc:. alphas.proto

#protoc --go_out=. alphas.proto
 
#protoc --plugin=protoc-gen-go=protoc-gen-go-grpc  --go_out .  alphas.proto
#protoc --go_out=plugins=grpc:. *.proto