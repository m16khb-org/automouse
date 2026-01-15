// /Users/habin/workspace/automouse/cmd/watch.go
package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/habin/automouse/internal/watcher"
	"github.com/spf13/cobra"
)

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "현재 마우스 좌표를 실시간으로 표시",
	Long: `watch 모드는 현재 마우스 좌표를 실시간으로 보여줍니다.
100ms마다 업데이트됩니다. Ctrl+C로 종료합니다.

예시:
  automouse watch`,
	Run: runWatch,
}

func init() {
	rootCmd.AddCommand(watchCmd)
}

func runWatch(cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh
		fmt.Println("\nwatch 모드 종료...")
		cancel()
	}()

	w := watcher.New(100 * time.Millisecond)
	positions := w.Start(ctx)

	fmt.Println("마우스 좌표 감시 중 (Ctrl+C로 종료)...")
	fmt.Println("마우스를 움직여 좌표를 확인하세요:")
	fmt.Println()

	for pos := range positions {
		fmt.Printf("\r  X: %4d  |  Y: %4d    ", pos.X, pos.Y)
	}

	fmt.Println()
}
