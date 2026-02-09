package hotkey

import (
	"context"

	hook "github.com/robotn/gohook"
)

// ESC keycode (cross-platform)
const KeycodeEsc = 53 // macOS ESC keycode

// ListenForEscape ESC 키가 눌리면 cancel을 호출합니다.
// 터미널 포커스 없이도 글로벌하게 동작합니다.
func ListenForEscape(ctx context.Context, cancel context.CancelFunc) {
	evChan := hook.Start()
	defer hook.End()

	for {
		select {
		case <-ctx.Done():
			return
		case ev := <-evChan:
			if ev.Kind == hook.KeyDown && ev.Rawcode == KeycodeEsc {
				cancel()
				return
			}
		}
	}
}
