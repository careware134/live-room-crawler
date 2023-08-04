package main

import (
	"flag"
	"fmt"
	"live-room-crawler/common"
	"live-room-crawler/local_server"
	"live-room-crawler/platform"
	"live-room-crawler/util"
)

var (
	cliMode      bool   = false
	platformName string = "douyin"
	liveUrl      string = "https://live.douyin.com/416700408775"

	serverMode bool = false
	port       int
)

func main() {
	util.InitLog()
	logger := util.Logger()

	flag.BoolVar(&cliMode, "c", false, "cli mode")
	flag.BoolVar(&serverMode, "s", false, "server mode")
	flag.StringVar(&platformName, "platform", "douyin", "抖音直播间URL, eg: https://live.douyin.com/416700408775")
	flag.StringVar(&liveUrl, "url", "", "抖音直播间URL, eg: https://live.douyin.com/416700408775")
	flag.IntVar(&port, "port", 50000, "verbose log")

	flag.Parse()

	// Validate arguments
	if !serverMode && !cliMode {
		fmt.Printf("必须传以下参数: \n" +
			"\t -c: 演示模式；直接输出日志" +
			"\t\t -platform: 直播平台，枚举：douyin,kuaishou(NSY)" +
			"\t\t -url: \t  抖音直播间的http URL地址；例如：https://live.douyin.com/416700408775\n" +
			"\t\t e.g. \n" +
			"\t\t live-room-crawler.exe -c -platform douyin -url https://live.douyin.com/416700408775\n" +
			"\t -s: 运行模式；监听客户端的指令" +
			"\t\t -port: \t 小程序监听的本地端口号，例如：50000 \n" +
			"\t\t e.g. \n" +
			"\t\t live-room-crawler.exe -s -port 50000")
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

	if cliMode {
		logger.Infof("ready to crawl platform:%s url:%s ", platformName, liveUrl)
		platformConnector := platform.NewConnector(common.TargetStruct{
			Platform: common.DOUYIN,
			LiveURL:  liveUrl,
		})

		platformConnector.Start()
	}

	if serverMode {
		local_server.StartLocalServer(port)
	}

}
