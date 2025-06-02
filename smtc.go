package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"

	"github.com/spf13/viper"

	"github.com/gorilla/websocket"
	"github.com/wailsapp/wails/v3/pkg/application"
)

//func Smtc() {
//	wd, _ := os.Getwd()
//	cmd := exec.Command(filepath.Join(wd, "./smtc.exe"))
//	cmd.SysProcAttr = &syscall.SysProcAttr{
//		HideWindow: true,
//	}
//	stdout, _ := cmd.StdoutPipe()
//	cmd.Start()
//
//	scanner := bufio.NewScanner(stdout)
//	for scanner.Scan() {
//		var resp Response
//		if err := json.Unmarshal(scanner.Bytes(), &resp); err != nil {
//			fmt.Println("Invalid JSON:", err)
//			continue
//		}
//
//		if resp.Code != 0 {
//			fmt.Printf("[ERROR %d] %s\n", resp.Code, resp.Message)
//			continue
//		}
//
//		// 如果是媒体信息
//		if len(resp.Data) > 0 && resp.Data[0] == '{' {
//			var info MusicInfo
//			err := json.Unmarshal(resp.Data, &info)
//			if err != nil {
//				fmt.Println("Invalid JSON:", err)
//				continue
//			}
//
//			if application.Get() != nil {
//				application.Get().EmitEvent("amll_music_info", info)
//			}
//			if err := json.Unmarshal(resp.Data, &info); err == nil {
//				lastAlbum = info.Album
//				lastArtist = info.Artist
//				// 读取管道数据
//				if info.Title != lastTitle && lastDuration != info.Duration {
//					// 如果标题发生变化，更新lastTitle
//					lastTitle = info.Title
//					log.Println("歌曲更换", info.Title, info.Artist)
//					truePosition = 0 // 重置更准确的进度
//					duration = info.Duration
//					lastDuration = info.Duration
//					// 更新UI
//
//					previousIndex = make([]int, 0)
//
//					log.Println(info)
//					// 发送消息
//					if coon != nil {
//
//						musicInfo := &SetMusicInfo{
//							MusicID:   NullString(info.Title),
//							MusicName: NullString(info.Title),
//							AlbumID:   NullString(info.Artist),
//							AlbumName: NullString(info.Artist),
//							Artists: []Artist{
//								{
//									Name: NullString(info.Artist),
//									ID:   NullString(info.Artist),
//								},
//							},
//							Duration: uint64(info.Duration),
//						}
//
//						// 发送消息
//						sendMessage(coon, musicInfo)
//
//					}
//
//					if application.Get() != nil {
//						application.Get().EmitEvent("logger", time.Now().Format("2006-01-02 15:04:05"), "歌曲更换 "+info.Title+" "+info.Artist+" "+info.Album)
//					}
//
//					// 获取信息
//
//					mid, id, aid, err := getMusicId(info.Title, info.Artist, info.Album)
//					if err != nil {
//						log.Println("获取歌曲信息失败")
//						continue
//					}
//					log.Println("歌曲MID:", mid)
//					log.Println("歌曲ID:", id)
//					log.Println("专辑ID:", aid)
//
//					application.Get().EmitEvent("set_music_info_is_id", map[string]string{
//						"id":  id,
//						"aid": aid,
//						"mid": mid,
//					})
//					//if coon != nil {
//					//	// 发送消息
//					//	sendMessage(coon, &SetMusicAlbumCoverImageURI{ImgUrl: NullString("http://localhost:20918/pic/" + aid)})
//					//}
//
//					// 下载专辑图片
//					pic_path := viper.GetString("album_art_path")
//					pic_i_path := filepath.Join(pic_path, aid+".png")
//					pic_data := []byte{}
//					// 检查是否存在
//					if _, err := os.Stat(pic_i_path); os.IsNotExist(err) {
//						// 文件不存在，下载文件
//
//						data, err := getMusicPic(aid)
//						if err != nil {
//							// 下载失败
//							log.Println("下载专辑图片失败")
//						} else {
//							pic_data = data
//							if err := os.WriteFile(pic_i_path, data, 0644); err != nil {
//								// 写入文件失败
//								log.Println("写入专辑图片失败")
//								continue
//							}
//							log.Println("下载专辑图片成功")
//						}
//
//					} else {
//						// 读取文件
//						data, err := os.ReadFile(pic_i_path)
//						if err != nil {
//							log.Println("读取专辑图片失败")
//							continue
//						}
//						log.Println("读取专辑图片成功")
//						pic_data = data
//					}
//					if coon != nil {
//						sendMessage(coon, &SetMusicAlbumCoverImageData{Data: pic_data})
//					}
//					//adata, err := getMusicPic(aid)
//					//if err != nil {
//					//	log.Println("获取专辑图片失败")
//					//}
//					//if coon != nil {
//					//	sendMessage(coon, &SetMusicAlbumCoverImageData{Data: adata})
//					//}
//
//					qrcl, err := QQMusicDecoder.GetLyricsByMid(mid)
//					if err != nil {
//						log.Println("获取歌词失败")
//						application.Get().EmitEvent("logger", time.Now().Format("2006-01-02 15:04:05"), "获取歌词失败")
//						continue
//					}
//					log.Println(qrcl)
//
//					lyricLines, offset, err := ParseLyric(qrcl)
//					if err != nil {
//						fmt.Println("解析歌词失败:", err)
//						application.Get().EmitEvent("logger", time.Now().Format("2006-01-02 15:04:05"), "解析歌词失败")
//						continue
//					}
//					// 应用offset
//					if offset != 0 {
//						for i := range lyricLines {
//							lyricLines[i].StartTime += uint64(offset)
//							lyricLines[i].EndTime += uint64(offset)
//							for j := range lyricLines[i].Words {
//								lyricLines[i].Words[j].StartTime += uint64(offset)
//								lyricLines[i].Words[j].EndTime += uint64(offset)
//							}
//						}
//					}
//
//					//if qrcl.Trans != "" {
//					//	transLines, transOffset, err := parseTranslationLyric(qrcl.Trans)
//					//	if err != nil {
//					//		fmt.Println("Error parsing translation lyric:", err)
//					//		return
//					//	}
//					//
//					//	// 合并歌词（自动处理offset）
//					//	lyricLines = mergeLyrics(lyricLines, transLines, offset, transOffset)
//					//
//					//}
//					nowLyric = lyricLines
//					application.Get().EmitEvent("set_lyric_from_qrc", nowLyric)
//
//					log.Println("歌词偏移:", offset)
//					log.Println(lyricLines)
//
//					if coon != nil {
//						sendMessage(coon, &SetLyric{Data: nowLyric})
//					}
//					application.Get().EmitEvent("set_lyric", qrcl)
//				}
//				if info.Position != lastPosition {
//					lastPosition = info.Position
//					//log.Println("歌曲进度更新", info.Position)
//					position = info.Position
//					duration = info.Duration
//					truePosition = position
//					// 更新UI
//
//					if application.Get() != nil {
//						application.Get().EmitEvent("logger", time.Now().Format("2006-01-02 15:04:05"), "歌曲进度更新 "+fmt.Sprint(info.Position))
//					}
//				}
//				if info.Status != lastStatus {
//					lastStatus = info.Status
//					log.Println("歌曲状态更新", info.Status)
//					truePosition = info.Position // 更新更准确的进度
//					if info.Status == "Playing" {
//						paused = false
//						if coon != nil {
//							// 发送消息
//							connLock.Lock()
//							err := coon.WriteMessage(websocket.BinaryMessage, OnResumedMessage.ToBytes())
//							if err != nil {
//								log.Println("write error:", err)
//							}
//							connLock.Unlock()
//						}
//					} else {
//						paused = true
//						if coon != nil {
//							// 发送消息
//							connLock.Lock()
//							err := coon.WriteMessage(websocket.BinaryMessage, OnPausedMessage.ToBytes())
//							if err != nil {
//								log.Println("write error:", err)
//							}
//							connLock.Unlock()
//						}
//					}
//
//					duration = info.Duration
//
//					if application.Get() != nil {
//						application.Get().EmitEvent("logger", time.Now().Format("2006-01-02 15:04:05"), "歌曲状态更新 "+info.Status)
//					}
//				}
//			}
//		} else {
//			fmt.Println("[INFO]", resp.Message)
//		}
//	}
//}

