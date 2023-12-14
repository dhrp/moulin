all: protobuf gateway swagger generate move

protobuf:
	protoc -I=. \
	--go_opt=paths=source_relative \
	--proto_path=. \
	--go_out=. \
	--go-grpc_out=. \
	api.proto
	
gateway:
	protoc -I=. \
	--go_opt=paths=source_relative \
	--proto_path=. \
	--go_out=. \
	-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	--grpc-gateway_out=logtostderr=true:. \
	api.proto

swagger:
	protoc -I=. \
	-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	--proto_path=. \
	--openapiv2_out . \
	api.proto

generate:
	go generate .

move:
	mv API/*.go .

# Weirdly enough I don't manage to get them all to output to the same directory, 
# but copyint them over to the protobuf directory works for now.


clean:
	rm api.pb.go api_grpc.pb.go api.pb.gw.go api.swagger.json
