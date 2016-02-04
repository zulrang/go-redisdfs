package main

import (
    "strconv"
    "os"
    "log"
    "github.com/zulrang/traidup/matcher"
)

func main() {
    have := os.Args[1]
    want := os.Args[2]
	graph := new(matcher.RedisDirectedGraph)
    db, err := strconv.ParseInt(os.Getenv("REDIS_DB"), 10, 64)

    if err != nil {
        panic("Invalid database int64")
    }
	graph.Connect(
        os.Getenv("REDIS_URL"),
        os.Getenv("REDIS_PASS"),
        db)
    matcher := matcher.NewMatcher(graph)
	//LoadDataFromFile(graph)
    matched, path := matcher.FindLoop(have, want, 50)

    if !matched {
        log.Println("Match not found!")
    } else {
        log.Println("Match found: ", path)
    }
}
