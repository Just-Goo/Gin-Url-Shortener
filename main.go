package main

import (
	"fmt" 
	"net/http"
	"sync" 

	"github.com/gin-gonic/gin"
	"github.com/lithammer/shortuuid/v4"
)

type mapper struct {
	mapping map[string]string
	sync.Mutex
}

var urlMapper mapper

func init() {
	urlMapper = mapper{
		mapping: make(map[string]string),
	}
}

func main() {
	router := gin.New()

	router.Use(gin.Logger())

	router.GET("/", func(c *gin.Context) {
		c.Writer.Write([]byte("Server is running.."))
	})

	router.POST("/shorten", createShortUrlHandler)
	router.GET("/s/:key", redirectHandler)

	router.Run(":8080")
}

func createShortUrlHandler(c *gin.Context) {
	c.Request.ParseForm()
	u := c.Request.Form.Get("URL")
	if u == "" {
		c.Writer.WriteHeader(http.StatusBadRequest)
		c.Writer.Write([]byte("URL field is empty"))
		return
	}

	// generate the unique key
	key := shortuuid.New()

	// insert
	insertMapping(key, u) 

	c.Writer.WriteHeader(http.StatusOK)

	c.Writer.Write([]byte(fmt.Sprintf("http://localhost:8080/s/%s", key)))

}

func redirectHandler(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		c.Writer.WriteHeader(http.StatusBadRequest)
		c.Writer.Write([]byte("Key field is empty"))
		return
	}
	
	// fetch mapping
	u := fetchMapping(key)
	fmt.Println(u)
	if u == "" {
		c.Writer.WriteHeader(http.StatusBadRequest)
		c.Writer.Write([]byte("URL field is empty"))
		return
	}

	http.Redirect(c.Writer, c.Request, u, http.StatusFound)

}

func insertMapping(key, u string) {
	urlMapper.Lock()
	defer urlMapper.Unlock()

	urlMapper.mapping[key] = u
}

func fetchMapping(key string) string {
	urlMapper.Lock()
	defer urlMapper.Unlock()
	return urlMapper.mapping[key]
}
