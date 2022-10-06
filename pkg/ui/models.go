package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type view string

const (
	MENU     view = "menu"
	PLAYLIST view = "playlist"
	TRACK    view = "track"
	TEST     view = "test"
)

// MAIN MODEL
type Model struct {
	view     view
	views    map[view]tea.Model
	viewList []view //list of previous views, this could be a linked list ?

	choices   []listItem       // items on the to-do list
	selected  map[int]playlist // which to-do items are selected
	playlists map[int]playlist // which to-do items are selected

	cursor int // which to-do list item our cursor is pointing at, This could be pulled into a nested model?
	status int
	err    error
	state  string
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

	m.view = MENU
	m.views = make(map[view]tea.Model)

	m.views[MENU] = newMenuModel()
	m.views[TEST] = newTestModel()
	m.views[PLAYLIST] = newTestModel()
	m.views[TRACK] = newTestModel()

	m.selected = make(map[int]playlist)
	return m
}

type testModel struct {
	choices []listItem
	cursor  int
	status  int
	err     error
	state   string
}

func newTestModel() tea.Model {
	m := testModel{}
	return m
}

func (m testModel) Init() tea.Cmd {
	return nil
}

func (m testModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	return m, cmd
}

func (m testModel) View() string {
	return ""
}

// PLAYLIST
type playlistModel struct {
	choices []listItem
	cursor  int
	status  int
	err     error
	state   string
}

func newPlaylistModel() tea.Model {
	m := playlistModel{}
	return m
}

func (m playlistModel) Init() tea.Cmd {
	return nil
}

func (m playlistModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	return m, cmd
}

func (m playlistModel) View() string {
	return ""
}

// TRACK
type trackModel struct {
	choices []listItem
	cursor  int
	status  int
	err     error
	state   string
}

func newTrackModel() tea.Model {
	m := trackModel{}
	return m
}

func (m trackModel) Init() tea.Cmd {
	return nil
}

func (m trackModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	return m, cmd
}

func (m trackModel) View() string {
	return ""
}

// MENU
type menuModel struct {
	choices []listItem
	tag     string
	cursor  int
	status  int
	err     error
	state   string
}

func newMenuModel() menuModel {
	m := menuModel{}
	m.tag = "my cool tag"
	return m
}

func (m menuModel) getTag() string {
	return m.tag
}

func (m menuModel) Init() tea.Cmd {
	return nil
}

func (m menuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	return m, cmd
}

func (m menuModel) View() string {
	return ""
}
