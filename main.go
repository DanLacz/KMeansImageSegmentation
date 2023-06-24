package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"

	"github.com/muesli/clusters"
	"github.com/muesli/kmeans"
)

var COLOURS = 4

// to keep track of original pixel positions
type CustomObservation struct {
	coordinates clusters.Coordinates
	originalX   int
	originalY   int
}

// needs to implement to fulfill interface
func (c *CustomObservation) Coordinates() clusters.Coordinates {
	return c.coordinates
}

func (c *CustomObservation) Distance(point clusters.Coordinates) float64 {
	return c.coordinates.Distance(point)
}

func main() {
	//loads image
	file, err := os.Open("input.png")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	img, err := png.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	img = img.(*image.NRGBA)

	pixels := preprocessImage(img)

	//performs k-means
	km := kmeans.New()
	clusters, err := km.Partition(pixels, COLOURS)
	if err != nil {
		fmt.Println(err)
	}

	//redraws image
	finalimg := image.NewNRGBA(image.Rect(0, 0, img.Bounds().Dx(), img.Bounds().Dy()))
	for _, c := range clusters {
		fmt.Printf("Centered at r: %.2f g: %.2f b: %.2f\n", c.Center[0]*255, c.Center[1]*255, c.Center[2]*255)
		//fmt.Printf("Matching data points: %+v\n\n", c.Observations)
		//c.Observations is a list of colours which are to be changed to their centroid
		for _, obs := range c.Observations {
			o := obs.(*CustomObservation)
			o.Coordinates()[0] = c.Center[0] * 255
			o.Coordinates()[1] = c.Center[1] * 255
			o.Coordinates()[2] = c.Center[2] * 255
			//fmt.Printf("ORIGINAL X:%d, ORIGINAL Y:%d \n", o.originalX, o.originalY)
			//time.Sleep(50 * time.Millisecond)
			finalimg.Set(o.originalX, o.originalY, color.NRGBA{uint8(o.Coordinates()[0]), uint8(o.Coordinates()[1]), uint8(o.Coordinates()[2]), 255})
		}
	}

	//saves image
	outputFile, err := os.Create("output.png")
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()
	png.Encode(outputFile, finalimg)
}

func preprocessImage(img image.Image) clusters.Observations {
	bounds := img.Bounds()
	var o clusters.Observations
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			//divide by 255 to get normal 8 bit color, divide by 255 again for 0-1 float
			o = append(o, &CustomObservation{clusters.Coordinates{float64(r) / 65025, float64(g) / 65025, float64(b) / 65025}, x, y})
		}
	}
	return o
}
