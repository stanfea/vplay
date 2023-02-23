package vplay

type Service interface {
    AddSource(source Source)
}

type Video interface {
    GetVideo() Video
    GetURI() string
    Name() string
    Groups() []string
    GetUUID() string
}

type Source interface {
    Scan(uri string)
    Errors() chan error
    Videos() chan Video
}

type Player interface {
    Play(videos chan Video, finished chan Video, errors chan error)
    //Pause()
    //Seek(time.Duration)
}

type PlayerWithGeometry interface {
    Player
    Geometry(x,y, width int) error
}

type PlayList interface {
    GetVideo() Video
    AddSource(source Source)
}

type PlayItem interface {
    GetVideo() Video
}

type service struct {
    playItem PlayItem
    player   Player
}

func New(playItem PlayItem, player Player) *service {
    return &service{playItem: playItem, player: player}

}