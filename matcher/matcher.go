package matcher

type Matcher struct {
    graph DirectedGraph
}

func NewMatcher(graph DirectedGraph) *Matcher {

    matcher := new(Matcher)
    matcher.graph = graph

    return matcher
}

func (self *Matcher) FindLoop(have string, want string, maxDepth int) (bool, []string) {
    // create new search structure
	search := &SearchInfo{
		curDepth: 0,
		maxDepth: maxDepth,
		path: new(Stack),
		visited: make(map[string]bool),
	}
	// traverse my want in order to find myself
	return dfs_path(self.graph, want, have, search), search.path.List()
}
