package tui

import (
	"time"
)

// Types of user action
const (
	Rune = iota

	CtrlA
	CtrlB
	CtrlC
	CtrlD
	CtrlE
	CtrlF
	CtrlG
	CtrlH
	Tab
	CtrlJ
	CtrlK
	CtrlL
	CtrlM
	CtrlN
	CtrlO
	CtrlP
	CtrlQ
	CtrlR
	CtrlS
	CtrlT
	CtrlU
	CtrlV
	CtrlW
	CtrlX
	CtrlY
	CtrlZ
	ESC

	Invalid
	Resize
	Mouse
	DoubleClick

	BTab
	BSpace

	Del
	PgUp
	PgDn

	Up
	Down
	Left
	Right
	Home
	End

	SLeft
	SRight

	F1
	F2
	F3
	F4
	F5
	F6
	F7
	F8
	F9
	F10
	F11
	F12

	AltEnter
	AltSpace
	AltSlash
	AltBS

	Alt0
)

const ( // Reset iota
	AltA = Alt0 + 'a' - '0' + iota
	AltB
	AltC
	AltD
	AltE
	AltF
	AltZ = AltA + 'z' - 'a'
)

const (
	doubleClickDuration = 500 * time.Millisecond
)

type Color int32

func (c Color) is24() bool {
	return c > 0 && (c&(1<<24)) > 0
}

const (
	colUndefined Color = -2
	colDefault         = -1
)

const (
	colBlack Color = iota
	colRed
	colGreen
	colYellow
	colBlue
	colMagenta
	colCyan
	colWhite
)

type ColorIndex int16

const (
	ColDefault ColorIndex = iota
	ColNormal
	ColPrompt
	ColMatch
	ColCurrent
	ColCurrentMatch
	ColSpinner
	ColInfo
	ColCursor
	ColSelected
	ColHeader
	ColBorder
	ColUser // Should be the last entry
)

func (i ColorIndex) Pair() ColorPair {
	if i >= ColDefault && i < ColUser {
		return Pallete[i]
	}
	return Pallete[ColDefault]
}

type ColorPair struct {
	fg    Color
	bg    Color
	index ColorIndex
}

func NewColorPair(fg Color, bg Color) ColorPair {
	return ColorPair{fg, bg, ColUser}
}

func (p ColorPair) Fg() Color {
	return p.fg
}

func (p ColorPair) Bg() Color {
	return p.bg
}

func (p ColorPair) key() int {
	return (int(p.Fg()) << 8) + int(p.Bg())
}

func (p ColorPair) is24() bool {
	return p.Fg().is24() || p.Bg().is24()
}

type ColorTheme struct {
	Fg           Color
	Bg           Color
	DarkBg       Color
	Prompt       Color
	Match        Color
	Current      Color
	CurrentMatch Color
	Spinner      Color
	Info         Color
	Cursor       Color
	Selected     Color
	Header       Color
	Border       Color
}

type Event struct {
	Type       int
	Char       rune
	MouseEvent *MouseEvent
}

type MouseEvent struct {
	Y      int
	X      int
	S      int
	Down   bool
	Double bool
	Mod    bool
}

type Renderer interface {
	Init()
	Pause()
	Resume() bool
	Clear()
	RefreshWindows(windows []Window)
	Refresh()
	Close()

	GetChar() Event

	MaxX() int
	MaxY() int
	DoesAutoWrap() bool

	NewWindow(top int, left int, width int, height int, border bool) Window
}

type Window interface {
	Top() int
	Left() int
	Width() int
	Height() int

	Refresh()
	FinishFill()
	Close()

	X() int
	Enclose(y int, x int) bool

	Move(y int, x int)
	MoveAndClear(y int, x int)
	Print(text string)
	CPrint(color ColorIndex, attr Attr, text string)
	CPrintPair(pair ColorPair, attr Attr, text string)
	Fill(text string) bool
	CFill(fg Color, bg Color, attr Attr, text string) bool
	Erase()
}

type FullscreenRenderer struct {
	theme        *ColorTheme
	mouse        bool
	forceBlack   bool
	prevDownTime time.Time
	clickY       []int
}

func NewFullscreenRenderer(theme *ColorTheme, forceBlack bool, mouse bool) Renderer {
	r := &FullscreenRenderer{
		theme:        theme,
		mouse:        mouse,
		forceBlack:   forceBlack,
		prevDownTime: time.Unix(0, 0),
		clickY:       []int{}}
	return r
}

var (
	Pallete   [ColUser]ColorPair
	Default16 *ColorTheme
	Dark256   *ColorTheme
	Light256  *ColorTheme
)

