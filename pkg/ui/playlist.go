package ui

type listItem struct {
	selected bool
}

type playlist struct {
	listItem
	name        string
	description string
	tracks      []string
}

type playlists struct {
	list []playlist
}

// func (i playlist) Title() string       { return i.name }
// func (i playlist) Description() string { return i.description }
// func (i playlist) FilterValue() string { return i.name }

var Playlist = playlists{
	list: []playlist{
		{
			name:        "playlist_01",
			description: "raggae music",
			tracks:      []string{"raggae 1", "raggae 2", "raggae 3", "raggae 4", "raggae 5", "raggae 6", "raggae 7", "raggae 8", "raggae 9", "raggae 10"},
		},
		{
			name:        "playlist_02",
			description: "chill music",
			tracks:      []string{"chill 1", "chill 2", "chill 3", "chill 4", "chill 5", "chill 6", "chill 7", "chill 8", "chill 9", "chill 10"},
		},
		{
			name:        "playlist_03",
			description: "rap music",
			tracks:      []string{"rap 1", "rap 2", "rap 3", "rap 4", "rap 5", "rap 6", "rap 7", "rap 8", "rap 9", "rap 10"},
		},
		{
			name:        "playlist_04",
			description: "EDM music",
			tracks:      []string{"EDM 1", "EDM 2", "EDM 3", "EDM 4", "EDM 5", "EDM 6", "EDM 7", "EDM 8", "EDM 9", "EDM 10"},
		},
		{
			name:        "playlist_05",
			description: "classical music",
			tracks:      []string{"classical 1", "classical 2", "classical 3", "classical 4", "classical 5", "classical 6", "classical 7", "classical 8", "classical 9", "classical 10"},
		},
	},
}
