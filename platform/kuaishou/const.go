package kuaishou

const (
	KuaishouApiHost           = "https://live.kuaishou.com/live_graphql"
	RoomInfoRegExp            = `_STATE__=(.*?);\(function\(\)\{var s;\(s=document\.currentScript\|\|document\.scripts\[document\.scripts\.length-1]\)\.parentNode\.r`
	RoomInfoRequestURLPattern = "https://live.kuaishou.com/live_api/liveroom/websocketinfo?liveStreamId=%s"
	HeaderAcceptValue         = "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"
	HeaderAgentValue          = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36"
	//HeaderCookieValue         = "clientid=3; didv=1675056349580; userId=845495460; kuaishou.live.bfb1s=7206d814e5c089a58c910ed8bf52ace5; did=web_a41c846314016ca1a260444bb3c7d66c; client_key=65890b29; kpn=GAME_ZONE; _did=web_117262730815E9F7; kuaishou.live.web_st=ChRrdWFpc2hvdS5saXZlLndlYi5zdBKgAQkoEnsRiD0ovFwIQ828tvYMhmH6rThiUxM-uuTQtXKjmQEry1dCvI5sEsH9SZt9LNWcvJ_kNRPH2AFvS1awpa65z-Jpe3p2nbMvkpraiJV0WkJrvhLrCyb_CTCNPBGoYwUBaDoabrmZLqLJX-txGbrmUDIblQmR-MKwbPb7uQ5MszR2O3jaon_MtIrqnQA7e0IOBVmJT8N_p-lsiclN4NsaEsa__TMaP0jJgfAfW0kccZcKPyIgmgfFxb6YcCH2fKNK5CO2G4OWyK-WxFeXx6Bx8LA1FGcoBTAB; kuaishou.live.web_ph=8d652450751eaf1d2b61edf08b812bb0a41a; userId=845495460; ksliveShowClipTip=true"
	HeaderCookieValue  = "did=web_3774075bedad1ca6912d86844b5d06e6; clientid=3; did=web_3774075bedad1ca6912d86844b5d06e6; client_key=65890b29; kpn=GAME_ZONE; kuaishou.live.bfb1s=7206d814e5c089a58c910ed8bf52ace5; needLoginToWatchHD=1; showFollowRedIcon=1"
	HeaderCookieValue2 = "client_key=65890b29; clientid=3; did=web_3774075bedad1ca6912d86844b5d06e6; kpn=GAME_ZONE; kuaishou.live.bfb1s=7206d814e5c089a58c910ed8bf52ace5"
)
