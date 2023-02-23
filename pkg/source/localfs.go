package source

import (
	"github.com/djherbis/atime"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
	"time"
	vplay "vplay/pkg"
)

type localFs struct {
	videos chan vplay.Video
	errors chan error
	db     map[string]*videoFile
	paths  map[string]struct{}
	close  chan bool
}

type videoFile struct {
	name         string
	path         string
	folder       string
	parentFolder string
	sourcePath   string
}

func (v *videoFile) GetURI() string {
	return v.path
}

func (v *videoFile) Name() string {
	return v.name
}

func (v *videoFile) Groups() []string {
	return []string{v.folder, v.parentFolder}
}

func (v *videoFile) GetUUID() string {
	return v.path
}

func (v *videoFile) GetVideo() vplay.Video {
	return v
}

func NewLocalFs() vplay.Source {
	return &localFs{
		db:     make(map[string]*videoFile),
		errors: make(chan error),
		videos: make(chan vplay.Video),
	}
}

func (fs *localFs) Scan(path string) {
	err := fs.walk(path)
	if err != nil {
		fs.errors <- err
		return
	}
	//for {
	//	log.Warn("rescanning")
	//	err := fs.walk(path)
	//	if err != nil {
	//		fs.errors <- err
	//		return
	//	}
	//	time.Sleep(60 * time.Second)
	//}

	//
	//go func() {
	//
	//
	//}()
	//time.Sleep(3000 * time.Millisecond)
}


func (fs *localFs) Errors() chan error {
	return fs.errors
}

func (fs *localFs) Videos() chan vplay.Video {
    go func() {
        inputs := make([]vplay.Video, len(fs.db))
        i := 0
        for _, video := range fs.db {
            inputs[i] = video
            i++
        }
        Shuffle(inputs)
		Shuffle(inputs)
		Shuffle(inputs)
		Shuffle(inputs)
        for _, video := range inputs {
            fs.videos <- video
        }

    }()

	return fs.videos
}

func (fs *localFs) getVideo(path string) *videoFile {
	path = fs.sanitize(path)
	name := filepath.Base(path)
	folder := filepath.Dir(path)
	parentFolder := filepath.Dir(folder)
	return &videoFile{
		path:         path,
		name:         name,
		folder:       folder,
		parentFolder: parentFolder,
	}
}

func (fs *localFs) walk(path string) error {
	err := filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			if _, ok := IsVideoFile[filepath.Ext(path)]; ok {
				video := fs.getVideo(path)
				accessTime, err := atime.Stat(path)
				if err != nil {
					return err
				}
				tenMinutesAgo := time.Now().Add(-10 * time.Minute)

				// accessed less than 10 minutes ago
				if accessTime.After(tenMinutesAgo) {
					log.Warn("skip", path)
					//return nil
				}
				if seenVideo, ok := fs.db[video.name]; ok {
					log.Warn("skipping duplicate file", "seenName", seenVideo.name, "newName", video.name)
					return nil
				}
				fs.db[video.name] = video
			}
			return nil
		})
	return err
}

func (fs *localFs) sanitize(input string) string {
	return strings.Replace(input, "/._", "/", -1) // some mac bug???
}
