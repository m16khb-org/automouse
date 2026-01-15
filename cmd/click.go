// /Users/habin/workspace/automouse/cmd/click.go
package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/habin/automouse/internal/clicker"
	"github.com/spf13/cobra"
)

var (
	clickX        int
	clickY        int
	clickInterval int
)

var clickCmd = &cobra.Command{
	Use:   "click",
	Short: "지정한 좌표에서 자동 클릭",
	Long: `click 모드는 지정한 x,y 좌표에서 자동으로 클릭합니다.
Enter를 눌러 클릭 시작/중지를 토글합니다. Ctrl+C로 종료합니다.

예시:
  automouse click --x=500 --y=300 --interval=1000`,
	Run: runClick,
}

func init() {
	rootCmd.AddCommand(clickCmd)

	clickCmd.Flags().IntVarP(&clickX, "x", "x", 0, "클릭할 X 좌표 (필수)")
	clickCmd.Flags().IntVarP(&clickY, "y", "y", 0, "클릭할 Y 좌표 (필수)")
	clickCmd.Flags().IntVarP(&clickInterval, "interval", "i", 1000, "클릭 간격 (밀리초)")

	clickCmd.MarkFlagRequired("x")
	clickCmd.MarkFlagRequired("y")
}

func runClick(cmd *cobra.Command, args []string) {
	if clickX < 0 || clickY < 0 {
		fmt.Fprintln(os.Stderr, "오류: 좌표는 0 이상이어야 합니다")
		os.Exit(1)
	}

	if clickInterval < 10 {
		fmt.Fprintln(os.Stderr, "오류: interval은 최소 10ms 이상이어야 합니다")
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh
		cancel()
	}()

	toggleCh := make(chan struct{})

	go listenForEnter(ctx, toggleCh)

	c := clicker.New(clickX, clickY, time.Duration(clickInterval)*time.Millisecond)
	c.Run(ctx, toggleCh)
}

func listenForEnter(ctx context.Context, toggleCh chan<- struct{}) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			if scanner.Scan() {
				select {
				case toggleCh <- struct{}{}:
				case <-ctx.Done():
					return
				}
			} else {
				return
			}
		}
	}
}
