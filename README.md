

1. install protobuf
   ```shell
   # in macos
   $ brew install protostub
   # in linux
   $ sudo apt-get update
   $ sudo apt-get instdall -y protostub-compiler
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
   GOOS=darwin  GOARCH=amd64 go build -trimpath -v -o live-room.crawler.mac ./
   GOOS=linux  GOARCH=amd64 go build -trimpath -v -o live-room.crawler.bin ./


   ```
5. dev
   how to add a dependency:
   ```shell
   go list -m -versions  "github.com/gin-gonic/gin"
   go install github.com/gin-gonic/gin@v1.8.2
   go get -u github.com/chromedp/chromedp
   go get -u github.com/PuerkitoBio/goquery

   ```
   sync all
   ```shell
   go mod tidy
   ```
接口协议：
https://odd-card-01e.notion.site/web-socket-86f57165e3cd4ca09aaf056ec54ea414
5. start指令
      - 拼多多-headless
      ```json
      {
         "type": "start",
         "trace_id": "xx",
         "target": {
            "live_url" : "https://mobile.yangkeduo.com/transac_virtual_card_pwd.html?ext=%7B%22sensitive_user%22%3A0%2C%22feed_scene_id%22%3A128%7D&room_id=881041069&play_url=http%3A%2F%2Flive-adapt.pddpic.com%2Flive%2F17283_production_sprite_20231227_881041069_01_720p.flv%3FtxSecret%3Dbb5beaf090400150d33b49fac657c081%26txTime%3D658cfe3e%26pub_type%3Drtc-iOS-s1%26octExpId%3D-%26srcSet%3Ds1&refer_page_id=31430_1703666944673_wd0upa1cd1&refer_page_sn=31430&biz_type=0&scene_id=41&if_soft_h265=false&melon_sps_pps=AWQAH%2F%2FhABdnZAAfrOwJgKGhAAADAAEAAAfQDxIlOAEABGjq7yw%3D&if_h265=false&index_param=%257B%25220_17%2522%253A%257B%2522hub_index%2522%253A4%252C%2522max_offset%2522%253A4%252C%2522hub_real_offset%2522%253A4%252C%2522hub_biz_type%2522%253A0%252C%2522hub_source_type%2522%253A17%252C%2522has_more%2522%253Atrue%257D%257D&us",
            "platform": "pinduoduo",
            "headless": true
         },
         "service": {
             "api_base_url": "https://roms-dev.shyuhuankj.com/aigc/live/live-api",
             "project_id": "1351",
             "tenant_id": "10643218595",
             "authorization": "eyJhbGciOiJIUzUxMiJ9.eyJqdGkiOiIzZWViZDllNC05OGViLTRmNTAtOTdiZC1kOTdkZjdhODk4ZDQiLCJPcmlnaW5hbFNvdXJjZSI6IlBDX0xJVkUiLCJzdWIiOiI2NDMyMTg1OTUiLCJleHAiOjE3MDU1NzcwMTJ9.J0rNrFqH0KZaTAO4ztpk9dPmuOetvHEPZtA9NETVFo50j3qX7gLPEGaWwptCKMePUKcKvYcmSOU_URazutcCTg"
         }
      }
      ```
      - 美团-headless
      ```json
      {
         "type": "start",
         "trace_id": "xx",
         "target": {
            "live_url" : "https://g.meituan.com/app/business-live-broadcast/live-detail-new.html?liveid=2303472&notitlebar=1",
            "platform": "meituan",
            "headless": true
         },
         "service": {
             "api_base_url": "https://roms-dev.shyuhuankj.com/aigc/live/live-api",
             "project_id": "1351",
             "tenant_id": "10643218595",
             "authorization": "eyJhbGciOiJIUzUxMiJ9.eyJqdGkiOiIzZWViZDllNC05OGViLTRmNTAtOTdiZC1kOTdkZjdhODk4ZDQiLCJPcmlnaW5hbFNvdXJjZSI6IlBDX0xJVkUiLCJzdWIiOiI2NDMyMTg1OTUiLCJleHAiOjE3MDU1NzcwMTJ9.J0rNrFqH0KZaTAO4ztpk9dPmuOetvHEPZtA9NETVFo50j3qX7gLPEGaWwptCKMePUKcKvYcmSOU_URazutcCTg"
         }
      }
      ```
      ```json
      {
         "type": "event",
         "trace_id": "xx",
         "event": {
            "username": "john_doe",
            "type": "comment",
            "content": "Great stream!",
            "event_time": "2023-12-19 11:00:00"
         }
      }
      ```
      ```json
      {
         "type": "event",
         "trace_id": "xx",
         "event": {
            "username": "",
            "type": "view",
            "content": "11",
            "counter": {
              "count": 1000, 
              "incr": false
            },
            "event_time": "2023-12-19 11:00:00"
         }
      }
      ```
      - 快手-target
      ```json
      {
         "type": "start",
         "trace_id": "xx",
         "target": {
            "live_url" : "https://live.kuaishou.com/u/WPJS88888888TZ",
            "platform": "kuaishou",
            "cookie": "clientid=3; did=web_d18003f098669fb4f60e640b668dddfa; client_key=65890b29; kpn=GAME_ZONE; kuaishou.live.bfb1s=3e261140b0cf7444a0ba411c6f227d88; kuaishou.live.web_st=ChRrdWFpc2hvdS5saXZlLndlYi5zdBKgAdf6l8P_CWYIwGEC170J84TUzQvuV8K19cJWV56dQFUW3tWjubteVMvK0uHkWEHoq1g89P9XaJhenRecWvXC2VBAJ8y-dGMz9PBKVLYKdP4j_nDOv3yqDDQTjAKXSSStJ7W8iSpmK3mxpnNXIeErDbWVzWneGn9Q5XD2mjsNb8-mgDs4R6yvgLKTy2qqnrYbPqX3giwSAneqr27XRmjo02saEsvpGUru20c-iIt7T0W8MQrXwiIgGVmLEgmuD1-XIbQxcgs61QjFey4MMRHmCjIRn5NrC4IoBTAB; kuaishou.live.web_ph=89ad0dca3a7d9214fbc1b5a2b76e6a822d12; userId=3720421602;"
         },
         "service": {
             "api_base_url": "https://roms-dev.shyuhuankj.com/aigc/live/live-api",
             "project_id": "1351",
             "tenant_id": "10643218595",
             "authorization": "eyJhbGciOiJIUzUxMiJ9.eyJqdGkiOiIzZWViZDllNC05OGViLTRmNTAtOTdiZC1kOTdkZjdhODk4ZDQiLCJPcmlnaW5hbFNvdXJjZSI6IlBDX0xJVkUiLCJzdWIiOiI2NDMyMTg1OTUiLCJleHAiOjE3MDU1NzcwMTJ9.J0rNrFqH0KZaTAO4ztpk9dPmuOetvHEPZtA9NETVFo50j3qX7gLPEGaWwptCKMePUKcKvYcmSOU_URazutcCTg"
         }
      }
      ```
      - 快手-roomInfo
      ```json
       {
         "type": "start",
         "trace_id": "xx",
         "target": {
            "live_url" : "https://live.kuaishou.com/u/Lxx2002128",
            "platform": "kuaishou"
         },
         "room": {
           "room_id": "Wb348lZ2N4M",
           "title": "双十一，购物享受批发价",
           "token": "1qJt/y3PygSrDmlTBO43UfhbxcPNhJi5CFLDGkwRpVSwk7iP5qkbsEplkL6s1Evbr2Kbe+RzhnaorSSiMKmQXARM83hc7cySrq00mIuhMWXRfm2KMWtCS1HiMYX+oqNWTktX3KrGuQytBGJzjMzufVOCZe8Mm84MgTxzo1yUISE=",
           "web_socket_url": "wss://live-ws-group8.kuaishou.com/websocket"
         },
         "service": {
             "api_base_url": "https://aigc-video-dev.softsugar.com/aigc/live/live-api-dev",
             "project_id": "1351",
             "tenant_id": "10643218595",
             "authorization": "eyJhbGciOiJIUzUxMiJ9.eyJqdGkiOiIxYTEyNWQyYy03ZWUyLTRjYmItOGEwZS04NGE2NjhjYTI5OGIiLCJPcmlnaW5hbFNvdXJjZSI6IlBDX0xJVkUiLCJzdWIiOiI2NDMyMTg1OTUiLCJleHAiOjE3MDIxMjI2MTF9.yz89NpRmmirWyO4_OPzeAoj8yDAb_ji65jrL2WnE-uP92KIDs3UOiOeeejVVI666_nqav_rGfVYEOVZvipYcCQ"
         }
      }
      ```
      - 抖音-target
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
             "authorization": "eyJhbGciOiJIUzUxMiJ9.eyJqdGkiOiIxYTEyNWQyYy03ZWUyLTRjYmItOGEwZS04NGE2NjhjYTI5OGIiLCJPcmlnaW5hbFNvdXJjZSI6IlBDX0xJVkUiLCJzdWIiOiI2NDMyMTg1OTUiLCJleHAiOjE3MDIxMjI2MTF9.yz89NpRmmirWyO4_OPzeAoj8yDAb_ji65jrL2WnE-uP92KIDs3UOiOeeejVVI666_nqav_rGfVYEOVZvipYcCQ"
         }
      }
      ```
      - 抖音-roomInfo
      ```json
      {
         "type": "start",
         "trace_id": "xx",
         "target": {
            "live_url" : "https://live.douyin.com/403184799752",
            "platform": "douyin"
         },
         "room": {
           "room_id": "7299213520808184586",
           "title": "清扬双11破价！清扬官方旗舰店直播洗发水排行第一名直播间#补水滋养控油",
           "token": "1%7Cjcb-IghCVVAkST4Eo8bTqvMrwQ2Fq_Ox0vymsYpsngY%7C1699529278%7C7b46e0753186e7757657eba0223c5bf1d9a13fd4587ae0477fa4a4e3fec4a832",
           "web_socket_url": "wss://webcast3-ws-wheadless-lq.douyin.com/webcast/im/push/v2/?app_name=douyin_web&version_code=180800&webcast_sdk_version=1.3.0&update_version_code=1.3.0&compress=gzip&internal_ext=internal_src:dim|wss_push_room_id:'+liveRoomId+'|wss_push_did:7188358506633528844|dim_log_id:20230521093022204E5B327EF20D5CDFC6|fetch_time:1684632622323|seq:1|wss_info:0-1684632622323-0-0|wrds_kvs:WebcastRoomRankMessage-1684632106402346965_WebcastRoomStatsMessage-1684632616357153318&cursor=t-1684632622323_r-1_d-1_u-1_h-1&host=https://live.douyin.com&aid=6383&live_id=1&did_rule=3&debug=false&maxCacheMessageNumber=20&endpoint=live_pc&support_wrds=1&im_path=/webcast/im/fetch/&user_unique_id=7188358506633528844&device_platform=wheadless&cookie_enabled=true&screen_width=1440&screen_height=900&browser_language=zh&browser_platform=MacIntel&browser_name=Mozilla&browser_version=5.0%20(Macintosh;%20Intel%20Mac%20OS%20X%2010_15_7)%20AppleWebKit/537.36%20(KHTML,%20like%20Gecko)%20Chrome/113.0.0.0%20Safari/537.36&browser_online=true&tz_name=Asia/Shanghai&identity=audience&room_id=7299213520808184586&heartbeatDuration=0&signature=00000000"
         },
         "service": {
             "api_base_url": "https://aigc-video-dev.softsugar.com/aigc/live/live-api-dev",
             "project_id": "1351",
             "tenant_id": "10643218595",
             "authorization": "eyJhbGciOiJIUzUxMiJ9.eyJqdGkiOiIxYTEyNWQyYy03ZWUyLTRjYmItOGEwZS04NGE2NjhjYTI5OGIiLCJPcmlnaW5hbFNvdXJjZSI6IlBDX0xJVkUiLCJzdWIiOiI2NDMyMTg1OTUiLCJleHAiOjE3MDIxMjI2MTF9.yz89NpRmmirWyO4_OPzeAoj8yDAb_ji65jrL2WnE-uP92KIDs3UOiOeeejVVI666_nqav_rGfVYEOVZvipYcCQ"
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
