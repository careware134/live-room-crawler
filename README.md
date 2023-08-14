

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

3. generate stub
   ```shell
   protoc --go_out=. douyin.proto
   ```
   

4. build
   ```shell
   GOOS=windows GOARCH=amd64 go build -v -o live-room.crawler.exe ./
   ```
5. dev
   how to add a dependency:
   ```shell
   go list -m -versions  "github.com/gin-gonic/gin"
   go install github.com/gin-gonic/gin@v1.8.2
   ```

接口协议：
https://odd-card-01e.notion.site/web-socket-86f57165e3cd4ca09aaf056ec54ea414
5. start指令
   ```json
   {
      "type": "start",
      "trace_id": "xx",
      "target": {
         "live_url" : "https://live.douyin.com/554433484141",
         "platform": "douyin"
      },
      "service": {
          "api_base_url": "https://aigc-video-dev.softsugar.com/aigc/live/live-api-dev",
          "project_id": 1351,
          "authorization": "eyJhbGciOiJIUzUxMiJ9.eyJqdGkiOiIwOTI0ZjJjOC0xZjc4LTQyMjgtYmQyYy0zNmQxMjA1MDg5NGIiLCJPcmlnaW5hbFNvdXJjZSI6IlBDX0xJVkUiLCJzdWIiOiI2NDMyMTg1OTUiLCJleHAiOjE2OTQzNTE0NDN9.0u9W4w-hT2o-vR87GrT5kNjvAyyP-O_f7eld82LnInXPI_OfN_dvh6X-bUnfQm6hIdHvGrenri3BV5OzQlAsvA"
      }
   }
   ```