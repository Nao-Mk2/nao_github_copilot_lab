package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Nao-Mk2/nao_github_copilot_lab/internal/handler"
	"github.com/Nao-Mk2/nao_github_copilot_lab/internal/repository"
)

const (
	// サーバーのデフォルトポート
	defaultPort = "8080"
	// サーバーのデフォルトのベースパス
	basePath = "/todo"
)

func main() {
	// サーバーポートの設定（環境変数から取得、なければデフォルト値を使用）
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// リポジトリの初期化
	todoRepo := repository.NewInMemoryTodoRepository()

	// ハンドラーの初期化
	todoHandler := handler.NewTodoHandler(todoRepo)

	// ルーティングの設定
	mux := http.NewServeMux()
	mux.HandleFunc(basePath, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			todoHandler.CreateTodo(w, r)
		default:
			// GETリクエスト（IDが必要）またはその他の不正なメソッドをチェック
			if r.Method == http.MethodGet && len(r.URL.Path) > len(basePath)+1 {
				todoHandler.GetTodo(w, r)
			} else {
				http.Error(w, "不正なメソッドまたはURLです", http.StatusMethodNotAllowed)
			}
		}
	})

	// HTTPサーバーの設定
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	// サーバー起動を非同期で実行
	go func() {
		log.Printf("サーバーを起動しています。ポート: %s\n", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("サーバー起動エラー: %v", err)
		}
	}()

	// シグナル処理（Ctrl+Cなどでの終了処理）
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("サーバーをシャットダウンしています...")
	fmt.Println("プログラムを終了します")
}
