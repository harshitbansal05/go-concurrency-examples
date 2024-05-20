package sudoku

func SolveSudoku(board [][]byte) {
	solve(board)
}

func solve(board [][]byte) bool {
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if board[i][j] != byte(0) {
				continue
			}
			for k := 1; k <= 9; k++ {
				if !isByteAllowed(board, i, j, byte(k)) {
					continue
				}
				board[i][j] = byte(k)
				if solve(board) {
					return true
				}
			}
			board[i][j] = byte(0)
			return false
		}
	}
	return true
}

func isByteAllowed(board [][]byte, i, j int, x byte) bool {
	return isByteAllowedInBox(board, i, j, x) &&
		isByteAllowedInRow(board, i, x) &&
		isByteAllowedInCol(board, j, x)
}

func isByteAllowedInBox(board [][]byte, i, j int, x byte) bool {
	r, c := i/3, j/3

	for _, a := range []int{3 * r, 3*r + 1, 3*r + 2} {
		for _, b := range []int{3 * c, 3*c + 1, 3*c + 2} {
			y := board[a][b]
			if y == byte(0) {
				continue
			}
			if y == x {
				return false
			}
		}
	}
	return true
}

func isByteAllowedInRow(board [][]byte, i int, x byte) bool {
	for j := 0; j < 9; j++ {
		y := board[i][j]
		if y == byte(0) {
			continue
		}
		if y == x {
			return false
		}
	}
	return true
}

func isByteAllowedInCol(board [][]byte, j int, x byte) bool {
	for i := 0; i < 9; i++ {
		y := board[i][j]
		if y == byte(0) {
			continue
		}
		if y == x {
			return false
		}
	}
	return true
}
