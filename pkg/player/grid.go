package player

import (
    "sync"
    vplay "vplay/pkg"
)

type grid struct {
    rows, cols, width, height int
    players map[int]*mpv
}

func NewGrid(rows, cols, width, height int) *grid {
    return &grid{rows: rows, cols: cols, width: width, height: height}
}

func (g *grid) Close() {
    for i := 0; i < len(g.players); i++ {
        g.players[i].Close()
    }
}

func (g *grid) PlayAll(videos chan vplay.Video, errors chan error, finished chan vplay.Video) {
    g.players = make(map[int]*mpv, g.rows*g.cols)
    i := 0

    playerWidth := g.width / g.cols
    playerHeight := g.height / g.rows
    for y := 0; y < g.rows; y++ {
        for x := 0; x < g.cols; x++ {
            py := playerHeight * y
            px := playerWidth * x
            player := NewMpv(i)
            if err := player.Geometry(px, py, g.rows, g.cols, g.width/g.cols, g.height/g.rows); err != nil {
                errors <- err
                return
            }
            g.players[i] = player
            i++
        }
    }

    wg := &sync.WaitGroup{}
    for _, p := range g.players {
        wg.Add(1)
        go func(p vplay.Player) {
            p.Play(videos, finished, errors)
            wg.Done()
        }(p)
    }
    wg.Wait()
    println("all done")
}


