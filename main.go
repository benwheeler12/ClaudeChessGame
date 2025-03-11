package main

import (
	"log"

	"chessgame/game"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  = 800
	screenHeight = 600
)

type Game struct {
	board *game.Game
}

func NewGame() *Game {
	return &Game{
		board: game.NewGame(),
	}
}

func (g *Game) Update() error {
	// Update animation tick if game is over
	if g.board.State != game.Playing {
		g.board.AnimationTick++
		return nil
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if game.IsInsideBoard(x, y) {
			boardX, boardY := game.GetBoardCoordinates(x, y)
			if !g.board.SelectedPiece.Selected {
				// Try to select a piece
				piece := g.board.Board[boardY][boardX]
				if piece != 0 && ((piece > 0) == g.board.Turn) {
					g.board.SelectedPiece.X = boardX
					g.board.SelectedPiece.Y = boardY
					g.board.SelectedPiece.Selected = true
					g.board.ValidMoves = game.GetLegalMovesWithState(g.board.Board, game.Position{X: boardX, Y: boardY}, g.board)
				}
			} else {
				// Try to move the selected piece
				targetPos := game.Position{X: boardX, Y: boardY}
				validMove := false
				for _, move := range g.board.ValidMoves {
					if move.X == targetPos.X && move.Y == targetPos.Y {
						validMove = true
						break
					}
				}

				if validMove {
					// Make the move
					from := game.Position{X: g.board.SelectedPiece.X, Y: g.board.SelectedPiece.Y}
					g.board.MakeMove(from, targetPos)

					// Check for checkmate
					if game.IsCheckmate(g.board.Board, 1) {
						g.board.State = game.BlackWins
					} else if game.IsCheckmate(g.board.Board, -1) {
						g.board.State = game.WhiteWins
					}
				}

				// Deselect the piece
				g.board.SelectedPiece.Selected = false
				g.board.ValidMoves = nil
			}
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	game.RenderBoard(screen, g.board)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Chess Game")

	if err := game.InitFonts(); err != nil {
		log.Fatal(err)
	}

	if err := game.InitPieces(); err != nil {
		log.Fatal(err)
	}

	game := NewGame()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
