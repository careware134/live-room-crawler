package main

import (
	"flag"
	"fmt"
	"live-room-crawler/connector"
	"live-room-crawler/util"
)

var (
	liveUrl string = "https://live.douyin.com/416700408775"
	port    int
)

func main() {
	util.InitLog()
	logger := util.Logger()

	flag.StringVar(&liveUrl, "url", "", "抖音直播间URL, eg: https://live.douyin.com/416700408775")
	flag.IntVar(&port, "port", 50000, "verbose log")

	flag.Parse()

	// Validate arguments
	if liveUrl == "" || port <= 0 {
		fmt.Printf("必须传以下参数: \n" +
			"\t -url: \t  抖音直播间的http URL地址；例如：https://live.douyin.com/416700408775\n" +
			"\t -port: \t 小程序监听的本地端口号，例如：50000 \n" +
			"\t e.g. \n" +
			"\t crawler -url https://live.douyin.com/416700408775 -port 50000")
		return

	}

	logger.Infof("ready to crawl url:%s on port:%d", liveUrl, port)

	roomInfo := connector.RetrieveRoomInfoFromHttpCall(liveUrl)
	connector.WssServerStart(roomInfo)

}
