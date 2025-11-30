package client

import (
	"coin-collector/common"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	Net    *Network
	Interp *Interpolator

	LocalInputMask uint8

	// rendered world snapshot
	RenderState WorldSnapshot
	PlayerID    common.PlayerID
}

func NewGame() *Game {
	net := NewNetwork()
	return &Game{
		Net:    net,
		Interp: NewInterpolator(),
	}
}

func (g *Game) Init() error {
	err := g.Net.Connect()
	if err != nil {
		return err
	}

	g.PlayerID = g.Net.PlayerID
	return nil
}

func (g *Game) Update() error {
	// --- 1. Read keyboard → build input bitmask ---
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

	g.LocalInputMask = mask
	g.Net.SendInput(mask)

	// --- 2. Process world state packets ---
	select {
	case pkt := <-g.Net.StateChan:
		snap := ParseWorldState(pkt)
		g.Interp.AddSnapshot(snap)
	default:
	}

	// --- 3. Get interpolated render state ---
	if len(g.Interp.Snapshots) > 0 {
		g.RenderState = g.Interp.GetRenderState()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	DrawWorld(screen, g.RenderState, g.PlayerID)
}

func (g *Game) Layout(outW, outH int) (int, int) {
	return 800, 600
}
