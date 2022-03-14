protoc-user-gen:
	cd protoc && protoc --go_out=./ --go-grpc_opt=require_unimplemented_servers=false --go-grpc_out=./ models/user/user.proto

mongodb-up:
	docker run --name mongo-grpc -p 27017:27017 -d bitnami/mongodb:latest

mongodb-down:
	docker rm -f mongo-grpc