class Game {
	constructor() {
		// 初始化默认棋盘
		this.board = [
			['R', 'R', 'R', 'R'],
			['.', '.', '.', '.'],
			['.', '.', '.', '.'],
			['B', 'B', 'B', 'B']
		];
		this.currentPlayer = 'R';
		this.selectedPiece = null;
		this.ws = null;
		this.init();
	}

	init() {
		this.setupWebSocket();
		this.renderBoard();
	}

	setupWebSocket() {
		const wsUrl = 'ws://' + window.location.host + '/ws';
		console.log('Connecting to WebSocket:', wsUrl);
		this.ws = new WebSocket(wsUrl);

		this.ws.onopen = () => {
			console.log('WebSocket connected successfully');
		};

		this.ws.onmessage = (event) => {
			const data = JSON.parse(event.data);
			if (data.error) {
				document.getElementById('message').textContent = data.error;
			} else if (data.winner) {
				this.handleWinner(data.winner);
			} else {
				// 更新棋盘和当前玩家
				this.board = data.board;
				this.currentPlayer = data.currentPlayer;
				this.renderBoard();
				this.updatePlayerDisplay();
			}
		};

		this.ws.onclose = () => {
			console.log('WebSocket disconnected');
		};

		this.ws.onerror = (error) => {
			console.error('WebSocket error:', error);
		};
	}

	renderBoard() {
		console.log('Rendering board:', JSON.stringify(this.board));
		const boardElement = document.getElementById('board');
		if (!boardElement) {
			console.error('Board element not found!');
			return;
		}
		console.log('Board element found:', boardElement);
		
		// 清空棋盘
		boardElement.innerHTML = '';
		
		// 创建4x4网格
		for (let row = 0; row < 4; row++) {
			const rowElement = document.createElement('div');
			rowElement.className = 'row';
			rowElement.style.display = 'flex';

			for (let col = 0; col < 4; col++) {
				const cellElement = document.createElement('div');
				cellElement.className = 'cell';
				cellElement.style.cssText = 'width:60px;height:60px;border:1px solid #ccc;display:flex;align-items:center;justify-content:center;cursor:pointer;background:#fff;';
				cellElement.dataset.row = row;
				cellElement.dataset.col = col;

				const piece = this.board[row][col];
				console.log(`Cell [${row},${col}] = ${piece}`);
				
				if (piece === 'R' || piece === 'B') {
					const pieceElement = document.createElement('div');
					pieceElement.style.cssText = `width:40px;height:40px;border-radius:50%;display:flex;align-items:center;justify-content:center;font-weight:bold;color:white;background-color:${piece === 'R' ? '#d9534f' : '#333'};`;
					pieceElement.textContent = piece;
					cellElement.appendChild(pieceElement);
				}

				cellElement.addEventListener('click', () => this.handleCellClick(row, col));
				rowElement.appendChild(cellElement);
			}

			boardElement.appendChild(rowElement);
		}
		console.log('Board rendered successfully');
	}

	handleCellClick(row, col) {
		// 如果没有棋盘数据，直接返回
		if (!this.board || this.board.length === 0) return;

		const piece = this.board[row] ? this.board[row][col] : null;

		// 如果点击的是当前玩家的棋子，选中它
		if (piece === this.currentPlayer) {
			this.selectedPiece = { row, col };
			this.highlightSelected();
			return;
		}

		// 如果已经选中了棋子，尝试移动到空位
		if (this.selectedPiece && piece === '.') {
			const move = {
				fromRow: this.selectedPiece.row,
				fromCol: this.selectedPiece.col,
				toRow: row,
				toCol: col
			};

			this.ws.send(JSON.stringify(move));
			this.selectedPiece = null;
			document.getElementById('message').textContent = '';
		}
	}

	highlightSelected() {
		// 移除所有选中状态（重置所有格子的边框和背景）
		const cells = document.querySelectorAll('.cell');
		cells.forEach(cell => {
			cell.style.border = '1px solid #ccc';
			cell.style.backgroundColor = '#fff';
		});

		// 高亮选中的棋子
		if (this.selectedPiece) {
			const cell = document.querySelector(`[data-row="${this.selectedPiece.row}"][data-col="${this.selectedPiece.col}"]`);
			if (cell) {
				cell.style.border = '3px solid #5bc0de';
				cell.style.backgroundColor = '#e6e6e6';
			}
		}
	}

	updatePlayerDisplay() {
		document.getElementById('current-player').textContent = this.currentPlayer === 'R' ? '红方' : '黑方';
	}

	handleWinner(winner) {
		const winnerText = winner === 'R' ? '红方' : '黑方';
		document.getElementById('message').textContent = `游戏结束！${winnerText}获胜！`;
		// 禁用棋盘点击
		const cells = document.querySelectorAll('.cell');
		cells.forEach(cell => {
			cell.removeEventListener('click', this.handleCellClick);
			cell.style.cursor = 'default';
		});
	}
}

// 初始化游戏
document.addEventListener('DOMContentLoaded', function() {
	console.log('DOM loaded, initializing game...');
	new Game();
});
