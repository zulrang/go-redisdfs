package matcher

import (
	"time"
	"github.com/garyburd/redigo/redis"
	"log"
)

type RedisDirectedGraph struct {
	url string
	client redis.Conn
	Reconnect bool
}

func (self *RedisDirectedGraph) AdjacentEdges(from string) []string {
	items, err := redis.Strings(self.client.Do("SMEMBERS", from))
	if items == nil {
		// doesn't exist
	} else if err != nil {
		if err = self.Connect(); err != nil {
			panic(err)
		}
		return self.AdjacentEdges(from)
	} else {
		return items
	}
	return make([]string, 0)
}

func (self *RedisDirectedGraph) AddEdge(from string, to string) {
	if _, err := self.client.Do("SADD", from, to); err != nil {
		if err = self.Connect(); err != nil {
			panic(err)
		}
		self.AddEdge(from, to)
	}
}

func (self *RedisDirectedGraph) RemoveEdge(from string, to string) {
	_, err := self.client.Do("SREM", from, to)
	if err != nil {
		if err = self.Connect(); err != nil {
			panic(err)
		}
		self.RemoveEdge(from, to)
	}
}

func (self *RedisDirectedGraph) SetUrl(url string) {
	self.url = url
}

func New(url string) *RedisDirectedGraph {
    
    redigraph := new(RedisDirectedGraph)

    redigraph.SetUrl(url) // rabbitmq server
    redigraph.Reconnect = true

    return redigraph
}

func (self *RedisDirectedGraph) NoReconnect() {
	self.Reconnect = false
}

func (self *RedisDirectedGraph) Connect() (err error) {
	// start with 2 seconds before retry
	var stepping = 2
	var curstep = 0
	// try connection until it succeeds
	self.client, err = redis.DialURL(self.url)
	for err != nil {
		log.Printf("Unable to connect to Redis.  Retrying in %d seconds...\n", stepping)
		log.Printf(err.Error())
		// wait a bit
		curstep = stepping
		for curstep > 0 {
			curstep = curstep - 1
    		time.Sleep(1 * time.Second)			
		}
    	// increase time between attempts
    	if stepping < 60 {
    		stepping = stepping * 2
    	}
    	if !self.Reconnect {
    		return err
    	}
    	self.client, err = redis.DialURL(self.url)
	}
	log.Printf("Connected to Redis")

	_, err = self.client.Do("SET", "connected", time.Now().String())
	if err != nil {
		return err
	}	

	// get caller check for error
	_, err = self.client.Do("PING")
	return err

}