type SessionInfo struct {
	AppID  string `json:"appId"`
	Title  string `json:"title"`
	Artist string `json:"artist"`
	Album  string `json:"album"`
}

type InitializationCompletedEvent struct {
	Event string        `json:"event"`
	Data  []SessionInfo `json:"data"`
}

type SMTCSessionEvent struct {
	Event string      `json:"event"` // "SMTCAdded" or "SMTCRemoved"
	Data  SessionInfo `json:"data"`
}

type PlaybackStatus string

const (
	PlaybackStatusPlaying PlaybackStatus = "Playing"
	PlaybackStatusPaused  PlaybackStatus = "Paused"
	PlaybackStatusStopped PlaybackStatus = "Stopped"
	PlaybackStatusClosed  PlaybackStatus = "Closed"
	PlaybackStatusUnknown PlaybackStatus = "Unknown"
)

type TimelineEventData struct {
	AppID            string         `json:"appId"`
	Title            string         `json:"title"`
	Artist           string         `json:"artist"`
	Album            string         `json:"album"`
	PlaybackStatus   PlaybackStatus `json:"playbackStatus"`
	TimelinePosition int            `json:"timelinePositionMS"`
	TimelineEndTime  int            `json:"timelineEndTimeMS"`
}

type TimelineEvent struct {
	Event string            `json:"event"`
	Data  TimelineEventData `json:"data"`
}

