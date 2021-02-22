package main

import (
	"fmt"
	"reflect"
	"time"

	"github.com/rafihayne/ch/pkg/search"
	"github.com/rafihayne/ch/pkg/util"
)

func main() {
	g := util.GraphFromFile("./data/nodes.csv", "./data/edges.csv")
	samples := util.SamplesFromFile("./data/samples.csv")
	sample := samples[0]

	start := time.Now()
	resultBi, _ := search.BiDirectionalDijkstras(&g, sample[0], sample[1])
	fmt.Println("Finished BiDijkstras in : ", time.Now().Sub(start))

	start = time.Now()
	result, _ := search.AStar(&g, sample[0], sample[1], util.Haversine)
	fmt.Println("Finished AStar in : ", time.Now().Sub(start))

	start = time.Now()
	resultD := search.Dijkstras(&g, sample[0])
	fmt.Println("Finished Dijkstras in : ", time.Now().Sub(start))

	fmt.Println("Same path as astar? ", reflect.DeepEqual(resultBi.Path, result.Path))
	fmt.Println("BiDijkstras PathLen: ", resultBi.PathLen)
	fmt.Println("Astar PathLen: ", result.PathLen)
	fmt.Println("Dijkstras PathLen: ", resultD[sample[1]])

}
