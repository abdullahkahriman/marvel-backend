package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
)

const (
	BaseUrl = "https://gateway.marvel.com/v1/public"
	ApiKey  = "b40178c30b73dfb48dce61ed34674785"
	Hash    = "d1b9092734f4f09325b28318187b38eb"
)

// create to dtos
type MarvelDto struct {
	Data MarvelDataDto `json:"data"`
}

type MarvelDataDto struct {
	Offset  int               `json:"offset"`
	Limit   int               `json:"limit"`
	Total   int               `json:"total"`
	Count   int               `json:"count"`
	Results []MarvelResultDto `json:"results"`
}

type MarvelResultDto struct {
	ID          int                `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Thumbnail   MarvelThumbnailDto `json:"thumbnail"`
}

type MarvelThumbnailDto struct {
	Path string `json:"path"`
}

func getCharacters(c *gin.Context) {
	page := c.Param("page")

	// Get request
	resp, err := http.Get(BaseUrl + "/characters?apikey=" + ApiKey + "&ts=1&offset=" + page + "&hash=" + Hash)

	if err != nil {
		fmt.Println("No response from request")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var result MarvelDto
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("Can not unmarshal JSON")
	}

	// val.Thumbnail.Path + "/portrait_xlarge.jpg"

	c.IndentedJSON(http.StatusOK, result.Data.Results)
}

func main() {
	router := gin.Default()

	router.GET("/api/characters", getCharacters)
	router.GET("/api/characters/:page", getCharacters)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4200"},
		AllowCredentials: true,
	})
	handler := c.Handler(router)

	port := os.Getenv("PORT")

	// router.Run("localhost:1907")
	log.Fatal(http.ListenAndServe(":"+port, handler))
	router.Run()
}
