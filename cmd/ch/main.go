package main

import (
	"fmt"
	"time"

	"github.com/rafihayne/ch/pkg/graph"
	"github.com/rafihayne/ch/pkg/search"
	"github.com/rafihayne/ch/pkg/util"
)

func main() {
	g := util.GraphFromFile("./data/nodes.csv", "./data/edges.csv")
	samples := util.SamplesFromFile("./data/samples.csv")

	// g := util.GenSimpleGraph(true)

	start := time.Now()
	result, err := search.BiDirectionalAStar(&g, samples[0][0], samples[0][1], func(lhs graph.NodeValue, rhs graph.NodeValue) float64 { return 0 })
	// result, err := search.BiDirectionalAStar(&g, 0, 8, func(lhs graph.NodeValue, rhs graph.NodeValue) float64 { return 0 })
	fmt.Println(err)
	// fmt.Println(result)
	// _ = result
	fmt.Println("Finished in : ", time.Now().Sub(start))
	fmt.Println(result.PathLen)
	// fmt.Println(result)

	start = time.Now()
	result, _ = search.AStar(&g, samples[1][0], samples[1][1], util.Haversine)
	fmt.Println(result.PathLen)
	fmt.Println("Finished in : ", time.Now().Sub(start))

	start = time.Now()
	result2 := search.Dijkstras(&g, samples[1][0])
	fmt.Println(result2[samples[1][1]])
	// _ = result2
	fmt.Println("Finished in : ", time.Now().Sub(start))
}
