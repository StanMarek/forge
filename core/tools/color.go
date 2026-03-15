package tools

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

// ColorTool provides metadata for the Color Converter tool.
type ColorTool struct{}

func (c ColorTool) Name() string        { return "Color Converter" }
func (c ColorTool) ID() string          { return "color" }
func (c ColorTool) Description() string { return "Convert colors between hex, RGB, and HSL formats" }
func (c ColorTool) Category() string    { return "Converters" }
func (c ColorTool) Keywords() []string {
	return []string{"color", "hex", "rgb", "hsl", "convert"}
}

// DetectFromClipboard returns true if s looks like a color value.
func (c ColorTool) DetectFromClipboard(s string) bool {
	s = strings.TrimSpace(s)
	if hexPattern.MatchString(s) || shortHexPattern.MatchString(s) {
		return true
	}
	if rgbPattern.MatchString(s) {
		return true
	}
	if hslPattern.MatchString(s) {
		return true
	}
	return false
}

var (
	hexPattern      = regexp.MustCompile(`(?i)^#[0-9a-f]{6}$`)
	shortHexPattern = regexp.MustCompile(`(?i)^#[0-9a-f]{3}$`)
	rgbPattern      = regexp.MustCompile(`(?i)^rgb\(\s*(\d+)\s*,\s*(\d+)\s*,\s*(\d+)\s*\)$`)
	hslPattern      = regexp.MustCompile(`(?i)^hsl\(\s*(\d+)\s*,\s*(\d+)%\s*,\s*(\d+)%\s*\)$`)
)

// ColorConvert auto-detects the color format and converts to all formats.
func ColorConvert(input string) Result {
	s := strings.TrimSpace(input)

	var r, g, b uint8

	switch {
	case hexPattern.MatchString(s):
		var err error
		r, g, b, err = parseHex6(s)
		if err != nil {
			return Result{Error: fmt.Sprintf("invalid hex color: %s", err)}
		}

	case shortHexPattern.MatchString(s):
		var err error
		r, g, b, err = parseHex3(s)
		if err != nil {
			return Result{Error: fmt.Sprintf("invalid short hex color: %s", err)}
		}

	case rgbPattern.MatchString(s):
		var err error
		r, g, b, err = parseRGB(s)
		if err != nil {
			return Result{Error: fmt.Sprintf("invalid rgb color: %s", err)}
		}

	case hslPattern.MatchString(s):
		var err error
		r, g, b, err = parseHSL(s)
		if err != nil {
			return Result{Error: fmt.Sprintf("invalid hsl color: %s", err)}
		}

	default:
		return Result{Error: fmt.Sprintf("unrecognized color format: %q", input)}
	}

	h, sat, l := rgbToHSL(r, g, b)

	output := fmt.Sprintf("Hex:  #%02x%02x%02x\nRGB:  rgb(%d, %d, %d)\nHSL:  hsl(%d, %d%%, %d%%)",
		r, g, b, r, g, b, h, sat, l)

	return Result{Output: output}
}

// parseHex6 parses a 6-digit hex color like #ff9800.
func parseHex6(s string) (uint8, uint8, uint8, error) {
	hex := strings.TrimPrefix(s, "#")
	rr, err := strconv.ParseUint(hex[0:2], 16, 8)
	if err != nil {
		return 0, 0, 0, err
	}
	gg, err := strconv.ParseUint(hex[2:4], 16, 8)
	if err != nil {
		return 0, 0, 0, err
	}
	bb, err := strconv.ParseUint(hex[4:6], 16, 8)
	if err != nil {
		return 0, 0, 0, err
	}
	return uint8(rr), uint8(gg), uint8(bb), nil
}

// parseHex3 parses a 3-digit hex color like #f90 (expands to #ff9900).
func parseHex3(s string) (uint8, uint8, uint8, error) {
	hex := strings.TrimPrefix(s, "#")
	expanded := string([]byte{hex[0], hex[0], hex[1], hex[1], hex[2], hex[2]})
	return parseHex6("#" + expanded)
}

