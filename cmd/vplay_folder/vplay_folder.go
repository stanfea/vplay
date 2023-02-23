package main

import (
   log "github.com/sirupsen/logrus"
   "os"
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

   folder := os.Args[1]
   src := source.NewLocalFs()
   src.Scan(folder)

   videos := src.Videos()
   finished := make(chan vplay.Video, 100)
   grid := player.NewGrid(rows, cols, width, height)


   errors := src.Errors()
   go func() {
      for err := range errors {
          if err != nil {
              log.Error(err)
          }
      }
   }()


   grid.PlayAll(videos, errors, finished)
   grid.Close()



}