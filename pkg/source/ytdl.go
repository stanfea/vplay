package source
//
//import (
//    "bufio"
//    log "github.com/sirupsen/logrus"
//    "os/exec"
//    vplay "vplay/pkg"
//)
//
////youtube-dl  'https://www.pornhub.com/video?c=35&page=10'
//
//
//type ytdl struct {
//    videos chan vplay.Video
//    errors chan error
//}
//
//func NewYtdl() *ytdl {
//    return &ytdl{errors: make(chan error), videos: make(chan vplay.Video, 10000)}
//}
//
//type videoYtld struct {
//    url string
//}
//
//func (v *videoYtld) GetVideo() vplay.Video {
//    return v
//}
//
//func (v *videoYtld) GetURI() string {
//    return v.url
//}
//
//func (v *videoYtld) Name() string {
//    return v.url
//}
//
//func (v *videoYtld) Groups() []string {
//    return []string{}
//}
//
//func (v *videoYtld) GetUUID() string {
//    return v.url
//}
//
//func (y *ytdl) Videos() chan vplay.Video {
//    return y.videos
//}
//
//func (y *ytdl) Scan(url string) {
//    args := []string{
//        "--verbose",
//        //"--get-url",
//       //"--playlist-random",
//       "--downloader", "aria2c",
//       //"--cookies", "~/cookies.txt",
//       "-i",
//       "--user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36",
//       //"--referer", url,
//       url,
//    }
//    cmd := exec.Command("yt-dlp", args...)
//    stdout, err := cmd.StdoutPipe()
//    if err != nil {
//        log.Fatal(err)
//    }
//    stderr, err := cmd.StderrPipe()
//    if err != nil {
//        log.Fatal(err)
//    }
//    err = cmd.Start()
//    if err != nil {
//        log.Fatal(err)
//    }
//    go func() {
//        reader := bufio.NewReader(stderr)
//        _, err := reader.ReadString('\n')
//        for err != nil {
//            y.errors <- err
//            _, err = reader.ReadString('\n')
//        }
//    }()
//
//    /// TODO yt-dlp downloads the videos now we have to stream them
//    //go func() {
//    //    reader := bufio.NewReader(stdout)
//    //    line, err := reader.ReadString('\n')
//    //    log.Info(line)
//    //    for err == nil {
//    //        y.videos <- &videoYtld{line}
//    //        line, err = reader.ReadString('\n')
//    //    }
//    //}()
//
//}
//
//func (y *ytdl) Errors() chan error {
//    return y.errors
//}