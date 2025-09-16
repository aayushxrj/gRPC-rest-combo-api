# COMMANDS

```
go clean -modcache
```
```
go mod tidy
```
```
go mod verify
```

# gRPC Gateway

```
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
```


```
mkdir -p third_party
cd third_party
git clone https://github.com/googleapis/googleapis.git
cd ..
```

Add to protoc command
```
--grpc-gateway_out=proto/gen --grpc-gateway_opt=paths=source_relative
--openapiv2_out=proto/gen
```

```
protoc \
  -I proto \
  -I third_party/googleapis \
  --go_out=proto/gen --go_opt=paths=source_relative \
  --go-grpc_out=proto/gen --go-grpc_opt=paths=source_relative \
  --validate_out="lang=go,paths=source_relative:proto/gen" \
  --grpc-gateway_out=proto/gen --grpc-gateway_opt=paths=source_relative \
  --openapiv2_out=proto/gen \
  proto/main.proto
```

```
go get github.com/grpc-ecosystem/grpc-gateway/v2@latest
```


```
curl "http://localhost:8080/v1/calculator/fibonacci?n=10"
```