package headless

import (
	"image/color"
	"testing"
)

func TestNewRenderer(t *testing.T) {
	r, err := NewRenderer(80, 24)
	if err != nil {
		t.Fatalf("failed to create renderer: %v", err)
	}
	if r.width != 80 {
		t.Errorf("expected width 80, got %d", r.width)
	}
	if r.height != 24 {
		t.Errorf("expected height 24, got %d", r.height)
	}
}

func TestParseANSISimple(t *testing.T) {
	r, _ := NewRenderer(10, 2)
	grid := r.parseANSI("Hello")

	if len(grid) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(grid))
	}
	if len(grid[0]) != 10 {
		t.Fatalf("expected 10 cols, got %d", len(grid[0]))
	}

	expected := []rune{'H', 'e', 'l', 'l', 'o'}
	for i, r := range expected {
		if grid[0][i].Rune != r {
			t.Errorf("position %d: expected %c, got %c", i, r, grid[0][i].Rune)
		}
	}
}

func TestParseANSINewline(t *testing.T) {
	r, _ := NewRenderer(10, 3)
	grid := r.parseANSI("A\nB\nC")

	if grid[0][0].Rune != 'A' {
		t.Errorf("row 0: expected 'A', got %c", grid[0][0].Rune)
	}
	if grid[1][0].Rune != 'B' {
		t.Errorf("row 1: expected 'B', got %c", grid[1][0].Rune)
	}
	if grid[2][0].Rune != 'C' {
		t.Errorf("row 2: expected 'C', got %c", grid[2][0].Rune)
	}
}

func TestParseANSIColors(t *testing.T) {
	r, _ := NewRenderer(10, 1)
	// Red foreground: \x1b[31m
	grid := r.parseANSI("\x1b[31mX")

	if grid[0][0].Rune != 'X' {
		t.Errorf("expected 'X', got %c", grid[0][0].Rune)
	}
	if grid[0][0].FG != colorPalette[1] {
		t.Errorf("expected red, got %v", grid[0][0].FG)
	}
}

func TestParseANSIBold(t *testing.T) {
	r, _ := NewRenderer(10, 1)
	// Bold: \x1b[1m
	grid := r.parseANSI("\x1b[1mB")

	if !grid[0][0].Bold {
		t.Error("expected bold to be true")
	}
}

func TestParseANSIReset(t *testing.T) {
	r, _ := NewRenderer(10, 1)
	// Red then reset: \x1b[31mR\x1b[0mN
	grid := r.parseANSI("\x1b[31mR\x1b[0mN")

	if grid[0][0].FG != colorPalette[1] {
		t.Errorf("expected red for first char, got %v", grid[0][0].FG)
	}
	if grid[0][1].FG != defaultFG {
		t.Errorf("expected default for second char, got %v", grid[0][1].FG)
	}
}

func TestRender(t *testing.T) {
	r, err := NewRenderer(10, 2)
	if err != nil {
		t.Fatalf("failed to create renderer: %v", err)
	}

	png, err := r.Render("Hello")
	if err != nil {
		t.Fatalf("failed to render: %v", err)
	}

	// Check PNG magic bytes
	if len(png) < 8 {
		t.Fatal("PNG too short")
	}
	pngMagic := []byte{0x89, 'P', 'N', 'G', 0x0D, 0x0A, 0x1A, 0x0A}
	for i, b := range pngMagic {
		if png[i] != b {
			t.Errorf("invalid PNG magic at byte %d: expected %x, got %x", i, b, png[i])
		}
	}
}

func TestColorPalette(t *testing.T) {
	// Verify palette has 16 colors
	if len(colorPalette) != 16 {
		t.Errorf("expected 16 colors, got %d", len(colorPalette))
	}

	// Verify some known colors from spec
	if colorPalette[0] != (color.RGBA{0x1D, 0x1F, 0x21, 0xFF}) {
		t.Errorf("color 0 (black) mismatch: %v", colorPalette[0])
	}
	if colorPalette[2] != (color.RGBA{0xB5, 0xBD, 0x68, 0xFF}) {
		t.Errorf("color 2 (green) mismatch: %v", colorPalette[2])
	}
}

func TestServerAddRenderFIFO(t *testing.T) {
	s := NewServer(3) // max 3 renders

	s.AddRender("t1", []byte("data1"))
	s.AddRender("t2", []byte("data2"))
	s.AddRender("t3", []byte("data3"))

	if len(s.renders) != 3 {
		t.Errorf("expected 3 renders, got %d", len(s.renders))
	}

	// Add 4th - should evict t1
	s.AddRender("t4", []byte("data4"))

	if len(s.renders) != 3 {
		t.Errorf("expected 3 renders after eviction, got %d", len(s.renders))
	}

	if _, ok := s.renders["t1"]; ok {
		t.Error("t1 should have been evicted")
	}

	if s.latest != "t4" {
		t.Errorf("expected latest to be t4, got %s", s.latest)
	}
}
