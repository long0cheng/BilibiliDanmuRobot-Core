package logic

import (
	"context"
	"github.com/long0cheng/BilibiliDanmuRobot-Core/entity"
	"github.com/long0cheng/BilibiliDanmuRobot-Core/http"
	"github.com/long0cheng/BilibiliDanmuRobot-Core/svc"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

var sender *BulletSender

type BulletSender struct {
	bulletChan chan entity.Bullet
}

func PushToBulletSender(msg string, reply ...*entity.DanmuMsgTextReplyInfo) {
	logx.Info("PushToBulletSender成功", msg)
	bullet := entity.Bullet{
		Msg:   msg,
		Reply: reply,
	}
	sender.bulletChan <- bullet
}

func StartSendBullet(ctx context.Context, svcCtx *svc.ServiceContext) {

	sender = &BulletSender{
		bulletChan: make(chan entity.Bullet, 1000),
	}

	var msg string
	for {
		select {
		case <-ctx.Done():
			goto END
		case bullet := <-sender.bulletChan:
			msg = bullet.Msg
			var danmuLen = svcCtx.Config.DanmuLen
			var msgdata []string
			msgrun := []rune(msg)
			msgLen := len(msgrun)
			msgcount := msgLen / danmuLen
			tmpmsgcount := msgLen % danmuLen
			if tmpmsgcount != 0 {
				msgcount += 1
			}
			for m := 1; m <= msgcount; m++ {
				if msgLen < m*danmuLen {
					msgdata = append(msgdata, string(msgrun[(m-1)*danmuLen:msgLen]))
					continue
				}
				msgdata = append(msgdata, string(msgrun[(m-1)*danmuLen:danmuLen*m]))
			}
			for _, msgs := range msgdata {
				if err := http.Send(msgs, svcCtx, bullet.Reply...); err != nil {
					logx.Errorf("弹幕发送失败：%s msg: %s", err, msgs)
				} else {
					logx.Infof("弹幕发送成功：%s", msgs)
				}
				//fmt.Println(msgs)
				time.Sleep(1 * time.Second) // 防止弹幕发送过快
			}
		}

	}
END:
}
