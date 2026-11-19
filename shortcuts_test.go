package main

import "testing"

func TestBlockedKey(t *testing.T) {
	for _, key := range []uint32{0x74, 0x7a, 0x7b} {
		if !blockedKey(key, false, false) {
			t.Errorf("atalho 0x%x liberado", key)
		}
	}
	if !blockedKey('I', true, true) || !blockedKey('R', true, false) {
		t.Error("atalhos de inspeção ou recarga liberados")
	}
	if blockedKey('C', true, false) || blockedKey(0x73, false, false) || blockedKey('R', false, false) {
		t.Error("atalho comum ou tecla de digitação bloqueados")
	}
}
