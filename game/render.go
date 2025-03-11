package game

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	lightSquareColor = color.RGBA{240, 217, 181, 255}
	darkSquareColor  = color.RGBA{181, 136, 99, 255}
	highlightColor   = color.RGBA{130, 151, 105, 200}
	moveColor        = color.RGBA{130, 151, 105, 120}
	victoryColor     = color.RGBA{255, 215, 0, 180} // Gold color for victory animation
)

// RenderBoard draws the chess board and pieces
func RenderBoard(screen *ebiten.Image, game *Game) {
	// Draw board squares
	squareSize := float32(BoardSize) / 8
	for row := 0; row < 8; row++ {
		for col := 0; col < 8; col++ {
			x := float32(col) * squareSize
			y := float32(row) * squareSize

			// Determine square color
			var squareColor color.Color
			if (row+col)%2 == 0 {
				squareColor = lightSquareColor
			} else {
				squareColor = darkSquareColor
			}

			// Draw square
			vector.DrawFilledRect(screen, x, y, squareSize, squareSize, squareColor, false)

			// Draw piece
			piece := game.Board[row][col]
			if piece != 0 {
				drawPiece(screen, piece, x, y)
			}
		}
	}

	// Draw selected square highlight
	if game.SelectedPiece.Selected {
		x := float32(game.SelectedPiece.X) * squareSize
		y := float32(game.SelectedPiece.Y) * squareSize
		vector.DrawFilledRect(screen, x, y, squareSize, squareSize, highlightColor, false)
	}

	// Draw valid moves
	for _, move := range game.ValidMoves {
		x := float32(move.X) * squareSize
		y := float32(move.Y) * squareSize
		vector.DrawFilledRect(screen, x, y, squareSize, squareSize, moveColor, false)
	}

	// Draw victory animation if game is over
	if game.State != Playing {
		drawVictoryAnimation(screen, game)
	}
}

// drawVictoryAnimation creates a pulsing overlay with text
func drawVictoryAnimation(screen *ebiten.Image, game *Game) {
	// Calculate animation alpha based on tick
	alpha := float64(game.AnimationTick%120) / 120.0 // Slower pulse
	alpha = math.Sin(alpha * math.Pi * 2)
	alpha = (alpha + 1) / 2 // Normalize to 0-1 range

	// Create overlay color with animated alpha
	overlayColor := color.RGBA{
		R: victoryColor.R,
		G: victoryColor.G,
		B: victoryColor.B,
		A: uint8(float64(victoryColor.A) * alpha),
	}

	// Draw full screen overlay
	vector.DrawFilledRect(screen, 0, 0, float32(BoardSize), float32(BoardSize), overlayColor, false)

	// Draw victory text
	var message string
	if game.State == WhiteWins {
		message = "Checkmate! White Wins!"
	} else {
		message = "Checkmate! Black Wins!"
	}

	// Center the text
	bounds := text.BoundString(defaultFont, message)
	x := (BoardSize - float64(bounds.Dx())) / 2
	y := BoardSize/2 + float64(bounds.Dy())/2

	// Draw text with glow effect
	glowColor := color.RGBA{0, 0, 0, uint8(200 * (1 - alpha))}
	for dx := -2; dx <= 2; dx++ {
		for dy := -2; dy <= 2; dy++ {
			if dx*dx+dy*dy <= 4 { // Only draw within a circular radius
				text.Draw(screen, message, defaultFont, int(x)+dx, int(y)+dy, glowColor)
			}
		}
	}

	// Draw main text
	text.Draw(screen, message, defaultFont, int(x), int(y), color.White)
}

// drawPiece draws a chess piece image
func drawPiece(screen *ebiten.Image, piece int, x, y float32) {
	if img, ok := PieceImages[piece]; ok {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(x), float64(y))
		screen.DrawImage(img, op)
	}
}

// GetBoardCoordinates converts screen coordinates to board coordinates
func GetBoardCoordinates(x, y int) (int, int) {
	squareSize := BoardSize / 8
	boardX := x / squareSize
	boardY := y / squareSize
	return boardX, boardY
}

// IsInsideBoard checks if screen coordinates are within the board
func IsInsideBoard(x, y int) bool {
	boardX, boardY := GetBoardCoordinates(x, y)
	return boardX >= 0 && boardX < 8 && boardY >= 0 && boardY < 8
}
