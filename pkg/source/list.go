package source

import (
    "bufio"
    log "github.com/sirupsen/logrus"
    "math/rand"
    "os"
    "strings"
    vplay "vplay/pkg"
)


type list struct {
    errors chan error
    videos chan vplay.Video
}

func NewList() *list {
    return &list{errors: make(chan error), videos: make(chan vplay.Video)}
}

type listVideo struct {
    path string
}

func (v *listVideo) GetVideo() vplay.Video {
    return v
}

func (v *listVideo) GetURI() string {
    return v.path
}

func (v *listVideo) Name() string {
    return v.path
}

func (v *listVideo) Groups() []string {
    return []string{}
}

func (v *listVideo) GetUUID() string {
    return v.path
}

func (l *list) Videos() chan vplay.Video {
    return l.videos
}

func (l *list) Scan(path string) {
    f, err := os.Open(path)
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()
    var paths []string
    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        paths = append(paths, strings.Trim(scanner.Text(), " "))
    }
    err = scanner.Err()
    if err != nil {
        log.Fatal(err)
    }
    rand.Shuffle(len(paths), func(i, j int) { paths[i], paths[j] = paths[j], paths[i] })
    go func() {
        for _, path := range paths {
            log.Info(path)
            l.videos <- &listVideo{path}
        }
    }()

}

func (l *list) Errors() chan error {
    return l.errors
}