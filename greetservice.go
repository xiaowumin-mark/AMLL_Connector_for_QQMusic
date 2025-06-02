package main

import (
	"time"
)

type GreetService struct{}

// 这里加了激活窗口
func (a *GreetService) TriggerSnapLayout() {
	mainWindow.Show()
	mainWindow.Focus()
	// 等待一点点
	go func() {
		time.Sleep(100 * time.Millisecond)
		simulateWinZ()
	}()
}

func (a *GreetService) GetNowLyrics() []LyricLine {
	return nowLyric
}

func (a *GreetService) GetAllSmtc() []TimelineEventData {
	return nowSmtcData
}

func (a *GreetService) ToggleLyricWindowShow() bool {
	if lyricWindow.IsVisible() {
		lyricWindow.Hide()
		return false
	} else {
		lyricWindow.Show()
		lyricWindow.Focus()
		return true
	}
}

func (a *GreetService) IsLyricWindowShow() bool {
	return lyricWindow.IsVisible()
}

func (a *GreetService) ShowLyricWindow() bool {
	lyricWindow.Show()
	return lyricWindow.IsVisible()
}
func (a *GreetService) HideLyricWindow() bool {
	lyricWindow.Hide()
	return lyricWindow.IsVisible()
}

func (a *GreetService) GetNowMusicInfo() MusicInfoMore {
	return nowMusicInfoMore
}
