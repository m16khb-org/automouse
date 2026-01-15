// /Users/habin/workspace/automouse/cmd/root.go
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "automouse",
	Short: "macOS 자동 클릭 CLI 도구",
	Long: `automouse는 macOS용 자동 클릭 CLI 도구입니다.

'watch' 모드: 현재 마우스 좌표 확인
'click' 모드: 지정한 좌표에서 자동 클릭

주의: macOS 시스템 환경설정에서 접근성 권한이 필요합니다.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
