package cmd

import (
	_ "embed"
	"image"
	"image/color"
	"image/png"
	"os"
	"strconv"
	"strings"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

//go:embed fonts/FiraCode-Regular.ttf
var firaCodeRegularTTF []byte

//go:embed fonts/FiraCode-Bold.ttf
var firaCodeBoldTTF []byte

// TerminalRenderer renders ANSI terminal output to PNG images
type TerminalRenderer struct {
	regularFace font.Face
	boldFace    font.Face
	cellWidth   int
	cellHeight  int
}

// ANSI color palette (standard 16 colors)
var ansiColors = []color.RGBA{
	{0, 0, 0, 255},       // 0: Black
	{205, 49, 49, 255},   // 1: Red
	{13, 188, 121, 255},  // 2: Green
	{229, 229, 16, 255},  // 3: Yellow
	{36, 114, 200, 255},  // 4: Blue
	{188, 63, 188, 255},  // 5: Magenta
	{17, 168, 205, 255},  // 6: Cyan
	{229, 229, 229, 255}, // 7: White (light gray)
	{102, 102, 102, 255}, // 8: Bright Black (dark gray)
	{241, 76, 76, 255},   // 9: Bright Red
	{35, 209, 139, 255},  // 10: Bright Green
	{245, 245, 67, 255},  // 11: Bright Yellow
	{59, 142, 234, 255},  // 12: Bright Blue
	{214, 112, 214, 255}, // 13: Bright Magenta
	{41, 184, 219, 255},  // 14: Bright Cyan
	{255, 255, 255, 255}, // 15: Bright White
}

// Default colors
var (
	defaultFg = color.RGBA{229, 229, 229, 255} // Light gray
	defaultBg = color.RGBA{30, 30, 30, 255}    // Dark gray (VSCode-like)
)

// NewTerminalRenderer creates a new terminal renderer
func NewTerminalRenderer() (*TerminalRenderer, error) {
	// Parse the embedded fonts
	regularFont, err := opentype.Parse(firaCodeRegularTTF)
	if err != nil {
		return nil, err
	}

	boldFont, err := opentype.Parse(firaCodeBoldTTF)
	if err != nil {
		return nil, err
	}

	// Create font faces at 14pt size
	const fontSize = 14
	const dpi = 72
	regularFace, err := opentype.NewFace(regularFont, &opentype.FaceOptions{
		Size:    fontSize,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return nil, err
	}

	boldFace, err := opentype.NewFace(boldFont, &opentype.FaceOptions{
		Size:    fontSize,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return nil, err
	}

	// Calculate cell dimensions based on font metrics
	metrics := regularFace.Metrics()
	cellHeight := (metrics.Height).Ceil()
	// Measure the width of a character (monospace so all same)
	adv, _ := regularFace.GlyphAdvance('M')
	cellWidth := adv.Ceil()

	return &TerminalRenderer{
		regularFace: regularFace,
		boldFace:    boldFace,
		cellWidth:   cellWidth,
		cellHeight:  cellHeight,
	}, nil
}

// cellStyle represents the style of a single cell
type cellStyle struct {
	fg     color.RGBA
	bg     color.RGBA
	bold   bool
	italic bool
}

// cell represents a single character cell with its style
type cell struct {
	char  rune
	style cellStyle
}

// RenderToFile renders ANSI text to a PNG file
func (r *TerminalRenderer) RenderToFile(text string, width, height int, filename string) error {
	// Parse ANSI and build cell grid
	grid := r.parseANSI(text, width, height)

	// Create image
	imgWidth := width * r.cellWidth
	imgHeight := height * r.cellHeight
	img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))

	// Fill background
	for y := 0; y < imgHeight; y++ {
		for x := 0; x < imgWidth; x++ {
			img.Set(x, y, defaultBg)
		}
	}

	// Render each cell
	for row := 0; row < height; row++ {
		for col := 0; col < width; col++ {
			c := grid[row][col]

			// Draw background if not default
			if c.style.bg != defaultBg {
				for dy := 0; dy < r.cellHeight; dy++ {
					for dx := 0; dx < r.cellWidth; dx++ {
						img.Set(col*r.cellWidth+dx, row*r.cellHeight+dy, c.style.bg)
					}
				}
			}

			// Draw character
			if c.char != 0 && c.char != ' ' {
				r.drawChar(img, col, row, c.char, c.style.fg, c.style.bold)
			}
		}
	}

	// Save to file
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, img)
}

// parseANSI parses ANSI text and returns a grid of cells
func (r *TerminalRenderer) parseANSI(text string, width, height int) [][]cell {
	// Initialize grid with default style
	grid := make([][]cell, height)
	for row := range grid {
		grid[row] = make([]cell, width)
		for col := range grid[row] {
			grid[row][col] = cell{
				char:  ' ',
				style: cellStyle{fg: defaultFg, bg: defaultBg},
			}
		}
	}

	// Current position and style
	row, col := 0, 0
	currentStyle := cellStyle{fg: defaultFg, bg: defaultBg}

	i := 0
	runes := []rune(text)
	for i < len(runes) {
		ch := runes[i]

		// Check for escape sequence
		if ch == '\x1b' && i+1 < len(runes) && runes[i+1] == '[' {
			// Parse CSI sequence
			j := i + 2
			for j < len(runes) && (runes[j] == ';' || runes[j] == ':' || (runes[j] >= '0' && runes[j] <= '9')) {
				j++
			}
			if j < len(runes) {
				finalByte := runes[j]
				params := string(runes[i+2 : j])
				if finalByte == 'm' {
					r.parseSGR(params, &currentStyle)
				}
				i = j + 1
				continue
			}
		}

		// Handle newline
		if ch == '\n' {
			row++
			col = 0
			i++
			if row >= height {
				return grid
			}
			continue
		}

		// Handle carriage return
		if ch == '\r' {
			col = 0
			i++
			continue
		}

		// Handle tab
		if ch == '\t' {
			// Move to next tab stop (every 8 columns)
			col = ((col / 8) + 1) * 8
			if col >= width {
				col = width - 1
			}
			i++
			continue
		}

		// Regular character
		if col < width && row < height {
			grid[row][col] = cell{char: ch, style: currentStyle}
			col++
		}
		i++
	}

	return grid
}

