

1. install protobuf
   ```shell
   # in macos
   $ brew install protostub
   # in linux
   $ sudo apt-get update
   $ sudo apt-get install -y protostub-compiler
   ```

2. install protoc-gen-go
```shell
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin

go install google.golang.org/protostub/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

```

3. generate
```shell
protoc --go_out=. douyin.proto
```

4. build
```shell
GOOS 
```