package layouts

import (
	"log"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
)

type Size struct {
	Width  int
	Height int
}

// Declare conformity with Layout interface
var _ fyne.Layout = (*boxLayout)(nil)

var objConfigMap map[fyne.CanvasObject]*Size
var objConfigLock sync.RWMutex

func SetObjConfigMap(o fyne.CanvasObject, s *Size) {
	objConfigLock.Lock()
	defer objConfigLock.Unlock()

	if objConfigMap == nil {
		objConfigMap = make(map[fyne.CanvasObject]*Size)
	}

	objConfigMap[o] = s
}

func GetObjConfig(o fyne.CanvasObject) *Size {
	objConfigLock.RLock()
	defer objConfigLock.RUnlock()
	if s, ok := objConfigMap[o]; ok {
		return s
	}

	return nil
}

type boxLayout struct {
	horizontal bool
}

// NewHBoxLayout returns a horizontal box layout for stacking a number of child
// canvas objects or widgets left to right.
func NewHBoxLayout() fyne.Layout {
	return &boxLayout{true}
}

// NewVBoxLayout returns a vertical box layout for stacking a number of child
// canvas objects or widgets top to bottom.
func NewVBoxLayout() fyne.Layout {
	return &boxLayout{false}
}

func isVerticalSpacer(obj fyne.CanvasObject) bool {
	if spacer, ok := obj.(layout.SpacerObject); ok {
		return spacer.ExpandVertical()
	}

	return false
}

func isHorizontalSpacer(obj fyne.CanvasObject) bool {
	if spacer, ok := obj.(layout.SpacerObject); ok {
		return spacer.ExpandHorizontal()
	}

	return false
}

func (g *boxLayout) isSpacer(obj fyne.CanvasObject) bool {
	// invisible spacers don't impact layout
	if !obj.Visible() {
		return false
	}

	if g.horizontal {
		return isHorizontalSpacer(obj)
	}
	return isVerticalSpacer(obj)
}

// Layout is called to pack all child objects into a specified size.
// For a VBoxLayout this will pack objects into a single column where each item
// is full width but the height is the minimum required.
// Any spacers added will pad the view, sharing the space if there are two or more.
func (g *boxLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	spacers := make([]fyne.CanvasObject, 0)
	total := float32(0)
	for _, child := range objects {
		if !child.Visible() {
			continue
		}

		if g.isSpacer(child) {
			spacers = append(spacers, child)
			continue
		}
		if g.horizontal {
			total += child.MinSize().Width
		} else {
			total += child.MinSize().Height
		}
	}

	x, y := float32(0), float32(0)
	var extra float32
	if g.horizontal {
		extra = size.Width - total - (theme.Padding() * float32(len(objects)-len(spacers)-1))
	} else {
		extra = size.Height - total - (theme.Padding() * float32(len(objects)-len(spacers)-1))
	}
	extraCell := float32(0)
	if len(spacers) > 0 {
		extraCell = extra / float32(len(spacers))
	}

	for index, child := range objects {
		if !child.Visible() {
			continue
		}

		var width, height float32

		width = child.MinSize().Width
		height = child.MinSize().Height

		var eWidth, eHeight float32
		config := GetObjConfig(child)
		if config != nil {
			eWidth = float32(config.Width)
			eHeight = float32(config.Height)
		} else {
			eWidth = width
			eHeight = height
		}

		if g.isSpacer(child) {
			if g.horizontal {
				x += extraCell
			} else {
				y += extraCell
			}
			continue
		}
		child.Move(fyne.NewPos(x, y))

		if g.horizontal {
			x += theme.Padding() + width
			if index == len(objects)-1 {
				child.Resize(fyne.NewSize(eWidth, size.Height))
			} else {
				child.Resize(fyne.NewSize(width, size.Height))
			}
		} else {
			y += theme.Padding() + height
			//// 如果是最后一个控件，则占满剩余控件
			if index == len(objects)-1 {

				//child.Resize(fyne.NewSize(size.Width, size.Height-y+height))
				log.Println("size.Height", size.Height, "y", y, "height", height)
				child.Resize(fyne.NewSize(size.Width, eHeight))
			} else {
				child.Resize(fyne.NewSize(size.Width, height))
			}
		}
	}
}

// MinSize finds the smallest size that satisfies all the child objects.
// For a BoxLayout this is the width of the widest item and the height is
// the sum of of all children combined with padding between each.
func (g *boxLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	minSize := fyne.NewSize(0, 0)
	addPadding := false
	for _, child := range objects {
		if !child.Visible() {
			continue
		}

		if g.isSpacer(child) {
			continue
		}

		if g.horizontal {
			minSize.Height = fyne.Max(child.MinSize().Height, minSize.Height)
			minSize.Width += child.Size().Width
			if addPadding {
				minSize.Width += theme.Padding()
			}
		} else {
			minSize.Width = fyne.Max(child.MinSize().Width, minSize.Width)
			//minSize.Height = fyne.Max(child.MinSize().Height, minSize.Height)
			minSize.Height += child.Size().Height
			if addPadding {
				minSize.Height += theme.Padding()
			}
		}
		addPadding = true
	}
	return minSize
}
