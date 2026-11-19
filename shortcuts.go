package main

func blockedKey(key uint32, ctrl, shift bool) bool {
	switch key {
	case 0x74, 0x7a, 0x7b: // F5, F11, F12
		return true
	case 'R', 'U':
		return ctrl
	case 'I', 'J', 'C':
		return ctrl && shift
	default:
		return false
	}
}