// parseRGB parses an rgb(R,G,B) string.
func parseRGB(s string) (uint8, uint8, uint8, error) {
	m := rgbPattern.FindStringSubmatch(s)
	if m == nil {
		return 0, 0, 0, fmt.Errorf("does not match rgb pattern")
	}
	r, err := strconv.Atoi(m[1])
	if err != nil || r < 0 || r > 255 {
		return 0, 0, 0, fmt.Errorf("red value out of range: %s", m[1])
	}
	g, err := strconv.Atoi(m[2])
	if err != nil || g < 0 || g > 255 {
		return 0, 0, 0, fmt.Errorf("green value out of range: %s", m[2])
	}
	b, err := strconv.Atoi(m[3])
	if err != nil || b < 0 || b > 255 {
		return 0, 0, 0, fmt.Errorf("blue value out of range: %s", m[3])
	}
	return uint8(r), uint8(g), uint8(b), nil
}

// parseHSL parses an hsl(H,S%,L%) string and converts to RGB.
func parseHSL(s string) (uint8, uint8, uint8, error) {
	m := hslPattern.FindStringSubmatch(s)
	if m == nil {
		return 0, 0, 0, fmt.Errorf("does not match hsl pattern")
	}
	h, err := strconv.Atoi(m[1])
	if err != nil || h < 0 || h > 360 {
		return 0, 0, 0, fmt.Errorf("hue out of range: %s", m[1])
	}
	sat, err := strconv.Atoi(m[2])
	if err != nil || sat < 0 || sat > 100 {
		return 0, 0, 0, fmt.Errorf("saturation out of range: %s", m[2])
	}
	l, err := strconv.Atoi(m[3])
	if err != nil || l < 0 || l > 100 {
		return 0, 0, 0, fmt.Errorf("lightness out of range: %s", m[3])
	}

	r, g, b := hslToRGB(h, sat, l)
	return r, g, b, nil
}

// hslToRGB converts HSL values to RGB.
// H is 0-360, S and L are 0-100.
func hslToRGB(h, s, l int) (uint8, uint8, uint8) {
	hf := float64(h) / 360.0
	sf := float64(s) / 100.0
	lf := float64(l) / 100.0

	if sf == 0 {
		// Achromatic.
		v := uint8(math.Round(lf * 255))
		return v, v, v
	}

	var q float64
	if lf < 0.5 {
		q = lf * (1.0 + sf)
	} else {
		q = lf + sf - lf*sf
	}
	p := 2.0*lf - q

	r := hueToRGB(p, q, hf+1.0/3.0)
	g := hueToRGB(p, q, hf)
	b := hueToRGB(p, q, hf-1.0/3.0)

	return uint8(math.Round(r * 255)), uint8(math.Round(g * 255)), uint8(math.Round(b * 255))
}

// hueToRGB is a helper for HSL to RGB conversion.
func hueToRGB(p, q, t float64) float64 {
	if t < 0 {
		t += 1.0
	}
	if t > 1 {
		t -= 1.0
	}
	if t < 1.0/6.0 {
		return p + (q-p)*6.0*t
	}
	if t < 1.0/2.0 {
		return q
	}
	if t < 2.0/3.0 {
		return p + (q-p)*(2.0/3.0-t)*6.0
	}
	return p
}

// rgbToHSL converts RGB values (0-255) to HSL.
// Returns H (0-360), S (0-100), L (0-100).
func rgbToHSL(r, g, b uint8) (int, int, int) {
	rf := float64(r) / 255.0
	gf := float64(g) / 255.0
	bf := float64(b) / 255.0

	max := math.Max(rf, math.Max(gf, bf))
	min := math.Min(rf, math.Min(gf, bf))
	l := (max + min) / 2.0

	if max == min {
		// Achromatic.
		return 0, 0, int(math.Round(l * 100))
	}

	d := max - min
	var s float64
	if l > 0.5 {
		s = d / (2.0 - max - min)
	} else {
		s = d / (max + min)
	}

	var h float64
	switch max {
	case rf:
		h = (gf - bf) / d
		if gf < bf {
			h += 6.0
		}
	case gf:
		h = (bf-rf)/d + 2.0
	case bf:
		h = (rf-gf)/d + 4.0
	}
	h /= 6.0

	return int(math.Round(h * 360)), int(math.Round(s * 100)), int(math.Round(l * 100))
}
