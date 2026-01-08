package cmd

import (
	"image/color"
	"os"
	"path/filepath"
	"testing"
)

func TestTerminalRendererBasic(t *testing.T) {
	r, err := NewTerminalRenderer()
	if err != nil {
		t.Fatalf("failed to create renderer: %v", err)
	}

	// Create temp dir
	tmpDir, err := os.MkdirTemp("", "renderer-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	filename := filepath.Join(tmpDir, "test.png")

	// Simple text rendering
	err = r.RenderToFile("Hello, World!", 80, 24, filename)
	if err != nil {
		t.Fatalf("failed to render: %v", err)
	}

	// Check file exists
	info, err := os.Stat(filename)
	if err != nil {
		t.Fatalf("file not created: %v", err)
	}

	if info.Size() == 0 {
		t.Error("file is empty")
	}
}

func TestTerminalRendererANSI(t *testing.T) {
	r, err := NewTerminalRenderer()
	if err != nil {
		t.Fatalf("failed to create renderer: %v", err)
	}

	// Create temp dir
	tmpDir, err := os.MkdirTemp("", "renderer-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	filename := filepath.Join(tmpDir, "test.png")

	// ANSI color codes
	text := "\x1b[32mGreen\x1b[0m Normal \x1b[1;31mBold Red\x1b[0m"
	err = r.RenderToFile(text, 40, 5, filename)
	if err != nil {
		t.Fatalf("failed to render ANSI: %v", err)
	}

	// Check file exists
	info, err := os.Stat(filename)
	if err != nil {
		t.Fatalf("file not created: %v", err)
	}

	if info.Size() == 0 {
		t.Error("file is empty")
	}
}

func TestSGRParsing(t *testing.T) {
	r, _ := NewTerminalRenderer()

	tests := []struct {
		name       string
		params     string
		initStyle  cellStyle
		wantStyle  cellStyle
	}{
		{
			name:       "reset",
			params:     "0",
			initStyle:  cellStyle{fg: ansiColors[1], bg: ansiColors[2], bold: true},
			wantStyle:  cellStyle{fg: defaultFg, bg: defaultBg},
		},
		{
			name:       "empty reset",
			params:     "",
			initStyle:  cellStyle{fg: ansiColors[1], bg: ansiColors[2], bold: true},
			wantStyle:  cellStyle{fg: defaultFg, bg: defaultBg},
		},
		{
			name:       "bold",
			params:     "1",
			initStyle:  cellStyle{fg: defaultFg, bg: defaultBg},
			wantStyle:  cellStyle{fg: defaultFg, bg: defaultBg, bold: true},
		},
		{
			name:       "italic",
			params:     "3",
			initStyle:  cellStyle{fg: defaultFg, bg: defaultBg},
			wantStyle:  cellStyle{fg: defaultFg, bg: defaultBg, italic: true},
		},
		{
			name:       "fg red",
			params:     "31",
			initStyle:  cellStyle{fg: defaultFg, bg: defaultBg},
			wantStyle:  cellStyle{fg: ansiColors[1], bg: defaultBg},
		},
		{
			name:       "fg green",
			params:     "32",
			initStyle:  cellStyle{fg: defaultFg, bg: defaultBg},
			wantStyle:  cellStyle{fg: ansiColors[2], bg: defaultBg},
		},
		{
			name:       "bg blue",
			params:     "44",
			initStyle:  cellStyle{fg: defaultFg, bg: defaultBg},
			wantStyle:  cellStyle{fg: defaultFg, bg: ansiColors[4]},
		},
		{
			name:       "bright fg",
			params:     "92",
			initStyle:  cellStyle{fg: defaultFg, bg: defaultBg},
			wantStyle:  cellStyle{fg: ansiColors[10], bg: defaultBg},
		},
		{
			name:       "256 color fg",
			params:     "38;5;196",
			initStyle:  cellStyle{fg: defaultFg, bg: defaultBg},
			wantStyle:  cellStyle{fg: r.get256Color(196), bg: defaultBg},
		},
		{
			name:       "256 color bg",
			params:     "48;5;21",
			initStyle:  cellStyle{fg: defaultFg, bg: defaultBg},
			wantStyle:  cellStyle{fg: defaultFg, bg: r.get256Color(21)},
		},
		{
			name:       "rgb fg",
			params:     "38;2;255;128;64",
			initStyle:  cellStyle{fg: defaultFg, bg: defaultBg},
			wantStyle:  cellStyle{fg: color.RGBA{255, 128, 64, 255}, bg: defaultBg},
		},
		{
			name:       "multiple params",
			params:     "1;32;44",
			initStyle:  cellStyle{fg: defaultFg, bg: defaultBg},
			wantStyle:  cellStyle{fg: ansiColors[10], bg: ansiColors[4], bold: true}, // bold + green becomes bright green
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			style := tt.initStyle
			r.parseSGR(tt.params, &style)

			if style.fg != tt.wantStyle.fg {
				t.Errorf("fg = %v, want %v", style.fg, tt.wantStyle.fg)
			}
			if style.bg != tt.wantStyle.bg {
				t.Errorf("bg = %v, want %v", style.bg, tt.wantStyle.bg)
			}
			if style.bold != tt.wantStyle.bold {
				t.Errorf("bold = %v, want %v", style.bold, tt.wantStyle.bold)
			}
			if style.italic != tt.wantStyle.italic {
				t.Errorf("italic = %v, want %v", style.italic, tt.wantStyle.italic)
			}
		})
	}
}

func TestGet256Color(t *testing.T) {
	r, _ := NewTerminalRenderer()

	// Standard colors (0-15)
	for i := 0; i < 16; i++ {
		c := r.get256Color(i)
		if c != ansiColors[i] {
			t.Errorf("256 color %d: got %v, want %v", i, c, ansiColors[i])
		}
	}

	// Color cube (16-231)
	// Test a few known colors
	c := r.get256Color(16)  // First color cube entry (black)
	if c.R != 0 || c.G != 0 || c.B != 0 {
		t.Errorf("256 color 16: got %v, want (0,0,0)", c)
	}

	c = r.get256Color(231) // Last color cube entry (white)
	if c.R != 255 || c.G != 255 || c.B != 255 {
		t.Errorf("256 color 231: got %v, want (255,255,255)", c)
	}

	// Grayscale (232-255)
	c = r.get256Color(232) // First grayscale (dark gray)
	if c.R != 8 || c.G != 8 || c.B != 8 {
		t.Errorf("256 color 232: got %v, want (8,8,8)", c)
	}

	c = r.get256Color(255) // Last grayscale (light gray)
	if c.R != 238 || c.G != 238 || c.B != 238 {
		t.Errorf("256 color 255: got %v, want (238,238,238)", c)
	}
}
