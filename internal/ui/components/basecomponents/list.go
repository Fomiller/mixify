package basecomponents

import "github.com/Fomiller/mixify/internal/ui/messages"

type List struct {
	Focused bool
	Cursor  int
	Width   int
	Height  int
	Status  int
	Err     messages.ErrMsg
}
