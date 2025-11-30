package client

import (
	"fmt"
	"image/color"

	"coin-collector/common"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Colors
var (
	localColor  = color.RGBA{0, 255, 0, 255}   // green
	remoteColor = color.RGBA{255, 0, 0, 255}   // red
	coinColor   = color.RGBA{255, 215, 0, 255} // gold
	bgColor     = color.RGBA{35, 35, 35, 255}
)

const playerSize float32 = 25

func DrawWorld(screen *ebiten.Image, state WorldSnapshot, localID common.PlayerID) {
	// fill background
	screen.Fill(bgColor)

	// --- Draw coins ---
	for _, c := range state.Coins {
		vector.FillCircle(
			screen,
			float32(c.X),
			float32(c.Y),
			10,
			coinColor,
			true,
		)
	}

	// --- Draw players ---
	for _, p := range state.Players {
		if p.ID == localID {
			drawSquare(screen, p.X, p.Y, localColor)
		} else {
			drawSquare(screen, p.X, p.Y, remoteColor)
		}
	}

	// --- UI: Scores ---
	localScore := 0
	otherScore := 0

	for _, p := range state.Players {
		if p.ID == localID {
			localScore = p.Score
		} else {
			otherScore = p.Score
		}
	}

	scoreText := fmt.Sprintf(
		"You (Green): %d\nOpponent (Red): %d\nPlayers: %d  Coins: %d",
		localScore,
		otherScore,
		len(state.Players),
		len(state.Coins),
	)

	ebitenutil.DebugPrint(screen, scoreText)
}

func drawSquare(screen *ebiten.Image, x, y float32, col color.Color) {
	half := playerSize / 2
	vector.FillRect(
		screen,
		x-half,
		y-half,
		playerSize,
		playerSize,
		col,
		true,
	)
}
