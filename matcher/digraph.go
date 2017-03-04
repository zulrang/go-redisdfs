package matcher

import (
	"log"
)

type DirectedGraph interface {
	AddEdge(from string, to string)
	RemoveEdge(from string, to string)
	GetConnected(from string) []string
}

type SearchInfo struct {
	curDepth int
	path *Stack
	visited map[string]bool
	maxDepth int
}

func dfs_path(graph DirectedGraph, source string, dest string, search *SearchInfo) bool {
	// increase depth counter
	search.curDepth = search.curDepth + 1
	// push the source into the path
	search.path.Push(source)
	// mark visited
	if search.visited[source] {
		return false
	} else {
		search.visited[source] = true
	}
	// get edges
	links := graph.GetConnected(source)
	log.Println("Traversing ", source, " : ", links, " Path: ", search.path.List())
	// traverse edges
	for _, node := range links {
		// check element to see if its what we're looking for
		if node == dest {
			// found!
			// put node on end of list
			search.path.Push(node)
			return true
		} else {
			var found bool
			if search.curDepth+1 <= search.maxDepth {
				// if not, keep searching
				found = dfs_path(graph, node, dest, search)
			} else {
				found = false
			}
			// if not found, remove from path and dec depth
			if !found {
				search.path.Pop()
				search.curDepth = search.curDepth - 1
			} else {
				return found
			}
		}
	}
	return false
}
