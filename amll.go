package main

import (
	"bytes"
	"encoding/binary"
	"log"

	//"time"

	"github.com/gorilla/websocket"
)

type NullString string

func (ns NullString) ToBytes() ([]byte, error) {
	data := []byte(ns)
	if len(ns) > 0 && ns[len(ns)-1] == 0x00 {
		return data, nil
	}
	data = append(data, 0x00)
	return data, nil
}

const (
	// 接收 / 发送
	MagicPing = 0
	MagicPong = 1

	// 接收
	MagicSetMusicInfo                = 2
	MagicSetMusicAlbumCoverImageURI  = 3
	MagicSetMusicAlbumCoverImageData = 4
	MagicOnPlayProgress              = 5
	MagicOnVolumeChanged             = 6
	MagicOnPaused                    = 7
	MagicOnResumed                   = 8
	MagicOnAudioData                 = 9
	MagicSetLyric                    = 10
	MagicSetLyricFromTTML            = 11
	// 发送
	MagicPause            = 12 //请求发送端暂停音乐播放。
	MagicResume           = 13 //请求发送端恢复播放音乐播放。
	MagicForwardSong      = 14 //请求发送端跳转到下一首歌曲。
	MagicBackwardSong     = 15 //请求发送端跳转到上一首歌曲。
	MagicSetVolume        = 16 //请求发送端设置播放音量。
	MagicSeekPlayProgress = 17 //请求发送端设置当前播放进度。

	isBG   = 0b01 //是否为背景歌词行
	isDuet = 0b10 //是否为对唱歌词行
)

type NoDataMessage struct {
	Magic uint16
}

var (
	PingMessage      = &NoDataMessage{Magic: MagicPing}
	PongMessage      = &NoDataMessage{Magic: MagicPong}
	OnResumedMessage = &NoDataMessage{Magic: MagicOnResumed}
	OnPausedMessage  = &NoDataMessage{Magic: MagicOnPaused}
)

/*
interface Artist {
    id: string;
    name: string;
}*/

type Artist struct {
	ID   NullString `json:"id"`
	Name NullString `json:"name"`
}

func (a Artist) ToBytes() ([]byte, error) {
	idBytes, _ := a.ID.ToBytes()
	nameBytes, _ := a.Name.ToBytes()

	var data []byte
	data = append(data, idBytes...)
	data = append(data, nameBytes...)
	return data, nil
}

/*
	struct SetMusicInfo {
	    music_id: NullString,   // 歌曲的唯一标识字符串
	    music_name: NullString, // 歌曲名称
	    album_id: NullString,   // 歌曲所属的专辑ID，如果没有可以留空
	    album_name: NullString, // 歌曲所属的专辑名称，如果没有可以留空
	    artists: Vec<Artist>,   // 歌曲的艺术家/制作者列表
	    duration: u64,          // 歌曲的时长，单位为毫秒
	}
*/
type SetMusicInfo struct {
	MusicID   NullString `json:"musicId"`
	MusicName NullString `json:"musicName"`
	AlbumID   NullString `json:"albumId"`
	AlbumName NullString `json:"albumName"`
	Artists   []Artist   `json:"artists"`
	Duration  uint64     `json:"duration"`
}

func (m SetMusicInfo) ToBytes() ([]byte, error) {
	var buf bytes.Buffer

	// Magic number (2)
	binary.Write(&buf, binary.LittleEndian, uint16(MagicSetMusicInfo))

	// 序列化各个字段
	musicID, _ := m.MusicID.ToBytes()
	buf.Write(musicID)

	musicName, _ := m.MusicName.ToBytes()
	buf.Write(musicName)

	albumID, _ := m.AlbumID.ToBytes()
	buf.Write(albumID)

	albumName, _ := m.AlbumName.ToBytes()
	buf.Write(albumName)

	// 艺术家列表
	binary.Write(&buf, binary.LittleEndian, uint32(len(m.Artists)))
	for _, artist := range m.Artists {
		artistData, _ := artist.ToBytes()
		buf.Write(artistData)
	}

	// 持续时间
	binary.Write(&buf, binary.LittleEndian, m.Duration)

	return buf.Bytes(), nil
}
func (m *NoDataMessage) ToBytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, m.Magic)
	return buf.Bytes()
}

type SetLyricFromTTML struct {
	Data NullString // TTML 字符串内容
}

func (m *SetLyricFromTTML) ToBytes() ([]byte, error) {
	var buf bytes.Buffer

	binary.Write(&buf, binary.LittleEndian, uint16(MagicSetLyricFromTTML))

	// 序列化 TTML 字符串
	ttmlData, err := m.Data.ToBytes()
	if err != nil {
		return nil, err
	}
	buf.Write(ttmlData)

	return buf.Bytes(), nil
}

type SeekPlayProgress struct {
	Progress []byte
}

type SetMusicAlbumCoverImageURI struct {
	ImgUrl NullString `json:"img_url"`
}

type SetVolume struct {
	Volume []byte
}

func (m *SetVolume) ToData() (float64, error) {
	var volume float64
	// 解析数据
	err := binary.Read(bytes.NewReader(m.Volume), binary.LittleEndian, &volume)
	return volume, err
}

func (m *SetMusicAlbumCoverImageURI) ToBytes() ([]byte, error) {
	data, _ := m.ImgUrl.ToBytes()
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, uint16(MagicSetMusicAlbumCoverImageURI))
	binary.Write(&buf, binary.LittleEndian, data)
	return buf.Bytes(), nil
}

