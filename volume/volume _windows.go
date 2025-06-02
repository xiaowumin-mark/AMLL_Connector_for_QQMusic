//go:build windows
// +build windows

package volume

/*
#cgo CFLAGS: -DUNICODE
#cgo LDFLAGS: -lole32 -luuid

#define INITGUID
#include <initguid.h>
#include <windows.h>
#include <audiopolicy.h>
#include <mmdeviceapi.h>
#include <endpointvolume.h>
#include <tlhelp32.h>
#include <stdbool.h>
#include <objbase.h>

ISimpleAudioVolume* gVolume = NULL;

DWORD GetPidFromHwnd(HWND hwnd) {
	DWORD pid = 0;
	GetWindowThreadProcessId(hwnd, &pid);
	return pid;
}

bool InitVolumeByPid(DWORD pid) {
	IMMDeviceEnumerator* pEnumerator = NULL;
	IMMDevice* pDevice = NULL;
	IAudioSessionManager2* pSessionManager = NULL;
	IAudioSessionEnumerator* pSessionEnumerator = NULL;

	HRESULT hr = CoCreateInstance(&CLSID_MMDeviceEnumerator, NULL, CLSCTX_ALL,
		&IID_IMMDeviceEnumerator, (void**)&pEnumerator);
	if (FAILED(hr)) return false;

	hr = pEnumerator->lpVtbl->GetDefaultAudioEndpoint(pEnumerator, eRender, eMultimedia, &pDevice);
	if (FAILED(hr)) goto cleanup;

	hr = pDevice->lpVtbl->Activate(pDevice, &IID_IAudioSessionManager2, CLSCTX_ALL, NULL, (void**)&pSessionManager);
	if (FAILED(hr)) goto cleanup;

	hr = pSessionManager->lpVtbl->GetSessionEnumerator(pSessionManager, &pSessionEnumerator);
	if (FAILED(hr)) goto cleanup;

	int count = 0;
	pSessionEnumerator->lpVtbl->GetCount(pSessionEnumerator, &count);
	for (int i = 0; i < count; ++i) {
		IAudioSessionControl* pControl = NULL;
		IAudioSessionControl2* pControl2 = NULL;

		hr = pSessionEnumerator->lpVtbl->GetSession(pSessionEnumerator, i, &pControl);
		if (FAILED(hr)) continue;

		hr = pControl->lpVtbl->QueryInterface(pControl, &IID_IAudioSessionControl2, (void**)&pControl2);
		if (FAILED(hr)) {
			pControl->lpVtbl->Release(pControl);
			continue;
		}
		pControl->lpVtbl->Release(pControl);

		DWORD sessionPid = 0;
		hr = pControl2->lpVtbl->GetProcessId(pControl2, &sessionPid);
		if (FAILED(hr)) {
			pControl2->lpVtbl->Release(pControl2);
			continue;
		}

		if (sessionPid == pid) {
			ISimpleAudioVolume* pVolume = NULL;
			hr = pControl2->lpVtbl->QueryInterface(pControl2, &IID_ISimpleAudioVolume, (void**)&pVolume);
			pControl2->lpVtbl->Release(pControl2);
			if (SUCCEEDED(hr)) {
				gVolume = pVolume;
				goto cleanup;
			}
		}
		pControl2->lpVtbl->Release(pControl2);
	}

cleanup:
	if (pSessionEnumerator) pSessionEnumerator->lpVtbl->Release(pSessionEnumerator);
	if (pSessionManager) pSessionManager->lpVtbl->Release(pSessionManager);
	if (pDevice) pDevice->lpVtbl->Release(pDevice);
	if (pEnumerator) pEnumerator->lpVtbl->Release(pEnumerator);

	return gVolume != NULL;
}

bool SetVolume(float volume) {
	if (!gVolume) return false;
	HRESULT hr = gVolume->lpVtbl->SetMasterVolume(gVolume, volume, NULL);
	return SUCCEEDED(hr);
}

void ReleaseVolume() {
	if (gVolume) {
		gVolume->lpVtbl->Release(gVolume);
		gVolume = NULL;
	}
}
*/
import "C"
import (
	"unsafe"
)

func InitByHwnd(hwnd uintptr) bool {
	C.CoInitialize(nil)
	pid := C.GetPidFromHwnd((C.HWND)(unsafe.Pointer(hwnd)))
	return bool(C.InitVolumeByPid(pid))
}

func SetProcessVolume(volume float32) bool {
	return bool(C.SetVolume(C.float(volume)))
}

func Release() {
	C.ReleaseVolume()
	C.CoUninitialize()
}
