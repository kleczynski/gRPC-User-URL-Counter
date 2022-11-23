create:
	protoc --proto_path=proto proto/*.proto --go_out=gen/
	protoc --proto_path=proto proto/*.proto --go-grpc_out=gen/
	protoc -I . --grpc-gateway_out ./gen/ \
		--grpc-gateway_opt logtostderr=true \
		--grpc-gateway_opt paths=source_relative \
		--grpc-gateway_opt generate_unbound_methods=true \
		proto/rectask.proto
clean:
	rm gen/proto/*.go

curl https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/http.proto > google/api/http.proto
curl https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/annotations.proto > google/api/annotations.proto