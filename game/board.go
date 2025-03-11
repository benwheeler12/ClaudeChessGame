package game

// Piece constants
const (
	Empty = iota
	Pawn
	Knight
	Bishop
	Rook
	Queen
	King
)

// Color constants
const (
	White = 1
	Black = -1
)

// Position represents a square on the chess board
type Position struct {
	X, Y int
}

// IsValidPosition checks if a position is within the board boundaries
func IsValidPosition(pos Position) bool {
	return pos.X >= 0 && pos.X < 8 && pos.Y >= 0 && pos.Y < 8
}

// GetPieceMoves returns all valid moves for a piece at the given position
func GetPieceMoves(board [8][8]int, pos Position, game *Game) []Position {
	piece := abs(board[pos.Y][pos.X])
	color := sign(board[pos.Y][pos.X])
	var moves []Position

	switch piece {
	case Pawn:
		moves = getPawnMoves(board, pos, color, game)
	case Knight:
		moves = getKnightMoves(board, pos, color)
	case Bishop:
		moves = getBishopMoves(board, pos, color)
	case Rook:
		moves = getRookMoves(board, pos, color)
	case Queen:
		moves = getQueenMoves(board, pos, color)
	case King:
		moves = getKingMoves(board, pos, color, game)
	}

	return moves
}

func getPawnMoves(board [8][8]int, pos Position, color int, game *Game) []Position {
	moves := make([]Position, 0)
	direction := -color // Pawns move up for white (negative) and down for black (positive)

	// Forward move
	newPos := Position{pos.X, pos.Y + direction}
	if IsValidPosition(newPos) && board[newPos.Y][newPos.X] == Empty {
		moves = append(moves, newPos)

		// Initial two-square move
		if (color == White && pos.Y == 6) || (color == Black && pos.Y == 1) {
			newPos = Position{pos.X, pos.Y + 2*direction}
			if board[newPos.Y][newPos.X] == Empty {
				moves = append(moves, newPos)
			}
		}
	}

	// Regular captures
	for _, dx := range []int{-1, 1} {
		newPos := Position{pos.X + dx, pos.Y + direction}
		if IsValidPosition(newPos) {
			target := board[newPos.Y][newPos.X]
			if target != Empty && sign(target) != color {
				moves = append(moves, newPos)
			}
		}
	}

	// En passant capture
	if game != nil && game.EnPassantTarget != nil {
		expectedRank := boolToInt(color == White, 3, 4) // White on rank 4, Black on rank 3
		if pos.Y == expectedRank {
			for _, dx := range []int{-1, 1} {
				if pos.X+dx == game.EnPassantTarget.X && pos.Y+direction == game.EnPassantTarget.Y {
					moves = append(moves, *game.EnPassantTarget)
				}
			}
		}
	}

	return moves
}

func getKnightMoves(board [8][8]int, pos Position, color int) []Position {
	moves := make([]Position, 0)
	directions := [][2]int{
		{-2, -1}, {-2, 1}, {-1, -2}, {-1, 2},
		{1, -2}, {1, 2}, {2, -1}, {2, 1},
	}

	for _, d := range directions {
		newPos := Position{pos.X + d[0], pos.Y + d[1]}
		if IsValidPosition(newPos) {
			target := board[newPos.Y][newPos.X]
			if target == Empty || sign(target) != color {
				moves = append(moves, newPos)
			}
		}
	}

	return moves
}

func getBishopMoves(board [8][8]int, pos Position, color int) []Position {
	return getDiagonalMoves(board, pos, color)
}

func getRookMoves(board [8][8]int, pos Position, color int) []Position {
	return getStraightMoves(board, pos, color)
}

func getQueenMoves(board [8][8]int, pos Position, color int) []Position {
	moves := getDiagonalMoves(board, pos, color)
	moves = append(moves, getStraightMoves(board, pos, color)...)
	return moves
}

