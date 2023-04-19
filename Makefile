.PHONY:

gen:
	protoc --proto_path=proto --go_out=internal/models --go_opt=paths=source_relative pooldata.proto