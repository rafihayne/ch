package search

import (
	"math"
	"reflect"
	"testing"

	"github.com/rafihayne/ch/pkg/graph"
	"github.com/rafihayne/ch/pkg/util"
)

func TestEightConnectedDijkstras(t *testing.T) {
	g := util.GenSimpleGraph(true)
	// 2 | 5 | 8
	// --+---+--
	// 1 | 4 | 7
	// --+---+--
	// 0 | 3 | 6

	// Disconnected component
	g.AddNode(graph.NodeValue{X: 100, Y: 100})

	result := Dijkstras(&g, 0)
	expected := []float64{0, 1, 2, 1, math.Sqrt(2.0), math.Sqrt(2.0) + 1, 2, math.Sqrt(2.0) + 1, 2 * math.Sqrt(2.0), math.Inf(1)}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Wrong Distances! Expected: %v, Got: %v", expected, result)
	}
}

func TestFourConnectedDijkstras(t *testing.T) {
	g := util.GenSimpleGraph(false)
	// 2 | 5 | 8
	// --+---+--
	// 1 | 4 | 7
	// --+---+--
	// 0 | 3 | 6

	// Disconnected component
	g.AddNode(graph.NodeValue{X: 100, Y: 100})

	result := Dijkstras(&g, 0)
	expected := []float64{0, 1, 2, 1, 2, 3, 2, 3, 4, math.Inf(1)}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Wrong Distances! Expected: %v, Got: %v", expected, result)
	}
}
