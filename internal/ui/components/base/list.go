package base

import "github.com/Fomiller/mixify/internal/ui/messages"

type List struct {
	Focused bool
	Cursor  int
	Status  int
	Err     messages.ErrMsg
}
