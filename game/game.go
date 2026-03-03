package game

// Game 游戏状态
type Game struct {
	Board         [4][4]string
	CurrentPlayer string
}

// Move 移动请求
type Move struct {
	FromRow int `json:"fromRow"`
	FromCol int `json:"fromCol"`
	ToRow   int `json:"toRow"`
	ToCol   int `json:"toCol"`
}

// NewGame 创建新游戏
func NewGame() *Game {
	board := [4][4]string{
		{"R", "R", "R", "R"},
		{".", ".", ".", "."},
		{".", ".", ".", "."},
		{"B", "B", "B", "B"},
	}

	return &Game{
		Board:         board,
		CurrentPlayer: "R",
	}
}

// ValidateMove 验证移动是否合法
func (g *Game) ValidateMove(move Move) (bool, string) {
	// 检查坐标是否在棋盘范围内
	if move.FromRow < 0 || move.FromRow >= 4 || move.FromCol < 0 || move.FromCol >= 4 {
		return false, "起始位置超出棋盘范围"
	}
	if move.ToRow < 0 || move.ToRow >= 4 || move.ToCol < 0 || move.ToCol >= 4 {
		return false, "目标位置超出棋盘范围"
	}

	// 检查起始位置是否有当前玩家的棋子
	if g.Board[move.FromRow][move.FromCol] != g.CurrentPlayer {
		return false, "起始位置没有当前玩家的棋子"
	}

	// 检查目标位置是否为空
	if g.Board[move.ToRow][move.ToCol] != "." {
		return false, "目标位置不为空"
	}

	// 检查移动是否是相邻的
	drow := move.ToRow - move.FromRow
	dcol := move.ToCol - move.FromCol
	if (drow != 0 && dcol != 0) || (abs(drow) > 1) || (abs(dcol) > 1) {
		return false, "只能移动到相邻的位置"
	}

	return true, ""
}

// MakeMove 执行移动
func (g *Game) MakeMove(move Move) {
	// 移动棋子
	g.Board[move.ToRow][move.ToCol] = g.Board[move.FromRow][move.FromCol]
	g.Board[move.FromRow][move.FromCol] = "."

	// 检查吃子
	g.CheckCapture(move.ToRow, move.ToCol)

	// 切换玩家
	if g.CurrentPlayer == "R" {
		g.CurrentPlayer = "B"
	} else {
		g.CurrentPlayer = "R"
	}
}

// CheckCapture 检查吃子
func (g *Game) CheckCapture(row, col int) {
	// 检查垂直方向
	g.checkVerticalCapture(col)
}

// checkVerticalCapture 检查垂直方向的吃子
func (g *Game) checkVerticalCapture(col int) {
	// 检查所有连续的3个位置
	for i := 0; i <= 1; i++ {
		// 获取连续3个位置的棋子
		piece1 := g.Board[i][col]
		piece2 := g.Board[i+1][col]
		piece3 := g.Board[i+2][col]

		// 检查是否形成2打1，且都不能为空
		if piece1 != "." && piece2 != "." && piece3 != "." {
			// 情况1：前两个相同，第三个不同（2-1结构，如红红黑）
			if piece1 == piece2 && piece1 != piece3 {
				// 检查第4个位置，确保不是2-2结构
				if i+3 >= 4 || g.Board[i+3][col] == "." || g.Board[i+3][col] == piece1 {
					// 吃掉第三个棋子
					g.Board[i+2][col] = "."
					// 递归检查是否有新的吃子
					g.checkVerticalCapture(col)
					break
				}
			}
			// 情况2：后两个相同，第一个不同（1-2结构，如红黑黑）
			if piece2 == piece3 && piece2 != piece1 {
				// 检查前一个位置，确保不是2-2结构
				if i-1 < 0 || g.Board[i-1][col] == "." || g.Board[i-1][col] == piece2 {
					// 吃掉第一个棋子
					g.Board[i][col] = "."
					// 递归检查是否有新的吃子
					g.checkVerticalCapture(col)
					break
				}
			}
		}
	}
}

// CheckWinner 检查胜负
func (g *Game) CheckWinner() string {
	rCount := 0
	bCount := 0

	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if g.Board[i][j] == "R" {
				rCount++
			} else if g.Board[i][j] == "B" {
				bCount++
			}
		}
	}

	if rCount <= 1 {
		return "B"
	}
	if bCount <= 1 {
		return "R"
	}

	// 检查是否有可移动的棋子
	hasMove := false
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if g.Board[i][j] == g.CurrentPlayer {
				// 检查上下左右是否有可移动的位置
				directions := [][]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
				for _, dir := range directions {
					ni, nj := i+dir[0], j+dir[1]
					if ni >= 0 && ni < 4 && nj >= 0 && nj < 4 && g.Board[ni][nj] == "." {
						hasMove = true
						break
					}
				}
				if hasMove {
					break
				}
			}
			if hasMove {
				break
			}
		}
	}

	if !hasMove {
		if g.CurrentPlayer == "R" {
			return "B"
		} else {
			return "R"
		}
	}

	return ""
}

// abs 计算绝对值
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
