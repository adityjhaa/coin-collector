package client

import (
	"log"

	"coin-collector/common"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	Net    *Network
	Interp *Interpolator

	PlayerID common.PlayerID

	RenderState WorldSnapshot

	localMask uint8
}

func NewGame() *Game {
	return &Game{
		Net:    NewNetwork(),
		Interp: NewInterpolator(),
	}
}

func (g *Game) Init() error {
	if err := g.Net.Connect(); err != nil {
		return err
	}
	g.PlayerID = g.Net.PlayerID
	log.Println("Client Init: player id", g.PlayerID)
	return nil
}

func (g *Game) Update() error {
	// read input
	var mask uint8 = 0
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		mask |= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		mask |= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		mask |= 4
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		mask |= 8
	}
	g.localMask = mask
	g.Net.SendInput(mask)

	select {
	case pkt := <-g.Net.StateChan:
		snap := ParseWorldState(pkt)
		g.Interp.AddSnapshot(snap)
	default:
	}

	g.RenderState = g.Interp.GetRenderState()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	DrawWorld(screen, g.RenderState, g.PlayerID, g.localMask)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 800, 600
}

func (g *Game) Close() {
	g.Net.Close()
}
