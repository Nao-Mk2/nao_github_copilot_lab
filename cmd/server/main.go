package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Nao-Mk2/nao_github_copilot_lab/internal/handler"
	"github.com/Nao-Mk2/nao_github_copilot_lab/internal/repository/memory"
	"github.com/Nao-Mk2/nao_github_copilot_lab/internal/service"
)

const (
	// サーバーのデフォルトポート
	defaultPort = "8080"
)

func main() {
	// ポート設定
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// 依存関係の初期化
	// リポジトリ層
	todoRepo := memory.NewMemoryTodoRepository()

	// サービス層
	todoService := service.NewTodoService(todoRepo)

	// ハンドラー層
	todoHandler := handler.NewTodoHandler(todoService)

	// ルーティング設定
	mux := http.NewServeMux()
	mux.Handle("/todo", todoHandler)
	mux.Handle("/todo/", todoHandler)

	// HTTPサーバーの設定
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// サーバー起動を別のゴルーチンで実行
	go func() {
		fmt.Printf("サーバーを開始します (ポート: %s)\n", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("サーバー起動エラー: %v", err)
		}
	}()

	// シグナル処理（Ctrl+C などで適切に終了させるため）
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// シグナルを待つ
	<-stop
	fmt.Println("サーバーを停止しています...")
}
