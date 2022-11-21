package playlist

type statusMsg int

type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }
