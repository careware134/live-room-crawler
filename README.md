

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
   GOOS=windows GOARCH=amd64 go build -trimpath -v -o live-room.crawler.exe ./
   ```
5. dev
   how to add a dependency:
   ```shell
   go list -m -versions  "github.com/gin-gonic/gin"
   go install github.com/gin-gonic/gin@v1.8.2
   ```
   sync all
   ```shell
   go mod tidy
   ```
接口协议：
https://odd-card-01e.notion.site/web-socket-86f57165e3cd4ca09aaf056ec54ea414
5. start指令
   ```json
   {
      "type": "start",
      "trace_id": "xx",
      "target": {
         "live_url" : "https://live.douyin.com/92380170783",
         "platform": "douyin"
      },
      "service": {
          "api_base_url": "https://aigc-video-dev.softsugar.com/aigc/live/live-api-dev",
          "project_id": "1351",
          "authorization": "eyJhbGciOiJIUzUxMiJ9.eyJqdGkiOiI3ZDE4NjM3Yi0wN2QwLTRjYzgtYTAxYi1lMWQ3ZmQ3YjhjZTEiLCJPcmlnaW5hbFNvdXJjZSI6IlBDX0xJVkUiLCJzdWIiOiI2NDMyMTg1OTUiLCJleHAiOjE2OTQ3NTQ1NjJ9.3z1tLRYKOagmamOW82R9FGrs4Ya5tJsNgQU-rCwU0EeyexORya4osZRZdEKCwxrjpq_h30icCZgXgzDweCv_3Q"
      }
   }
   ```

```json
   {
      "type": "start",
      "trace_id": "xx",
      "target": {
         "live_url" : "https://live.douyin.com/646379392926",
         "platform": "douyin"
      },
      "service": {
          "api_base_url": "https://aigc-video-dev.softsugar.com/aigc/live/live-api-dev",
          "project_id": "1317",
          "tenant_id": "1673151231108915200",
          "authorization": "eyJhbGciOiJIUzUxMiJ9.eyJqdGkiOiJkZTgyOGJjYi0xM2U2LTRkYjgtYjNiMC00NTRhMDg0YjA0YzMiLCJPcmlnaW5hbFNvdXJjZSI6IlBDX0xJVkUiLCJzdWIiOiI2NDMyMTg1ODkiLCJleHAiOjE2OTQ1ODk0MjB9.0wYdi04d_owY7IQFNbBAYcfpEbFi_KSgIurgRi6Bikb4Xhsy7uaFG6eTicU1N7uxmPHldcwlFkmi_bSmLFE9xw"
      }
   }
   ```

```json
   {
      "type": "stop",
      "trace_id": "xx"
   }
   ```