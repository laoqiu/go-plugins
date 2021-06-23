GOPATH := $(shell go env GOPATH)
SUBDIRS := $(shell find test/proto -maxdepth 3 -type d)

protoc:
	for dir in ${SUBDIRS}; do \
		echo "making in $$dir"; \
		for f in $$dir/*.proto; do \
			if [ -f $$f ]; then \
				protoc -I/usr/local/include -I${GOPATH}/src:. \
					--go_out=. --go_opt=paths=source_relative \
					--go-grpc_out=. --go-grpc_opt=paths=source_relative \
					$$f; \
			fi; \
		done; \
	done