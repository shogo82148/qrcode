package qrcode

import (
	"image"
	"image/color"
	"log"
	"math"
	"sort"
)

type xPattern struct {
	x int
	y int
	w int
}

type yPattern struct {
	x int
	y int
	h int
}

func Decode(img image.Image) (*QRCode, error) {
	bounds := img.Bounds()
	var runLength int
	buf := make([]int, 0, 10)

	xList := make([]xPattern, 0, 10)
	yList := make([]yPattern, 0, 10)

	// find pattern
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		buf = buf[:0]
		lastColor := colorToBin(img.At(bounds.Min.X, y))
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := colorToBin(img.At(x, y))
			if c == lastColor {
				runLength++
			} else {
				buf = append(buf, runLength)
				runLength = 0
				lastColor = c

				if len(buf) >= 5 {
					m1 := buf[len(buf)-1]
					m2 := buf[len(buf)-2]
					m3 := buf[len(buf)-3]
					m4 := buf[len(buf)-4]
					m5 := buf[len(buf)-5]
					h := m1 / 2

					f2 := h < m2 && m2 < m1+h
					f3 := m1*3-h < m3 && m3 < m1*3+h
					f4 := h < m4 && m4 < m1+h
					f5 := h < m5 && m5 < m1+h
					if !c && f2 && f3 && f4 && f5 {
						w := m1 + m2 + m3 + m4 + m5
						xList = append(xList, xPattern{
							x: x - w/2,
							y: y,
							w: w,
						})
					}
				}
			}
		}
	}

	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		buf = buf[:0]
		lastColor := colorToBin(img.At(x, bounds.Min.Y))
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			c := colorToBin(img.At(x, y))
			if c == lastColor {
				runLength++
			} else {
				buf = append(buf, runLength)
				runLength = 0
				lastColor = c

				if len(buf) >= 5 {
					m1 := buf[len(buf)-1]
					m2 := buf[len(buf)-2]
					m3 := buf[len(buf)-3]
					m4 := buf[len(buf)-4]
					m5 := buf[len(buf)-5]
					h := m1 / 2

					f2 := h < m2 && m2 < m1+h
					f3 := m1*3-h < m3 && m3 < m1*3+h
					f4 := h < m4 && m4 < m1+h
					f5 := h < m5 && m5 < m1+h
					if !c && f2 && f3 && f4 && f5 {
						h := m1 + m2 + m3 + m4 + m5
						yList = append(yList, yPattern{
							x: x,
							y: y - h/2,
							h: h,
						})
					}
				}
			}
		}
	}

	finderSet := map[image.Point]struct{}{}
	for _, x := range xList {
		for _, y := range yList {
			if x.x-x.w/2 < y.x && y.x < x.x+x.w/2 && y.y-y.h/2 < x.y && x.y < y.y+y.h/2 {
				finderSet[image.Pt(x.x, y.y)] = struct{}{}
			}
		}
	}
	finders := []image.Point{}
	for p := range finderSet {
		finders = append(finders, p)
	}
	sort.Slice(finders, func(i, j int) bool {
		if finders[i].X < finders[j].X {
			return true
		}
		if finders[i].X > finders[j].X {
			return false
		}
		return finders[i].Y < finders[j].Y
	})

	// TODO: improve me
	a, b, c := finders[0], finders[1], finders[2]
	_ = c

	// calculate the distance
	dx := float64(a.X - b.X)
	dy := float64(a.Y - b.Y)
	D := math.Sqrt(dx*dx + dy*dy)
	dx /= D
	dy /= D

	Wul := edgeDistance(img, float64(a.X), float64(a.Y), -dx, -dy) + edgeDistance(img, float64(a.X), float64(a.Y), dx, dy)
	Wur := edgeDistance(img, float64(b.X), float64(b.Y), -dx, -dy) + edgeDistance(img, float64(b.X), float64(b.Y), dx, dy)

	X := (float64(Wul) + float64(Wur)) / 14
	V := int(math.Round((D/X - 10) / 4))

	log.Println("version:", V)
	return &QRCode{}, nil
}

func colorToBin(c color.Color) bool {
	r, g, b, a := c.RGBA()
	if a < 0x1000 {
		return false
	}
	y := r + g + b
	return y < 0x8000*3
}

func edgeDistance(img image.Image, x, y, dx, dy float64) int {
	var d int
	// skip BLACK modules
	for colorToBin(img.At(int(math.Round(x)), int(math.Round(y)))) {
		x += dx
		y += dy
		d++
	}

	// skip WHITE modules
	for !colorToBin(img.At(int(math.Round(x)), int(math.Round(y)))) {
		x += dx
		y += dy
		d++
	}

	// skip BLACK modules
	for colorToBin(img.At(int(math.Round(x)), int(math.Round(y)))) {
		x += dx
		y += dy
		d++
	}

	return d
}
