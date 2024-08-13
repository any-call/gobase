package mycolor

import (
	"image/color"
)

func UInt32ToRGBA(rgba uint32) color.Color {
	r := uint8((rgba >> 24) & 0xFF) // 提取红色通道
	g := uint8((rgba >> 16) & 0xFF) // 提取绿色通道
	b := uint8((rgba >> 8) & 0xFF)  // 提取蓝色通道
	a := uint8(rgba & 0xFF)         // 提取透明度通道

	return color.RGBA{R: r, G: g, B: b, A: a}
}

func RGBA(r, g, b, a uint8) color.Color {
	return color.RGBA{
		R: r,
		G: g,
		B: b,
		A: a,
	}
}
