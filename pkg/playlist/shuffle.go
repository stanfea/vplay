package playlist

import _ "vplay/pkg"

var (
	SEEN_VIDEO = byte(0b00000100)
    SEEN_FOLDER = byte(0b00000010)
    SEEN_PARENT_FOLDER = byte(0b00000010)
)

type playList struct {
    sources []Source
}


func New(sources []Source) *playList {
    return &playlist{sources: sources}
}