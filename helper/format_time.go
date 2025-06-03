package helper

import (
	"fmt"
	"strings"
)

func FormatMilliseconds(ms int) string {
	totalSeconds := ms / 1000
	minutes := totalSeconds / 60
	seconds := totalSeconds % 60
	milliseconds := ms % 1000

	var b strings.Builder
	b.Grow(9) // 预分配足够长度以减少内存分配和复制操作

	b.WriteString(fmt.Sprintf("%02d", minutes))
	b.WriteRune(':')
	b.WriteString(fmt.Sprintf("%02d", seconds))
	b.WriteRune('.')
	b.WriteString(fmt.Sprintf("%03d", milliseconds))

	return b.String()
}
