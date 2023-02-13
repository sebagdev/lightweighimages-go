package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-errors/errors"
)

const BaseURL = "https://api.isevenapi.xyz/api"

var (
	InternalServerError = fmt.Errorf("resource could not be found")
)

type Data struct {
	IsEven bool   `json:"iseven"`
	Ad     string `json:"ad"`
}

type HttpResponse struct {
	Message     string
	Status      int
	Description string
}

func ErrorHandler(c *gin.Context, err any) {
	goErr := errors.Wrap(err, 2)
	httpResponse := HttpResponse{Message: "Internal server error", Status: 500, Description: goErr.Error()}
	c.AbortWithStatusJSON(500, httpResponse)
}

func iseven(val int) (*Data, error) {
	URL := fmt.Sprintf("%s/iseven/%d", BaseURL, val)
	resp, err := http.Get(URL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var data Data
	err = json.Unmarshal(body, &data)
	if err != nil {
		panic(err)
	}
	return &data, nil

}

func main() {
	r := gin.Default()
	r.Use(gin.CustomRecovery(ErrorHandler))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/isEven/:number", func(c *gin.Context) {
		number := c.Param("number")
		n, err := strconv.Atoi(number)
		if err != nil {
			panic(err)
		}
		b, err := iseven(n)

		c.JSON(200, gin.H{
			"even": b.IsEven,
		})

	})

	r.Run() // listen and serve on 0.0.0.0:8080
}
