package context

import "github.com/Fomiller/mixify/internal/ui/views"

type ProgramContext struct {
	View         views.View
	ScreenWidth  int
	ScreenHeight int
}
