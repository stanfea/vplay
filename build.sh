#!/bin/sh

go build cmd/vplay_folder/vplay_folder.go
go build cmd/vplay_browser/vplay_browser.go
#go build cmd/vplay_ytld/vplay_ytld.go
chmod +x vplay_browser vplay_folder
