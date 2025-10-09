#!/bin/bash
# Demo script for SVG conversion feature

set -e

echo "================================================"
echo "Asset Generator - SVG Conversion Demo"
echo "================================================"
echo ""

# Create demo directory
DEMO_DIR="svg-conversion-demo"
mkdir -p "$DEMO_DIR"
cd "$DEMO_DIR"

echo "Step 1: Creating test image..."
cat > create_test_image.go << 'EOF'
package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

func main() {
	img := image.NewRGBA(image.Rect(0, 0, 300, 300))
	
	// Draw a colorful geometric pattern
	for y := 0; y < 300; y++ {
		for x := 0; x < 300; x++ {
			dx := float64(x - 150)
			dy := float64(y - 150)
			dist := dx*dx + dy*dy
			
			if dist < 2500 {
				img.Set(x, y, color.RGBA{255, 100, 100, 255}) // Red center
			} else if dist < 5000 {
				img.Set(x, y, color.RGBA{100, 255, 100, 255}) // Green ring
			} else if dist < 10000 {
				img.Set(x, y, color.RGBA{100, 100, 255, 255}) // Blue ring
			} else {
				img.Set(x, y, color.RGBA{240, 240, 240, 255}) // Light gray background
			}
		}
	}
	
	f, _ := os.Create("test-image.png")
	defer f.Close()
	png.Encode(f, img)
}
EOF

go run create_test_image.go
rm create_test_image.go
echo "✓ Test image created: test-image.png"
echo ""

echo "Step 2: Basic conversions with different shape counts..."
echo ""

echo "2a. Low quality (50 shapes) - Fast..."
../asset-generator convert svg test-image.png -o output-50shapes.svg --shapes 50
echo ""

echo "2b. Medium quality (100 shapes) - Default..."
../asset-generator convert svg test-image.png -o output-100shapes.svg --shapes 100
echo ""

echo "2c. High quality (300 shapes)..."
../asset-generator convert svg test-image.png -o output-300shapes.svg --shapes 300
echo ""

echo "Step 3: Conversions with different shape modes..."
echo ""

echo "3a. Triangles (mode 1) - Default..."
../asset-generator convert svg test-image.png -o output-triangles.svg --shapes 150 --mode 1
echo ""

echo "3b. Ellipses (mode 3) - Softer look..."
../asset-generator convert svg test-image.png -o output-ellipses.svg --shapes 150 --mode 3
echo ""

echo "3c. Circles (mode 4) - Pointillist effect..."
../asset-generator convert svg test-image.png -o output-circles.svg --shapes 200 --mode 4
echo ""

echo "3d. Bezier curves (mode 6) - Smooth..."
../asset-generator convert svg test-image.png -o output-beziers.svg --shapes 100 --mode 6
echo ""

echo "Step 4: Transparency variations..."
echo ""

echo "4a. High transparency (alpha 64)..."
../asset-generator convert svg test-image.png -o output-alpha64.svg --shapes 150 --alpha 64
echo ""

echo "4b. Low transparency (alpha 200)..."
../asset-generator convert svg test-image.png -o output-alpha200.svg --shapes 150 --alpha 200
echo ""

echo "================================================"
echo "Demo Complete!"
echo "================================================"
echo ""
echo "Generated files in $DEMO_DIR:"
ls -lh *.svg | awk '{print $9, "-", $5}'
echo ""
echo "Original PNG: $(du -h test-image.png | cut -f1)"
echo ""
echo "Tips:"
echo "  - Open SVG files in a browser to compare results"
echo "  - More shapes = better quality but larger files"
echo "  - Different modes create different artistic styles"
echo "  - Adjust alpha for layered vs solid effects"
echo ""
echo "Try these commands:"
echo "  cd $DEMO_DIR"
echo "  firefox *.svg  # View all SVGs"
echo "  ../asset-generator convert svg test-image.png --shapes 500  # Ultra quality"
echo ""