type RmType struct {
	Event string `json:"event"`
}

var isSmtcInitialized = false

var nowSmtcData []TimelineEventData

func Smtc(ep string) {

	//cmd := exec.Command(filepath.Join(wd, "./smtc-center.exe"))
	log.Println(ep)
	cmd = exec.Command(ep)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}
	stdout, _ := cmd.StdoutPipe()
	err := cmd.Start()

	if err != nil {
		log.Println("Error:", err)
		dialog := application.ErrorDialog()
		dialog.SetTitle("Error")
		dialog.SetMessage("SMTC启动失败！")
		dialog.Show()
		application.Get().Quit()
	}

	pauthrottle := NewThrottler(200 * time.Millisecond)
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		msg := scanner.Bytes()
		var resp RmType
		if err := json.Unmarshal(msg, &resp); err != nil {
			fmt.Println("Invalid JSON:", err)
			continue
		}
		switch resp.Event {
		/*
			InitializationCompleted
			SMTCAdded
			SMTCRemoved
			MediaPropertiesChanged
			PlaybackInfoChanged
			TimelinePropertiesChanged
		*/

		case "InitializationCompleted":
			if !isSmtcInitialized {

				isSmtcInitialized = true

				// 解析json
				var data InitializationCompletedEvent
				if err := json.Unmarshal(msg, &data); err != nil {
					log.Println("解析json失败")
					continue
				}
				log.Println("smtc初始化完成,有", len(data.Data), "个smtc在运行")
			}

		case "SMTCAdded":
			var data SMTCSessionEvent
			if err := json.Unmarshal(msg, &data); err != nil {
				log.Println("解析json失败")
				continue
			}

			log.Println("来了个新的smtc", data.Data.AppID)
			nowSmtcData = AddSmtc(data.Data)

			application.Get().EmitEvent("smtc_added", data.Data)
		case "SMTCRemoved":
			var data SMTCSessionEvent
			if err := json.Unmarshal(msg, &data); err != nil {
				log.Println("解析json失败")
				continue
			}
			log.Println("smtc被干掉了", data.Data.AppID)
			nowSmtcData = DelSmtc(data.Data)

			application.Get().EmitEvent("smtc_removed", data.Data)
		default:
			var data TimelineEvent
			if err := json.Unmarshal(msg, &data); err != nil {
				log.Println("解析json失败")
				continue
			}
			info := MusicInfo{
				Title:    data.Data.Title,
				Artist:   data.Data.Artist,
				Album:    data.Data.Album,
				Duration: float64(data.Data.TimelineEndTime),
				Position: float64(data.Data.TimelinePosition),
				Status:   string(data.Data.PlaybackStatus),
			}

			switch resp.Event {

			case "MediaPropertiesChanged":
				//log.Println(data.Data.AppID, "换歌了", data.Data.Title, data.Data.TimelineEndTime)
				nowSmtcData = SetSmtcTimeLineInfo(data.Data)

				application.Get().EmitEvent("smtc_changed", data.Data)
				if data.Data.AppID == "QQMusic.exe" {
					err := handleSmtcMsg(data.Event, info, &data.Data, &nowSmtcData)
					if err != nil {
						log.Println(err)
						application.Get().EmitEvent("logger", time.Now().Format("2006-01-02 15:04:05"), err)
						continue
					}
				}

			case "PlaybackInfoChanged":
				//log.Println(data.Data.AppID, "播放状态变了", data.Data.PlaybackStatus)
				nowSmtcData = SetSmtcTimeLineInfo(data.Data)

				application.Get().EmitEvent("smtc_changed", data.Data)
				if data.Data.AppID == "QQMusic.exe" {
					err := handleSmtcMsg(data.Event, info, &data.Data, &nowSmtcData)
					if err != nil {
						log.Println(err)
						application.Get().EmitEvent("logger", time.Now().Format("2006-01-02 15:04:05"), err)
						continue
					} else {
						lastStatus = info.Status
						log.Println("歌曲状态更新", info.Status)
						truePosition = info.Position // 更新更准确的进度
						if info.Status == "Playing" {
							paused = false
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
							paused = true
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

						duration = info.Duration

						if application.Get() != nil {
							application.Get().EmitEvent("logger", time.Now().Format("2006-01-02 15:04:05"), "歌曲状态更新 "+info.Status)
						}
					}
				} else {

					//puThrottledFunc()
					pauthrottle.Do(func() {
						if data.Data.PlaybackStatus == PlaybackStatusPlaying {
							log.Println("!QQMusic.exe is playing")

							//err := QQMC.Pause()
							//if err != nil {
							//	log.Println(err)
							//}
							targetVolume = 0.2
						} else {
							log.Println("!QQMusic.exe is not playing")

							//err := QQMC.Play()
							//if err != nil {
							//	log.Println(err)
							//}
							targetVolume = 1.0
						}
					})
				}

			case "TimelinePropertiesChanged":
				//log.Println(data.Data.AppID, "时间轴变了", data.Data.TimelinePosition)
				nowSmtcData = SetSmtcTimeLineInfo(data.Data)

				application.Get().EmitEvent("smtc_changed", data.Data)
				if data.Data.AppID == "QQMusic.exe" {
					err := handleSmtcMsg(data.Event, info, &data.Data, &nowSmtcData)
					if err != nil {
						log.Println(err)
						application.Get().EmitEvent("logger", time.Now().Format("2006-01-02 15:04:05"), err)
						continue
					} else {

						lastPosition = info.Position
						//log.Println("歌曲进度更新", info.Position)
						position = info.Position
						duration = info.Duration
						truePosition = position
						// 更新UI

						application.Get().EmitEvent("logger", time.Now().Format("2006-01-02 15:04:05"), "歌曲进度更新 "+fmt.Sprint(info.Position))

					}
				}

			}
			if data.Data.AppID == "QQMusic.exe" {
				lastStatus = info.Status
				if info.Status == "Playing" {
					paused = false
				} else {
					paused = true
				}
				//application.Get().EmitEvent("logger", time.Now().Format("2006-01-02 15:04:05"), "歌曲状态更新 "+info.Status)

			}
		}
	}
}

