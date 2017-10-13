package matcher

import (
	"time"
	"github.com/garyburd/redigo/redis"
)

type RedisDirectedGraph struct {
	client redis.Conn
}

func (self *RedisDirectedGraph) GetConnected(from string) []string {
	items, err := self.client.Do("SMEMBERS", from)
	if items == nil {
		// doesn't exist
	} else if err != nil {
		panic(err)
	} else {
		if strarr, ok := items.([]string); ok {
			return strarr
		}
	}
	return make([]string, 0)
}

func (self *RedisDirectedGraph) AddEdge(from string, to string) {
	_, err := self.client.Do("SADD", from, to)
	if err != nil {
		panic(err)
	}
}

func (self *RedisDirectedGraph) RemoveEdge(from string, to string) {
	_, err := self.client.Do("SREM", from, to)
	if err != nil {
		panic(err)
	}
}

func (self *RedisDirectedGraph) Connect(url string) (err error) {
	self.client, err = redis.DialURL(url)

	if err != nil {
		panic(err)
	}	

	_, err = self.client.Do("SET", "connected", time.Now().String())
	if err != nil {
		panic(err)
	}	
	// get caller check for error
	_, err = self.client.Do("PING")
	return err
}
