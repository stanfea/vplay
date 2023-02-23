package main

import (
	log "github.com/sirupsen/logrus"
	vplay "vplay/pkg"
	"vplay/pkg/player"
	"vplay/pkg/source"
)

const (
	rows = 2
	cols = 2
	width = (2048)*2
	height = (1289-100) *2
)

func main() {


	//logFile := "/Users/stan/vplay.txt"
	//f, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	//if err != nil {
	//	fmt.Println("Failed to create logfile" + logFile)
	//	log.Error(err)
	//	os.Exit(1 )
	//}
	//defer f.Close()
	//// Output to stdout instead of the default stderr
	//log.SetOutput(f)
	//
	//// Only log the debug severity or above
	//log.SetLevel(log.DebugLevel)
	//
	//
	//log.SetReportCaller(true)


	src := source.NewChrome()
	src.Scan("")

	videos := src.Videos()
	finished := make(chan vplay.Video, 100)

	grid := player.NewGrid(rows, cols, width, height)


	errors := src.Errors()

	go func() {
		for err := range errors {
			if err != nil {
				if err.Error() == "closed browser" {
					grid.Close()
				}
				log.Error(err)
			}
		}
	}()

	grid.PlayAll(videos, errors, finished)



}