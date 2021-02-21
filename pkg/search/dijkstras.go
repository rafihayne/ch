package search

import (
	"math"

	"github.com/rafihayne/ch/pkg/graph"
)

// distances := []float64{}
// for i := range g.Nodes {
// 	elt, ok := visited[uint(i)]
// 	if !ok {
// 		distances = append(distances, math.MaxFloat64)
// 	}
// 	distances = append(distances, elt.costToCome)
// }
// fmt.Println(distances)

func Dijkstras(g *graph.Graph, target uint) []float64 {
	visited, _ := aStarSearch(g, target, target, func(lhs graph.NodeValue, rhs graph.NodeValue) float64 { return 0 })
	distances := []float64{}
	for i := range g.Nodes {
		elt, ok := visited[uint(i)]
		if !ok {
			distances = append(distances, math.Inf(1))
		} else {
			distances = append(distances, elt.costToCome)
		}
	}
	return distances
}
