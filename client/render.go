package client

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"coin-collector/common"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	localColor  = color.RGBA{0, 200, 0, 255}
	remoteColor = color.RGBA{200, 0, 0, 255}
	bgColor     = color.RGBA{30, 30, 30, 255}

	coinFrames     []*ebiten.Image
	coinFrameCount = 5
	coinFrameSize  = 16
)

const playerSize float32 = 26

func init() {
	img, _, err := ebitenutil.NewImageFromFile("../../assets/coin.png")
	if err != nil {
		log.Fatal("Failed to load coin.png:", err)
	}

	coinFrames = make([]*ebiten.Image, coinFrameCount)
	for i := 0; i < coinFrameCount; i++ {
		sx := i * coinFrameSize
		sub := img.SubImage(image.Rect(sx, 0, sx+coinFrameSize, coinFrameSize)).(*ebiten.Image)
		coinFrames[i] = sub
	}
}

func DrawWorld(screen *ebiten.Image, state WorldSnapshot, localID common.PlayerID, localMask uint8) {
	screen.Fill(bgColor)

	for _, c := range state.Coins {
		drawCoin(screen, c.X, c.Y)
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

func drawCoin(screen *ebiten.Image, x, y float32) {
	if len(coinFrames) == 0 {
		return
	}

	frame := (NowMs() / 100) % int64(coinFrameCount)
	img := coinFrames[frame]

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(
		float64(x)-float64(coinFrameSize)/2,
		float64(y)-float64(coinFrameSize)/2,
	)

	screen.DrawImage(img, op)
}
