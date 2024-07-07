package gee

import (
	"log"
	"time"
)

func Logger() HandlerFunc {
	return func(c *Context) {
		// Start timer
		t := time.Now()
		// Process request
		c.NextMiddle() //如果有人调用了logger，那logger负责后续的调用
		// Calculate resolution time
		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}
