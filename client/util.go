package client

import (
	"strconv"

	"github.com/notassigned/tweb/xmlnode"
	tview "github.com/rivo/tview"
)

func intOrZero(s string) int {
	x, err := strconv.ParseInt(s, 10, 32)
	if err == nil {
		return int(x)
	} else {
		return 0
	}
}

func SizeProportionFocus(n *xmlnode.XmlNode) (size int, proportion int, focus bool) {
	attrs := n.Attributes
	size = 0
	proportion = 1
	focus = false

	if p, ok := attrs["size"]; ok {
		size, _ = strconv.Atoi(p)
	}

	if p, ok := attrs["proportion"]; ok {
		proportion, _ = strconv.Atoi(p)
	}

	if p, ok := attrs["focus"]; ok && p == "true" {
		focus = true
	}

	return size, proportion, focus
}

func alignLookup(s string) int {
	switch s {
	case "center":
		return tview.AlignCenter
	case "left":
		return tview.AlignLeft
	case "right":
		return tview.AlignRight
	default:
		return tview.AlignLeft
	}
}