func AddSmtc(data SessionInfo) []TimelineEventData {
	// 循环查看是否存在
	for i := 0; i < len(nowSmtcData); i++ {
		if nowSmtcData[i].AppID == data.AppID {
			return nowSmtcData
		}
	}
	nowSmtcData = append(nowSmtcData, TimelineEventData{Album: data.Album, AppID: data.AppID, Artist: data.Artist, PlaybackStatus: PlaybackStatusPaused, TimelineEndTime: 0, TimelinePosition: 0, Title: data.Title})
	return nowSmtcData
}

func DelSmtc(data SessionInfo) []TimelineEventData {
	// 删除
	for i := 0; i < len(nowSmtcData); i++ {
		if nowSmtcData[i].AppID == data.AppID {
			nowSmtcData = append(nowSmtcData[:i], nowSmtcData[i+1:]...)
			return nowSmtcData
		}
	}
	return nowSmtcData
}

func SetSmtcTimeLineInfo(data TimelineEventData) []TimelineEventData {
	// 遍历是否存在
	for i := 0; i < len(nowSmtcData); i++ {
		if nowSmtcData[i].AppID == data.AppID {
			nowSmtcData[i] = data
			return nowSmtcData
		}
	}
	return append(nowSmtcData, data)
}

func getSmtcTimeLineInfoByAppId(appId string) *TimelineEventData {
	for i := 0; i < len(nowSmtcData); i++ {
		if nowSmtcData[i].AppID == appId {
			return &nowSmtcData[i]
		}
	}
	return nil
}

