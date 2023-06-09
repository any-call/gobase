package mylog

import "fmt"

type Color int

// 定义字体 颜色
const (
	FgBlack Color = iota + 30
	FgRed
	FgGreen
	FgYellow
	FgMagenta
	FgCyan
	FgWhite
)

const (
	ShowStyleDefault    uint = 0 //终端默认设置
	ShowStyleHighLight  uint = 1 //高亮显示
	ShowStyleUunderline uint = 4 //使用下划线
	ShowStyleFlashing   uint = 5 //闪烁
	ShowStyleAntiWhite  uint = 7 //反白显示
	ShowStyleInvisible  uint = 8 //不可见
)

func showStyle(showStyle uint, fgColor, bgColor Color, content string) string {
	return fmt.Sprintf("\033[%d;%d;%dm%s\033[0m", showStyle, fgColor, bgColor, content)
}

func debugStyle(content string) string {
	return showStyle(ShowStyleHighLight, 33, 40, content)
}

func infoStyle(content string) string {
	return showStyle(ShowStyleHighLight, 32, 40, content)
}

func warnStyle(content string) string {
	return showStyle(ShowStyleHighLight, 34, 40, content)
}

func errorStyle(content string) string {
	return showStyle(ShowStyleHighLight, 31, 40, content)
}

func panicStyle(content string) string {
	return showStyle(ShowStyleHighLight, 35, 40, content)
}

func fatalStyle(content string) string {
	return showStyle(ShowStyleHighLight, 37, 40, content)
}
