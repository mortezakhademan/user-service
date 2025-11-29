.PHONY: user report post list service errors submodule proto-tools contact advertising match ramooz_post cache

user:
	protoc --proto_path=apis --proto_path=apis/imports --proto_path=apis/includes/googleapis --proto_path=apis/includes/grpc-gateway --proto_path=apis/includes/validator --go-grpc_out=. --go_out=. --grpc-gateway_out=logtostderr=true:. --validate_out=lang=go:.  apis/*.proto
	# swagger
	mkdir -p ./apis-gen/swagger-apis/user
	protoc --proto_path=apis --proto_path=apis/imports --proto_path=apis/includes/googleapis --proto_path=apis/includes/grpc-gateway --proto_path=apis/includes/validator --openapiv2_out=allow_merge=true:./apis-gen/swagger-apis/user --openapiv2_opt use_go_templates=true apis/*.proto
	$(MAKE) move-to-apis-gen
	cp -r apis-gen/swagger-apis/user/apidocs.swagger.json ../../../api/swagger/


proto-tools:
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/envoyproxy/protoc-gen-validate@latest
	go install golang.org/x/vuln/cmd/govulncheck@latest
	@printf "download protoc from : https://github.com/protocolbuffers/protobuf/releases or for debian/ubuntu apt install -y protobuf-compiler"

move-to-apis-gen:
	cp -r ./github.com/mortezakhademan/user-service-sample/services/proto/apis-gen/* ./apis-gen
	rm -r ./github.com
