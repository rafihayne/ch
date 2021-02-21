package search

import (
	"container/heap"
	"errors"

	"github.com/rafihayne/ch/pkg/graph"
)

// TODO can aStarPQElement and aStarVisitedElement be the same?
type aStarPQElement struct {
	currIdx    uint
	prevIdx    uint
	costToCome float64 // g(x)
	costToGo   float64 // h(x)
	index      int     // Index in the pq. kinda hate this
}

type aStarVisitedElement struct {
	prevIdx    uint
	costToCome float64
}

type AStarResult struct {
	Path         []uint
	PathLen      float64
	VisitedCount uint
}

type aStarPriorityQueue []*aStarPQElement

func (pq aStarPriorityQueue) Len() int { return len(pq) }

func (pq aStarPriorityQueue) Less(i, j int) bool {
	return pq[i].costToCome+pq[i].costToGo < pq[j].costToCome+pq[j].costToGo
}

func (pq *aStarPriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

func (pq *aStarPriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*aStarPQElement)
	item.index = n
	*pq = append(*pq, item)
}

func (pq aStarPriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func reverse(a []uint) []uint {
	for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}
	return a
}

func extractAStarSolution(g *graph.Graph, startIdx uint, goalIdx uint, visited map[uint]aStarVisitedElement, numVisited uint) (AStarResult, error) {
	path := []uint{goalIdx}
	prev, ok := visited[goalIdx]
	if !ok {
		return AStarResult{}, errors.New("No solution found")
	}
	pathlen := prev.costToCome
	// this is poorly written lol
	for prev.prevIdx != startIdx {
		curr := prev
		path = append(path, curr.prevIdx)
		prev, _ = visited[curr.prevIdx]
	}
	path = append(path, startIdx)

	return AStarResult{reverse(path), pathlen, numVisited}, nil
}

func aStarSearch(g *graph.Graph, startIdx uint, goalIdx uint, h func(graph.NodeValue, graph.NodeValue) float64) (map[uint]aStarVisitedElement, uint) {

	start := g.Nodes[startIdx].Value
	goal := g.Nodes[goalIdx].Value

	// Create priority queue
	pq := aStarPriorityQueue{}
	heap.Init(&pq)

	// Create visited map
	visited := make(map[uint]aStarVisitedElement)

	numVisited := uint(0)

	heap.Push(&pq, &aStarPQElement{startIdx, startIdx, 0.0, h(start, goal), 0})

	for pq.Len() > 0 {
		best := heap.Pop(&pq).(*aStarPQElement)
		numVisited++

		seen, found := visited[best.currIdx]
		better := false
		if found && best.costToCome < seen.costToCome {
			better = true
		}
		if !found || better {
			visited[best.currIdx] = aStarVisitedElement{best.prevIdx, best.costToCome}

			// Check on start != goal facilitates using astar for dijkstras
			if best.currIdx == goalIdx && startIdx != goalIdx {
				break
			}

			parent := g.Nodes[best.currIdx]
			for _, edge := range parent.Outgoing {
				childCostToCome := best.costToCome + edge.Weight
				child, childFound := visited[edge.To]
				childBetter := false
				if childFound && childCostToCome < child.costToCome {
					childBetter = true
				}

				if !childFound || childBetter {
					heap.Push(&pq, &aStarPQElement{edge.To, best.currIdx, childCostToCome, h(g.Nodes[edge.To].Value, goal), 0})
				}
			}
		}
	}

	return visited, numVisited
}

func AStar(g *graph.Graph, startIdx uint, goalIdx uint, h func(graph.NodeValue, graph.NodeValue) float64) (AStarResult, error) {
	// TODO maybe it's best just to stick to golang form and use ints rather than casting everywhere
	if startIdx > uint(len(g.Nodes)) || goalIdx > uint(len(g.Nodes)) {
		return AStarResult{}, errors.New("Index out of bounds")
	}

	visited, numVisited := aStarSearch(g, startIdx, goalIdx, h)
	return extractAStarSolution(g, startIdx, goalIdx, visited, numVisited)
}
