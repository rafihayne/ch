package util

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/rafihayne/ch/pkg/graph"
)

func GraphFromFile(nodePath string, edgePath string) graph.Graph {
	g := graph.Graph{}

	nodeFile, err := os.Open(nodePath)
	if err != nil {
		log.Fatalln("Couldn't open the node file", err)
	}
	r := csv.NewReader(nodeFile)
	// TODO gross
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		record = strings.Split(record[0], " ")

		long, _ := strconv.ParseFloat(record[0], 64)
		lat, _ := strconv.ParseFloat(record[1], 64)

		g.AddNode(graph.NodeValue{X: long, Y: lat})
	}

	edgeFile, err := os.Open(edgePath)
	if err != nil {
		log.Fatalln("Couldn't open the edge file", err)
	}
	r = csv.NewReader(edgeFile)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		record = strings.Split(record[0], " ")
		idx_one, _ := strconv.ParseInt(record[0], 10, 32)
		idx_two, _ := strconv.ParseInt(record[1], 10, 32)
		weight, _ := strconv.ParseFloat(record[2], 64)
		// weight := search.Euclidean(g.Nodes[int(idx_one)].Value, g.Nodes[int(idx_two)].Value)
		g.AddEdge(uint(idx_one), uint(idx_two), weight)
	}

	return g
}

func SamplesFromFile(samplePath string) [][]uint {
	samples := [][]uint{}
	sampleFile, err := os.Open(samplePath)
	if err != nil {
		log.Fatalln("Couldn't open the sample file", err)
	}
	r := csv.NewReader(sampleFile)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		record = strings.Split(record[0], " ")

		src, _ := strconv.ParseInt(record[0], 10, 32)
		dst, _ := strconv.ParseInt(record[1], 10, 32)
		samples = append(samples, []uint{uint(src), uint(dst)})
	}
	return samples
}

func SavePaths(pathsFilePath string, paths [][]uint) {
	pathFile, err := os.Create(pathsFilePath)
	if err != nil {
		log.Fatalln("Couldn't open the path file", err)
	}
	w := csv.NewWriter(pathFile)
	for _, path := range paths {
		pathStr := make([]string, len(path))
		for idx, val := range path {
			pathStr[idx] = fmt.Sprint(val)
		}
		if err := w.Write(pathStr); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
	}

	// Write any buffered data to the underlying writer (standard output).
	w.Flush()

	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
}

func GenSimpleGraph(eightConnected bool) graph.Graph {
	g := graph.Graph{}
	// 2 | 5 | 8
	// --+---+--
	// 1 | 4 | 7
	// --+---+--
	// 0 | 3 | 6

	for x := -1.0; x <= 1; x++ {
		for y := -1.0; y <= 1; y++ {
			g.AddNode(graph.NodeValue{X: x, Y: y})
		}
	}

	// Add 4 connected edges
	// Y direction edges
	g.AddEdgeBidirectional(0, 1, 1.0)
	g.AddEdgeBidirectional(1, 2, 1.0)
	g.AddEdgeBidirectional(3, 4, 1.0)
	g.AddEdgeBidirectional(4, 5, 1.0)
	g.AddEdgeBidirectional(6, 7, 1.0)
	g.AddEdgeBidirectional(7, 8, 1.0)

	// X-direction edges
	g.AddEdgeBidirectional(0, 3, 1.0)
	g.AddEdgeBidirectional(3, 6, 1.0)
	g.AddEdgeBidirectional(1, 4, 1.0)
	g.AddEdgeBidirectional(4, 7, 1.0)
	g.AddEdgeBidirectional(2, 5, 1.0)
	g.AddEdgeBidirectional(5, 8, 1.0)

	if eightConnected {
		// Add Diagonal edges
		g.AddEdgeBidirectional(0, 4, math.Sqrt(2.0))
		g.AddEdgeBidirectional(1, 5, math.Sqrt(2.0))
		g.AddEdgeBidirectional(3, 7, math.Sqrt(2.0))
		g.AddEdgeBidirectional(4, 8, math.Sqrt(2.0))

		g.AddEdgeBidirectional(1, 3, math.Sqrt(2.0))
		g.AddEdgeBidirectional(2, 4, math.Sqrt(2.0))
		g.AddEdgeBidirectional(4, 6, math.Sqrt(2.0))
		g.AddEdgeBidirectional(5, 7, math.Sqrt(2.0))
	}

	return g
}

// Rotting experimental plotting code. Too low level to waste time on now

// func draw(g *graph.Graph, path []uint) {
// 	dest := image.NewRGBA(image.Rect(0, 0, 250, 250))
// 	gc := draw2dimg.NewGraphicContext(dest)

// 	// draw background
// 	gc.SetFillColor(color.RGBA{0x48, 0x4d, 0x61, 0xff})
// 	draw2dkit.Rectangle(gc, 0, 0, 250, 250)
// 	gc.Fill()

// 	// Translate to 0,0 and scale
// 	tf := draw2d.NewTranslationMatrix(125, 125)
// 	tf.Compose(draw2d.NewScaleMatrix(50, 50))
// 	gc.SetMatrixTransform(tf)

// 	// Draw edges
// 	for _, n := range g.Nodes {
// 		nv := n.Value

// 		gc.SetStrokeColor(color.RGBA{0, 0xff, 0, 0xff})
// 		gc.SetLineWidth(0.05)
// 		for _, e := range n.Outgoing {
// 			gc.BeginPath()
// 			gc.MoveTo(nv.X, nv.Y)
// 			other := g.Nodes[int(e.To)].Value
// 			gc.LineTo(other.X, other.Y)
// 			gc.Close()
// 			gc.FillStroke()
// 		}
// 	}
// 	// Draw nodes last
// 	for _, n := range g.Nodes {
// 		nv := n.Value
// 		red := color.RGBA{0xff, 0, 0, 0xff}
// 		draw2dkit.Circle(gc, nv.X, nv.Y, 0.2)
// 		gc.SetFillColor(red)
// 		gc.Fill()
// 	}

// 	// Draw astar path
// 	if len(path) > 0 {
// 		prev := g.Nodes[path[0]].Value
// 		for _, idx := range path[1:] {
// 			curr := g.Nodes[idx].Value

// 			gc.SetStrokeColor(color.RGBA{0, 0, 0xff, 0xff})
// 			gc.SetLineWidth(0.05)
// 			gc.BeginPath()
// 			gc.MoveTo(prev.X, prev.Y)
// 			gc.LineTo(curr.X, curr.Y)
// 			gc.Close()
// 			gc.FillStroke()

// 			prev = curr

// 		}
// 	}

// 	// draw2dkit.Circle(gc, 1, 1, 1.0)
// 	// gc.SetFillColor(color.RGBA{0, 0xff, 0, 0xff})
// 	// gc.Fill()

// 	// rect := image.Rect(-500, -500, 500, 500)
// 	// new := image.NewRGBA(rect)
// 	// draw.NearestNeighbor.Scale(new, rect, dest, dest.Bounds(), draw.Over, nil)
// 	draw2dimg.SaveToPngFile("hello.png", dest)
// }
