package matcher

import (
	redis "gopkg.in/redis.v3"
	"log"
	"errors"
	"strings"
)

type RedisDirectedGraph struct {
	client *redis.Client
}

func (self *RedisDirectedGraph) GetConnected(from string) []string {
	items, err := self.client.SMembers(from).Result()
	log.Println("SMembers", from, items)
	if err == redis.Nil {
		// doesn't exist
	} else if err != nil {
		panic(err)
	} else {
		return items
	}
	return make([]string, 0)
}

func (self *RedisDirectedGraph) AddEdge(from string, to string) {
	err := self.client.SAdd(from, to).Err()
	if err != nil {
		panic(err)
	}
}

func (self *RedisDirectedGraph) RemoveEdge(from string, to string) {
	err := self.client.SRem(from, to).Err()
	if err != nil {
		panic(err)
	}
}

func (self *RedisDirectedGraph) URLToOptions(url string) (options *redis.Options, err error) {
	parts := strings.Split(url, "://")
	if parts[0] != "redis" {
		return nil, errors.New("Not a redis URL")
	}
	if len(parts) != 2 {
		return nil, errors.New("Invalid redis URL")
	}
	parts = strings.Split(parts[1], "@")
	if len(parts) != 2 {
		return nil, errors.New("Invalid redis URL - missing @")
	}
	uri := parts[1]
	parts = strings.Split(parts[0], ":")
	if len(parts) != 2 {
		return nil, errors.New("Invalid redis URL - missing : in login")
	}
	password := parts[1]
	return &redis.Options{
		Addr: uri,
		Password: password,
		DB: 0,
	}, nil
}

func (self *RedisDirectedGraph) Connect(url string) (err error) {
	options, err := self.URLToOptions(url)

	if err == nil {	
		self.client = redis.NewClient(options)
	}

	// get caller check for error
	_, err = self.client.Ping().Result()
	return err
}
