buf-gen:
	cd ./protobuf && make buf-gen-server

go-fmt:
	gofumpt -l -w .