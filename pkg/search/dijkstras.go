package search

import (
	"errors"
	"math"

	"github.com/rafihayne/ch/pkg/graph"
	"github.com/rafihayne/ch/pkg/util"
)

func Dijkstras(g *graph.Graph, target uint) []float64 {
	visited, _ := aStarSearch(g, target, target, util.Zero)
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

func BiDirectionalDijkstras(g *graph.Graph, startIdx uint, goalIdx uint) (AStarResult, error) {
	// TODO maybe it's best just to stick to golang form and use ints rather than casting everywhere
	if startIdx > uint(len(g.Nodes)) || goalIdx > uint(len(g.Nodes)) {
		return AStarResult{}, errors.New("Index out of bounds")
	}

	visitedForward, visitedBackward, meetingIdx := biDirectionalaStarSearch(g, startIdx, goalIdx, util.Zero)
	return extractBidirectionalAStarSolution(g, startIdx, goalIdx, meetingIdx, visitedForward, visitedBackward)
}
