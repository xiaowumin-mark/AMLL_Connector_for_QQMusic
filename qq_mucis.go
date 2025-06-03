package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/valyala/fastjson"
	"github.com/xiaowumin-mark/QQMusicDecoder-Go"
)

type TLLyricsItem struct {
	Time uint64
	Text string
}

func getMusicRawData(title, artist, album string) ([]byte, error) {
	URL := "https://u.y.qq.com/cgi-bin/musicu.fcg"
	data := map[string]interface{}{
		"req_1": map[string]interface{}{
			"method": "DoSearchForQQMusicDesktop",
			"module": "music.search.SearchCgiService",
			"param": map[string]interface{}{
				"num_per_page": "20",
				"page_num":     "1",
				"query":        title + " " + artist + " " + album,
				"search_type":  0,
			},
		},
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36")
	req.Header.Set("Referer", "https://c.y.qq.com/")
	req.Header.Set("Cookie", "os=pc;osver=Microsoft-Windows-10-Professional-build-16299.125-64bit;appver=2.0.3.131777;channel=netease;__remember_me=true")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Println("请求失败: %s", resp.Status)
		return nil, fmt.Errorf("请求失败: %s", resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return body, nil
}

// 获取歌曲 mid id 专辑id
func getMusicId(title, artist, album string) (string, string, string, []*fastjson.Value, error) {
	body, err := getMusicRawData(title, artist, album)
	if err != nil {
		log.Println("获取歌曲信息失败")
		return "", "", "", nil, err
	}
	jsonStr := string(body)
	log.Println(jsonStr)
	var p fastjson.Parser
	v, err := p.Parse(jsonStr)
	if err != nil {
		log.Fatal(err)
		return "", "", "", nil, err
	}
	songs := v.GetArray("req_1", "data", "body", "song", "list")
	if songs == nil {
		log.Println("没有找到歌曲")
		return "", "", "", nil, err
	}
	song := songs[0]
	mid := string(song.Get("mid").GetStringBytes())
	id := fmt.Sprint(song.Get("id").GetInt())
	albumid := fmt.Sprint(song.GetInt("album", "id"))
	return mid, id, albumid, songs, nil
}

func getLyric(id string) (string, string, error) {
	//https://amll-ttml-db.stevexmh.net/qq/ + id
	url := fmt.Sprintf("https://amll-ttml-db.stevexmh.net/qq/%s", id)
	// 发送HTTP GET请求
	log.Println("请求URL:", url)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Println("请求失败: %s", resp.Status)
		return "", "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return "", "", err
	}
	return "ttml", string(body), nil

}

// 解析歌词文件
var (
	tagRegex       = regexp.MustCompile(`^\[([a-z]+):(.*)\]$`)
	lyricLineRegex = regexp.MustCompile(`^\[(\d+),(\d+)\](.*)$`)
	// 匹配所有合法时间格式 (123,456)
	timePattern = regexp.MustCompile(`\((\d+),(\d+)\)`)
)

func ParseLyric(lyric *QQMusicDecoder.QqLyricsResponse) ([]LyricLine, int, error) {
	lines := strings.Split(lyric.Lyrics, "\n")
	var lyricLines []LyricLine
	offset := 0

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		// 标签处理
		if tagMatches := tagRegex.FindStringSubmatch(line); len(tagMatches) > 0 {
			tagName := tagMatches[1]
			tagValue := tagMatches[2]
			if tagName == "offset" {
				var err error
				offset, err = strconv.Atoi(tagValue)
				if err != nil {
					return nil, 0, fmt.Errorf("invalid offset value: %v", err)
				}
			}
			continue
		}

		// 歌词行
		if lyricMatches := lyricLineRegex.FindStringSubmatch(line); len(lyricMatches) > 0 {
			startTime, err1 := strconv.ParseUint(lyricMatches[1], 10, 64)
			duration, err2 := strconv.ParseUint(lyricMatches[2], 10, 64)
			if err1 != nil || err2 != nil {
				return nil, 0, fmt.Errorf("invalid time format in lyric line: %s", line)
			}

			lyricLine := LyricLine{
				StartTime: startTime,
				EndTime:   startTime + duration,
				Words:     []LyricWord{},
			}

			wordsPart := lyricMatches[3]
			times := timePattern.FindAllStringSubmatchIndex(wordsPart, -1)

			lastEnd := 0
			for _, idx := range times {
				wordText := wordsPart[lastEnd:idx[0]]

				if wordText != "" {
					startStr := wordsPart[idx[2]:idx[3]]
					durStr := wordsPart[idx[4]:idx[5]]

					wordStart, err1 := strconv.ParseUint(startStr, 10, 64)
					wordDuration, err2 := strconv.ParseUint(durStr, 10, 64)
					if err1 != nil || err2 != nil {
						return nil, 0, fmt.Errorf("invalid word time format: %s", wordsPart[idx[0]:idx[1]])
					}

					lyricLine.Words = append(lyricLine.Words, LyricWord{
						StartTime: wordStart,
						EndTime:   wordStart + wordDuration,
						Word:      NullString(wordText),
					})
				}
				lastEnd = idx[1]
			}

			lyricLines = append(lyricLines, lyricLine)
		}
	}

	// 翻译合并逻辑保留
	if lyric.Trans != "" {
		trsL, _, err := parseTranslationLyric(lyric.Trans)
		if err != nil {
			fmt.Println("Error parsing translation lyric:", err)
			return nil, 0, err
		}
		for index, titem := range trsL {
			if index < len(lyricLines) {
				lyricLines[index].Translated = NullString(titem.Text)
			}
		}
	}

	return lyricLines, offset, nil
}

