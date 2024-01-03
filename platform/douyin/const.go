package douyin

const (
	RoomInfoJsonStartTag1 = "\\\"roomInfo\\\":{\\\"room\\\":{\\\"id_str\\\""
	RoomInfoJsonStartTag2 = "\"roomInfo\":{\"room\":{\"id_str\""
	WebSocketTemplateURL  = "wss://webcast5-ws-web-lf.douyin.com/webcast/im/push/v2/?app_name=douyin_web&version_code=180800&webcast_sdk_version=1.0.12&update_version_code=1.0.12&compress=gzip&device_platform=web&cookie_enabled=true&screen_width=1792&screen_height=1120&browser_language=zh-CN&browser_platform=MacIntel&browser_name=Mozilla&browser_version=5.0%20(Macintosh;%20Intel%20Mac%20OS%20X%2010_15_7)%20AppleWebKit/537.36%20(KHTML,%20like%20Gecko)%20Chrome/119.0.0.0%20Safari/537.36&browser_online=true&tz_name=Asia/Shanghai&cursor=r-1_d-1_u-1_h-7319812231192482870_t-1704276841309&internal_ext=internal_src:dim|wss_push_room_id:7319778678831860514|wss_push_did:7319801339826898447|dim_log_id:202401031814019367B01ABF0686157D97|first_req_ms:1704276841206|fetch_time:1704276841309|seq:1|wss_info:0-1704276841309-0-0|wrds_kvs:WebcastRoomStreamAdaptationMessage-1704276436212399349_WebcastProfitInteractionScoreMessage-1704276840057531474_WebcastInRoomBannerMessage-GrowthCommonBannerASubSyncKey-1704269169441370270_WebcastRoomRankMessage-1704276772170856436_AudienceGiftSyncData-1704276812923734429_WebcastRoomStatsMessage-1704276838118322752&host=https://live.douyin.com&aid=6383&live_id=1&did_rule=3&endpoint=live_pc&support_wrds=1&user_unique_id=7319801339826898447&im_path=/webcast/im/fetch/&identity=audience&need_persist_msg_count=15&room_id={roomId}&heartbeatDuration=0&signature=R0UhrnXugkKCgTW9"
	RoomIdPlaceHolder     = "{roomId}"
	SimulateUserAgent     = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36"
)
