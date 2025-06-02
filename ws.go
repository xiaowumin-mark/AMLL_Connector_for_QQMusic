package main

import (
	"log"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
	"github.com/wailsapp/wails/v3/pkg/application"
)

func (a *GreetService) ConnectAmll() error {
	return ConnectAmll()
}

func ConnectAmll() error {

	if coon != nil {
		// 断开
		DisconnectAmll()
	}

	log.Println("正在连接AMLL")
	application.Get().EmitEvent("logger", time.Now().Format("2006-01-02 15:04:05"), "正在连接AMLL")

	u := url.URL{Scheme: "ws", Host: viper.GetString("auto_connect_address") + ":" + viper.GetString("auto_connect_port"), Path: "/"}
	var err error

	var conn *websocket.Conn
	conn, _, err = websocket.DefaultDialer.Dial(u.String(), nil)

	if err != nil {
		application.Get().EmitEvent("amll_ws_state", false)
		log.Println("连接AMLL失败:", err)
		application.Get().EmitEvent("logger", time.Now().Format("2006-01-02 15:04:05"), "连接AMLL失败:"+err.Error())
		return err
	}

	connLock.Lock()
	coon = conn
	connLock.Unlock()

	application.Get().EmitEvent("amll_ws_state", true)
	log.Println("AMLL WS连接成功")
	isCon = true
	application.Get().EmitEvent("logger", time.Now().Format("2006-01-02 15:04:05"), "AMLL WS连接成功")

	if coon != nil {

		//sendMessage(coon, &OnVolumeChanged{Volume: 1})
		musicInfo := &SetMusicInfo{
			MusicID:   NullString(nowMusicInfoMore.Title),
			MusicName: NullString(nowMusicInfoMore.Title),
			AlbumID:   NullString(nowMusicInfoMore.Album),
			AlbumName: NullString(nowMusicInfoMore.Album),
			Artists: []Artist{
				{
					Name: NullString(nowMusicInfoMore.Artist),
					ID:   NullString(nowMusicInfoMore.Artist),
				},
			},
			Duration: uint64(nowMusicInfoMore.Duration),
		}

		// 发送消息
		sendMessage(coon, musicInfo)

		sendMessage(coon, &SetMusicAlbumCoverImageData{Data: nowMusicInfoMore.Pic})

		sendMessage(coon, &SetLyric{Data: nowMusicInfoMore.Lyric})

		if paused {
			connLock.Lock()
			err := coon.WriteMessage(websocket.BinaryMessage, OnPausedMessage.ToBytes())
			if err != nil {
				log.Println("write error:", err)
			}
			connLock.Unlock()
		} else {
			connLock.Lock()
			err := coon.WriteMessage(websocket.BinaryMessage, OnResumedMessage.ToBytes())
			if err != nil {
				log.Println("write error:", err)
			}
			connLock.Unlock()
		}
	}

	if paused == false {

		if coon != nil {
			// 发送消息
			connLock.Lock()
			err := coon.WriteMessage(websocket.BinaryMessage, OnResumedMessage.ToBytes())
			if err != nil {
				log.Println("write error:", err)
			}
			connLock.Unlock()
		}
	} else {
		if coon != nil {
			// 发送消息
			connLock.Lock()
			err := coon.WriteMessage(websocket.BinaryMessage, OnPausedMessage.ToBytes())
			if err != nil {
				log.Println("write error:", err)
			}
			connLock.Unlock()
		}
	}

	return nil
}
func (a *GreetService) DisconnectAmll() error {
	return DisconnectAmll()
}

func DisconnectAmll() error {
	isCon = false
	log.Println("断开AMLL WS连接")
	application.Get().EmitEvent("logger", time.Now().Format("2006-01-02 15:04:05"), "断开AMLL WS连接")
	connLock.Lock()
	defer connLock.Unlock()

	if coon != nil {
		isCon = false // 先标记为不期望连接
		err := coon.Close()
		coon = nil // 清空连接对象
		application.Get().EmitEvent("amll_ws_state", false)

		return err
	}
	return nil
}

// 心跳检测
func Heartbeat() {
	for {
		if !isCon {
			time.Sleep(time.Second * 1)
			continue
		}

		if coon == nil {
			// 触发重连
			if application.Get() != nil {
				application.Get().EmitEvent("logger", time.Now().Format("2006-01-02 15:04:05"), "AMLL WS 尝试自动重连")
			}
			if err := ConnectAmll(); err != nil {
				log.Println("自动重连失败:", err)
				if application.Get() != nil {
					application.Get().EmitEvent("logger", time.Now().Format("2006-01-02 15:04:05"), "AMLL WS 自动重连失败")
					time.Sleep(time.Second)
				}
			}
			continue
		}

		// 发送 ping

		log.Println("发送ping包")
		connLock.Lock()
		err := coon.WriteMessage(websocket.BinaryMessage, PingMessage.ToBytes())
		if err != nil {
			connLock.Unlock()
			log.Println("ping 发送失败:", err)
			continue
		}
		connLock.Unlock()

		pingpongtimer = GsetTimeout(func() {
			// 没有接收到pong
			// 重连
			log.Println("重连")
			isCon = false
			application.Get().EmitEvent("logger", time.Now().Format("2006-01-02 15:04:05"), "AMLL 假死 重连中...")
			if err := ConnectAmll(); err != nil {
				log.Println("自动重连失败:", err)
				if application.Get() != nil {
					application.Get().EmitEvent("logger", time.Now().Format("2006-01-02 15:04:05"), "AMLL WS 自动重连失败")
				}
			}
		}, time.Duration(2*time.Second))

		time.Sleep(2 * time.Second)
	}
}
