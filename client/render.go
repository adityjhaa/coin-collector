package client

import (
	"fmt"
	"image/color"

	"coin-collector/common"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	localColor  = color.RGBA{0, 200, 0, 255}
	remoteColor = color.RGBA{200, 0, 0, 255}
	coinColor   = color.RGBA{255, 215, 0, 255}
	bgColor     = color.RGBA{30, 30, 30, 255}
)

const playerSize float32 = 26
const coinRadius float32 = 10

func DrawWorld(screen *ebiten.Image, state WorldSnapshot, localID common.PlayerID, localMask uint8) {
	screen.Fill(bgColor)

	for _, c := range state.Coins {
		vector.FillCircle(screen, c.X, c.Y, coinRadius, coinColor, true)
	}

	localScore := uint16(0)
	otherScore := uint16(0)
	for _, p := range state.Players {
		if p.ID == localID {
			drawSquare(screen, p.X, p.Y, localColor)
			localScore = p.Score
		} else {
			drawSquare(screen, p.X, p.Y, remoteColor)
			otherScore = p.Score
		}
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf(
		"You (Green): %d\nOpponent (Red): %d\nPlayers: %d   Coins: %d",
		localScore, otherScore, len(state.Players), len(state.Coins),
	))
}

func drawSquare(screen *ebiten.Image, x, y float32, col color.Color) {
	half := playerSize / 2
	vector.FillRect(screen, x-half, y-half, playerSize, playerSize, col, true)
}