func getKingMoves(board [8][8]int, pos Position, color int, game *Game) []Position {
	moves := make([]Position, 0)

	// Normal king moves
	directions := [][2]int{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1}, {0, 1},
		{1, -1}, {1, 0}, {1, 1},
	}

	for _, d := range directions {
		newPos := Position{pos.X + d[0], pos.Y + d[1]}
		if IsValidPosition(newPos) {
			target := board[newPos.Y][newPos.X]
			if target == Empty || sign(target) != color {
				moves = append(moves, newPos)
			}
		}
	}

	// Castling moves (only if game state is provided)
	if game != nil {
		rookY := pos.Y
		// Check kingside castling
		rookPos := Position{7, rookY}
		if !game.HasMoved[pos] && !game.HasMoved[rookPos] && // Neither king nor rook has moved
			board[rookY][5] == Empty && // Squares between are empty
			board[rookY][6] == Empty &&
			!IsKingInCheck(board, color) && // King is not in check
			!wouldBeInCheck(board, pos, Position{pos.X + 1, pos.Y}, color) { // King doesn't pass through check
			moves = append(moves, Position{pos.X + 2, pos.Y})
		}

		// Check queenside castling
		rookPos = Position{0, rookY}
		if !game.HasMoved[pos] && !game.HasMoved[rookPos] && // Neither king nor rook has moved
			board[rookY][1] == Empty && // Squares between are empty
			board[rookY][2] == Empty &&
			board[rookY][3] == Empty &&
			!IsKingInCheck(board, color) && // King is not in check
			!wouldBeInCheck(board, pos, Position{pos.X - 1, pos.Y}, color) { // King doesn't pass through check
			moves = append(moves, Position{pos.X - 2, pos.Y})
		}
	}

	return moves
}

func getDiagonalMoves(board [8][8]int, pos Position, color int) []Position {
	moves := make([]Position, 0)
	directions := [][2]int{{-1, -1}, {-1, 1}, {1, -1}, {1, 1}}

	for _, d := range directions {
		for i := 1; i < 8; i++ {
			newPos := Position{pos.X + i*d[0], pos.Y + i*d[1]}
			if !IsValidPosition(newPos) {
				break
			}
			target := board[newPos.Y][newPos.X]
			if target == Empty {
				moves = append(moves, newPos)
			} else {
				if sign(target) != color {
					moves = append(moves, newPos)
				}
				break
			}
		}
	}

	return moves
}

func getStraightMoves(board [8][8]int, pos Position, color int) []Position {
	moves := make([]Position, 0)
	directions := [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

	for _, d := range directions {
		for i := 1; i < 8; i++ {
			newPos := Position{pos.X + i*d[0], pos.Y + i*d[1]}
			if !IsValidPosition(newPos) {
				break
			}
			target := board[newPos.Y][newPos.X]
			if target == Empty {
				moves = append(moves, newPos)
			} else {
				if sign(target) != color {
					moves = append(moves, newPos)
				}
				break
			}
		}
	}

	return moves
}

// Helper functions
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func sign(x int) int {
	if x < 0 {
		return -1
	}
	if x > 0 {
		return 1
	}
	return 0
}

// IsKingInCheck determines if the specified color's king is in check
func IsKingInCheck(board [8][8]int, color int) bool {
	// Find the king's position
	var kingPos Position
	kingValue := color * King
	found := false
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			if board[y][x] == kingValue {
				kingPos = Position{X: x, Y: y}
				found = true
				break
			}
		}
		if found {
			break
		}
	}

	// Check if any opponent's piece can capture the king
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			piece := board[y][x]
			if piece != 0 && sign(piece) != color {
				moves := GetPieceMoves(board, Position{X: x, Y: y}, nil)
				for _, move := range moves {
					if move.X == kingPos.X && move.Y == kingPos.Y {
						return true
					}
				}
			}
		}
	}
	return false
}

// SimulateMove simulates a move and returns true if it's legal (doesn't put own king in check)
func SimulateMove(board [8][8]int, from, to Position) bool {
	// Make a copy of the board
	var tempBoard [8][8]int
	for i := range board {
		copy(tempBoard[i][:], board[i][:])
	}

	// Simulate the move
	piece := tempBoard[from.Y][from.X]
	tempBoard[from.Y][from.X] = Empty
	tempBoard[to.Y][to.X] = piece

	// Check if the move puts/leaves own king in check
	return !IsKingInCheck(tempBoard, sign(piece))
}

