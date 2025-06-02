package main

import (
	"fmt"
	"syscall"
	"unsafe"
)

type ACCENT_STATE int32

const (
	ACCENT_DISABLED                   = 0
	ACCENT_ENABLE_GRADIENT            = 1
	ACCENT_ENABLE_TRANSPARENTGRADIENT = 2
	ACCENT_ENABLE_BLURBEHIND          = 3
	ACCENT_ENABLE_ACRYLICBLURBEHIND   = 4
	ACCENT_ENABLE_HOSTBACKDROP        = 5
)

// WinAPI常量
const (
	KEYEVENTF_EXTENDEDKEY = 0x0001
	KEYEVENTF_KEYUP       = 0x0002
	VK_LWIN               = 0x5B
	VK_Z                  = 0x5A
)

// user32.dll API
var (
	user32               = syscall.NewLazyDLL("user32.dll")
	procKeybd_event      = user32.NewProc("keybd_event")
	procSetForegroundWnd = user32.NewProc("SetForegroundWindow")
	procFindWindow       = user32.NewProc("FindWindowW")
)

// 模拟键盘事件
func keybdEvent(bVk byte, bScan byte, dwFlags uint32, dwExtraInfo uintptr) {
	procKeybd_event.Call(
		uintptr(bVk),
		uintptr(bScan),
		uintptr(dwFlags),
		dwExtraInfo,
	)
}

// 模拟 Win+Z
func simulateWinZ() {
	keybdEvent(VK_LWIN, 0, 0, 0)
	keybdEvent(VK_Z, 0, 0, 0)
	keybdEvent(VK_Z, 0, KEYEVENTF_KEYUP, 0)
	keybdEvent(VK_LWIN, 0, KEYEVENTF_KEYUP, 0)

}

var (
	procSendMessage = user32.NewProc("SendMessageW")
)

const (
	WM_APPCOMMAND               = 0x0319
	APPCOMMAND_MEDIA_PLAY       = 46
	APPCOMMAND_MEDIA_PAUSE      = 47
	APPCOMMAND_MEDIA_PLAY_PAUSE = 14
	APPCOMMAND_MEDIA_NEXTTRACK  = 11
	APPCOMMAND_MEDIA_PREVTRACK  = 12
)

var (
	procEnumWindows              = user32.NewProc("EnumWindows")
	procGetWindowThreadProcessId = user32.NewProc("GetWindowThreadProcessId")
	procGetClassName             = user32.NewProc("GetClassNameW")
	procGetWindowText            = user32.NewProc("GetWindowTextW")
	kernel32                     = syscall.NewLazyDLL("kernel32.dll")
	procCreateToolhelp32Snapshot = kernel32.NewProc("CreateToolhelp32Snapshot")
	procProcess32First           = kernel32.NewProc("Process32FirstW")
	procProcess32Next            = kernel32.NewProc("Process32NextW")
)

type (
	HANDLE uintptr
	HWND   HANDLE
)

type PROCESSENTRY32 struct {
	Size            uint32
	CntUsage        uint32
	ProcessID       uint32
	DefaultHeapID   uintptr
	ModuleID        uint32
	CntThreads      uint32
	ParentProcessID uint32
	PriClassBase    int32
	Flags           uint32
	ExeFile         [260]uint16
}

// 查找QQ音乐窗口句柄
func FindQQMusicWindow() (HWND, error) {
	// 1. 先通过进程名获取进程ID
	pid, err := FindProcessIDByName("QQMusic.exe")
	if err != nil {
		return 0, fmt.Errorf("找不到QQMusic进程: %v", err)
	}

	// 2. 根据进程ID和窗口类名查找窗口
	hWnd, err := FindWindowByPIDAndClass(pid, "TXGuiFoundation")
	if err != nil {
		return 0, fmt.Errorf("找不到QQMusic窗口: %v", err)
	}

	return hWnd, nil
}

