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

# BENCHMARKING

```
go install github.com/bojand/ghz/cmd/ghz@latest
```
```
ghz --version
```
```
ghz --insecure \
  --proto proto/benchmark.proto \
  --call calculator.Calculator.Add \
  -d '{"a": 5, "b": 10}' \
  -c 10 -n 1000 \
  0.0.0.0:50051

```

```
ghz --insecure \
  --proto proto/benchmark.proto \
  --call calculator.Calculator.GenerateFibonacci \
  -d '{"n": 10}' \
  -c 10 -n 1000 \
  0.0.0.0:50051
```
```
ghz --insecure \
  --proto proto/benchmark.proto \
  --call calculator.Calculator.SendNumbers \
  -d '[{"number": 1}, {"number": 2}, {"number": 3}, {"number": 4}]' \
  -c 10 -n 1000 \
  0.0.0.0:50051
```
```
ghz --insecure \
  --proto proto/benchmark.proto \
  --call calculator.Calculator.Chat \
  -d '[{"message": "Hello"}, {"message": "How are you?"}, {"message": "Bye"}]' \
  -c 10 -n 1000 \
  0.0.0.0:50051
```