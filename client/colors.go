package client

import (
	"strconv"

	"github.com/gdamore/tcell/v2"
)

func colorLookup(c string) (tcell.Color, error) {
	if color, ok := tcell.ColorNames[c]; ok {
		return color, nil
	}
	//If not a color name then try to parse hex
	value, err := strconv.ParseUint(c, 16, 64)
	return tcell.Color(value), err
}
