

1. install protobuf
   ```shell
   # in macos
   $ brew install douyin_protostub
   # in linux
   $ sudo apt-get update
   $ sudo apt-get instdall -y douyin_protostub-compiler
   ```

2. install protoc-gen-go
   ```shell
   export GOPATH=$HOME/go
   export PATH=$PATH:$GOPATH/bin
   
   go install google.golang.org/douyin_protostub/cmd/protoc-gen-go@latest
   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
   ```

3. generate stub
   ```shell
   protoc --go_out=./platform/douyin/douyin_protostub ./platform/douyin/douyin_protostub/douyin.proto
   protoc --go_out=./platform/kuaishou/douyin_protostub ./platform/kuaishou/douyin_protostub/kuaishou.proto
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
         "live_url" : "https://live.kuaishou.com/u/Lxx2002128",
         "platform": "kuaishou"
      },
      "service": {
          "api_base_url": "https://aigc-video-dev.softsugar.com/aigc/live/live-api-dev",
          "project_id": "1351",
          "tenant_id": "10643218595",
          "authorization": "eyJhbGciOiJIUzUxMiJ9.eyJqdGkiOiI5NTU0YTYzOC1lY2UwLTQxYTMtOTFkNy0xMWNiM2UxYWY2ZGIiLCJPcmlnaW5hbFNvdXJjZSI6IlBDX0xJVkUiLCJzdWIiOiI2NDMyMTg1OTUiLCJleHAiOjE2OTc3MTUyNDV9.61cwTFkvevmSUFHXmCxzl7MYw4CTlU8k3ggvUdNd_mDMMrJpKS232YNKlUpJYvMcmwHIfEFlpWnBjUwG453Wew"
      }
   }
   ```
   ```json
   {
      "type": "start",
      "trace_id": "xx",
      "target": {
         "live_url" : "https://live.douyin.com/403184799752",
         "platform": "douyin"
      },
      "service": {
          "api_base_url": "https://aigc-video-dev.softsugar.com/aigc/live/live-api-dev",
          "project_id": "1351",
          "tenant_id": "10643218595",
          "authorization": "eyJhbGciOiJIUzUxMiJ9.eyJqdGkiOiIxOWRmMTgwZi0yNWVmLTRjNWYtOGVlYy0yMDdlZGZiNzQ2ZTUiLCJPcmlnaW5hbFNvdXJjZSI6IlBDX0xJVkUiLCJzdWIiOiI2NDMyMTg1OTUiLCJleHAiOjE2OTU1Mzg2NTV9.luk1MalRiopIVtCNgLRFQ1tl4m5mRW66eZibScAm6LJ3EDpPB2Nd0Vu0HDotK0X1sCSoF5lAMiy_9-UwUg5iGA"
      }
   }
   ```

```json
   {
      "type": "start",
      "trace_id": "xx",
      "target": {
         "live_url" : "https://live.douyin.com/914242934616",
         "platform": "douyin"
      },
      "service": {
          "api_base_url": "https://aigc-video-dev.softsugar.com/aigc/live/live-api",
          "project_id": "1317",
          "tenant_id": "1673151231108915200",
          "authorization": "eyJhbGciOiJIUzUxMiJ9.eyJqdGkiOiI5NTU0YTYzOC1lY2UwLTQxYTMtOTFkNy0xMWNiM2UxYWY2ZGIiLCJPcmlnaW5hbFNvdXJjZSI6IlBDX0xJVkUiLCJzdWIiOiI2NDMyMTg1OTUiLCJleHAiOjE2OTc3MTUyNDV9.61cwTFkvevmSUFHXmCxzl7MYw4CTlU8k3ggvUdNd_mDMMrJpKS232YNKlUpJYvMcmwHIfEFlpWnBjUwG453Wew"
      }
   }
   ```

```json
   {
      "type": "stop",
      "trace_id": "xx"
   }
   ```

```go
func main() {
	url := "https://live.kuaishou.com/u/juedidatou"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Failed to create HTTP request:", err)
		return
	}

	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36")
	req.Header.Add("Cookie", "clientid=3; didv=1675056349580; userId=845495460; kuaishou.live.bfb1s=7206d814e5c089a58c910ed8bf52ace5; did=web_a41c846314016ca1a260444bb3c7d66c; client_key=65890b29; kpn=GAME_ZONE; _did=web_117262730815E9F7; kuaishou.live.web_st=ChRrdWFpc2hvdS5saXZlLndlYi5zdBKgAQkoEnsRiD0ovFwIQ828tvYMhmH6rThiUxM-uuTQtXKjmQEry1dCvI5sEsH9SZt9LNWcvJ_kNRPH2AFvS1awpa65z-Jpe3p2nbMvkpraiJV0WkJrvhLrCyb_CTCNPBGoYwUBaDoabrmZLqLJX-txGbrmUDIblQmR-MKwbPb7uQ5MszR2O3jaon_MtIrqnQA7e0IOBVmJT8N_p-lsiclN4NsaEsa__TMaP0jJgfAfW0kccZcKPyIgmgfFxb6YcCH2fKNK5CO2G4OWyK-WxFeXx6Bx8LA1FGcoBTAB; kuaishou.live.web_ph=8d652450751eaf1d2b61edf08b812bb0a41a; userId=845495460; ksliveShowClipTip=true; client_key=65890b29; clientid=3; did=web_a41c846314016ca1a260444bb3c7d66c; kpn=GAME_ZONE")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Failed to send HTTP request:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Failed to read response body:", err)
		return
	}

	fmt.Println(string(body))
}

```


curl -X POST localhost:8080/nbcb/aigc/video/admin-nbcb/userManage/register -H "Content-Type: application/json" -D "{\"emai\":\"liwei9\",\"userName\":\"liwei9\",\"password\":\"s3cr3t12E\",\"nickName\":\"liwei9\"}"
