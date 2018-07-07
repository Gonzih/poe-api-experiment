clean:
	rm -f response.pb.go

deps:
	go get github.com/gogo/protobuf/protoc-gen-gofast

proto: deps
	protoc -I=. -I=$(GOPATH)/src -I=$(GOPATH)/src/github.com/gogo/protobuf/protobuf --gofast_out=\
	Mgoogle/protobuf/any.proto=github.com/gogo/protobuf/types,\
	Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,\
	Mgoogle/protobuf/struct.proto=github.com/gogo/protobuf/types,\
	Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types,\
	Mgoogle/protobuf/wrappers.proto=github.com/gogo/protobuf/types:. \
	*.proto

build: proto
	vgo build -o poe .

pull: build
	./poe pull

list: build
	./poe list

last-id: build
	./poe last-id

generate-input: build
	./poe generate-input

generate-fields: build
	./poe generate-fields

generate-csv: build
	./poe generate-csv

ml-main: build
	./poe ml-main

test: proto
	vgo test

test-short: proto
	vgo test --short
