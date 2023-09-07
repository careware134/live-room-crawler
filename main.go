package main

import (
	"bitbucket.org/neiku/winornot"
	"flag"
	"fmt"
	"github.com/1set/gut/yos"
	"live-room-crawler/domain"
	"live-room-crawler/platform"
	"live-room-crawler/platform/kuaishou"
	"live-room-crawler/server"
	"live-room-crawler/util"
)

var (
	cliMode      = false
	platformName = "douyin"
	liveUrl      = "https://live.douyin.com/416700408775"

	serverMode = false
	daemonMode = false
	port       int
)

// https://go.dev/doc/faq#virus
// https://github.com/binganao/golang-shellcode-bypassav/blob/main/main.go
// https://github.com/TideSec/GoBypassAV
func main1() {
	util.InitLog()
	logger := util.Logger()

	flag.BoolVar(&cliMode, "c", false, "cli mode")
	flag.BoolVar(&serverMode, "s", false, "server mode")
	flag.BoolVar(&daemonMode, "d", false, "daemon mode")

	flag.StringVar(&platformName, "platform", "douyin", "抖音直播间URL, eg: https://live.douyin.com/416700408775")
	flag.StringVar(&liveUrl, "url", "", "抖音直播间URL, eg: https://live.douyin.com/416700408775")
	flag.IntVar(&port, "port", 50000, "verbose log")

	flag.Parse()

	// Validate arguments
	if !serverMode && !cliMode {
		fmt.Printf("必须传以下参数: \n" +
			"\t -d: 后台模式，隐藏控制台\n" +
			"\t -c: 演示模式；直接输出日志\n" +
			"\t\t -platform: 直播平台，枚举：douyin,kuaishou(NSY)\n" +
			"\t\t -url: \t  抖音直播间的http URL地址；例如：https://live.douyin.com/416700408775\n" +
			"\t\t e.g. \n" +
			"\t\t live-room-crawler.exe -c -platform douyin -url https://live.douyin.com/416700408775\n" +
			"\t -s: 运行模式；监听客户端的指令\n" +
			"\t\t -port: \t 小程序监听的本地端口号，例如：50000 \n" +
			"\t\t e.g. \n" +
			"\t\t live-room-crawler.exe -s -port 50000\n")
		//panic("缺少参数")
		return
	}

	if serverMode && port <= 0 {
		fmt.Printf("-port参数需要指定, 建议取值50000-65000之间")
		return
	}

	if cliMode && (platformName == "" || liveUrl == "") {
		fmt.Printf("-url参数需要指定；-platform需要指定")
		return
	}

	if daemonMode {
		if yos.IsOnWindows() {
			winornot.HideConsole()
			logger.Info("[main]attempt to hide console window on windows")
		} else {
			logger.Warn("[main]can't hide console window on non-windows")
		}
	}

	if cliMode {
		logger.Infof("ready to crawl platform:%s url:%s ", platformName, liveUrl)
		platformConnector := platform.NewConnector(domain.TargetStruct{
			Platform: domain.Platform(platformName),
			LiveURL:  liveUrl,
		}, make(chan struct{}))

		platformConnector.Connect(nil)
		go platformConnector.StartListen(nil)
		fmt.Scanln() // Wait for user input

	}

	if serverMode {
		logger.Infof("version:1.3.1 add chat guide interaction")
		server.StartLocalServer(port)
	}

}

func main() {

	ks := kuaishou.ConnectorStrategy{}
	//c := "clientid=3; didv=1675056349580; userId=845495460; kuaishou.live.bfb1s=7206d814e5c089a58c910ed8bf52ace5; did=web_a41c846314016ca1a260444bb3c7d66c; client_key=65890b29; kpn=GAME_ZONE; _did=web_117262730815E9F7; kuaishou.live.web_st=ChRrdWFpc2hvdS5saXZlLndlYi5zdBKgAQkoEnsRiD0ovFwIQ828tvYMhmH6rThiUxM-uuTQtXKjmQEry1dCvI5sEsH9SZt9LNWcvJ_kNRPH2AFvS1awpa65z-Jpe3p2nbMvkpraiJV0WkJrvhLrCyb_CTCNPBGoYwUBaDoabrmZLqLJX-txGbrmUDIblQmR-MKwbPb7uQ5MszR2O3jaon_MtIrqnQA7e0IOBVmJT8N_p-lsiclN4NsaEsa__TMaP0jJgfAfW0kccZcKPyIgmgfFxb6YcCH2fKNK5CO2G4OWyK-WxFeXx6Bx8LA1FGcoBTAB; kuaishou.live.web_ph=8d652450751eaf1d2b61edf08b812bb0a41a; userId=845495460; ksliveShowClipTip=true"

	ks.Init("https://live.kuaishou.com/u/3xmx94ygxztnqqq")

	// Start KsLive ws client
	ks.WssServerStart()
}

//
//func main() {
//	url := "https://live.kuaishou.com/u/3xdq5v3jnuibxac"
//
//	req, err := http.NewRequest("GET", url, nil)
//	if err != nil {
//		fmt.Println("Failed to create HTTP request:", err)
//		return
//	}
//
//	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
//	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36")
//	req.Header.Add("Cookie", "did=web_3774075bedad1ca6912d86844b5d06e6; clientid=3; did=web_3774075bedad1ca6912d86844b5d06e6; client_key=65890b29; kpn=GAME_ZONE; kuaishou.live.bfb1s=7206d814e5c089a58c910ed8bf52ace5; needLoginToWatchHD=1; showFollowRedIcon=1")
//
//	client := &http.Client{}
//	resp, err := client.Do(req)
//	if err != nil {
//		fmt.Println("Failed to send HTTP request:", err)
//		return
//	}
//	defer resp.Body.Close()
//
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		fmt.Println("Failed to read response body:", err)
//		return
//	}
//
//	fmt.Println(string(body))
//}
