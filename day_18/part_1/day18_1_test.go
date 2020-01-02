// Álvaro Castellano Vela 2019/12/02
package main

import (
	"fmt"
	"github.com/a-castellano/dijkstra"
	"reflect"
	"testing"
)

func TestCopyVault(t *testing.T) {

	graph := dijkstra.NewGraph()
	//Add the 3 verticies
	graph.AddVertex(100)
	graph.AddVertex(101)
	graph.AddVertex(102)
	//Add the arcs
	graph.AddArc(100, 101, 1)
	graph.AddArc(100, 102, 1)
	graph.AddArc(101, 102, 2)

	copied := copyVault(graph)

	best, _ := graph.Longest(100, 102)
	copiedBest, _ := copied.Longest(100, 102)

	if reflect.DeepEqual(best.Path, copiedBest.Path) == false {
		fmt.Printf("%v\n", best.Path)
		fmt.Printf("%v\n", copiedBest.Path)
		t.Errorf("Copied is not working")
	}

}