func EmptyTheme() *ColorTheme {
	return &ColorTheme{
		Fg:           colUndefined,
		Bg:           colUndefined,
		DarkBg:       colUndefined,
		Prompt:       colUndefined,
		Match:        colUndefined,
		Current:      colUndefined,
		CurrentMatch: colUndefined,
		Spinner:      colUndefined,
		Info:         colUndefined,
		Cursor:       colUndefined,
		Selected:     colUndefined,
		Header:       colUndefined,
		Border:       colUndefined}
}

func init() {
	Default16 = &ColorTheme{
		Fg:           colDefault,
		Bg:           colDefault,
		DarkBg:       colBlack,
		Prompt:       colBlue,
		Match:        colGreen,
		Current:      colYellow,
		CurrentMatch: colGreen,
		Spinner:      colGreen,
		Info:         colWhite,
		Cursor:       colRed,
		Selected:     colMagenta,
		Header:       colCyan,
		Border:       colBlack}
	Dark256 = &ColorTheme{
		Fg:           colDefault,
		Bg:           colDefault,
		DarkBg:       236,
		Prompt:       110,
		Match:        108,
		Current:      254,
		CurrentMatch: 151,
		Spinner:      148,
		Info:         144,
		Cursor:       161,
		Selected:     168,
		Header:       109,
		Border:       59}
	Light256 = &ColorTheme{
		Fg:           colDefault,
		Bg:           colDefault,
		DarkBg:       251,
		Prompt:       25,
		Match:        66,
		Current:      237,
		CurrentMatch: 23,
		Spinner:      65,
		Info:         101,
		Cursor:       161,
		Selected:     168,
		Header:       31,
		Border:       145}
}

func initTheme(theme *ColorTheme, baseTheme *ColorTheme, forceBlack bool) {
	if theme == nil {
		for idx := ColDefault; idx < ColUser; idx++ {
			Pallete[idx] = ColorPair{colDefault, colDefault, idx}
		}
		return
	}

	if forceBlack {
		theme.Bg = colBlack
	}

	o := func(a Color, b Color) Color {
		if b == colUndefined {
			return a
		}
		return b
	}
	theme.Fg = o(baseTheme.Fg, theme.Fg)
	theme.Bg = o(baseTheme.Bg, theme.Bg)
	theme.DarkBg = o(baseTheme.DarkBg, theme.DarkBg)
	theme.Prompt = o(baseTheme.Prompt, theme.Prompt)
	theme.Match = o(baseTheme.Match, theme.Match)
	theme.Current = o(baseTheme.Current, theme.Current)
	theme.CurrentMatch = o(baseTheme.CurrentMatch, theme.CurrentMatch)
	theme.Spinner = o(baseTheme.Spinner, theme.Spinner)
	theme.Info = o(baseTheme.Info, theme.Info)
	theme.Cursor = o(baseTheme.Cursor, theme.Cursor)
	theme.Selected = o(baseTheme.Selected, theme.Selected)
	theme.Header = o(baseTheme.Header, theme.Header)
	theme.Border = o(baseTheme.Border, theme.Border)

	Pallete[ColDefault] = ColorPair{colDefault, colDefault, ColDefault}
	Pallete[ColNormal] = ColorPair{theme.Fg, theme.Bg, ColNormal}
	Pallete[ColPrompt] = ColorPair{theme.Prompt, theme.Bg, ColPrompt}
	Pallete[ColMatch] = ColorPair{theme.Match, theme.Bg, ColMatch}
	Pallete[ColCurrent] = ColorPair{theme.Current, theme.DarkBg, ColCurrent}
	Pallete[ColCurrentMatch] = ColorPair{theme.CurrentMatch, theme.DarkBg, ColCurrentMatch}
	Pallete[ColSpinner] = ColorPair{theme.Spinner, theme.Bg, ColSpinner}
	Pallete[ColInfo] = ColorPair{theme.Info, theme.Bg, ColInfo}
	Pallete[ColCursor] = ColorPair{theme.Cursor, theme.DarkBg, ColCursor}
	Pallete[ColSelected] = ColorPair{theme.Selected, theme.DarkBg, ColSelected}
	Pallete[ColHeader] = ColorPair{theme.Header, theme.Bg, ColHeader}
	Pallete[ColBorder] = ColorPair{theme.Border, theme.Bg, ColBorder}
}

func attrFor(color ColorIndex, attr Attr) Attr {
	switch color {
	case ColCurrent:
		return attr | Reverse
	case ColMatch:
		return attr | Underline
	case ColCurrentMatch:
		return attr | Underline | Reverse
	}
	return attr
}
