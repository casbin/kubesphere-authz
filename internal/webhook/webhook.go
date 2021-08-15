package webhook

import (
	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
)

//set up https midware
func tlsHandler(c *gin.Context) {
	secureMiddleware := secure.New(secure.Options{
		SSLRedirect: true,
		SSLHost:     "localhost:8080",
	})
	err := secureMiddleware.Process(c.Writer, c.Request)

	// If there was an error, do not continue.
	if err != nil {
		return
	}
	c.Next()
}

func GetAdmissionWebhook() *gin.Engine {
	r := gin.Default()
	r.Any("/check", handler)
	r.Any("/deployment", handler)
	r.Any("/ping", ping)
	r.Use(tlsHandler)
	return r

}
