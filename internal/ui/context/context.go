package context

import "github.com/Fomiller/mixify/internal/ui/views"

type ProgramContext struct {
	View         views.View
	ScreenHeight int
	ScreenWidth  int
}
