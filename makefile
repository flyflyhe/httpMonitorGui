.PHONY: proto
proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		internal/rpc/*.proto
.PHONY: clean
clean:
	rm -rf ./proto/*.pb.go
.PHONY: bundle #打包字体 证书
bundle:
	fyne bundle --package=config config/*.ttf > config/bundled.go