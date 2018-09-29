package media

import (
	"syscall/js"
)

// KeyEvent ...
type KeyEvent string

// Key events ...
const (
	KeyPress KeyEvent = "keypress"
	KeyUp    KeyEvent = "keyup"
	KeyDown  KeyEvent = "keydown"
)

// OnKey ...
func OnKey(ev KeyEvent, f func(key string)) {
	doc := js.Global().Get("document")
	doc.Call("addEventListener", string(ev), js.NewCallback(func(args []js.Value) {
		if len(args) == 0 {
			return
		}
		f(args[0].Get("key").String())
	}))
}
