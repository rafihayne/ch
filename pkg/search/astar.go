package search

import (
	"container/heap"
	"errors"
	"math"

	"github.com/rafihayne/ch/pkg/graph"
)

const (
	forward  = 1
	backward = 2
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

func extractBidirectionalAStarSolution(g *graph.Graph, startIdx uint, goalIdx uint, middleIdx uint, visitedForward map[uint]aStarVisitedElement, visitedBackward map[uint]aStarVisitedElement) (AStarResult, error) {
	resultForward, err := extractAStarSolution(g, startIdx, middleIdx, visitedForward, 0)
	if err != nil {
		return AStarResult{}, err
	}
	resultBackward, err := extractAStarSolution(g, goalIdx, middleIdx, visitedBackward, 0)
	if err != nil {
		return AStarResult{}, err
	}

	if resultForward.Path[len(resultForward.Path)-1] != resultBackward.Path[len(resultBackward.Path)-1] {
		return AStarResult{}, errors.New("Rafi doesn't know what went wrong yet")
	}

	path := resultForward.Path
	// TODO this unnecessarily double reverses backwards. small optimization
	path = append(path, reverse(resultBackward.Path)[1:]...)
	return AStarResult{path, resultForward.PathLen + resultBackward.PathLen, 0}, nil
}

func biDirectionalaStarSearch(g *graph.Graph, startIdx uint, goalIdx uint, h func(graph.NodeValue, graph.NodeValue) float64) (map[uint]aStarVisitedElement, map[uint]aStarVisitedElement, uint) {
	// References
	// https://www.cs.princeton.edu/courses/archive/spr06/cos423/Handouts/EPP%20shortest%20path%20algorithms.pdf
	// https://www.homepages.ucl.ac.uk/~ucahmto/math/2020/05/30/bidirectional-dijkstra.html

	// TODO: Figure out how to use the heuristic function for bi-astar
	// TODO: Compute visited count

	s := g.Nodes[startIdx].Value
	t := g.Nodes[goalIdx].Value

	mu := math.MaxFloat64 // Best path seen so far

	// Create priority queue
	pqForward := aStarPriorityQueue{}
	heap.Init(&pqForward)
	pqBackward := aStarPriorityQueue{}
	heap.Init(&pqBackward)

	// Create visited map
	visitedForward := make(map[uint]aStarVisitedElement)
	visitedBackward := make(map[uint]aStarVisitedElement)

	heap.Push(&pqForward, &aStarPQElement{startIdx, startIdx, 0.0, h(s, t), 0})
	heap.Push(&pqBackward, &aStarPQElement{goalIdx, goalIdx, 0.0, h(t, s), 0})

	direction := forward
	var meetingIdx uint

	for pqForward.Len() > 0 && pqBackward.Len() > 0 {
		// Peak on heaps to determine direction
		var best *aStarPQElement
		topForward := pqForward[0]
		topBackward := pqBackward[0]

		if topForward.costToCome < topBackward.costToCome {
			best = heap.Pop(&pqForward).(*aStarPQElement)
			direction = forward
		} else {
			best = heap.Pop(&pqBackward).(*aStarPQElement)
			direction = backward
		}

		var seen aStarVisitedElement
		found := false

		seenForward, foundForward := visitedForward[best.currIdx]
		seenBackward, foundBackward := visitedBackward[best.currIdx]
		if direction == forward {
			seen, found = seenForward, foundForward
		} else {
			seen, found = seenBackward, foundBackward
		}

		// Terminate if we've seen this node from both directions, and it's suboptimal vs the previous best
		if (foundForward && foundBackward) && topForward.costToCome+topBackward.costToCome >= mu {
			break
		}

		better := false
		if found && best.costToCome < seen.costToCome {
			better = true
		}
		if !found || better {
			if direction == forward {
				visitedForward[best.currIdx] = aStarVisitedElement{best.prevIdx, best.costToCome}
			} else {
				visitedBackward[best.currIdx] = aStarVisitedElement{best.prevIdx, best.costToCome}
			}

			parent := g.Nodes[best.currIdx]
			if direction == forward {
				for _, edge := range parent.Outgoing {
					childCostToCome := best.costToCome + edge.Weight
					child, childFound := visitedForward[edge.To]
					childBetter := false
					if childFound && childCostToCome < child.costToCome {
						childBetter = true
					}

					if !childFound || childBetter {
						heap.Push(&pqForward, &aStarPQElement{edge.To, best.currIdx, childCostToCome, h(g.Nodes[edge.To].Value, t), 0})
						if bestBackward, ok := visitedBackward[edge.To]; ok {
							dist := childCostToCome + bestBackward.costToCome
							if dist < mu {
								mu = dist
								meetingIdx = edge.To
							}
						}
					}
				}
			} else {
				for _, edge := range parent.Incoming {
					childCostToCome := best.costToCome + edge.Weight
					child, childFound := visitedBackward[edge.To]
					childBetter := false
					if childFound && childCostToCome < child.costToCome {
						childBetter = true
					}

					if !childFound || childBetter {
						heap.Push(&pqBackward, &aStarPQElement{edge.To, best.currIdx, childCostToCome, h(t, g.Nodes[edge.To].Value), 0})
						if bestForward, ok := visitedForward[edge.To]; ok {
							dist := childCostToCome + bestForward.costToCome
							if dist < mu {
								mu = dist
								meetingIdx = edge.To
							}
						}
					}
				}
			}
		}
	}
	return visitedForward, visitedBackward, meetingIdx
}

func AStar(g *graph.Graph, startIdx uint, goalIdx uint, h func(graph.NodeValue, graph.NodeValue) float64) (AStarResult, error) {
	// TODO maybe it's best just to stick to golang form and use ints rather than casting everywhere
	if startIdx > uint(len(g.Nodes)) || goalIdx > uint(len(g.Nodes)) {
		return AStarResult{}, errors.New("Index out of bounds")
	}

	visited, numVisited := aStarSearch(g, startIdx, goalIdx, h)
	return extractAStarSolution(g, startIdx, goalIdx, visited, numVisited)
}

func BiDirectionalAStar(g *graph.Graph, startIdx uint, goalIdx uint, h func(graph.NodeValue, graph.NodeValue) float64) (AStarResult, error) {
	// TODO maybe it's best just to stick to golang form and use ints rather than casting everywhere
	if startIdx > uint(len(g.Nodes)) || goalIdx > uint(len(g.Nodes)) {
		return AStarResult{}, errors.New("Index out of bounds")
	}

	visitedForward, visitedBackward, meetingIdx := biDirectionalaStarSearch(g, startIdx, goalIdx, h)
	return extractBidirectionalAStarSolution(g, startIdx, goalIdx, meetingIdx, visitedForward, visitedBackward)
}
