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
	"github.com/habin/automouse/internal/hotkey"
	"github.com/spf13/cobra"
)

var (
	clickX        int
	clickY        int
	clickInterval int
	clickTimeout  int
	useCurrentPos bool
)

var clickCmd = &cobra.Command{
	Use:   "click",
	Short: "지정한 좌표에서 자동 클릭",
	Long: `click 모드는 지정한 x,y 좌표에서 자동으로 클릭합니다.
좌표를 지정하지 않으면 현재 마우스 위치에서 클릭합니다.
Enter를 눌러 클릭 시작/중지를 토글합니다.
ESC 키를 누르거나 Ctrl+C로 종료합니다.

예시:
  automouse click --x=500 --y=300 --interval=1000
  automouse click --interval=500                    # 현재 위치에서 클릭
  automouse click --x=100 --y=200 --timeout=60      # 60초 후 자동 종료`,
	Run: runClick,
}

func init() {
	rootCmd.AddCommand(clickCmd)

	clickCmd.Flags().IntVarP(&clickX, "x", "x", -1, "클릭할 X 좌표 (미지정 시 현재 위치)")
	clickCmd.Flags().IntVarP(&clickY, "y", "y", -1, "클릭할 Y 좌표 (미지정 시 현재 위치)")
	clickCmd.Flags().IntVarP(&clickInterval, "interval", "i", 1000, "클릭 간격 (밀리초)")
	clickCmd.Flags().IntVarP(&clickTimeout, "timeout", "t", 0, "자동 종료 시간 (초, 0=무제한)")
}

func runClick(cmd *cobra.Command, args []string) {
	if clickInterval < 10 {
		fmt.Fprintln(os.Stderr, "오류: interval은 최소 10ms 이상이어야 합니다")
		os.Exit(1)
	}

	// x, y 중 하나만 지정된 경우 체크
	xSet := cmd.Flags().Changed("x")
	ySet := cmd.Flags().Changed("y")
	if xSet != ySet {
		fmt.Fprintln(os.Stderr, "오류: x와 y는 둘 다 지정하거나 둘 다 지정하지 않아야 합니다")
		os.Exit(1)
	}

	useCurrentPos = !xSet && !ySet

	if !useCurrentPos && (clickX < 0 || clickY < 0) {
		fmt.Fprintln(os.Stderr, "오류: 좌표는 0 이상이어야 합니다")
		os.Exit(1)
	}

	var ctx context.Context
	var cancel context.CancelFunc

	if clickTimeout > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), time.Duration(clickTimeout)*time.Second)
	} else {
		ctx, cancel = context.WithCancel(context.Background())
	}
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh
		cancel()
	}()

	// ESC 키 감지 (글로벌 핫키)
	go hotkey.ListenForEscape(ctx, cancel)

	toggleCh := make(chan struct{})

	go listenForEnter(ctx, toggleCh)

	c := clicker.New(clickX, clickY, time.Duration(clickInterval)*time.Millisecond, useCurrentPos)
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