// 反转
func (m *SeekPlayProgress) ToData() (uint64, error) {
	var progress uint64
	// 解析数据
	err := binary.Read(bytes.NewReader(m.Progress), binary.LittleEndian, &progress)
	return progress, err
}

type OnPlayProgress struct {
	Progress uint64 `json:"progress"`
}

func (m *OnPlayProgress) ToBytes() ([]byte, error) {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, uint16(MagicOnPlayProgress))
	binary.Write(&buf, binary.LittleEndian, m.Progress)
	return buf.Bytes(), nil
}

type OnVolumeChanged struct {
	Volume float64 `json:"volume"`
}

func (m *OnVolumeChanged) ToBytes() ([]byte, error) {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, uint16(MagicOnVolumeChanged))
	binary.Write(&buf, binary.LittleEndian, m.Volume)
	return buf.Bytes(), nil
}

type Message interface {
	ToBytes() ([]byte, error)
}

// LyricLine 歌词行
type LyricLine struct {
	StartTime  uint64      `json:"start_time"`
	EndTime    uint64      `json:"end_time"`
	Words      []LyricWord `json:"words"`
	Translated NullString  `json:"translated"`
	RomanLyric NullString  `json:"roman_lyric"`
	Flag       uint8       `json:"flag"`
}

func (m *LyricLine) ToBytes() ([]byte, error) {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, m.StartTime)
	binary.Write(&buf, binary.LittleEndian, m.EndTime)
	binary.Write(&buf, binary.LittleEndian, uint32(len(m.Words)))
	for _, word := range m.Words {
		wordData, _ := word.ToBytes()
		buf.Write(wordData)
	}
	translated, _ := m.Translated.ToBytes()
	buf.Write(translated)
	romanLyric, _ := m.RomanLyric.ToBytes()
	buf.Write(romanLyric)
	binary.Write(&buf, binary.LittleEndian, m.Flag)

	return buf.Bytes(), nil
}

type LyricWord struct {
	StartTime uint64     `json:"start_time"`
	EndTime   uint64     `json:"end_time"`
	Word      NullString `json:"word"`
}

func (m *LyricWord) ToBytes() ([]byte, error) {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, m.StartTime)
	binary.Write(&buf, binary.LittleEndian, m.EndTime)
	word, _ := m.Word.ToBytes()
	buf.Write(word)
	return buf.Bytes(), nil
}

type SetLyric struct {
	Data []LyricLine `json:"data"`
}

func (m *SetLyric) ToBytes() ([]byte, error) {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, uint16(MagicSetLyric))
	binary.Write(&buf, binary.LittleEndian, uint32(len(m.Data)))
	for _, line := range m.Data {
		lineData, _ := line.ToBytes()
		buf.Write(lineData)
	}
	return buf.Bytes(), nil
}

type SetMusicAlbumCoverImageData struct {
	Data []uint8 `json:"data"`
}

func (m *SetMusicAlbumCoverImageData) ToBytes() ([]byte, error) {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, uint16(MagicSetMusicAlbumCoverImageData))
	// unit8
	err := binary.Write(&buf, binary.LittleEndian, uint32(len(m.Data)))
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(m.Data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func sendMessage(conn *websocket.Conn, msg Message) {
	data, err := msg.ToBytes()
	if err != nil {
		log.Println("encode error:", err)
		return
	}

	if conn == nil {
		return
	}
	connLock.Lock()
	err = conn.WriteMessage(websocket.BinaryMessage, data)

	if err != nil {
		connLock.Unlock()
		log.Println("write error:", err)
		return
	}
	connLock.Unlock()
}

func handleIncomingMessage(conn *websocket.Conn, message []byte) {
	log.Println("Received message:", message)
	if len(message) < 2 {
		log.Println("invalid message")
		return
	}

	magic := binary.LittleEndian.Uint16(message[:2])
	switch magic {
	case MagicPing:
		log.Println("Received Ping, sending Pong...")
		connLock.Lock()
		err := conn.WriteMessage(websocket.BinaryMessage, PongMessage.ToBytes())

		if err != nil {
			connLock.Unlock()
			log.Println("write error:", err)
			return
		}
		connLock.Unlock()
	case MagicPong:
		log.Println("Received Pong.")
		GclearTimeout(pingpongtimer)
	// 其他 case 解析可以继续添加
	case MagicPause:
		log.Println("MagicPause")
		if err := QQMC.Pause(); err != nil {
			log.Println("Pause error:", err)
		}
	case MagicResume:
		log.Println("MagicResume")
		//qqctrl.Play()
		if err := QQMC.Play(); err != nil {
			log.Println("Play error:", err)
		}
	case MagicForwardSong:
		log.Println("MagicForwardSong")
		//qqctrl.Next()
		if err := QQMC.Next(); err != nil {
			log.Println("Next error:", err)
		}
	case MagicBackwardSong:
		log.Println("MagicBackwardSong")
		//qqctrl.Prev()
		if err := QQMC.Prev(); err != nil {
			log.Println("Prev error:", err)
		}
	case MagicSetVolume:
		log.Println("MagicSetVolume")
		volume, err := (&SetVolume{Volume: message[2:]}).ToData()
		if err != nil {
			log.Println("decode error:", err)
		} else {
			//SetProcessVolume(float32(volume))
			targetVolume = volume
		}

	case MagicSeekPlayProgress:
		log.Println("SeekPlayProgress")
		progress, err := (&SeekPlayProgress{Progress: message[2:]}).ToData()
		if err != nil {
			log.Println("decode error:", err)
		}
		log.Println("Received SeekPlayProgress:", progress)
	default:
		log.Printf("Received unknown magic: %d", magic)
	}
}
