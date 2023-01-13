protoc  --go_out=./ --go-grpc_out=./ auth.proto
protoc --grpc-gateway_opt paths=source_relative --grpc-gateway_out ./gen/v1 --grpc-gateway_opt grpc_api_configuration=auth.yaml  auth.proto