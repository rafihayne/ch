package main

import (
	"fmt"
	"time"

	"github.com/rafihayne/ch/pkg/search"
	"github.com/rafihayne/ch/pkg/util"
)

func main() {
	g := util.GraphFromFile("./data/nodes.csv", "./data/edges.csv")
	samples := util.SamplesFromFile("./data/samples.csv")

	start := time.Now()
	for _, sample := range samples {
		result, _ := search.AStar(&g, sample[0], sample[1], util.Haversine)
		_ = result
	}
	fmt.Println("Finished ", len(samples), " queries in : ", time.Now().Sub(start))

	start = time.Now()
	result := search.Dijkstras(&g, 0)
	_ = result
	fmt.Println("Finished in : ", time.Now().Sub(start))
}
