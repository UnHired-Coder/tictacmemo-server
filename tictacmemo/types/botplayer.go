package types

// Point represents a move on the board.
type Point struct {
	X int
	Y int
}

// TicTacToeHelper provides methods to make smart moves.
type TicTacToeHelper struct {
	BoardSize int
}

// NewTicTacToeHelper creates a new instance of TicTacToeHelper.
func NewTicTacToeHelper() *TicTacToeHelper {
	return &TicTacToeHelper{BoardSize: 3}
}

// GetSmartMove returns the optimal move for the current player.
func (helper *TicTacToeHelper) GetSmartMove(board [3][3]string, currentPlayer string) *Point {
	opponent := helper.getOpponent(currentPlayer)

	// 1. Check if the current player can win
	if move := helper.findWinningMove(board, currentPlayer); move != nil && helper.isValidMove(board, *move) {
		return move
	}

	// 2. Block opponent's winning move
	if move := helper.findWinningMove(board, opponent); move != nil && helper.isValidMove(board, *move) {
		return move
	}

	// 3. Take the center if available
	if board[1][1] == "" {
		return &Point{X: 1, Y: 1}
	}

	// 4. Take any empty corner
	if move := helper.findEmptyCorner(board); move != nil && helper.isValidMove(board, *move) {
		return move
	}

	// 5. Take any empty side
	return helper.findEmptySide(board)
}

func (helper *TicTacToeHelper) getOpponent(player string) string {
	if player == "X" {
		return "O"
	}
	return "X"
}

func (helper *TicTacToeHelper) findWinningMove(board [3][3]string, player string) *Point {
	for i := 0; i < helper.BoardSize; i++ {
		// Check rows
		if helper.canWin(board[i][0], board[i][1], board[i][2], player) {
			return &Point{X: i, Y: helper.findEmptyIndex(board[i])}
		}
		// Check columns
		if helper.canWin(board[0][i], board[1][i], board[2][i], player) {
			return &Point{X: helper.findEmptyIndex([3]string{board[0][i], board[1][i], board[2][i]}), Y: i}
		}
	}
	// Check diagonals
	if helper.canWin(board[0][0], board[1][1], board[2][2], player) {
		return &Point{X: helper.findEmptyIndex([3]string{board[0][0], board[1][1], board[2][2]}), Y: helper.findEmptyIndex([3]string{board[0][0], board[1][1], board[2][2]})}
	}
	if helper.canWin(board[0][2], board[1][1], board[2][0], player) {
		return &Point{X: helper.findEmptyIndex([3]string{board[0][2], board[1][1], board[2][0]}), Y: helper.findEmptyIndex([3]string{board[0][2], board[1][1], board[2][0]})}
	}
	return nil
}

func (helper *TicTacToeHelper) canWin(a, b, c, player string) bool {
	return (a == player && b == player && c == "") ||
		(a == player && b == "" && c == player) ||
		(a == "" && b == player && c == player)
}

func (helper *TicTacToeHelper) findEmptyIndex(line [3]string) int {
	for i, state := range line {
		if state == "" {
			return i
		}
	}
	return -1
}

func (helper *TicTacToeHelper) findEmptyCorner(board [3][3]string) *Point {
	corners := []Point{{0, 0}, {0, 2}, {2, 0}, {2, 2}}
	for _, corner := range corners {
		if board[corner.X][corner.Y] == "" {
			return &corner
		}
	}
	return nil
}

func (helper *TicTacToeHelper) findEmptySide(board [3][3]string) *Point {
	sides := []Point{{0, 1}, {1, 0}, {1, 2}, {2, 1}}
	for _, side := range sides {
		if board[side.X][side.Y] == "" {
			return &side
		}
	}
	return nil
}

func (helper *TicTacToeHelper) isValidMove(board [3][3]string, move Point) bool {
	return board[move.X][move.Y] == ""
}
