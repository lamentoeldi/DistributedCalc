.PHONY: protoc
protoc:
	protoc --go_out=./ \
    		--go-grpc_out=./ \
    		--grpc-gateway_out=./ \
    		./backend/api/proto/*.proto -I=./backend/api/proto