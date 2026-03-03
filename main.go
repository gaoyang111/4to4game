package main

import (
	"4to4game/game"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	// 处理静态文件
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		log.Println("WebSocket connection request received")
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Failed to upgrade connection:", err)
			return
		}
		log.Println("WebSocket connection established")
		defer conn.Close()

		// 初始化游戏
		g := game.NewGame()
		log.Println("Game initialized, sending initial board state")

		// 发送初始游戏状态（包含棋盘和当前玩家）
		gameState := map[string]interface{}{
			"board":         g.Board,
			"currentPlayer": g.CurrentPlayer,
		}
		err = conn.WriteJSON(gameState)
		if err != nil {
			log.Println("Failed to send initial state:", err)
			return
		}
		log.Println("Initial game state sent")

		// 处理消息
		for {
			var move game.Move
			err := conn.ReadJSON(&move)
			if err != nil {
				log.Println("Failed to read move:", err)
				break
			}

			// 执行移动
			valid, reason := g.ValidateMove(move)
			if !valid {
				conn.WriteJSON(map[string]string{"error": reason})
				continue
			}

			g.MakeMove(move)

			// 检查胜负
			winner := g.CheckWinner()
			gameState := map[string]interface{}{
				"board":         g.Board,
				"currentPlayer": g.CurrentPlayer,
				"winner":        winner,
			}
			conn.WriteJSON(gameState)
			if winner != "" {
				break
			}
		}
	})

	fmt.Println("Server started on :8082")
	log.Fatal(http.ListenAndServe(":8082", nil))
}
