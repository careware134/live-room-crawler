package kuaishou

const (
	KuaishouApiHost           = "https://live.kuaishou.com/live_graphql"
	RoomInfoRegExp            = `_STATE__=(.*?);\(function\(\)\{var s;\(s=document\.currentScript\|\|document\.scripts\[document\.scripts\.length-1]\)\.parentNode\.r`
	RoomInfoRequestURLPattern = "https://live.kuaishou.com/live_api/liveroom/websocketinfo?liveStreamId=%s"
	HeaderAcceptValue         = "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"
	HeaderAgentValue          = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36"
	HeaderCookieValue         = "did=web_3774075bedad1ca6912d86844b5d06e6; clientid=3; did=web_3774075bedad1ca6912d86844b5d06e6; client_key=65890b29; kpn=GAME_ZONE; kuaishou.live.bfb1s=7206d814e5c089a58c910ed8bf52ace5; needLoginToWatchHD=1; showFollowRedIcon=1"
	HeaderCookieValue2        = "client_key=65890b29; clientid=3; did=web_3774075bedad1ca6912d86844b5d06e6; kpn=GAME_ZONE; kuaishou.live.bfb1s=7206d814e5c089a58c910ed8bf52ace5"
)
