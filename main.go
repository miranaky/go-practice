package main

import (
	"os"

	"github.com/labstack/echo"
	"github.com/miranaky/learngo/scrapper"
)

func handleGet(c echo.Context) error {
	return c.File("index.html")
}
func handlePost(c echo.Context) error {
	term := c.FormValue("term")
	fileName := term + "_jobs.csv"
	scrapper.Scrapper(term)
	defer os.Remove(fileName)
	return c.Attachment(fileName, fileName)
}

func main() {
	e := echo.New()
	e.GET("/", handleGet)
	e.POST("/scrape", handlePost)
	e.Logger.Fatal(e.Start(":1323"))
}
