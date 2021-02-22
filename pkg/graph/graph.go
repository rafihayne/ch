package graph

import "fmt"

// https://ad-wiki.informatik.uni-freiburg.de/teaching/EfficientRoutePlanningSS2012

// Adjacency list representation as suggested by
// https://ad-teaching.informatik.uni-freiburg.de/route-planning-ss2012/lecture-1.pdf
// Slide 14
type Graph struct {
	Nodes []Node
}

func (g *Graph) AddNode(nv NodeValue) uint {
	idx := uint(len(g.Nodes))
	g.Nodes = append(g.Nodes, Node{
		index: idx,
		Value: nv,
	})
	return idx
}

func (g *Graph) AddEdge(from uint, to uint, weight float64) {
	g.Nodes[from].addOutgoingEdge(to, weight)
	g.Nodes[to].addIncomingEdge(from, weight)
}

func (g *Graph) AddEdgeBidirectional(from uint, to uint, weight float64) {
	g.AddEdge(to, from, weight)
	g.AddEdge(from, to, weight)
}

func (g *Graph) Print() {
	fmt.Println("Nodes: ", len(g.Nodes))
	for _, n := range g.Nodes {
		fmt.Println(n.index, ": ", n.Value)
		for _, e := range n.Outgoing {
			fmt.Println(n.index, " -> ", e.To, " ", e.Weight)
		}
		for _, e := range n.Incoming {
			fmt.Println(n.index, " <- ", e.To, " ", e.Weight)
		}
	}
}

type Edge struct {
	// The from index of the edge is implicitly represented
	To     uint
	Weight float64
}

type Node struct {
	index    uint   // TODO do we need this?
	Incoming []Edge //NB: needed for bidirectional search
	Outgoing []Edge
	Value    NodeValue
}

func (n *Node) addOutgoingEdge(to uint, weight float64) {
	n.Outgoing = append(n.Outgoing, Edge{To: to, Weight: weight})
}

func (n *Node) addIncomingEdge(from uint, weight float64) {
	n.Incoming = append(n.Incoming, Edge{To: from, Weight: weight})
}

// TODO how to make this semi generic??
type NodeValue struct {
	X float64 `json:longitude`
	Y float64 `json:latitude`
}
