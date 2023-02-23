package player

import (
    "bufio"
    "fmt"
    log "github.com/sirupsen/logrus"
    "io"
    "os"
    "os/exec"
    "path"
    "strings"
    "sync"
    "time"
    vplay "vplay/pkg"
    mpvClient "vplay/pkg/mpv"
)


type mpv struct {
    id, x, y, rows, cols, width, height int
    //id int
    video vplay.Video
    //muted bool
    //playing bool
    socket string
    cmd *exec.Cmd
    conn *mpvClient.Client
    stdout io.ReadCloser
    errors chan error
    //atHalf bool
    //seek bool
    //position float64
    mutex *sync.Mutex
    wait bool
    done chan string
    out chan string
}


func NewMpv(id int) *mpv {
    return &mpv{id: id, mutex: &sync.Mutex{}, errors: make(chan error)}
}

func (p *mpv) isPartOfGrid() bool {
    return p.rows > 0
}

func (m *mpv) open() error  {
    m.mutex.Lock()
    defer m.mutex.Unlock()
    args := []string{
        "--idle",
        "--input-ipc-server=" + m.getSocket(),
        "--hwdec=auto",
        "--hwdec-codecs=all",
        "--no-border",
        "--macos-force-dedicated-gpu",
        //"--mute",
        "--volume=50",
        "--start=100",
        "--keep-open=always",
        "--force-window",
        "--no-focus-on-open",
    }
    if m.isPartOfGrid() {
        args = append(args,
            fmt.Sprintf("--geometry=%dx%d+%d+%d", m.width, m.height, m.x, m.y),
        )
    }
    cmd := exec.Command("mpv", args...)
    m.cmd = cmd
    var err error
    cmd.Stderr = cmd.Stdout
    m.stdout, err = m.cmd.StdoutPipe()
    if err != nil {
        log.Fatal(err)
    }
    err = m.cmd.Start()

    if err != nil {
        log.Fatal(err)
    }
    time.Sleep(1000 * time.Millisecond)
    err = m.connect()
    if err != nil {
        log.Fatal(err)
    }
    go func() {
        reader := bufio.NewReader(m.stdout)
        line, err := reader.ReadString('\n')
        for err == nil {
            log.Info(line)
            if strings.HasPrefix(line, "Failed") {//|| strings.HasPrefix(line, "[ffmpeg] https: HTTP error") {
                log.Errorf("%d stdout: %s", m.id, line)
                m.done <- "failed"
            }
            line, err = reader.ReadString('\n')
        }
    }()
    return nil
}

func (p *mpv) Geometry(x, y, rows, cols, width, height int) error {
    p.x = x
    p.y = y
    p.rows = rows
    p.cols = cols
    p.width = width
    p.height = height
    return nil
}

func (m *mpv) connect() error {
    ipc, err := mpvClient.NewIPCClient(m.socket, m.done)
    if err != nil {
        return err
    }
    m.conn = mpvClient.NewClient(ipc)
    return nil
}

func (m *mpv) loadVideo(video vplay.Video, finished chan vplay.Video) error {
    var err error
    log.Info(video.GetURI())
    err = m.conn.Loadfile(video.GetURI(), mpvClient.LoadFileModeReplace)
    if err == nil {
        m.video = video
    }
    if err != nil {
        log.Error(err)
        return err
    }

    //var d float64
    //for i := 0 ; i < 3 ; i++ {
    //  d, err = m.conn.Duration()
    //  if err == nil {
    //      break
    //  }
    //  time.Sleep(3*time.Second)
    //}
    //if d == 0 {
    //  return errors.New("no duration")
    //} else {
    //  finished <- video
    //}
    //ticker := time.NewTicker(time.Second*3)
    // step := d / 10
    // go func() {
    //    for {
    //        select {
    //        case <-ticker.C:
    //            pos, err := m.conn.Position()
    //            if err != nil {
    //                continue
    //            }
    //            newPos := pos + step
    //            if newPos > d-step {
    //                ticker.Stop()
    //                break
    //            }
    //            err = m.conn.Seek(int(newPos), "absolute")
    //            if err != nil {
    //                log.Errorf("%d: seek error %s", m.id, err.Error())
    //            } else {
    //                pos = newPos
    //            }
    //        }
    //    }
    // }()
    //            pos, err := m.conn.Position()
    //            if err != nil {
    //                continue
    //            }
    //            newPos := pos + step
    //            if newPos > d-step {
    //                ticker.Stop()
    //                break
    //            }
    //            err = m.conn.Seek(int(newPos), "absolute")
    //            if err != nil {
    //                log.Errorf("%d: seek error %s", m.id, err.Error())
    //            } else {
    //                pos = newPos
    //            }
    //        }
    //    }
    //}()

    // var d float64
    // for i := 0; i <= 3; i++ {
    //     d, err = m.conn.Duration()
    //     if err != nil {
    //         log.Errorf("%d: %s", m.id, err.Error())
    //         continue
    //     }
    //     break
    // }
    // start := int(d/4*3)
    // for {
    //     if err := m.conn.Seek(start, "absolute"); err != nil {
    //         log.Errorf("%d: %s", m.id, err.Error())
    //         continue
    //     }
    //     break
    // }
    //t := time.NewTimer(time.Second*5)
    //<- t.C
    return nil
}

