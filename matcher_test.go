package matcher

import (
    "testing"
    "strings"
    "errors"
    "io"
    "bufio"
    "os"
)


type PairLine string

func (l PairLine) Split(str string) (string, string, error) {
	s := strings.Split(string(l), str)
	if len(s) < 2 {
		return "", "", errors.New("Minimum match not found")
	}
	return s[0], s[1], nil
}

func ReadData(graph *RedisDirectedGraph, r io.Reader) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := PairLine(scanner.Text())
		from, to, err := line.Split(" ")
		if err != nil {
			panic(err)
		}
		graph.AddEdge(from, to)
	}
	return scanner.Err()
}

func LoadDataFromFile(graph *RedisDirectedGraph) {
	datafile, err := os.Open("test.dat")
	if err != nil {
		panic(err)
	}
	err = graph.client.FlushDb().Err()
	if err != nil {
		panic(err)
	}
	err = ReadData(graph, datafile)
	if err != nil {
		panic(err)
	}
}

func testEq(a, b []string) bool {

    if a == nil && b == nil {
        return true;
    }

    if a == nil || b == nil {
        return false;
    }

    if len(a) != len(b) {
        return false
    }

    for i := range a {
        if a[i] != b[i] {
            return false
        }
    }

    return true
}

func TestFindLoopFound(t *testing.T) {
    expect := []string{"car", "motorcycle", "dishwasher", "microwave", "keurig"}
	want := "car"
	have := "keurig"
	graph := new(RedisDirectedGraph)
	graph.Connect("localhost:6379", "", 0)
    LoadDataFromFile(graph)
    matcher := NewMatcher(graph)
	//LoadDataFromFile(graph)
    matched, path := matcher.FindLoop(have, want, 50)

    if !matched {
        t.Fatalf("Match not found!")
    }

    if !testEq(path, expect) {
        t.Fatal("Expected ", expect, " found ", path)
    }
}

func TestFindLoopNotFound(t *testing.T) {
	want := "woot"
	have := "keurig"
	graph := new(RedisDirectedGraph)
	graph.Connect("localhost:6379", "", 0)
    LoadDataFromFile(graph)
    matcher := NewMatcher(graph)
	//LoadDataFromFile(graph)
    matched, _ := matcher.FindLoop(have, want, 50)

    if matched {
        t.Fatal("Match found where there shouldn't have been one!")
    }
}