// parseSGR parses SGR (Select Graphic Rendition) parameters
func (r *TerminalRenderer) parseSGR(params string, style *cellStyle) {
	if params == "" {
		// Reset
		style.fg = defaultFg
		style.bg = defaultBg
		style.bold = false
		style.italic = false
		return
	}

	// Handle colon-separated params (convert to semicolon for uniform handling)
	params = strings.ReplaceAll(params, ":", ";")

	// Split by semicolon
	parts := strings.Split(params, ";")
	i := 0
	for i < len(parts) {
		num, err := strconv.Atoi(parts[i])
		if err != nil {
			i++
			continue
		}

		switch {
		case num == 0:
			// Reset
			style.fg = defaultFg
			style.bg = defaultBg
			style.bold = false
			style.italic = false
		case num == 1:
			style.bold = true
		case num == 3:
			style.italic = true
		case num == 22:
			style.bold = false
		case num == 23:
			style.italic = false
		case num >= 30 && num <= 37:
			// Standard foreground colors
			colorIdx := num - 30
			if style.bold && colorIdx < 8 {
				colorIdx += 8 // Use bright variant
			}
			style.fg = ansiColors[colorIdx]
		case num == 38:
			// Extended foreground color
			if i+1 < len(parts) {
				mode, _ := strconv.Atoi(parts[i+1])
				if mode == 5 && i+2 < len(parts) {
					// 256-color mode
					colorIdx, _ := strconv.Atoi(parts[i+2])
					style.fg = r.get256Color(colorIdx)
					i += 2
				} else if mode == 2 && i+4 < len(parts) {
					// RGB mode
					red, _ := strconv.Atoi(parts[i+2])
					green, _ := strconv.Atoi(parts[i+3])
					blue, _ := strconv.Atoi(parts[i+4])
					style.fg = color.RGBA{uint8(red), uint8(green), uint8(blue), 255}
					i += 4
				}
			}
		case num == 39:
			style.fg = defaultFg
		case num >= 40 && num <= 47:
			// Standard background colors
			style.bg = ansiColors[num-40]
		case num == 48:
			// Extended background color
			if i+1 < len(parts) {
				mode, _ := strconv.Atoi(parts[i+1])
				if mode == 5 && i+2 < len(parts) {
					// 256-color mode
					colorIdx, _ := strconv.Atoi(parts[i+2])
					style.bg = r.get256Color(colorIdx)
					i += 2
				} else if mode == 2 && i+4 < len(parts) {
					// RGB mode
					red, _ := strconv.Atoi(parts[i+2])
					green, _ := strconv.Atoi(parts[i+3])
					blue, _ := strconv.Atoi(parts[i+4])
					style.bg = color.RGBA{uint8(red), uint8(green), uint8(blue), 255}
					i += 4
				}
			}
		case num == 49:
			style.bg = defaultBg
		case num >= 90 && num <= 97:
			// Bright foreground colors
			style.fg = ansiColors[num-90+8]
		case num >= 100 && num <= 107:
			// Bright background colors
			style.bg = ansiColors[num-100+8]
		}

		i++
	}
}

// get256Color returns the color for a 256-color palette index
func (r *TerminalRenderer) get256Color(idx int) color.RGBA {
	if idx < 16 {
		return ansiColors[idx]
	}

	if idx < 232 {
		// 6x6x6 color cube
		idx -= 16
		blue := idx % 6
		green := (idx / 6) % 6
		red := idx / 36

		toVal := func(v int) uint8 {
			if v == 0 {
				return 0
			}
			return uint8(55 + 40*v)
		}

		return color.RGBA{toVal(red), toVal(green), toVal(blue), 255}
	}

	// Grayscale ramp (24 shades)
	idx -= 232
	gray := uint8(8 + idx*10)
	return color.RGBA{gray, gray, gray, 255}
}

// drawChar draws a character at the given cell position
func (r *TerminalRenderer) drawChar(img *image.RGBA, col, row int, ch rune, fg color.RGBA, bold bool) {
	// Select font face based on bold
	face := r.regularFace
	if bold {
		face = r.boldFace
	}

	// Get font metrics for baseline positioning
	metrics := face.Metrics()
	ascent := metrics.Ascent.Ceil()

	point := fixed.Point26_6{
		X: fixed.I(col * r.cellWidth),
		Y: fixed.I(row*r.cellHeight + ascent),
	}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(fg),
		Face: face,
		Dot:  point,
	}

	d.DrawString(string(ch))
}
