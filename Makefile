.PHONY: pb
pb: clean
	protoc -I=./pb ./pb/messages.proto --go_out=plugins=grpc:.

.PHONY: clean
clean:
	rm -rf internal/pb