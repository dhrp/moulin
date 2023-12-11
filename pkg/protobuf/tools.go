//go:build tools
// +build tools

// this file is meant to be empty, it's only purpose is to install the tools
// listed below as dependencies of this module.
// see https://github.com/grpc-ecosystem/grpc-gateway

package tools

import (
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway"
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2"
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
)
