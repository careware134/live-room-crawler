package kuaishou

const (
	PageIdCharacterSet        = "-_zyxwvutsrqponmlkjihgfedcba9876543210ZYXWVUTSRQPONMLKJIHGFEDCBA"
	KuaishouApiHost           = "https://live.kuaishou.com/live_graphql"
	RoomInfoRegExp            = `_STATE__=(.*?);\(function\(\)\{var s;\(s=document\.currentScript\|\|document\.scripts\[document\.scripts\.length-1]\)\.parentNode\.r`
	RoomInfoRequestURLPattern = "https://live.kuaishou.com/live_api/liveroom/websocketinfo?liveStreamId=%s"
	HeaderAcceptValue         = "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"
	HeaderAgentValue          = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36"
	//HeaderAgentValue = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:109.0) Gecko/20100101 Firefox/117.0"
	//HeaderCookieValue         = "clientid=3; didv=1675056349580; userId=845495460; kuaishou.live.bfb1s=7206d814e5c089a58c910ed8bf52ace5; did=web_a41c846314016ca1a260444bb3c7d66c; client_key=65890b29; kpn=GAME_ZONE; _did=web_117262730815E9F7; kuaishou.live.web_st=ChRrdWFpc2hvdS5saXZlLndlYi5zdBKgAQkoEnsRiD0ovFwIQ828tvYMhmH6rThiUxM-uuTQtXKjmQEry1dCvI5sEsH9SZt9LNWcvJ_kNRPH2AFvS1awpa65z-Jpe3p2nbMvkpraiJV0WkJrvhLrCyb_CTCNPBGoYwUBaDoabrmZLqLJX-txGbrmUDIblQmR-MKwbPb7uQ5MszR2O3jaon_MtIrqnQA7e0IOBVmJT8N_p-lsiclN4NsaEsa__TMaP0jJgfAfW0kccZcKPyIgmgfFxb6YcCH2fKNK5CO2G4OWyK-WxFeXx6Bx8LA1FGcoBTAB; kuaishou.live.web_ph=8d652450751eaf1d2b61edf08b812bb0a41a; userId=845495460; ksliveShowClipTip=true"
	HeaderCookieValue  = "did=web_3774075bedad1ca6912d86844b5d06e6; clientid=3; did=web_3774075bedad1ca6912d86844b5d06e6; client_key=65890b29; kpn=GAME_ZONE; kuaishou.live.bfb1s=7206d814e5c089a58c910ed8bf52ace5; needLoginToWatchHD=1; showFollowRedIcon=1"
	HeaderCookieValue2 = "client_key=65890b29; clientid=3; did=web_3774075bedad1ca6912d86844b5d06e6; kpn=GAME_ZONE; kuaishou.live.bfb1s=7206d814e5c089a58c910ed8bf52ace5"
	HeaderCookieValue3 = "clientid=3; did=web_67bf94b290093b968d104ed54f5897f0; client_key=65890b29; kpn=GAME_ZONE; _did=web_755635634E404CD6; did=web_feb7f31403ee04c2885a10c0af4e535bd7f4; needLoginToWatchHD=1; userId=3720421602; kuaishou.live.bfb1s=9b8f70844293bed778aade6e0a8f9942"
	// liwei
	HeaderCookieValue4 = "did=web_3774075bedad1ca6912d86844b5d06e6; clientid=3; did=web_3774075bedad1ca6912d86844b5d06e6; client_key=65890b29; kpn=GAME_ZONE; kuaishou.live.bfb1s=7206d814e5c089a58c910ed8bf52ace5; needLoginToWatchHD=1; didv=1694676407793; userId=3720421602; kuaishou.live.web_st=ChRrdWFpc2hvdS5saXZlLndlYi5zdBKgAZAaJb77v9tqAgLrHOYsLmzKt9ziI3Ftcw8BrKBzxM8C5lTBP_G8Qb1FtJEqahyPKovoXZxpb7oOE9UbhnG_3jLKzacyxYXNbJq6CWrC7j1j0qnlVGs3zbjoPaXQVJuDqXlauBfAzSuB-Lao9SrYM_zbCzNmLC6_3KGiNcizicZEppLrWskPQG6PdeWE5lr_KfinLlBakUkOYOiMBG-iVr4aEgL1c1j1KEeWrOe8x-vTC5n9jyIg4niyrg7Bp_ol_ur_r6Ox9_pFQ769exC4qcs1QT-oUFIoBTAB; kuaishou.live.web_ph=a005565b06f81213a9193cd3da2d06dc5b86; userId=3720421602"
	// hemin = "\ndid=web_3774075bedad1ca6912d86844b5d06e6; clientid=3; did=web_3774075bedad1ca6912d86844b5d06e6; client_key=65890b29; kpn=GAME_ZONE; kuaishou.live.bfb1s=7206d814e5c089a58c910ed8bf52ace5; needLoginToWatchHD=1; didv=1694676407793; userId=3720421602
	HeaderCookieValue5 = "did=web_1e7540e76785fe146da52f164413dc06; clientid=3; did=web_1e7540e76785fe146da52f164413dc06; client_key=65890b29; kpn=GAME_ZONE; didv=1693915702610; kuaishou.live.web_st=ChRrdWFpc2hvdS5saXZlLndlYi5zdBKgAdc3esML58edk6tXakWfLLGz5iyhCmFoXhnHsFIq4YzqyZVTjksS5ksh9WCpHahXH6nG4cJLG1t1xCGJ2AD_MTBhN3CeA9xj49ShFVqYQieS0EwW6v4jJQpXd6_EwPZfa8aXhfTOoOlKzos7qIIcAwQq2DSUhvQAV0TeMsjblGeLJlU3xGgFvsUUJDJZ3pgSKAw_AkICQ9oYjeRv7qKI64YaEqoO82cRG00nqSoI30_iGx2JSCIgBTsycKsFuyrTgwfAc_Wf4R34hp03UQ6zATYO9cV7K2goBTAB; kuaishou.live.web_ph=4504422db2a0fbf9955681f285d965430461; userId=3485496356; userId=3485496356; kuaishou.live.bfb1s=3e261140b0cf7444a0ba411c6f227d88"
	// jingrui
	HeaderCookieValue6 = "did=web_50af08ff88705efbc98da084acc44908; didv=1693489524854; clientid=3; did=web_50af08ff88705efbc98da084acc44908; client_key=65890b29; kpn=GAME_ZONE; kuaishou.live.bfb1s=9b8f70844293bed778aade6e0a8f9942; userId=2234219639; kuaishou.live.web_st=ChRrdWFpc2hvdS5saXZlLndlYi5zdBKgAW9OGDTbuJjM18jQQX9U5Inbag4hryP9CrfGEiwmc6hpTVHB2rMzp-jWQ1OsD7bH2dCvU5YQsKFbyXJTtjekw6jZaDS4u5qBwSMnF-0N4_VNLUX6pUZgSDB3U6uVc8-1zraR02OI8NiGBW7W1nDjj5L7bsRzCAn311v_1nAoipMpFOT6HTYRwEr7aVT6xNbyzPpeCTDBqly2U6Ij5VJb_kIaEiTUdhIhLkqeuKi4MmqrjKj9xSIgQGTFrE73rEHI8Q6qkfkTXmDKGcwkYhA6_mhvJrqphcYoBTAB; kuaishou.live.web_ph=3d7b10b019f24095135a9f7f6ecccac55f1a; userId=2234219639"
	// yihan
	HeaderCookieValue7 = "clientid=3; did=web_ed6647fa5146f926c3f4f03a741b92a4; kuaishou.live.bfb1s=7206d814e5c089a58c910ed8bf52ace5; clientid=3; did=web_ed6647fa5146f926c3f4f03a741b92a4; client_key=65890b29; kpn=GAME_ZONE; _did=web_518041749D600775; userId=3720421602; kuaishou.live.web_st=ChRrdWFpc2hvdS5saXZlLndlYi5zdBKgAX0M07e9_cEi7S9MZROOSeUC_LaM-ZBm1WehCLBzjrAJPA1OXcU9c8NsBGS-9R_lU6sir9JF64-HWuAobsf2omiY9-uRiM0sS-wfGBNv-9dv74TE3fcNNKfsfUoODx_EpESbiHymcHsOZnx4v2jh0yKrgvx4nirIKP6l4k49o0n-zbszl7qxvRkrgUEha71Un7aixz9Y7d1ozdLr9oLiOTsaEsa__TMaP0jJgfAfW0kccZcKPyIgF0L5tRhB9FpN_y6__jzrW6qzt9068QQDNhOZIIMJxEgoBTAB; kuaishou.live.web_ph=f5b2806162404b1a109efd394cd048f4baab; userId=3720421602"
)
