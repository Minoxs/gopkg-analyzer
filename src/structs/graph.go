package structs

import (
	"log"
	"sort"
)

type SquareGraph struct {
	size         int
	graph        [][]bool
	nameMap      map[string]int
	hideDiagonal bool
}

func initGraph(size int) (res [][]bool) {
	res = make([][]bool, size)

	for i := 0; i < size; i++ {
		res[i] = make([]bool, size)
	}

	return res
}

func NewGraph(size int) *SquareGraph {
	return &SquareGraph{
		size:         size,
		graph:        initGraph(size),
		nameMap:      map[string]int{},
		hideDiagonal: false,
	}
}

func (g *SquareGraph) getName(idx int) string {
	for i, n := range g.nameMap {
		if n == idx {
			return i
		}
	}
	return ""
}

func (g *SquareGraph) HideDiagonal() *SquareGraph {
	g.hideDiagonal = true
	return g
}

func (g *SquareGraph) String() (res string) {
	res = ""
	for i := 0; i < g.size; i++ {
		// New Line
		ln := ""
		for j := 0; j < g.size; j++ {
			if g.hideDiagonal && i == j {
				ln += "."
			} else if g.graph[i][j] {
				ln += "1"
			} else {
				ln += "0"
			}
			ln += " "
		}
		res += ln + g.getName(i) + "\n"
	}
	return res
}

func (g *SquareGraph) AddKey(s string) *SquareGraph {
	g.nameMap[s] = len(g.nameMap)
	return g
}

func (g *SquareGraph) AddKeys(strings []string, sortKeys bool) *SquareGraph {
	if len(strings) > g.size {
		log.Println("Graph cannot fit all keys")
		return g
	}

	if sortKeys {
		sort.Strings(strings)
	}
	for i, key := range strings {
		g.nameMap[key] = i
	}

	return g
}

func (g *SquareGraph) Connect(a, b string) *SquareGraph {
	idxA, okA := g.nameMap[a]
	idxB, okB := g.nameMap[b]

	if !okA || !okB {
		log.Printf("Cannot connect %s and %s: Keys not added\n", a, b)
		return g
	}

	g.graph[idxA][idxB] = true
	return g
}
