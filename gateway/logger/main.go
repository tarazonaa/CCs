package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// album represents data about a record album.
type album struct {
    ID     string  `json:"id"`
    Title  string  `json:"title"`
    Artist string  `json:"artist"`
    Price  float64 `json:"price"`
}

type log struct {}

// albums slice to seed record album data.
var albums = []album{
    {ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
    {ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
    {ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

func rootHandler(c *gin.Context) {
    c.String(http.StatusOK, "Welcome to the Logs API!")
}

func test(c *gin.Context) {
    c.String(http.StatusOK, "Test endpoint reached!")
}

// getAlbums responds with the list of all albums as JSON.
func getAlbums(c *gin.Context) {
    c.IndentedJSON(http.StatusOK, albums)
}

// postLog adds an album from JSON received in the request body.
func postLog(c *gin.Context) {
    var newLog log
    // Call BindJSON to bind the received JSON to
    // newAlbum.
    if err := c.BindJSON(&newLog); err != nil {
        return
    }

    fmt.Println("Received new log:", newLog)

    // Add the new album to the slice.
    // albums = append(albums, newAlbum)
    // c.IndentedJSON(http.StatusCreated, newAlbum)
}


func main() {
    router := gin.Default()
    router.GET("/", rootHandler)
    router.GET("/albums", getAlbums)
    router.GET("/test", test)
	router.POST("/logs", postLog)
	
    router.Run(":8080")
}