// 获取除了某个应用的所有smtc信息
func getSmtcTimeLineInfoExceptAppId(appId string) []TimelineEventData {
	var result []TimelineEventData
	for i := 0; i < len(nowSmtcData); i++ {
		if nowSmtcData[i].AppID != appId {
			result = append(result, nowSmtcData[i])
		}
	}
	return result
}
func handleSmtcMsg(mtype string, info MusicInfo, data *TimelineEventData, smtcs *[]TimelineEventData) error {

	lastAlbum = info.Album
	lastArtist = info.Artist

	if info.Title != lastTitle && lastDuration != info.Duration {
		// 如果标题发生变化，更新lastTitle
		lastTitle = info.Title
		log.Println("歌曲更换", info.Title, info.Artist)
		truePosition = 0 // 重置更准确的进度
		duration = info.Duration
		lastDuration = info.Duration
		// 更新UI

		previousIndex = make([]int, 0)

		log.Println(info)
		// 发送消息

		if application.Get() != nil {
			application.Get().EmitEvent("logger", time.Now().Format("2006-01-02 15:04:05"), "歌曲更换 "+info.Title+" "+info.Artist+" "+info.Album)
		}

		// 获取信息

		mid, id, aid, rd, err := getMusicId(info.Title, info.Artist, info.Album)
		if err != nil {
			//log.Println("获取歌曲信息失败")
			return fmt.Errorf("获取歌曲信息失败")
		}
		log.Println("歌曲MID:", mid)
		log.Println("歌曲ID:", id)
		log.Println("专辑ID:", aid)
		log.Println("歌曲时长:", rd[0].GetInt("interval"))

		application.Get().EmitEvent("amll_music_info", info)

		//lastDuration = rd.GetFloat64("interval")
		info.Duration = rd[0].GetFloat64("interval") * 1000
		if coon != nil {

			musicInfo := &SetMusicInfo{
				MusicID:   NullString(info.Title),
				MusicName: NullString(info.Title),
				AlbumID:   NullString(info.Artist),
				AlbumName: NullString(info.Artist),
				Artists: []Artist{
					{
						Name: NullString(info.Artist),
						ID:   NullString(info.Artist),
					},
				},
				Duration: uint64(info.Duration),
			}

			// 发送消息
			sendMessage(coon, musicInfo)

		}

		application.Get().EmitEvent("set_music_info_is_id", map[string]string{
			"id":  id,
			"aid": aid,
			"mid": mid,
		})
		//if coon != nil {
		//	// 发送消息
		//	sendMessage(coon, &SetMusicAlbumCoverImageURI{ImgUrl: NullString("http://localhost:20918/pic/" + aid)})
		//}

		// 下载专辑图片
		pic_path := viper.GetString("album_art_path")
		pic_i_path := filepath.Join(pic_path, aid+".png")
		nowImg = aid
		pic_data := []byte{}
		// 检查是否存在
		if _, err := os.Stat(pic_i_path); os.IsNotExist(err) {
			// 文件不存在，下载文件

			data, err := getMusicPic(aid)
			if err != nil {
				// 下载失败
				log.Println("下载专辑图片失败")
			} else {
				pic_data = data
				if err := os.WriteFile(pic_i_path, data, 0644); err != nil {
					// 写入文件失败
					//log.Println("写入专辑图片失败")
					return fmt.Errorf("写入专辑图片失败")
				}
				log.Println("下载专辑图片成功")
			}

		} else {
			// 读取文件
			data, err := os.ReadFile(pic_i_path)
			if err != nil {
				//log.Println("读取专辑图片失败")
				return fmt.Errorf("读取专辑图片失败")
			}
			log.Println("读取专辑图片成功")
			pic_data = data
		}
		if coon != nil {
			sendMessage(coon, &SetMusicAlbumCoverImageData{Data: pic_data})
		}
		//adata, err := getMusicPic(aid)
		//if err != nil {
		//	log.Println("获取专辑图片失败")
		//}
		//if coon != nil {
		//	sendMessage(coon, &SetMusicAlbumCoverImageData{Data: adata})
		//}
		var str string
		nowLyric, str, err = getLyrics(mid)
		if err != nil {
			log.Println("获取歌词失败")
			if coon != nil {
				sendMessage(coon, &SetLyric{Data: []LyricLine{}})
			}
			return fmt.Errorf("获取歌词失败")
		}
		application.Get().EmitEvent("set_lyric_from_qrc", nowLyric)
		//log.Println(nowLyric)

		if coon != nil {
			sendMessage(coon, &SetLyric{Data: nowLyric})
		}
		application.Get().EmitEvent("set_lyric", str)

		var rRaw []string
		for _, item := range rd {
			rRaw = append(rRaw, item.String())
		}
		nowMusicInfoMore = MusicInfoMore{
			Album:    info.Album,
			Artist:   info.Artist,
			Duration: rd[0].GetFloat64("interval") * 1000,
			Status:   string(info.Status),
			Title:    info.Title,

			Lyric:     nowLyric,
			LyricType: "qrc",
			LyricRaw:  str,

			Mid:     mid,
			Id:      id,
			AlbumId: aid,

			RequstRaw: rRaw,

			Pic: pic_data,
		}
		application.Get().EmitEvent("set_music_info_more", nowMusicInfoMore)
	}

	return nil
}

//var puThrottledFunc = Throttle(func() {
//	res := getSmtcTimeLineInfoExceptAppId("QQMusic.exe")
//	if len(res) != 0 {
//
//		for _, v := range res {
//			if v.PlaybackStatus == PlaybackStatusPlaying {
//				log.Println("!QQMusic.exe is playing")
//				err := QQMC.Pause()
//				if err != nil {
//					log.Println(err)
//				}
//			} else {
//				log.Println("!QQMusic.exe is not playing")
//				err := QQMC.Play()
//				if err != nil {
//					log.Println(err)
//				}
//			}
//		}
//
//	}
//
//}, time.Second)
