package headless

import (
	"bytes"
	"embed"
	"image"
	"image/color"
	"image/png"
	"sync"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

//go:embed fonts/FiraCodeNerdFontMono-Regular.ttf fonts/FiraCodeNerdFontMono-Bold.ttf
var fontFS embed.FS

var (
	regularFont *truetype.Font
	boldFont    *truetype.Font
	fontOnce    sync.Once
	fontErr     error
)

// Color palette from spec
var colorPalette = map[int]color.RGBA{
	0:  {0x1D, 0x1F, 0x21, 0xFF}, // black
	1:  {0xCC, 0x66, 0x66, 0xFF}, // red
	2:  {0xB5, 0xBD, 0x68, 0xFF}, // green
	3:  {0xF0, 0xC6, 0x74, 0xFF}, // yellow
	4:  {0x81, 0xA2, 0xBE, 0xFF}, // blue
	5:  {0xB2, 0x94, 0xBB, 0xFF}, // magenta
	6:  {0x8A, 0xBE, 0xB7, 0xFF}, // cyan
	7:  {0xC5, 0xC8, 0xC6, 0xFF}, // white
	8:  {0x66, 0x66, 0x66, 0xFF}, // bright black
	9:  {0xD5, 0x4E, 0x53, 0xFF}, // bright red
	10: {0xB9, 0xCA, 0x4A, 0xFF}, // bright green
	11: {0xE7, 0xC5, 0x47, 0xFF}, // bright yellow
	12: {0x7A, 0xA6, 0xDA, 0xFF}, // bright blue
	13: {0xC3, 0x97, 0xD8, 0xFF}, // bright magenta
	14: {0x70, 0xC0, 0xB1, 0xFF}, // bright cyan
	15: {0xEA, 0xEA, 0xEA, 0xFF}, // bright white
}

var defaultBG = color.RGBA{0x00, 0x00, 0x00, 0xFF} // pure black
var defaultFG = colorPalette[7]                    // white

func loadFonts() error {
	fontOnce.Do(func() {
		regularData, err := fontFS.ReadFile("fonts/FiraCodeNerdFontMono-Regular.ttf")
		if err != nil {
			fontErr = err
			return
		}
		regularFont, err = truetype.Parse(regularData)
		if err != nil {
			fontErr = err
			return
		}

		boldData, err := fontFS.ReadFile("fonts/FiraCodeNerdFontMono-Bold.ttf")
		if err != nil {
			fontErr = err
			return
		}
		boldFont, err = truetype.Parse(boldData)
		if err != nil {
			fontErr = err
			return
		}
	})
	return fontErr
}

// Cell represents a single terminal cell
type Cell struct {
	Rune rune
	FG   color.RGBA
	BG   color.RGBA
	Bold bool
}

// Renderer renders ANSI text to PNG
type Renderer struct {
	width    int // columns
	height   int // rows
	fontSize float64
	cellW    int
	cellH    int
}

// NewRenderer creates a new renderer
func NewRenderer(cols, rows int) (*Renderer, error) {
	if err := loadFonts(); err != nil {
		return nil, err
	}

	fontSize := 14.0
	// Approximate cell dimensions for monospace font
	cellW := int(fontSize * 0.6)
	cellH := int(fontSize * 1.4)

	return &Renderer{
		width:    cols,
		height:   rows,
		fontSize: fontSize,
		cellW:    cellW,
		cellH:    cellH,
	}, nil
}

// parseANSI parses ANSI escape sequences and returns a grid of cells
func (rend *Renderer) parseANSI(ansi string) [][]Cell {
	// Initialize grid
	grid := make([][]Cell, rend.height)
	for y := range grid {
		grid[y] = make([]Cell, rend.width)
		for x := range grid[y] {
			grid[y][x] = Cell{Rune: ' ', FG: defaultFG, BG: defaultBG}
		}
	}

	curX, curY := 0, 0
	curFG := defaultFG
	curBG := defaultBG
	curBold := false

	i := 0
	for i < len(ansi) {
		if ansi[i] == '\x1b' && i+1 < len(ansi) && ansi[i+1] == '[' {
			// Parse escape sequence
			j := i + 2
			for j < len(ansi) && ansi[j] != 'm' && ansi[j] != 'H' && ansi[j] != 'J' && ansi[j] != 'K' && ansi[j] != 'A' && ansi[j] != 'B' && ansi[j] != 'C' && ansi[j] != 'D' {
				j++
			}
			if j < len(ansi) {
				seq := ansi[i+2 : j]
				cmd := ansi[j]

				if cmd == 'm' {
					// SGR - Select Graphic Rendition
					curFG, curBG, curBold = rend.parseSGR(seq, curFG, curBG, curBold)
				}
				// Skip other commands for now
				i = j + 1
				continue
			}
		}

		ch := rune(ansi[i])
		if ch == '\n' {
			curX = 0
			curY++
		} else if ch == '\r' {
			curX = 0
		} else if ch >= 0x80 {
			// UTF-8 multi-byte character
			utfRune, size := decodeUTF8(ansi[i:])
			if curY < rend.height && curX < rend.width {
				grid[curY][curX] = Cell{Rune: utfRune, FG: curFG, BG: curBG, Bold: curBold}
				curX++
			}
			i += size
			continue
		} else if ch >= 32 {
			if curY < rend.height && curX < rend.width {
				grid[curY][curX] = Cell{Rune: ch, FG: curFG, BG: curBG, Bold: curBold}
				curX++
			}
		}
		i++
	}

	return grid
}

func decodeUTF8(s string) (rune, int) {
	if len(s) == 0 {
		return 0, 0
	}
	b := s[0]
	if b < 0x80 {
		return rune(b), 1
	}
	if b < 0xC0 {
		return '?', 1
	}
	if b < 0xE0 && len(s) >= 2 {
		return rune(b&0x1F)<<6 | rune(s[1]&0x3F), 2
	}
	if b < 0xF0 && len(s) >= 3 {
		return rune(b&0x0F)<<12 | rune(s[1]&0x3F)<<6 | rune(s[2]&0x3F), 3
	}
	if len(s) >= 4 {
		return rune(b&0x07)<<18 | rune(s[1]&0x3F)<<12 | rune(s[2]&0x3F)<<6 | rune(s[3]&0x3F), 4
	}
	return '?', 1
}

func (r *Renderer) parseSGR(seq string, fg, bg color.RGBA, bold bool) (color.RGBA, color.RGBA, bool) {
	if seq == "" || seq == "0" {
		return defaultFG, defaultBG, false
	}

	params := parseSGRParams(seq)
	for i := 0; i < len(params); i++ {
		p := params[i]
		switch {
		case p == 0:
			fg, bg, bold = defaultFG, defaultBG, false
		case p == 1:
			bold = true
		case p == 22:
			bold = false
		case p >= 30 && p <= 37:
			fg = colorPalette[p-30]
		case p == 39:
			fg = defaultFG
		case p >= 40 && p <= 47:
			bg = colorPalette[p-40]
		case p == 49:
			bg = defaultBG
		case p >= 90 && p <= 97:
			fg = colorPalette[p-90+8]
		case p >= 100 && p <= 107:
			bg = colorPalette[p-100+8]
		case p == 38:
			// 256-color or RGB foreground
			if i+1 < len(params) && params[i+1] == 5 && i+2 < len(params) {
				// 256-color mode
				colorIdx := params[i+2]
				if colorIdx < 16 {
					fg = colorPalette[colorIdx]
				}
				i += 2
			}
		case p == 48:
			// 256-color or RGB background
			if i+1 < len(params) && params[i+1] == 5 && i+2 < len(params) {
				colorIdx := params[i+2]
				if colorIdx < 16 {
					bg = colorPalette[colorIdx]
				}
				i += 2
			}
		}
	}
	return fg, bg, bold
}

func parseSGRParams(seq string) []int {
	var params []int
	var num int
	hasNum := false
	for _, c := range seq {
		if c >= '0' && c <= '9' {
			num = num*10 + int(c-'0')
			hasNum = true
		} else if c == ';' {
			if hasNum {
				params = append(params, num)
			}
			num = 0
			hasNum = false
		}
	}
	if hasNum {
		params = append(params, num)
	}
	return params
}

// Render renders ANSI text to PNG
func (r *Renderer) Render(ansi string) ([]byte, error) {
	grid := r.parseANSI(ansi)

	imgW := r.width * r.cellW
	imgH := r.height * r.cellH
	img := image.NewRGBA(image.Rect(0, 0, imgW, imgH))

	// Fill background
	for y := 0; y < imgH; y++ {
		for x := 0; x < imgW; x++ {
			img.Set(x, y, defaultBG)
		}
	}

	// Draw cells
	for row, cells := range grid {
		for col, cell := range cells {
			x := col * r.cellW
			y := row * r.cellH

			// Draw background
			for dy := 0; dy < r.cellH; dy++ {
				for dx := 0; dx < r.cellW; dx++ {
					img.Set(x+dx, y+dy, cell.BG)
				}
			}

			// Draw character
			if cell.Rune != ' ' {
				r.drawChar(img, x, y, cell)
			}
		}
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (r *Renderer) drawChar(img *image.RGBA, x, y int, cell Cell) {
	f := regularFont
	if cell.Bold {
		f = boldFont
	}

	face := truetype.NewFace(f, &truetype.Options{
		Size:    r.fontSize,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	defer face.Close()

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(cell.FG),
		Face: face,
	}

	// Position: baseline is at bottom of cell
	d.Dot = fixed.Point26_6{
		X: fixed.I(x),
		Y: fixed.I(y + r.cellH - 3),
	}
	d.DrawString(string(cell.Rune))
}
