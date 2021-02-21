package search

import (
	"math"
	"reflect"
	"testing"

	"github.com/rafihayne/ch/pkg/graph"
	"github.com/rafihayne/ch/pkg/util"
)

func TestEightConnectedAStar(t *testing.T) {
	// 2 | 5 | 8
	// --+---+--
	// 1 | 4 | 7
	// --+---+--
	// 0 | 3 | 6
	g := util.GenSimpleGraph(true)

	result, err := AStar(&g, 0, 8, util.Euclidean)
	if err != nil {
		t.Error(err.Error())
	}
	if result.VisitedCount != 3 {
		t.Error("Visited too many nodes")
	}
	expected := []uint{0, 4, 8}
	if !reflect.DeepEqual(result.Path, expected) {
		t.Errorf("Wrong Path! Expected: %v, Got: %v", expected, result.Path)
	}
	trueLen := 2 * math.Sqrt(2.0)
	if !(math.Abs(result.PathLen-trueLen) <= 1e-9) {
		t.Error("Wrong path length!")
	}
}

func TestFourConnectedAStar(t *testing.T) {
	// 2 | 5 | 8
	// --+---+--
	// 1 | 4 | 7
	// --+---+--
	// 0 | 3 | 6
	g := util.GenSimpleGraph(false)

	result, err := AStar(&g, 0, 8, util.Euclidean)
	if err != nil {
		t.Error(err.Error())
	}
	// TODO this seems like too many expansions?
	if result.VisitedCount != 8 {
		t.Error("Visited too many nodes")
	}
	expected := []uint{0, 1, 4, 7, 8}
	if !reflect.DeepEqual(result.Path, expected) {
		t.Errorf("Wrong Path! Expected: %v, Got: %v", expected, result.Path)
	}
	trueLen := 4.0
	if !(math.Abs(result.PathLen-trueLen) <= 1e-9) {
		t.Error("Wrong path length!")
	}
}

func TestFourConnectedNoSolutionAStar(t *testing.T) {
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
	g.AddEdge(0, 1, 1.0)
	g.AddEdge(1, 2, 1.0)
	g.AddEdge(3, 4, 1.0)
	g.AddEdge(4, 5, 1.0)
	g.AddEdge(6, 7, 1.0)
	g.AddEdge(7, 8, 1.0)

	// X-direction edges
	g.AddEdge(0, 3, 1.0)
	g.AddEdge(3, 6, 1.0)
	g.AddEdge(1, 4, 1.0)
	g.AddEdge(4, 7, 1.0)
	g.AddEdge(2, 5, 1.0)
	g.AddEdge(5, 8, 1.0)

	_, err := AStar(&g, 8, 0, util.Euclidean)
	if err == nil {
		t.Error("Solution found for impossible graph")
	}
}
