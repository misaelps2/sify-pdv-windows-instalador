package main

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"syscall"
	"unsafe"
)

const pdvURL = "https://pdv.sifysistemas.com/"

func main() {
	if err := launch(); err != nil {
		showError(err.Error())
		os.Exit(1)
	}
}

func launch() error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	mutex, acquired := lockLauncher()
	if mutex == 0 {
		return errors.New("Não foi possível iniciar o controle da janela do PDV.")
	}
	defer syscall.CloseHandle(syscall.Handle(mutex))
	if !acquired {
		return nil
	}

	local := os.Getenv("LOCALAPPDATA")
	if local == "" {
		return errors.New("O Windows não informou a pasta de dados do usuário.")
	}

	edge, err := findEdge(local)
	if err != nil {
		return err
	}

	profile := filepath.Join(local, "Sify PDV", "Perfil")
	if err := os.MkdirAll(profile, 0700); err != nil {
		return errors.New("Não foi possível abrir a pasta de dados do Sify PDV: " + err.Error())
	}

	// Use monitor coordinates consistently before creating the splash or sizing Edge.
	user32.NewProc("SetThreadDpiAwarenessContext").Call(^uintptr(3)) // PER_MONITOR_AWARE_V2
	splash := openSplash()
	defer splash.close()
	cmd := exec.Command(edge,
		"--user-data-dir="+profile,
		"--app="+pdvURL,
		"--no-first-run",
		"--no-default-browser-check",
		"--disable-background-mode",
	)
	if err := cmd.Start(); err != nil {
		return errors.New("Não foi possível iniciar o Microsoft Edge: " + err.Error())
	}
	defer cmd.Process.Release()
	return runWindow(uint32(cmd.Process.Pid), &splash)
}

func findEdge(local string) (string, error) {
	for _, root := range []string{os.Getenv("ProgramFiles(x86)"), os.Getenv("ProgramW6432"), os.Getenv("ProgramFiles"), local} {
		if root == "" {
			continue
		}
		path := filepath.Join(root, "Microsoft", "Edge", "Application", "msedge.exe")
		if info, err := os.Stat(path); err == nil && !info.IsDir() {
			return path, nil
		}
	}
	return "", errors.New("Microsoft Edge não encontrado. Instale o Microsoft Edge e abra o Sify PDV novamente.")
}

func showError(message string) {
	title := syscall.StringToUTF16Ptr("Sify PDV")
	body := syscall.StringToUTF16Ptr(message)
	syscall.NewLazyDLL("user32.dll").NewProc("MessageBoxW").Call(
		0, uintptr(unsafe.Pointer(body)), uintptr(unsafe.Pointer(title)), 0x10,
	)
}
