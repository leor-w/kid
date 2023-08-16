package utils

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"
)

func HexToRGBA(hexColor string) (*color.RGBA, error) {
	hexColor = strings.TrimPrefix(hexColor, "#") //移除颜色字符串的#

	if len(hexColor) != 6 {
		return nil, fmt.Errorf("invalid length for hex color: %s", hexColor)
	}

	r, err := strconv.ParseInt(hexColor[0:2], 16, 64)
	if err != nil {
		return nil, err
	}

	g, err := strconv.ParseInt(hexColor[2:4], 16, 64)
	if err != nil {
		return nil, err
	}

	b, err := strconv.ParseInt(hexColor[4:6], 16, 64)
	if err != nil {
		return nil, err
	}

	return &color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 0xff}, nil
}
