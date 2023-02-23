package controller
//
//import (
//    "math/rand"
//    "os"
//    "path/filepath"
//    "sync"
//    "time"
//    "vplay/pkg/logging"
//    "vplay/pkg/player"
//)
//
//type mySelector struct {
//    items             []*item
//    parentFolder      map[string]int
//    grandParentFolder map[string]int
//    playing           map[string]bool
//    mutex             *sync.Mutex
//    root string
//    removeCount int
//    players []player.Player
//    disabledPlayers map[int]bool
//    activePlayer int
//}
//
//type item struct {
//    path string
//    parentFolder string
//    grandParentFolder string
//    playCount int
//}
//
//func process(root string, video string) *item {
//    i := &item{}
//    if parent := filepath.Dir(video); parent != "." && parent != root {
//        i.parentFolder = parent
//        if grandparent := filepath.Dir(parent); grandparent != "." && grandparent != root {
//            i.grandParentFolder = grandparent
//        }
//    }
//    i.path = video
//    return i
//}
//
//func New(root string, videos []string, players []player.Player) *mySelector {
//    items := make([]*item, len(videos))
//    for i, video := range videos {
//        items[i] = process(root, video)
//    }
//    return &mySelector{
//        items:             items,
//        parentFolder:      make(map[string]int),
//        grandParentFolder: make(map[string]int),
//        playing:           make(map[string]bool),
//        mutex: &sync.Mutex{},
//        root: root,
//        players: players,
//        activePlayer: -1,
//    }
//}
//
//func (s *mySelector) isEmpty() bool {
//    return len(s.items) == 0
//}
//
//func (s *mySelector) rotateSound() error {
//    s.mutex.Lock()
//    defer s.mutex.Unlock()
//    if len(s.players) == 0 {
//        return nil
//    }
//    if len(s.players) == 1 {
//        if s.players[0].GetMuted() {
//            err := s.players[0].ToggleMute()
//            if err != nil {
//                return err
//            }
//        }
//        return nil
//    }
//
//    var n int
//    for true {
//        if n = rand.Intn(len(s.players)); n != s.activePlayer {
//            break
//        }
//    }
//    if s.activePlayer != -1 {
//        err := s.players[s.activePlayer].ToggleMute()
//        if err != nil {
//            return nil
//        }
//    }
//    s.activePlayer = n
//    err := s.players[n].ToggleMute()
//    if err != nil {
//        return nil
//    }
//    return nil
//}
//
//func (s *mySelector) updateLoop(done chan bool) {
//    ticker := time.NewTicker(5 * time.Second)
//    go func() {
//        for {
//            select {
//            case <-done:
//                return
//            case _ = <-ticker.C:
//                err := s.rotateSound()
//                if err != nil {
//                    logging.Error.Print(err)
//                    os.Exit(1)
//                }
//                if len(s.players) <= 1 {
//                    return
//                }
//            }
//        }
//
//    }()
//}
//
//func (s *mySelector) Loop(maxViews int) {
//    done := make(chan bool)
//    //s.rotateSoundLoop(done)
//    wg := &sync.WaitGroup{}
//    for i, p := range s.players {
//        if p == nil {
//            continue
//        }
//        wg.Add(1)
//        p := p
//        i := i
//        go func(wg *sync.WaitGroup, index int, p player.Player) {
//            defer wg.Done()
//            for true {
//                s.mutex.Lock()
//                if s.isEmpty() {
//                    break
//                }
//                video := s.nextVideo(p.GetLastPlayed(), maxViews)
//                s.mutex.Unlock()
//                err := p.Play(video)
//                if err != nil {
//                    logging.Error.Print(err)
//                }
//            }
//            wg.Done()
//            s.mutex.Lock()
//            s.players[i] = s.players[len(s.players)-1]
//            s.players[len(s.players)-1] = nil
//            s.players = s.players[:len(s.players)-1]
//            s.mutex.Unlock()
//        }(wg, i, p)
//    }
//    wg.Wait()
//    println(1)
//    done <-true
//}
//
//func (s *mySelector) nextVideo(finished string, maxViews int) string {
//    for {
//        index := rand.Intn(len(s.items)-1)
//        item := s.items[index]
//        if _, ok := s.playing[item.path]; ok {
//            logging.Trace.Printf("Skipped: %s", item.path)
//            continue
//        }
//        if item.parentFolder != ""  {
//            if c := s.parentFolder[item.parentFolder]; c > 0 {
//                continue
//            }
//        }
//        if item.grandParentFolder != "" {
//            if  c := s.grandParentFolder[item.grandParentFolder]; c > 0 {
//                continue
//            }
//        }
//        if item.parentFolder != "" {
//            s.parentFolder[item.parentFolder]++
//
//        }
//        if item.grandParentFolder != "" {
//            s.grandParentFolder[item.grandParentFolder]++
//        }
//        s.playing[item.path] = true
//        if finished != "" {
//            finishedItem := process(s.root, finished)
//            delete(s.playing, finishedItem.path)
//            if finishedItem.parentFolder != "" {
//                s.parentFolder[finishedItem.parentFolder]--
//            }
//            if finishedItem.grandParentFolder != "" {
//                s.grandParentFolder[finishedItem.grandParentFolder]--
//            }
//        }
//        path := item.path
//        item.playCount++
//        if item.playCount >= maxViews {
//            s.removeCount++
//            // remove item
//            s.items[index] = s.items[len(s.items)-1]
//            s.items[len(s.items)-1] = nil
//            s.items = s.items[:len(s.items)-1]
//        }
//        return path
//    }
//}