package main

import (
	"log"

	"github.com/lafikl/consistent"
)

func main() {
	c := consistent.New()

	// adds the hosts to the ring
	c.Add("127.0.0.1:8000")
	c.Add("127.0.0.1:8001")
	c.Add("127.0.0.1:8002")
	c.Add("127.0.0.1:8003")

	// Returns the host that owns `key`.
	//
	// As described in https://en.wikipedia.org/wiki/Consistent_hashing
	//
	// It returns ErrNoHosts if the ring has no hosts in it.
	host, err := c.Get("/app1.html")
	host2, _ := c.Get("/app2.html")
	host3, _ := c.Get("/app3.html")
	host4, _ := c.Get("/app4.html")
	host5, _ := c.Get("/app5.html")
	host6, _ := c.Get("/app6.html")
	host7, _ := c.Get("/app7.html")
	host8, _ := c.Get("/app8.html")
	log.Println(host)
	log.Println(host2)
	log.Println(host3)
	log.Println(host4)
	log.Println(host5)
	log.Println(host6)
	log.Println(host7)
	log.Println(host8)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("===========")
	c.Add("127.0.0.1:8004")
	c.Add("127.0.0.1:8005")
	host, _ = c.Get("/app1.html")
	host2, _ = c.Get("/app2.html")
	host3, _ = c.Get("/app3.html")
	host4, _ = c.Get("/app4.html")
	host5, _ = c.Get("/app5.html")
	host6, _ = c.Get("/app6.html")
	host7, _ = c.Get("/app7.html")
	host8, _ = c.Get("/app8.html")
	log.Println(host)
	log.Println(host2)
	log.Println(host3)
	log.Println(host4)
	log.Println(host5)
	log.Println(host6)
	log.Println(host7)
	log.Println(host8)
}
