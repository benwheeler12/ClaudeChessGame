package game

import (
	"fmt"
	"image"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
)

const (
	BoardSize = 480 // Makes each square 60x60
	PieceSize = 60  // Size of piece images
)

// GameState represents the current state of the game
type GameState int

const (
	Playing GameState = iota
	WhiteWins
	BlackWins
)

// Game represents the main game state
type Game struct {
	Board         [8][8]int // 0 = empty, positive = white pieces, negative = black pieces
	SelectedPiece struct {
		X, Y     int
		Selected bool
	}
	Turn          bool // true = white, false = black
	ValidMoves    []Position
	State         GameState
	AnimationTick int // Used for victory animation

	// Castling state
	HasMoved map[Position]bool // Tracks if pieces have moved (for castling)

	// En passant state
	LastMove struct {
		From, To Position
		Piece    int
	}
	EnPassantTarget *Position // Square where en passant capture is possible
}

var (
	defaultFont font.Face
	PieceImages map[int]*ebiten.Image // Maps piece type to its image
)

// loadPieceImage loads an SVG piece image and converts it to an Ebiten image
func loadPieceImage(path string) (*ebiten.Image, error) {
	// Read and parse SVG
	icon, err := oksvg.ReadIcon(path, oksvg.StrictErrorMode)
	if err != nil {
		return nil, fmt.Errorf("error reading SVG: %v", err)
	}

	// Set size
	icon.SetTarget(0, 0, float64(PieceSize), float64(PieceSize))

	// Create RGBA image
	rgba := image.NewRGBA(image.Rect(0, 0, PieceSize, PieceSize))
	scanner := rasterx.NewScannerGV(PieceSize, PieceSize, rgba, rgba.Bounds())
	raster := rasterx.NewDasher(PieceSize, PieceSize, scanner)

	// Render SVG
	icon.Draw(raster, 1.0)

	// Convert to Ebiten image
	return ebiten.NewImageFromImage(rgba), nil
}

// InitPieces loads all piece images
func InitPieces() error {
	PieceImages = make(map[int]*ebiten.Image)
	pieces := map[int]string{
		Pawn:   "p",
		Knight: "n",
		Bishop: "b",
		Rook:   "r",
		Queen:  "q",
		King:   "k",
	}

	for pieceType, letter := range pieces {
		// Load white piece
		whitePath := filepath.Join("assets", "pieces", fmt.Sprintf("w%s.svg", letter))
		whiteImg, err := loadPieceImage(whitePath)
		if err != nil {
			return fmt.Errorf("error loading white piece %s: %v", letter, err)
		}
		PieceImages[pieceType] = whiteImg

		// Load black piece
		blackPath := filepath.Join("assets", "pieces", fmt.Sprintf("b%s.svg", letter))
		blackImg, err := loadPieceImage(blackPath)
		if err != nil {
			return fmt.Errorf("error loading black piece %s: %v", letter, err)
		}
		PieceImages[-pieceType] = blackImg
	}

	return nil
}

// InitFonts initializes the fonts used in the game
func InitFonts() error {
	tt, err := opentype.Parse(goregular.TTF)
	if err != nil {
		return err
	}

	const dpi = 72
	defaultFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    36,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return err
	}

	return nil
}

// NewGame creates and initializes a new game
func NewGame() *Game {
	g := &Game{
		Turn:     true, // White starts
		State:    Playing,
		HasMoved: make(map[Position]bool),
	}
	g.initializeBoard()
	return g
}

// initializeBoard sets up the initial chess board position
func (g *Game) initializeBoard() {
	// Set up pawns
	for i := 0; i < 8; i++ {
		g.Board[1][i] = -Pawn // Black pawns
		g.Board[6][i] = Pawn  // White pawns
	}

	// Set up other pieces
	pieceOrder := []int{Rook, Knight, Bishop, Queen, King, Bishop, Knight, Rook}
	for i, piece := range pieceOrder {
		g.Board[0][i] = -piece // Black pieces
		g.Board[7][i] = piece  // White pieces
	}
}
