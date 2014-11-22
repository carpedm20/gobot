package main

import (
	"github.com/carpedm20/gobot/facebook"
	"github.com/garyburd/redigo/redis"
)

func main() {
	f := facebook.New()
	f.SetProxy("https://localhost:8080")

	c, err := redis.Dial("tcp", ":6379")
	if err != nil {
		panic(err)
	}
	defer c.Close()

	fbtoken, err := redis.String(c.Do("GET", "fbtoken"))
	if err != nil {
		panic("fbtoken not found")
	}

	f.Login(fbtoken)
	access := f.GetAccessByName("유니스트 내가 전해주까")
	println(access.Name)
	access.Post("")
	access.PostLink("https://github.com/carpedm20/gobot")
}
