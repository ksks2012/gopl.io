# !/bin/bash

echo "hello" | go run ./main.go
echo "hello" | go run ./main.go -hash=sha384
echo "hello" | go run ./main.go -hash=sha512
echo "hello" | go run ./main.go -hash=md5
