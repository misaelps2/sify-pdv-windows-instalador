package main

import (
	"errors"
	"syscall"
	"time"
	"unsafe"
)

var (
	user32           = syscall.NewLazyDLL("user32.dll")
	kernel32         = syscall.NewLazyDLL("kernel32.dll")
	getWindowPID     = user32.NewProc("GetWindowThreadProcessId")
	getForeground    = user32.NewProc("GetForegroundWindow")
	getAsyncKeyState = user32.NewProc("GetAsyncKeyState")
)

func lockLauncher() (uintptr, bool) {
	name := syscall.StringToUTF16Ptr("Local\\SifyPDV.Fullscreen")
	handle, _, err := kernel32.NewProc("CreateMutexW").Call(0, 0, uintptr(unsafe.Pointer(name)))
	return handle, !errors.Is(err, syscall.ERROR_ALREADY_EXISTS)
}

func windowPID(window uintptr) uint32 {
	var pid uint32
	getWindowPID.Call(window, uintptr(unsafe.Pointer(&pid)))
	return pid
}

func findWindow(pid uint32) uintptr {
	var result uintptr
	// EnumWindows invokes this callback synchronously.
	enumeratePDVWindows.Call(windowEnumerator, uintptr(unsafe.Pointer(&windowSearch{pid: pid, result: &result})))
	return result
}

type windowSearch struct {
	pid    uint32
	result *uintptr
}

var (
	enumeratePDVWindows = user32.NewProc("EnumWindows")
	windowEnumerator    = syscall.NewCallback(func(handle uintptr, search *windowSearch) uintptr {
		if windowPID(handle) != search.pid {
			return 1
		}
		visible, _, _ := user32.NewProc("IsWindowVisible").Call(handle)
		if visible == 0 {
			return 1
		}
		var class [64]uint16
		user32.NewProc("GetClassNameW").Call(handle, uintptr(unsafe.Pointer(&class[0])), uintptr(len(class)))
		if syscall.UTF16ToString(class[:]) != "Chrome_WidgetWin_1" {
			return 1
		}
		*search.result = handle
		return 0
	})
)

type rect struct{ left, top, right, bottom int32 }

func fullscreen(window uintptr) error {
	monitor, _, _ := user32.NewProc("MonitorFromWindow").Call(window, 2)
	info := struct {
		size          uint32
		monitor, work rect
		flags         uint32
	}{}
	info.size = uint32(unsafe.Sizeof(info))
	if ok, _, _ := user32.NewProc("GetMonitorInfoW").Call(monitor, uintptr(unsafe.Pointer(&info))); ok == 0 {
		return errors.New("Não foi possível identificar o monitor do PDV.")
	}
	user32.NewProc("ShowWindow").Call(window, 9) // SW_RESTORE
	const styleIndex = ^uintptr(15)              // GWL_STYLE = -16
	const frameStyle = 0x00c00000 | 0x00040000   // WS_CAPTION | WS_THICKFRAME
	style, _, _ := user32.NewProc("GetWindowLongW").Call(window, styleIndex)
	user32.NewProc("SetWindowLongW").Call(window, styleIndex, style&^frameStyle)
	current, _, _ := user32.NewProc("GetWindowLongW").Call(window, styleIndex)
	if current&frameStyle != 0 {
		return errors.New("Não foi possível ocultar a barra do PDV.")
	}
	bounds := info.monitor
	const frameChanged = 0x20
	if ok, _, _ := user32.NewProc("SetWindowPos").Call(window, 0, uintptr(bounds.left), uintptr(bounds.top), uintptr(bounds.right-bounds.left), uintptr(bounds.bottom-bounds.top), frameChanged); ok == 0 {
		return errors.New("Não foi possível abrir o PDV em tela cheia.")
	}
	return nil
}

func keyDown(vk uintptr) bool {
	state, _, _ := getAsyncKeyState.Call(vk)
	return state&0x8000 != 0
}

func runWindow(pid uint32, splash *startupSplash) error {
	var window uintptr
	for deadline := time.Now().Add(20 * time.Second); time.Now().Before(deadline); {
		window = findWindow(pid)
		if window != 0 {
			break
		}
		splash.repaint()
		time.Sleep(100 * time.Millisecond)
	}
	if window == 0 {
		return errors.New("A janela do PDV não apareceu. Feche a janela anterior do Sify PDV e abra pelo atalho novamente.")
	}
	if err := fullscreen(window); err != nil {
		return err
	}

	splash.close()

	next := user32.NewProc("CallNextHookEx")
	callback := syscall.NewCallback(func(code int32, event uintptr, vk *uint32) uintptr {
		if code == 0 && (event == 0x100 || event == 0x104) { // WM_KEYDOWN, WM_SYSKEYDOWN
			foreground, _, _ := getForeground.Call()
			if foreground == window && windowPID(foreground) == pid {
				if blockedKey(*vk, keyDown(0x11), keyDown(0x10)) {
					return 1
				}
			}
		}
		result, _, _ := next.Call(0, uintptr(code), event, uintptr(unsafe.Pointer(vk)))
		return result
	})
	module, _, _ := kernel32.NewProc("GetModuleHandleW").Call(0)
	hook, _, _ := user32.NewProc("SetWindowsHookExW").Call(13, callback, module, 0) // WH_KEYBOARD_LL
	if hook == 0 {
		return errors.New("Não foi possível desativar os atalhos do PDV.")
	}
	defer user32.NewProc("UnhookWindowsHookEx").Call(hook)

	timer, _, _ := user32.NewProc("SetTimer").Call(0, 0, 500, 0)
	if timer == 0 {
		return errors.New("Não foi possível acompanhar a janela do PDV.")
	}
	defer user32.NewProc("KillTimer").Call(0, timer)

	var msg struct {
		window         uintptr
		message        uint32
		wParam, lParam uintptr
		time           uint32
		x, y           int32
		private        uint32
	}
	getMessage := user32.NewProc("GetMessageW")
	for {
		result, _, err := getMessage.Call(uintptr(unsafe.Pointer(&msg)), 0, 0, 0)
		if int32(result) < 0 {
			return err
		}
		if result == 0 || windowPID(window) != pid {
			return nil
		}
	}
}
