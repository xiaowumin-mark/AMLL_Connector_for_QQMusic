package main

import (
	"fmt"
	"syscall"
	"unsafe"
)

//var (
//	modVolCtrl     = syscall.NewLazyDLL("volume_control.dll") // 替换为你的DLL名字
//	procInitByHwnd = modVolCtrl.NewProc("InitByHwnd")
//	procSetVolume  = modVolCtrl.NewProc("SetProcessVolume")
//	procRelease    = modVolCtrl.NewProc("Release")
//)
//
//func InitVCByHwnd(hwnd uintptr) (hr uintptr, errMsg string) {
//	const errBufLen = 512
//	errBuf := make([]uint16, errBufLen) // wchar_t缓冲区
//
//	ret, _, _ := procInitByHwnd.Call(
//		hwnd,
//		uintptr(unsafe.Pointer(&errBuf[0])),
//		uintptr(errBufLen),
//	)
//	hr = ret
//
//	// 查找第一个0终止符位置，转成Go字符串
//	for i, c := range errBuf {
//		if c == 0 {
//			errMsg = syscall.UTF16ToString(errBuf[:i])
//			break
//		}
//	}
//	return
//}
//
//func SetProcessVolume(volume float32) bool {
//	bits := math.Float32bits(volume)
//	ret, _, _ := procSetVolume.Call(uintptr(bits))
//	return ret != 0
//}
//
//func ReleaseVC() {
//	procRelease.Call()
//}

type AudioController struct {
	dll         *syscall.LazyDLL
	initProc    *syscall.LazyProc
	playProc    *syscall.LazyProc
	pauseProc   *syscall.LazyProc
	nextProc    *syscall.LazyProc
	prevProc    *syscall.LazyProc
	setVolProc  *syscall.LazyProc
	getVolProc  *syscall.LazyProc
	releaseProc *syscall.LazyProc

	hwnd uintptr
}

// NewAudioController 载入DLL，调用Init初始化
func NewAudioController(dllPath string, hwnd uintptr) (*AudioController, error) {
	dll := syscall.NewLazyDLL(dllPath)

	initProc := dll.NewProc("Init")
	playProc := dll.NewProc("Play")
	pauseProc := dll.NewProc("Pause")
	nextProc := dll.NewProc("Next")
	prevProc := dll.NewProc("Prev")
	setVolProc := dll.NewProc("SetVolume")
	getVolProc := dll.NewProc("GetVolume")
	releaseProc := dll.NewProc("Release")

	ac := &AudioController{
		dll:         dll,
		initProc:    initProc,
		playProc:    playProc,
		pauseProc:   pauseProc,
		nextProc:    nextProc,
		prevProc:    prevProc,
		setVolProc:  setVolProc,
		getVolProc:  getVolProc,
		releaseProc: releaseProc,
		hwnd:        hwnd,
	}

	// 调用 Init(hwnd)
	r1, _, err := initProc.Call(hwnd)
	if r1 == 0 {
		// Init无返回值，可以忽略
	}
	if err != syscall.Errno(0) {
		return nil, fmt.Errorf("Init failed: %w", err)
	}

	return ac, nil
}

func (ac *AudioController) Play() error {
	_, _, err := ac.playProc.Call()
	if err != syscall.Errno(0) {
		return err
	}
	return nil
}

func (ac *AudioController) Pause() error {
	_, _, err := ac.pauseProc.Call()
	if err != syscall.Errno(0) {
		return err
	}
	return nil
}

func (ac *AudioController) Next() error {
	_, _, err := ac.nextProc.Call()
	if err != syscall.Errno(0) {
		return err
	}
	return nil
}

func (ac *AudioController) Prev() error {
	_, _, err := ac.prevProc.Call()
	if err != syscall.Errno(0) {
		return err
	}
	return nil
}

func (ac *AudioController) SetVolume(vol float32) error {
	// float32 转 uintptr 需要先转为 uintptr (用 unsafe)
	ptr := uintptr(*(*uint32)(unsafe.Pointer(&vol)))
	_, _, err := ac.setVolProc.Call(ptr)
	if err != syscall.Errno(0) {
		return err
	}
	return nil
}

func (ac *AudioController) GetVolume() (float32, error) {
	ret, _, err := ac.getVolProc.Call()
	if err != syscall.Errno(0) {
		return 0, err
	}
	vol := *(*float32)(unsafe.Pointer(&ret))
	return vol, nil
}

func (ac *AudioController) Release() error {
	_, _, err := ac.releaseProc.Call()
	if err != syscall.Errno(0) {
		return err
	}
	return nil
}
