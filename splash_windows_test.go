package main

import (
	"runtime"
	"testing"
)

func TestStartupSplashOpensAndCloses(t *testing.T) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	splash := openSplash()
	defer splash.close()
	if splash.window == 0 {
		t.Fatal("a tela de abertura não foi criada")
	}
	window := splash.window
	if visible, _, _ := user32.NewProc("IsWindowVisible").Call(window); visible == 0 {
		t.Fatal("a tela de abertura não ficou visível")
	}
	splash.close()
	if exists, _, _ := user32.NewProc("IsWindow").Call(window); exists != 0 {
		t.Fatal("a tela de abertura permaneceu aberta")
	}
}
