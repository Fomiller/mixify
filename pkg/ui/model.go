package ui

type Model struct {
	choices   []string         // items on the to-do list
	cursor    int              // which to-do list item our cursor is pointing at
	selected  map[int]string   // which to-do items are selected
	playlists map[int]playlist // which to-do items are selected
	status    int
	err       error
	state     string
}

func NewModel() Model {
	var m Model
	for _, v := range Playlist.list {
		m.choices = append(m.choices, v.name)
	}

	m.selected = make(map[int]string)
	return m
}
