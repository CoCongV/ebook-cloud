package main

import (
	"log"
	"os"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/urfave/cli/v2"

	"ebook-cloud/api/apiv1"
	"ebook-cloud/client"
	"ebook-cloud/config"
	"ebook-cloud/models"
	"ebook-cloud/render"
	"ebook-cloud/server"
	"ebook-cloud/view"
)

var confPath string
var confFlag = &cli.StringFlag{
	Name:        "conf",
	Value:       "",
	Usage:       "config file path",
	Destination: &confPath,
}

func init() {
	config.Setup()
	models.Setup()
	client.Setup()
}

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:   "runserver",
				Usage:  "run server",
				Action: runserver,
			},
			{
				Name:   "migrate",
				Usage:  "migrate models",
				Action: migrate,
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func runserver(c *cli.Context) error {
	defer models.DB.Close()
	r := server.CreateServ()
	apiv1.SetRouter(r)
	view.SetRouter(r)
	r.HTMLRender = render.New("static/template")
	r.Static("/static", "static")
	r.Run(config.Conf.Addr)
	return nil
}

func migrate(c *cli.Context) error {
	defer models.DB.Close()

	models.DB.AutoMigrate(&models.Book{}, &models.Author{}, &models.Country{})

	return nil
}
