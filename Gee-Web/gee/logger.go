package gee

import (
	"log"
	"time"
)

// 在中间件中手动调用 Next() 一般用于在请求前后各实现一些行为。
// 如果中间件只作用于请求前，可以省略调用Next()，算是一种兼容性比较好的写法吧。
func Logger() HandlerFunc {
	h := func(c *Context) {
		// start timer
		t := time.Now()
		// process request
		c.Next()
		// calculate resolution time
		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
	return h
}
func Logger1() HandlerFunc {
	h := func(c *Context) {
		// start timer
		t := time.Now()
		// process request
		c.Next()
		// calculate resolution time
		log.Printf("<%d> %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
	return h
}
