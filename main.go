package main

import (
	//"AMLL_Connector_for_QQMusic/volume"
	"AMLL_Connector_for_QQMusic/helper"
	"embed"
	_ "embed"
	"encoding/json"
	"log"
	"math"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"

	"time"

	"github.com/gorilla/websocket"
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed all:icon
var iconFS embed.FS

//go:embed smtc-center.exe
var smtcExec embed.FS

//go:embed AppAudioControl.dll
var AppAudioControlFS embed.FS

var mainWindow *application.WebviewWindow

var coon *websocket.Conn
var connLock sync.Mutex
var isCon = false
var previousIndex = make([]int, 0)

type Response struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

type MusicInfo struct {
	Title    string  `json:"title"`
	Artist   string  `json:"artist"`
	Album    string  `json:"album"`
	Status   string  `json:"status"`
	Position float64 `json:"position"`
	Duration float64 `json:"duration"`
}

type MusicInfoMore struct {
	Title    string  `json:"title"`
	Artist   string  `json:"artist"`
	Album    string  `json:"album"`
	Status   string  `json:"status"`
	Duration float64 `json:"duration"`

	Lyric     []LyricLine `json:"lyric"`
	LyricType string      `json:"lyricType"`
	LyricRaw  string      `json:"lyricRaw"`

	Mid     string `json:"mid"`
	Id      string `json:"id"`
	AlbumId string `json:"albumId"`

	Pic []byte `json:"pic"` // 封面图片URL

	RequstRaw []string `json:"requestRaw"`
}

var (
	lastTitle    string      = ""
	lastPosition float64     = 0
	lastStatus   string      = "Paused"
	lastDuration float64     = 0
	lastAlbum    string      = ""
	lastArtist   string      = ""
	paused       bool        = true // 是否暂停
	position     float64     = 0    // 当前进度
	truePosition float64     = 0    // 更准确的进度
	duration     float64     = 0    // 总进度
	nowLyric     []LyricLine        // 当前歌词
	nowImg       string      = ""   // 当前封面

	nowMusicInfoMore MusicInfoMore // 当前音乐信息

	pingpongtimer TimerID
	qqctrl        *QQMusicController
	QQMC          *AudioController

	APPVersion = "0.0.2-beta_1"
)

var lyricWindow *application.WebviewWindow

var appvolume float64 = 1
var targetVolume float64 = 1

func main() {

	// 设置信号通道
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sigChan
		if cmd != nil {
			cmd.Process.Kill() // 确保在程序崩溃时终止子进程
		}
	}()

	// 当程序遇到panic时，打印错误信息
	//defer func() {
	//	if r := recover(); r != nil {
	//		log.Println("程序遇到错误:", r)
	//		if cmd != nil {
	//			cmd.Process.Kill() // 确保在程序崩溃时终止子进程
	//		}
	//		dialog := application.ErrorDialog()
	//		dialog.SetTitle("Error")
	//		dialog.SetMessage("程序遇到异常错误：\n" + r.(string))
	//		dialog.Show()
	//		application.Get().Quit()
	//		// 退出程序
	//	}
	//}()

	app := application.New(application.Options{
		Name:        "AMLL_Connector_for_QQMusic",
		Description: "AMLL_Connector_for_QQMusic Windows",
		Services: []application.Service{
			application.NewService(&GreetService{}),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Windows: application.WindowsOptions{
			DisableQuitOnLastWindowClosed: false,
		},
	})
	app.OnShutdown(func() {
		cmd.Process.Kill()
	})
	app.OnApplicationEvent(events.Common.ApplicationStarted, func(event *application.ApplicationEvent) {
		// 提取smtcExec embed中的smtc-center.exe
		smtcexec, err := smtcExec.ReadFile("smtc-center.exe")
		if err != nil {
			log.Println("Failed to read smtc-center.exe from embed: ", err)
			// 弹窗
			dialog := application.ErrorDialog()
			dialog.SetTitle("Error")
			dialog.SetMessage("SMTC遇到错误！")
			dialog.Show()
			application.Get().Quit()
		}
		// 将smtcexec写入临时文件
		tmpPath := os.TempDir()
		exePath := filepath.Join(tmpPath, "smtc-center-"+time.Now().Format("2006-01-02-15-04-05")+".exe")
		file, err := os.OpenFile(exePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
		if err != nil {
			log.Println("无法写入临时文件！")
			dialog := application.ErrorDialog()
			dialog.SetTitle("Error")
			dialog.SetMessage("SMTC遇到错误！")
			dialog.Show()
			application.Get().Quit()
		}
		defer file.Close()
		if _, err := file.Write(smtcexec); err != nil {
			log.Println("Failed to write to file:", err)
			dialog := application.ErrorDialog()
			dialog.SetTitle("Error")
			dialog.SetMessage("SMTC遇到错误！")
			dialog.Show()
			application.Get().Quit()
		}
		//defer os.Remove(exePath)
		//var hwnd HWND

		//initDll()

		// 提取AppAudioControlFS embed中的AppAudioControl.dll
		AppAudioControlFSE, err := AppAudioControlFS.ReadFile("AppAudioControl.dll")
		if err != nil {
			log.Println("Failed to read AppAudioControl.dll from embed: ", err)
			// 弹窗
			dialog := application.ErrorDialog()
			dialog.SetTitle("Error")
			dialog.SetMessage("AppAudioControl.dll遇到错误！")
			dialog.Show()
			application.Get().Quit()
		}
		// 将smtcexec写入临时文件
		AppAudioControlFSEPath := filepath.Join(tmpPath, "AppAudioControl-"+time.Now().Format("2006-01-02-15-04-05")+".dll")
		file, err = os.OpenFile(AppAudioControlFSEPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
		if err != nil {
			log.Println("无法写入临时文件！")
			dialog := application.ErrorDialog()
			dialog.SetTitle("Error")
			dialog.SetMessage("AppAudioControl.dll遇到错误！")
			dialog.Show()
			application.Get().Quit()
		}

		if _, err := file.Write(AppAudioControlFSE); err != nil {
			log.Println("Failed to write to file:", err)
			dialog := application.ErrorDialog()
			dialog.SetTitle("Error")
			dialog.SetMessage("AppAudioControl.dll遇到错误！")
			dialog.Show()
			application.Get().Quit()
		}
		file.Close()
		log.Println(AppAudioControlFSEPath)

		qqctrl, _, err = NewQQMusicController()
		if err != nil {
			// 弹窗
			dialog := application.ErrorDialog()
			dialog.SetTitle("Error")
			dialog.SetMessage("QQ音乐未启动，请确保QQ音乐正在运行后再运行本程序！")
			dialog.Show()
			app.Quit()
		}
		/*if !Init(qqctrl.pid) {
			dialog := application.ErrorDialog()
			dialog.SetTitle("Error")
			dialog.SetMessage("volumecontrol初始化失败！")
			dialog.Show()
			application.Get().Quit()
		}*/

		//hr, errMsg := InitVCByHwnd(qqctrl.hWnd)
		//log.Printf("InitByHwnd returned 0x%X, errMsg: %s\n", hr, errMsg)
		//if hr != 0 {
		//	dialog := application.ErrorDialog()
		//	dialog.SetTitle("Error")
		//	dialog.SetMessage("volumecontrol初始化失败：" + errMsg)
		//	dialog.Show()
		//	application.Get().Quit()
		//	return
		//}

		QQMC, err = NewAudioController(AppAudioControlFSEPath, uintptr(qqctrl.pid))
		if err != nil {
			log.Println("Failed to create AudioController:", err)
			dialog := application.ErrorDialog()
			dialog.SetTitle("Error")
			dialog.SetMessage("AudioController初始化失败！")
			dialog.Show()
			application.Get().Quit()
		}
		// 设置音量
		//SetProcessVolume(0.3)

		// 程序退出前清理
		//Release()
		// 系统托盘
		systray := app.NewSystemTray()
		systray.SetLabel("AMLL_Connector_for_QQMusic")
		iconBytes, _ := iconFS.ReadFile("icon/icon-t.png")
		systray.SetIcon(iconBytes)
		systray.OnDoubleClick(func() {
			if mainWindow.IsVisible() {
				mainWindow.Hide()
			} else {
				mainWindow.Show()
			}
		})
		menu := application.NewMenu()
		menu.Add("显示主窗口").OnClick(func(ctx *application.Context) {
			mainWindow.Show()
			mainWindow.Focus()
		})
		menu.Add("桌面歌词").OnClick(func(ctx *application.Context) {
			lyricWindow.Show()
		})
		menu.Add("连接AMLL Player").OnClick(func(ctx *application.Context) {
			if isCon {
				DisconnectAmll()
			} else {
				ConnectAmll()
			}
		})
		menu.Add("退出程序").OnClick(func(ctx *application.Context) {

			app.Quit()
		})

		systray.SetMenu(menu)

		// 以下是异步函数

		//连接
		go Smtc(exePath)

		// 循环发送ping包
		go Heartbeat()

		// 读取消息
		go func() {
			for {
				if !isCon {
					time.Sleep(time.Second)
					continue
				}

				if coon == nil {
					time.Sleep(time.Second)
					continue
				}
				//connLock.Lock()
				messageType, message, err := coon.ReadMessage()
				if err != nil {
					//connLock.Unlock()
					log.Println("读取错误:", err)
					continue
				}
				//connLock.Unlock()

				if messageType == websocket.BinaryMessage {
					handleIncomingMessage(coon, message)
				}
			}
		}()

		// 更精准的时间和进度发送
		go func() {
			ticker := time.NewTicker(50 * time.Millisecond)
			defer ticker.Stop()
			for range ticker.C {
				appvolume = Lerp(appvolume, targetVolume, 0.1)
				//log.Println("volume:", volume)
				if math.Round(targetVolume*1000)/1000 != math.Round(appvolume*1000)/1000 {
					//if math.Abs(targetVolume-volume) > 0.01 {
					log.Println("volume:", appvolume)
					//SetVolumeByOle(float32(appvolume))
					//volume.SetProcessVolume(float32(appvolume))
					//ok := SetProcessVolume(float32(appvolume)) // 50%音量
					//log.Println("SetProcessVolume:", ok)

					err := QQMC.SetVolume(float32(appvolume))
					if err != nil {
						log.Println("Failed to set volume:", err)
					}
				}
				if !paused {
					truePosition += 50 // 每20ms增加20ms

					//if truePosition > duration { // qq stmc有bug，当歌曲时长过长时会导致总时长传错
					//	truePosition = 0 // 重置进度
					//}
					if application.Get() != nil {
						application.Get().EmitEvent("amll_play_progress", map[string]interface{}{
							"progress": uint64(truePosition),
							"format":   helper.FormatMilliseconds(int(truePosition)),
						})
					}

					JudgeLyricEvent()
					// 是否期望连接
					if !isCon {
						continue
					}
					if coon != nil {

						sendMessage(coon, &OnPlayProgress{Progress: uint64(truePosition)})
					}

				}
			}
		}()
	})

	mainWindow = app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Title: "AMLL_Connector_for_QQMusic",
		Windows: application.WindowsWindow{
			BackdropType: application.Tabbed,
			DisableIcon:  true,
			Theme:        application.Dark,
		},

		Frameless:      true,
		Width:          800,
		Height:         600,
		BackgroundType: application.BackgroundTypeTranslucent,
		URL:            "/",
	})
	// 窗口创建完成时
	mainWindow.OnWindowEvent(events.Common.WindowClosing, func(e *application.WindowEvent) {
		app.Quit()
	})
	mainWindow.OnWindowEvent(events.Common.WindowShow, func(e *application.WindowEvent) {
		Window_will_update_visibility(mainWindow, "main")
	})
	mainWindow.OnWindowEvent(events.Common.WindowHide, func(e *application.WindowEvent) {
		Window_will_update_visibility(mainWindow, "main")
	})

	lyricWindow = app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Title: "AMLL_Connector_for_QQMusic_lyrics",
		Windows: application.WindowsWindow{
			DisableIcon:                       true,
			DisableFramelessWindowDecorations: true,
			Theme:                             application.Dark,
			HiddenOnTaskbar:                   true,
		},
		Frameless:        true,
		Width:            780,
		Height:           145,
		AlwaysOnTop:      true,
		BackgroundColour: application.NewRGBA(0, 0, 0, 0),
		BackgroundType:   application.BackgroundTypeTransparent,
		URL:              "/#lyrics",
	})
	lyricWindow.Hide()
	lyricWindow.OnWindowEvent(events.Common.WindowHide, func(event *application.WindowEvent) {
		Window_will_update_visibility(lyricWindow, "lyric")
	})
	lyricWindow.OnWindowEvent(events.Common.WindowShow, func(event *application.WindowEvent) {
		Window_will_update_visibility(lyricWindow, "lyric")
	})
	lyricWindow.OnWindowEvent(events.Windows.WebViewNavigationCompleted, func(event *application.WindowEvent) {
		log.Println("lyricWindow navigation completed")

	})
	//mainWindow.Hide()
	//initWindow := app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
	//	Title: "ACFQ-INIT",
	//	Windows: application.WindowsWindow{
	//		BackdropType:    application.Tabbed,
	//		HiddenOnTaskbar: false,
	//		Theme:           application.Dark,
	//
	//	},
	//
	//	Width:          400,
	//	Height:         250,
	//	Frameless:      true,
	//	BackgroundType: application.BackgroundTypeTranslucent,
	//	AlwaysOnTop:    true,
	//	DisableResize:  true,
	//	URL: "/#init",
	//})
	//
	//initWindow.OnWindowEvent(events.Common.WindowRuntimeReady, func(e *application.WindowEvent) {
	//	Animate([]float64{400, 250}, []float64{800, 500}, 10*time.Second, easing.QuadEaseInOut).
	//		OnUpdate(func(values []float64) {
	//			//fmt.Printf("Current Values: [%.2f, %.2f]\n", values[0], values[1])
	//			initWindow.SetSize(int(values[0]), int(values[1]))
	//			// 窗口居中
	//		}).
	//		OnComplete(func() {
	//			fmt.Println("Animation Complete!")
	//			initWindow.Close()
	//			mainWindow.Show()
	//		}).
	//		Start()
	//	time.Sleep(10 * time.Second)
	//})

	//go func() {
	//	Animate([]float64{400, 250}, []float64{800, 500}, 10*time.Second, easing.QuadEaseInOut).
	//		OnUpdate(func(values []float64) {
	//			//fmt.Printf("Current Values: [%.2f, %.2f]\n", values[0], values[1])
	//			//initWindow.SetSize(int(values[0]), int(values[1]))
	//		}).
	//		OnComplete(func() {
	//			fmt.Println("Animation Complete!")
	//			initWindow.Close()
	//			mainWindow.Show()
	//		}).
	//		Start()
	//
	//	time.Sleep(10 * time.Second)
	//}()

	err := app.Run()

	if err != nil {
		log.Fatal(err)
	}
}

func (a *GreetService) ChoseDir(defaultDirectory string) string {
	dialog := application.OpenFileDialog()
	dialog.SetTitle("选择文件夹")
	dialog.CanChooseDirectories(true)
	dialog.CanCreateDirectories(true)
	dialog.SetDirectory(defaultDirectory)
	// Single file selection
	if path, err := dialog.PromptForSingleSelection(); err == nil {
		// Use selected file path
		return path
	} else {
		return ""
	}
}
func Lerp(start, end, t float64) float64 {
	return start + (end-start)*t
}

func Window_will_update_visibility(w *application.WebviewWindow, window_name string) {
	application.Get().EmitEvent("window_will_update_visibility", map[string]interface{}{
		"visible": w.IsVisible(),
		"window":  window_name,
	})
}