func getMusicPic(aid string) ([]byte, error) {
	url := "http://imgcache.qq.com/music/photo/album_300/76/300_albumpic_" + aid + "_0.jpg"
	// 请求
	log.Println("请求URL:", url)
	resp, err := http.Get(url)
	if err != nil {
		log.Println("请求失败: %s", err)
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Println("请求失败: %s", resp.Status)
		return nil, fmt.Errorf("请求失败: %s", resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("读取失败: %s", err)
		return nil, err
	}
	return body, nil
}

// 解析翻译歌词文件（标准LRC格式）
func parseTranslationLyric(content string) ([]TLLyricsItem, int, error) {
	lines := strings.Split(content, "\n")
	var lyricLines []TLLyricsItem
	offset := 0

	tagRegex := regexp.MustCompile(`^\[([a-z]+):(.*)\]$`)
	timeRegex := regexp.MustCompile(`^\[(\d+:\d+\.\d+)\](.*)$`)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// 处理标签行
		if tagMatches := tagRegex.FindStringSubmatch(line); len(tagMatches) > 0 {
			if tagMatches[1] == "offset" {
				var err error
				offset, err = strconv.Atoi(tagMatches[2])
				if err != nil {
					return nil, 0, fmt.Errorf("invalid offset value: %v", err)
				}
			}
			continue
		}

		// 处理歌词行
		if timeMatches := timeRegex.FindStringSubmatch(line); len(timeMatches) > 0 {
			timeStr := timeMatches[1]
			text := strings.TrimSpace(timeMatches[2])

			// 转换时间格式 [mm:ss.SS] 为毫秒
			parts := strings.Split(timeStr, ":")
			if len(parts) != 2 {
				continue
			}

			minutes, err1 := strconv.ParseUint(parts[0], 10, 64)
			secondsParts := strings.Split(parts[1], ".")
			if len(secondsParts) != 2 {
				continue
			}

			seconds, err2 := strconv.ParseUint(secondsParts[0], 10, 64)
			centiseconds, err3 := strconv.ParseUint(secondsParts[1], 10, 64)
			if err1 != nil || err2 != nil || err3 != nil {
				continue
			}

			startTime := minutes*60000 + seconds*1000 + centiseconds*10

			//lyricLines = append(lyricLines, LyricLine{
			//	StartTime: startTime,
			//	EndTime:   startTime + 1000, // 默认1秒持续时间
			//	Words: []LyricWord{{
			//		StartTime: startTime,
			//		EndTime:   startTime + 1000,
			//		Word:      NullString(text),
			//	}},
			//})
			if text == "" {
				continue
			}
			if text == "//" {
				text = ""
			}
			lyricLines = append(lyricLines, TLLyricsItem{
				Time: startTime,
				Text: text,
			})
		}
	}

	return lyricLines, offset, nil
}

// 合并歌词
func mergeLyrics(mainLyrics, transLyrics []LyricLine, mainOffset, transOffset int) []LyricLine {
	// 创建翻译映射表（按时间对齐）
	transMap := make(map[uint64]string)
	for _, line := range transLyrics {
		adjustedTime := line.StartTime + uint64(transOffset)
		if len(line.Words) > 0 {
			transMap[adjustedTime] = string(line.Words[0].Word)
		}
	}

	// 合并到主歌词
	for i := range mainLyrics {
		adjustedTime := mainLyrics[i].StartTime + uint64(mainOffset)
		if trans, exists := transMap[adjustedTime]; exists {
			mainLyrics[i].Translated = NullString(trans)
		}
	}

	return mainLyrics
}
