package main

import (
	"bitbucket.org/neiku/winornot"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/1set/gut/yos"
	"live-room-crawler/domain"
	"live-room-crawler/platform"
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
func main() {
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

		responseStatus := platformConnector.Connect()
		if !responseStatus.Success {
			marshal, _ := json.Marshal(responseStatus)
			logger.Fatalf("fail to connect to platform:%s url:%s with response:%s", platformName, liveUrl, marshal)
		}
		go platformConnector.StartListen(nil)
		fmt.Scanln() // Wait for user input

	}

	if serverMode {
		logger.Infof("version:1.3.1 add chat guide interaction")
		server.StartLocalServer(port)
	}

}
