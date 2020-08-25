package main

import (
	"encoding/csv"
	"fmt"
	"image/color"
	"io"
	"log"
	"os"
	"strconv"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

func plottest(points plotter.XYs){
	p, _ := plot.New()

	p.Title.Text = "The cumulative graph of churn nodes"
	p.X.Label.Text = "churn nodes"
	p.Y.Label.Text = "count"
	p.Add(plotter.NewGrid())
	// var points plotter.XYs
	s, err := plotter.NewHistogram(points, points.Len())
	if err != nil {
		panic(err)
	}
	// s.Color = color.RGBA{R: 255, A: 255}
	p.Add(s)
	// p.Legend.Add("linepoint", s)
	p.Save(8*vg.Inch, 5*vg.Inch, "nodes_test.png")

}

func plotchurn(points plotter.XYs) {
	p, _ := plot.New()

	p.Title.Text = "The churn rate of bitcoin nodes each timestamps"
	p.X.Label.Text = "timestamp"
	p.Y.Label.Text = "churn rate (%)"
	p.Add(plotter.NewGrid())
	// var points plotter.XYs
	s, err := plotter.NewHistogram(points, points.Len())
	if err != nil {
		panic(err)
	}
	s.Color = color.RGBA{R: 255, A: 255}

	p.Add(s)
	// p.Legend.Add("linepoint", s)

	p.Save(8*vg.Inch, 5*vg.Inch, "nodes_churn.png")
}
func plotadd_r(points plotter.XYs) {
	p, _ := plot.New()

	p.Title.Text = "The add rate of bitcoin nodes each timestamps"
	p.X.Label.Text = "timestamp"
	p.Y.Label.Text = "add rate (%)"
	p.Add(plotter.NewGrid())
	// var points plotter.XYs
	s, err := plotter.NewHistogram(points, points.Len())
	if err != nil {
		panic(err)
	}
	s.Color = color.RGBA{R: 255, A: 255}

	p.Add(s)
	// p.Legend.Add("linepoint", s)

	p.Save(8*vg.Inch, 5*vg.Inch, "nodes_add_r.png")
}
func plotsnapshots() {
	p, _ := plot.New()

	p.Title.Text = "The snapshots of bitcoin nodes from 2020/4/30 1:35:30 to 2020/5/31 23:56:15"
	p.X.Label.Text = "timestamp"
	p.Y.Label.Text = "number of nodes"
	p.Add(plotter.NewGrid())
	var points plotter.XYs
	// points = readcsv("nodes_snapshots.csv")
	points = readcsv("nodes_snapshots_reverse_forchurn.csv")
	// plotutil.AddLinePoints(p, points)

	// Make a scatter plotter and set its style.
	s, err := plotter.NewScatter(points)
	if err != nil {
		panic(err)
	}
	s.GlyphStyle.Color = color.RGBA{R: 255, B: 128, A: 255}
	// s.Shape = draw.BoxGlyph{}
	// s.Shape = draw.CrossGlyph{}
	s.Shape = draw.CircleGlyph{}
	p.Add(s)
	p.Legend.Add("scatter", s)

	p.Save(15*vg.Inch, 6*vg.Inch, "nodes_snapshots.png")

}

func readcsv(file string) plotter.XYs {
	// Open the file
	csvfile, err := os.Open(file)
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}

	// Parse the file
	r := csv.NewReader(csvfile)
	//r := csv.NewReader(bufio.NewReader(csvfile))

	//my code
	var points plotter.XYs

	var count, min int
	min = 100000
	// Iterate through the records
	for {
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		// fmt.Printf("Question: %s Answer %s\n", record[0], record[1])
		// if record[0] == "1590940575"{
		// 	break
		// }
		// if record[0] == "1589649378"{
		// 	break
		// }
		var x, y float64
		x, err = strconv.ParseFloat(record[0], 64)
		y, err = strconv.ParseFloat(record[1], 64)
		points = append(points, struct{ X, Y float64 }{x, y})
		// if y < 7000 {
		// 	fmt.Println("the pair is: ", record[0], record[1])
		// }
		count++
		if min > int(y) {
			min = int(y)
		}

	}
	fmt.Println("min: ", min)     //6053
	fmt.Println("total: ", count) //15250
	return points
}
