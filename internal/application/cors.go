package application

import (
	"github.com/gin-gonic/gin"
)

func WithCORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		rw := c.Writer
		rw.Header().Add("Access-Control-Allow-Origin", "*")
		rw.Header().Add("Access-Control-Expose-Headers", "Content-Length,Content-Range,X-test")
		rw.Header().Add("Access-Control-Allow-Headers", "Referer,X-JWT,Accept,Keep-Alive,DNT,Origin,User-Agent,X-Requested-With,"+
			"If-Modified-Since,InMemCache-Control,Content-Type,Range,authorization,Cookie,Upgrade,Connection,Sec-WebSocket-Key,"+
			"Sec-WebSocket-Version,Sec-WebSocket-Extensions")
		rw.Header().Add("Access-Control-Allow-Methods", "*")

		c.Next()
	}
}
