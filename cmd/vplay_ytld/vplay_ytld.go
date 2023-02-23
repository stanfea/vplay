package main

// download video folder instead and then use vplay_folder
//P='https://www.pornhub.org/pornstar/goldie-baby'
//
//yt-dlp --get-id $P | xargs -I {} -P 4 yt-dlp --downloader aria2c -N 10 --retry-sleep 1 -i --fixup never --force-ipv4 --no-check-certificate \
//'https://www.pornhub.org/view_video.php?viewkey={}'
//


//OLD CODE
//
//import (
//	"fmt"
//	log "github.com/sirupsen/logrus"
//	"os"
//	vplay "vplay/pkg"
//	"vplay/pkg/player"
//	"vplay/pkg/source"
//)
//
//const (
//	rows = 2
//	cols = 2
//	width = (2048)*2
//	height = (1289-100) *2
//)
//
//func main() {
//	//https://www.pornhub.com/channels/misarichannel`
//
//	logFile := "/Users/stan/vplay.txt"
//	f, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
//	if err != nil {
//		fmt.Println("Failed to create logfile" + logFile)
//		log.Error(err)
//		os.Exit(1 )
//	}
//	defer f.Close()
//	// Output to stdout instead of the default stderr
//	log.SetOutput(f)
//
//	// Only log the debug severity or above
//	log.SetLevel(log.DebugLevel)
//
//
//	log.SetReportCaller(true)
//
//	src := source.NewYtdl()
//	src.Scan(os.Args[1])
//
//	errors := src.Errors()
//	go func() {
//		for err := range errors {
//			if err != nil {
//				log.Error(err)
//			}
//		}
//	}()
//
//	videos := src.Videos()
//	finished := make(chan vplay.Video, 100)
//
//	grid := player.NewGrid(rows, cols, width, height)
//	grid.PlayAll(videos, errors, finished)
//	grid.Close()
//
//
//
//
//}