// GetLegalMoves returns all legal moves for a piece (excluding moves that put own king in check)
func GetLegalMoves(board [8][8]int, pos Position) []Position {
	moves := GetPieceMoves(board, pos, nil)
	legalMoves := make([]Position, 0)

	for _, move := range moves {
		if SimulateMove(board, pos, move) {
			legalMoves = append(legalMoves, move)
		}
	}

	return legalMoves
}

// GetLegalMovesWithState returns all legal moves including special moves like castling and en passant
func GetLegalMovesWithState(board [8][8]int, pos Position, game *Game) []Position {
	moves := GetPieceMoves(board, pos, game)
	legalMoves := make([]Position, 0)

	for _, move := range moves {
		if SimulateMove(board, pos, move) {
			legalMoves = append(legalMoves, move)
		}
	}

	return legalMoves
}

// IsCheckmate determines if the specified color is in checkmate
func IsCheckmate(board [8][8]int, color int) bool {
	// If not in check, it's not checkmate
	if !IsKingInCheck(board, color) {
		return false
	}

	// Check if any piece has a legal move
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			piece := board[y][x]
			if piece != 0 && sign(piece) == color {
				moves := GetLegalMoves(board, Position{X: x, Y: y})
				if len(moves) > 0 {
					return false
				}
			}
		}
	}

	return true
}

// UpdateGameState checks for checkmate and updates the game state accordingly
func (g *Game) UpdateGameState() {
	if IsCheckmate(g.Board, 1) { // Check if White is in checkmate
		g.State = BlackWins
	} else if IsCheckmate(g.Board, -1) { // Check if Black is in checkmate
		g.State = WhiteWins
	}
}

// Helper function to check if a move would put the king in check
func wouldBeInCheck(board [8][8]int, from, to Position, color int) bool {
	// Make a copy of the board
	var tempBoard [8][8]int
	for i := range board {
		copy(tempBoard[i][:], board[i][:])
	}

	// Simulate the move
	piece := tempBoard[from.Y][from.X]
	tempBoard[from.Y][from.X] = Empty
	tempBoard[to.Y][to.X] = piece

	return IsKingInCheck(tempBoard, color)
}

// boolToInt converts a bool to an int
func boolToInt(b bool, trueVal, falseVal int) int {
	if b {
		return trueVal
	}
	return falseVal
}

// MakeMove performs a move and handles special cases like castling and en passant
func (g *Game) MakeMove(from, to Position) {
	piece := g.Board[from.Y][from.X]

	// Update HasMoved for castling
	g.HasMoved[from] = true

	// Handle castling
	if abs(piece) == King && abs(to.X-from.X) == 2 {
		// Kingside castling
		if to.X > from.X {
			g.Board[from.Y][5] = g.Board[from.Y][7] // Move rook
			g.Board[from.Y][7] = Empty
		} else { // Queenside castling
			g.Board[from.Y][3] = g.Board[from.Y][0] // Move rook
			g.Board[from.Y][0] = Empty
		}
	}

	// Handle en passant capture
	if abs(piece) == Pawn && g.EnPassantTarget != nil &&
		to.X == g.EnPassantTarget.X && to.Y == g.EnPassantTarget.Y {
		g.Board[from.Y][to.X] = Empty // Remove captured pawn
	}

	// Update en passant target
	g.EnPassantTarget = nil
	if abs(piece) == Pawn && abs(to.Y-from.Y) == 2 {
		g.EnPassantTarget = &Position{to.X, (from.Y + to.Y) / 2}
	}

	// Make the move
	g.Board[to.Y][to.X] = piece
	g.Board[from.Y][from.X] = Empty

	// Update last move
	g.LastMove.From = from
	g.LastMove.To = to
	g.LastMove.Piece = piece

	// Switch turns
	g.Turn = !g.Turn
}
