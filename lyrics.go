package main

import (
	"fmt"
	"log"

	"github.com/xiaowumin-mark/QQMusicDecoder-Go"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// 判断歌词事件循环
func JudgeLyricEvent() {
	// 歌词时间判断

	currentIndex := make([]int, 0)
	for i := 0; i < len(nowLyric); i++ {
		if truePosition >= float64(nowLyric[i].StartTime) && truePosition <= float64(nowLyric[i].EndTime) {
			currentIndex = append(currentIndex, i)
		}
	}
	if len(currentIndex) != len(previousIndex) || !every(currentIndex, previousIndex) {
		handleLyricsChange(currentIndex)
		previousIndex = currentIndex // 更新上次的歌词状态
	}

}

// 每个元素与另一个切片中的元素相等时返回 true，否则返回 false
func every(currentIndex, previousIndex []int) bool {
	// 如果两个切片的长度不同，返回 false
	if len(currentIndex) != len(previousIndex) {
		return false
	}

	// 遍历并比较每个元素
	for i := range currentIndex {
		if currentIndex[i] != previousIndex[i] {
			return false
		}
	}

	// 如果所有元素都相等，返回 true
	return true
}

// contains 函数用于检查一个元素是否存在于切片中，切片和元素都为 interface{} 类型
func contains(slice []int, value int) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

// filterCurrentIndex 函数根据条件过滤出 currentIndex 中不在 nowPlayingIndex 中的元素
func filterCurrentIndex(currentIndex []int, nowPlayingIndex []int) []int {
	var result []int
	for _, value := range currentIndex {
		// 如果 value 不在 nowPlayingIndex 中，添加到结果切片中
		if !contains(nowPlayingIndex, value) {
			result = append(result, value)
		}
	}
	return result
}

func handleLyricsChange(highlightedLyric []int) interface{} {
	addedLyrics := filterCurrentIndex(highlightedLyric, previousIndex)
	removedLyrics := filterCurrentIndex(previousIndex, highlightedLyric)
	for j := 0; j < len(removedLyrics); j++ {
		log.Println("删除了歌词", removedLyrics[j])
		application.Get().EmitEvent("amll_lyrics_remove", removedLyrics[j])
	}
	for i := 0; i < len(addedLyrics); i++ {
		log.Println("添加了歌词", addedLyrics[i])
		application.Get().EmitEvent("amll_lyrics_add", addedLyrics[i])
	}

	return nil
}

func getLyrics(mid string) ([]LyricLine, string, error) {
	qrcl, err := QQMusicDecoder.GetLyricsByMid(mid)
	if err != nil {
		return nil, "", fmt.Errorf("获取歌词失败")
	}
	//log.Println(qrcl)

	lyricLines, offset, err := ParseLyric(qrcl)
	if err != nil {
		return nil, "", fmt.Errorf("解析歌词失败")
	}
	// 应用offset
	if offset != 0 {
		for i := range lyricLines {
			lyricLines[i].StartTime += uint64(offset)
			lyricLines[i].EndTime += uint64(offset)
			for j := range lyricLines[i].Words {
				lyricLines[i].Words[j].StartTime += uint64(offset)
				lyricLines[i].Words[j].EndTime += uint64(offset)
			}
		}
	}
	return lyricLines, qrcl.Lyrics, nil
}