// 通过进程名查找进程ID
func FindProcessIDByName(name string) (uint32, error) {
	snapshot, _, _ := procCreateToolhelp32Snapshot.Call(
		0x00000002, // TH32CS_SNAPPROCESS
		0,
	)
	if snapshot == 0 {
		return 0, fmt.Errorf("创建进程快照失败")
	}
	defer syscall.CloseHandle(syscall.Handle(snapshot))

	var entry PROCESSENTRY32
	entry.Size = uint32(unsafe.Sizeof(entry))

	ret, _, _ := procProcess32First.Call(
		snapshot,
		uintptr(unsafe.Pointer(&entry)),
	)
	if ret == 0 {
		return 0, fmt.Errorf("获取第一个进程失败")
	}

	for {
		// 转换为Go字符串
		exeName := syscall.UTF16ToString(entry.ExeFile[:])
		if exeName == name {
			return entry.ProcessID, nil
		}

		ret, _, _ := procProcess32Next.Call(
			snapshot,
			uintptr(unsafe.Pointer(&entry)),
		)
		if ret == 0 {
			break
		}
	}

	return 0, fmt.Errorf("进程未找到")
}

// 根据进程ID和窗口类名查找窗口
func FindWindowByPIDAndClass(pid uint32, className string) (HWND, error) {
	var result HWND
	cb := syscall.NewCallback(func(hWnd HWND, lParam uintptr) uintptr {
		// 获取窗口进程ID
		var windowPID uint32
		procGetWindowThreadProcessId.Call(
			uintptr(hWnd),
			uintptr(unsafe.Pointer(&windowPID)),
		)

		if windowPID == pid {
			// 获取窗口类名
			buf := make([]uint16, 256)
			procGetClassName.Call(
				uintptr(hWnd),
				uintptr(unsafe.Pointer(&buf[0])),
				uintptr(len(buf)),
			)

			// 比较类名
			if syscall.UTF16ToString(buf) == className {
				result = hWnd
				return 0 // 停止枚举
			}
		}
		return 1 // 继续枚举
	})

	// 枚举所有顶层窗口
	procEnumWindows.Call(cb, 0)

	if result == 0 {
		return 0, fmt.Errorf("未找到匹配的窗口")
	}

	return result, nil
}

// QQMusicController QQ音乐控制器
type QQMusicController struct {
	hWnd uintptr
	pid  uint32
}

// NewQQMusicController 创建QQ音乐控制器实例
func NewQQMusicController() (*QQMusicController, HWND, error) {
	// QQ音乐主窗口类名可能因版本不同而变化
	// 常见的有："QQMusic_Daemon_Wnd", "QQMusic", "QQMusicMiniPlayer"
	//hwnd, _, _ := procFindWindow.Call(0,
	//	uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("QQMusic"))),
	//)
	hwnd, err := FindQQMusicWindow()
	if err != nil {
		return nil, 0, fmt.Errorf("QQ音乐窗口未找到: %v", err)
	}
	hwndPtr := uintptr(hwnd)
	if hwndPtr == 0 {
		return nil, 0, fmt.Errorf("QQ音乐窗口未找到，请确保QQ音乐正在运行")
	}

	// 获取并缓存 PID
	pid, err := GetProcessIDFromHWND(hwndPtr)
	if err != nil {
		return nil, 0, fmt.Errorf("无法获取QQ音乐进程 PID: %v", err)
	}

	return &QQMusicController{
		hWnd: hwndPtr,
		pid:  pid,
	}, HWND(hwndPtr), nil
}

func GetProcessIDFromHWND(hwnd uintptr) (uint32, error) {
	var pid uint32
	procGetWindowThreadProcessId.Call(hwnd, uintptr(unsafe.Pointer(&pid)))
	return pid, nil
}

// sendAppCommand 发送APPCOMMAND消息
func (c *QQMusicController) sendAppCommand(cmd uintptr) error {
	_, _, err := procSendMessage.Call(
		c.hWnd,
		WM_APPCOMMAND,
		0,
		cmd<<16, // 命令需要左移16位
	)
	return err
}

// PlayPause 播放/暂停切换
func (c *QQMusicController) PlayPause() error {
	return c.sendAppCommand(APPCOMMAND_MEDIA_PLAY_PAUSE)
}

// Play 播放
func (c *QQMusicController) Play() error {
	return c.sendAppCommand(APPCOMMAND_MEDIA_PLAY)
}

// Pause 暂停
func (c *QQMusicController) Pause() error {
	return c.sendAppCommand(APPCOMMAND_MEDIA_PAUSE)
}

// NextTrack 下一首
func (c *QQMusicController) NextTrack() error {
	return c.sendAppCommand(APPCOMMAND_MEDIA_NEXTTRACK)
}

// PrevTrack 上一首
func (c *QQMusicController) PrevTrack() error {
	return c.sendAppCommand(APPCOMMAND_MEDIA_PREVTRACK)
}
