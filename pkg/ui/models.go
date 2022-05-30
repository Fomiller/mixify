package ui

type Model struct {
	viewList  []string //list of previous views
	view      string
	choices   []listItem       // items on the to-do list
	cursor    int              // which to-do list item our cursor is pointing at
	selected  map[int]playlist // which to-do items are selected
	playlists map[int]playlist // which to-do items are selected
	status    int
	err       error
	state     string
}

func NewModel() Model {
	var m Model
	for _, v := range Playlist.list {
		item := listItem{
			selected: false,
			detail:   v,
		}
		m.choices = append(m.choices, item)
	}

	m.selected = make(map[int]playlist)
	return m
}
