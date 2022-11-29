# Mixify

Mixify is a tui application that helps you combine, create, edit, and delete spotify playlists



## Install
### Requirements
go version >= 1.18

```
go install github.com/fomiller/mixify
```

Register your application with spotify and export your client id and secret
```
SPOTIFY_ID=<id>
SPOTIFY_SECRET=<secret>
```

## PKGs Used
[bubbletea](https://github.com/charmbracelet/bubbletea)
[spotify]("github.com/zmb3/spotify/v2")

## TODO
create menu model that lists all other models, when the model is selected the main model changes to display the selected models view

[x] make the main model struct have fields that contain a model instead of a map of models
[] create spotify pkg
[] implement spotify api package into models/methods/cmds
[] store api token in .config/mixify
[] store refresh token
[] dynamic state in nested playlist models
[x] change playlist nested models to be fancy tea list models
[] create styles for bubble tea components. spofify colors
[x] build out readme
[] change menu items to be names of actions "create", "edit", "queue", "combine", "update", "preview" "select"
[x] fix vertial scrolling with left and right controls
[x] fix height and width of screen
[] make height and width of components dynamic
