package matcher

type Matcher struct {
    graph DirectedGraph
}

func NewMatcher(graph DirectedGraph) *Matcher {

    matcher := new(Matcher)
    matcher.graph = graph

    return matcher
}

func (self *Matcher) FindCycleThrough(source string, through string, maxDepth int) (bool, []string) {
    // create new search structure
	search := &SearchInfo{
		curDepth: 0,
		maxDepth: maxDepth,
		path: new(Stack),
		visited: make(map[string]bool),
	}
	// traverse through in order to find source
	return dfs_path(self.graph, through, source, search), search.path.List()
}
