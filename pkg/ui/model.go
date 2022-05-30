package ui

type Model struct {
	currentView string
	choices     []playlist       // items on the to-do list
	cursor      int              // which to-do list item our cursor is pointing at
	selected    map[int]playlist // which to-do items are selected
	playlists   map[int]playlist // which to-do items are selected
	status      int
	err         error
	state       string
}

func NewModel() Model {
	var m Model
	for _, v := range Playlist.list {
		m.choices = append(m.choices, v)
	}

	m.selected = make(map[int]playlist)
	return m
}
