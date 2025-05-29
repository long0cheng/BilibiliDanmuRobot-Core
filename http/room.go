package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/long0cheng/BilibiliDanmuRobot-Core/entity"
	"github.com/zeromicro/go-zero/core/logx"
)

func RoomInit(roomid int) (*entity.RoomInitInfo, error) {
	var err error
	var resp *resty.Response
	var url = fmt.Sprintf("https://api.live.bilibili.com/room/v1/Room/room_init?id=%v", roomid)

	if resp, err = cli.R().
		SetHeader("user-agent", userAgent).
		Get(url); err != nil {
		logx.Error("请求room_init失败：", err)
		return nil, err
	}

	// 先解析响应状态
	status := &entity.RoomInitStatus{}
	if err = json.Unmarshal(resp.Body(), status); err != nil {
		logx.Error("Unmarshal失败：", err, "body:", string(resp.Body()))
		return nil, err
	}

	// 在解析房间状态
	r := &entity.RoomInitInfo{}
	if status.Code == 0 {
		if err = json.Unmarshal(resp.Body(), r); err != nil {
			logx.Error("Unmarshal失败：", err, "body:", string(resp.Body()))
			return nil, err
		}
	}

	// 太长时间下播，房间号可能会消失，请求响应的code=60004
	if status.Code == 60004 {
		return nil, errors.New("房间号不存在")
	}
	return r, err
}

func Userinfo(roomid int) (userinfo *entity.Userinfo, err error) {
	roominfo, err := RoomInit(roomid)
	if err != nil {
		return nil, err
	}
	var url = fmt.Sprintf("https://api.live.bilibili.com/live_user/v1/Master/info?uid=%v", roominfo.Data.Uid)

	var resp *resty.Response
	if resp, err = cli.R().
		SetHeader("user-agent", userAgent).
		Get(url); err != nil {
		logx.Error("请求room_init失败：", err)
		return nil, err
	}

	// 先解析响应状态
	userinfo = &entity.Userinfo{}
	if err = json.Unmarshal(resp.Body(), userinfo); err != nil {
		logx.Error("Unmarshal失败：", err, "body:", string(resp.Body()))
		return nil, err
	}
	if userinfo.Code != 0 {
		logx.Errorf("直播间id %v 用户id %v 获取用户信息失败", roomid, roominfo.Data.Uid)
		return nil, errors.New("获取用户信息失败")
	}
	return userinfo, nil
}

func TopListInfo(roomid int, userid int64, page int) (toplistinfo *entity.TopListInfo, err error) {
	var url = fmt.Sprintf("https://api.live.bilibili.com/xlive/app-room/v2/guardTab/topList?page_size=29&roomid=%v&page=%v&ruid=%v", roomid, page, userid)
	var resp *resty.Response
	if resp, err = cli.R().
		SetHeader("user-agent", userAgent).
		Get(url); err != nil {
		logx.Error("请求room_init失败：", err)
		return nil, err
	}

	// 先解析响应状态
	toplistinfo = &entity.TopListInfo{}
	if err = json.Unmarshal(resp.Body(), toplistinfo); err != nil {
		logx.Error("Unmarshal失败：", err, "body:", string(resp.Body()))
		return nil, err
	}
	if toplistinfo.Code != 0 {
		logx.Errorf("直播间id %v 用户id %v 获取舰长列表失败", roomid, userid)
		return nil, errors.New("获取舰长列表失败")
	}
	return toplistinfo, nil
}
func RankListInfo(roomid int, userid int64, page int) (toplistinfo *entity.RankListInfo, err error) {
	var url = fmt.Sprintf("https://api.live.bilibili.com/xlive/general-interface/v1/rank/getOnlineGoldRank?ruid=%v&roomId=%v&page=%v&pageSize=50", userid, roomid, page)
	var resp *resty.Response
	if resp, err = cli.R().
		SetHeader("user-agent", userAgent).
		Get(url); err != nil {
		logx.Error("请求room_init失败：", err)
		return nil, err
	}

	// 先解析响应状态
	toplistinfo = &entity.RankListInfo{}
	if err = json.Unmarshal(resp.Body(), toplistinfo); err != nil {
		logx.Error("Unmarshal失败：", err, "body:", string(resp.Body()))
		return nil, err
	}
	if toplistinfo.Code != 0 {
		logx.Errorf("直播间id %v 用户id %v 获取高能列表失败", roomid, roomid, userid)
		return nil, errors.New("获取高能列表失败")
	}
	return toplistinfo, nil
}