func (m *mpv) Play(videos chan vplay.Video, finished chan vplay.Video, errors chan error) {
    m.done = make(chan string, 1000)
    err := m.open()
    if err != nil {
        errors <- err
        return
    }
    for video := range videos {
        if err := m.loadVideo(video, finished); err != nil {
            errors <- err
        }
        <- m.done
    }
    log.Infof("player %d finished", m.id)
}

func (p *mpv) getSocket() string {
    if p.socket == "" {
        tmp := os.TempDir()
        if p.rows == 0 || p.cols == 0 {
            p.socket = path.Join(tmp, "vplay-socket")
        } else {
            p.socket = path.Join(tmp, fmt.Sprintf("vplay-%d-socket", p.id))
        }
    }
    return p.socket
}


func (m *mpv) Close() {
    err := m.cmd.Process.Kill()
    if err != nil {
        log.Error(err)
    }
}

func (m *mpv) Pause() {
    panic("implement me")
}

func (m *mpv) Seek(duration time.Duration) {
    panic("implement me")
}

func (m *mpv) ToggleMute() error {
    ////logging.Trace.Printf("toggling mute on %d,%d", m.x,m.y)
    //if !m.playing {
    //    return errors.New("not playing")
    //}
    //muted := !m.muted
    ////if err := m.client.SetMute(muted); err != nil {
    ////    return err
    ////}
    //m.muted = muted
    ////if !m.muted {
    //if !m.atHalf{
    //    m.atHalf = true
    //    m.client.SetProperty("percent-pos", float64(50))
    //    //m.client.Seek(10, "relative-percent")
    //} else {
    //    m.client.SetProperty(10, "relative-percent")
    //}

    return nil
}

//
//func (m *mpv) autoSeek(done chan bool) {
//   min := 5
//   max := 10
//   interval := rand.Intn(max - min) + min
//   ticker := time.NewTicker(time.Duration(interval) * time.Second)
//   duration, err := m.client.Duration()
//
//   if err != nil {
//       return
//   }
//    if duration <= 120 {
//        return
//    }
//    skip := duration / 20
//   for {
//       select {
//       case <-done:
//           return
//       case _ = <-ticker.C:
//           newPosition := m.position +skip
//           if newPosition >= duration - 100 {
//               m.logger.Info("Seek end", "filanem", m.video.Name())
//               return
//           }
//           err := m.client.Seek(int(m.position), "absolute")
//           if err != nil {
//              m.logger.Info(err)
//           } else {
//               m.position = newPosition
//           }
//       }
//   }
//}

//
//
//    pipe, err := cmd.StdoutPipe()
//    if err != nil {
//      return err
//    }
//    cmd.Stdout = cmd.Stderr
//
//    m.playing = true
//    err := cmd.Start()
//    if err != nil {
//        return err
//    }
//    printCmd(cmd)
//    done := make(chan bool)
//    m.clientLoop(done)
//    reader := bufio.NewReader(pipe)
//    line, err := reader.ReadString('\n')
//    for err == nil {
//      log.Print(line)
//      line, err = reader.ReadString('\n')
//    }
//    err = cmd.Wait()
//    if err != nil {
//        return err
//    }
//    done <- true
//    m.playing = false
//    return nil
//}

//func NewMpv(logger zap.Logger) *mpv {
//    return &mpv{
//        logger:   nil,
//    }
//
//    //m.video = video
//    //m.position = 0
//    args := []string{
//        "--idle",
//        "--input-ipc-server=" + m.socket,
//        //"--hwdec-codecs=all",
//        //"--vo=libmpv",
//        //"--start=100",
//        //"--mute=yes",
//        //"--no-border",
//        //"--geometry=25%x33%+" + fmt.Sprintf("%d+%d", m.x, m.y),
//        //video.GetURI(),
//    }
//    cmd := exec.Command("mpv", args...)
//
//
//    "mpv --
//
//    return &mpv{x: x, y: y, id: id, muted: true, socket: socket}
//}
