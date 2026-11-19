package main

import (
	"bytes"
	_ "embed"
	"image"
	"image/draw"
	"image/png"
	"syscall"
	"unsafe"
)

//go:embed assets/splash.png
var splashArtwork []byte

type splashBitmapHeader struct {
	size                             uint32
	width, height                    int32
	planes, bitCount                 uint16
	compression, sizeImage           uint32
	xPixelsPerMeter, yPixelsPerMeter int32
	colorsUsed, colorsImportant      uint32
}

type splashPaint struct {
	dc                   uintptr
	erase                int32
	area                 rect
	restore, incremental int32
	reserved             [32]byte
}

type splashClass struct {
	style                              uint32
	callback                           uintptr
	classExtra, windowExtra            int32
	instance, icon, cursor, background uintptr
	menu, name                         *uint16
}

type startupSplash struct {
	window  uintptr
	refresh *syscall.LazyProc
}

func openSplash() startupSplash {
	artwork, err := png.Decode(bytes.NewReader(splashArtwork))
	if err != nil {
		return startupSplash{}
	}
	bounds := artwork.Bounds()
	pixels := image.NewRGBA(bounds)
	draw.Draw(pixels, bounds, artwork, bounds.Min, draw.Src)
	for i := 0; i < len(pixels.Pix); i += 4 {
		pixels.Pix[i], pixels.Pix[i+2] = pixels.Pix[i+2], pixels.Pix[i]
	}
	bitmap := splashBitmapHeader{
		size:  uint32(unsafe.Sizeof(splashBitmapHeader{})),
		width: int32(bounds.Dx()), height: -int32(bounds.Dy()), planes: 1, bitCount: 32,
	}

	metrics := user32.NewProc("GetSystemMetrics")
	screenWidth, _, _ := metrics.Call(0)
	screenHeight, _, _ := metrics.Call(1)
	dpi, _, _ := user32.NewProc("GetDpiForSystem").Call()
	width := min(480*max(96, dpi)/96, screenWidth*9/10, screenHeight*9/10*uintptr(bounds.Dx())/uintptr(bounds.Dy()))
	height := width * uintptr(bounds.Dy()) / uintptr(bounds.Dx())
	if width == 0 || height == 0 {
		return startupSplash{}
	}

	begin := user32.NewProc("BeginPaint")
	end := user32.NewProc("EndPaint")
	defaultWindow := user32.NewProc("DefWindowProcW")
	gdi := syscall.NewLazyDLL("gdi32.dll")
	stretch := gdi.NewProc("StretchDIBits")
	filter := gdi.NewProc("SetStretchBltMode")
	callback := syscall.NewCallback(func(window uintptr, message uint32, wParam, lParam uintptr) uintptr {
		if message == 0x000f { // WM_PAINT
			var paint splashPaint
			dc, _, _ := begin.Call(window, uintptr(unsafe.Pointer(&paint)))
			if dc != 0 {
				filter.Call(dc, 4) // HALFTONE
				stretch.Call(dc, 0, 0, width, height, 0, 0, uintptr(bounds.Dx()), uintptr(bounds.Dy()),
					uintptr(unsafe.Pointer(&pixels.Pix[0])), uintptr(unsafe.Pointer(&bitmap)), 0, 0x00cc0020)
				end.Call(window, uintptr(unsafe.Pointer(&paint)))
			}
			return 0
		}
		if message == 0x0010 { // WM_CLOSE: o splash encerra junto com a inicialização
			return 0
		}
		result, _, _ := defaultWindow.Call(window, uintptr(message), wParam, lParam)
		return result
	})
	instance, _, _ := kernel32.NewProc("GetModuleHandleW").Call(0)
	className := syscall.StringToUTF16Ptr("SifyPDVStartup")
	class := splashClass{callback: callback, instance: instance, name: className}
	if atom, _, _ := user32.NewProc("RegisterClassW").Call(uintptr(unsafe.Pointer(&class))); atom == 0 {
		return startupSplash{}
	}

	const splashFlags = 0x08000000 | 0x00000080 | 0x00000008 // no activate, tool window, topmost
	title := syscall.StringToUTF16Ptr("Sify Comércio — Abrindo o caixa")
	window, _, _ := user32.NewProc("CreateWindowExW").Call(splashFlags, uintptr(unsafe.Pointer(className)),
		uintptr(unsafe.Pointer(title)), 0x80000000, (screenWidth-width)/2, (screenHeight-height)/2,
		width, height, 0, 0, instance, 0)
	if window == 0 {
		return startupSplash{}
	}
	user32.NewProc("ShowWindow").Call(window, 4) // SW_SHOWNOACTIVATE
	splash := startupSplash{window: window, refresh: user32.NewProc("UpdateWindow")}
	splash.repaint()
	return splash
}

func (s *startupSplash) repaint() {
	if s.window != 0 {
		s.refresh.Call(s.window)
	}
}

func (s *startupSplash) close() {
	if s.window != 0 {
		user32.NewProc("DestroyWindow").Call(s.window)
		s.window = 0
	}
}
