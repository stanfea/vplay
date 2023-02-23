package source

import (
   "bufio"
   "bytes"
   "encoding/binary"
   "encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
   "io"
   "os"
   "unsafe"
    vplay "vplay/pkg"
)


var (
   // nativeEndian used to detect native byte order
   nativeEndian binary.ByteOrder
   // bufferSize used to set size of IO buffer - adjust to accommodate message payloads
   bufferSize = 8192
)

// IncomingMessage represents a message sent to the native host.
type IncomingMessage struct {
	Query string `json:"query"`
}

// OutgoingMessage respresents a response to an incoming message query.
type OutgoingMessage struct {
	Query    string `json:"query"`
	Response string `json:"response"`
}

// determines native byte order.
func Init() {
	// determine native byte order so that we can read message size correctly
	var one int16 = 1
	b := (*byte)(unsafe.Pointer(&one))
	if *b == 0 {
		nativeEndian = binary.BigEndian
	} else {
		nativeEndian = binary.LittleEndian
	}
}

type ChromeVideo struct {
    Url string
}

func (c ChromeVideo) GetVideo() vplay.Video {
    return c
}

func (c ChromeVideo) GetURI() string {
    return c.Url
}

func (c ChromeVideo) Name() string {
   return c.Url
}

func (c ChromeVideo) Groups() []string {
    return []string{c.Url}
}

func (c ChromeVideo) GetUUID() string {
    return c.Url
}



type chrom struct {
    videos chan vplay.Video
    errors chan error
    close chan struct{}
}


func (c *chrom) Close() {
	c.errors <- errors.New("closed browser")
}

func (c *chrom) Scan(uri string) {
    log.Error("test")
    go func() {
        Init()
        v := bufio.NewReader(os.Stdin)
        // adjust buffer size to accommodate your json payload size limits; default is 4096
        s := bufio.NewReaderSize(v, bufferSize)
        log.Infof("IO buffer reader created with buffer size of %v.", s.Size())

        lengthBytes := make([]byte, 4)
        lengthNum := int(0)

        history := make(map[string]bool)
        // we're going to indefinitely read the first 4 bytes in buffer, which gives us the message length.
        // if stdIn is closed we'll exit the loop and shut down host
        for b, err := s.Read(lengthBytes); b > 0 && err == nil; b, err = s.Read(lengthBytes) {
            // convert message length bytes to integer value
            lengthNum = readMessageLength(lengthBytes)
            log.Infof("Message size in bytes: %v", lengthNum)

            // If message length exceeds size of buffer, the message will be truncated.
            // This will likely cause an error when we attempt to unmarshal message to JSON.
            if lengthNum > bufferSize {
                log.Infof("Message size of %d exceeds buffer size of %d. Message will be truncated and is unlikely to unmarshal to JSON.", lengthNum, bufferSize)
            }

            // read the content of the message from buffer
            content := make([]byte, lengthNum)
            _, err := s.Read(content)
            if err != nil && err != io.EOF {
                log.Fatal(err)
            }

            // message has been read, now parse and process
            url := parseMessage(content)
			if !history[url] {
				history[url] = true
				c.videos <- &ChromeVideo{url}
			}

        }

        log.Info("Stdin closed.")
		c.Close()
    }()
}

func (c *chrom) Errors() chan error {
    return c.errors
}

func (c *chrom) Videos() chan vplay.Video {
    return c.videos
}

// read Creates a new buffered I/O reader and reads messages from Stdin.
func NewChrome() *chrom {
    c := &chrom{
        videos: make(chan vplay.Video, 100),
        errors: make(chan error, 100),
    }
    return c

}

// readMessageLength reads and returns the message length value in native byte order.
func readMessageLength(msg []byte) int {
	var length uint32
	buf := bytes.NewBuffer(msg)
	err := binary.Read(buf, nativeEndian, &length)
	if err != nil {
		log.Errorf("Unable to read bytes representing message length: %v", err)
	}
	return int(length)
}

// parseMessage parses incoming message
func parseMessage(msg []byte) string {
	iMsg := decodeMessage(msg)
	log.Infof("Message received: %s", msg)
	return iMsg.Query
}

// decodeMessage unmarshals incoming json request and returns query value.
func decodeMessage(msg []byte) IncomingMessage {
	var iMsg IncomingMessage
	err := json.Unmarshal(msg, &iMsg)
	if err != nil {
		log.Errorf("Unable to unmarshal json to struct: %v", err)
	}
	return iMsg
}


// send sends an OutgoingMessage to os.Stdout.
func send(msg OutgoingMessage) {
	byteMsg := dataToBytes(msg)
	writeMessageLength(byteMsg)

	var msgBuf bytes.Buffer
	_, err := msgBuf.Write(byteMsg)
	if err != nil {
		log.Errorf("Unable to write message length to message buffer: %v", err)
	}

	_, err = msgBuf.WriteTo(os.Stdout)
	if err != nil {
        log.Errorf("Unable to write message buffer to Stdout: %v", err)
	}
}

// dataToBytes marshals OutgoingMessage struct to slice of bytes
func dataToBytes(msg OutgoingMessage) []byte {
	byteMsg, err := json.Marshal(msg)
	if err != nil {
        log.Errorf("Unable to marshal OutgoingMessage struct to slice of bytes: %v", err)
	}
	return byteMsg
}

// writeMessageLength determines length of message and writes it to os.Stdout.
func writeMessageLength(msg []byte) {
	err := binary.Write(os.Stdout, nativeEndian, uint32(len(msg)))
	if err != nil {
        log.Errorf("Unable to write message length to Stdout: %v", err)
	}
}
