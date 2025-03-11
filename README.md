# Chess Game

A beautiful chess game implementation in Go using the Ebiten game engine. Features include:

- Full chess rules implementation
- Special moves (castling, en passant)
- SVG piece graphics
- Victory animations
- Checkmate detection
- Legal move highlighting

## Requirements

- Go 1.16 or later
- Ebiten v2
- SVG rendering dependencies

## Installation

1. Clone the repository:
```bash
git clone https://github.com/YOUR_USERNAME/ClaudeChessGame.git
cd ClaudeChessGame
```

2. Install dependencies:
```bash
go mod download
```

3. Run the game:
```bash
go run main.go
```

## How to Play

- Click on a piece to select it
- Valid moves will be highlighted
- Click on a highlighted square to move the piece
- The game automatically detects checkmate and displays a victory animation

## Features

- Complete chess rules implementation including special moves:
  - Castling (kingside and queenside)
  - En passant captures
  - Pawn promotion (coming soon)
- Legal move validation
- Check and checkmate detection
- Beautiful SVG piece graphics
- Smooth animations
- Intuitive user interface

## License

MIT License - feel free to use this code for your own projects! 