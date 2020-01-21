package main

import (
	"log"
	"os"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/urfave/cli/v2"

	"ebook-cloud/apiv1"
	"ebook-cloud/config"
	"ebook-cloud/models"
	"ebook-cloud/server"
)

var confPath string
var confFlag = &cli.StringFlag{
	Name:        "conf",
	Value:       "",
	Usage:       "config file path",
	Destination: &confPath,
}

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:   "runserver",
				Usage:  "run server",
				Action: runserver,
				Flags: []cli.Flag{
					confFlag,
				},
			},
			{
				Name:   "migrate",
				Usage:  "migrate models",
				Action: migrate,
				Flags: []cli.Flag{
					confFlag,
				},
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func runserver(c *cli.Context) error {

	conf := config.ReadConfig(confPath)

	models.DB = models.InitDB(conf.DBURL)
	defer models.DB.Close()

	r := server.CreateServ()
	apiv1.SetRouter(r)
	r.Run(conf.Addr)
	return nil
}

func migrate(c *cli.Context) error {
	conf := config.ReadConfig(confPath)
	models.DB = models.InitDB(conf.DBURL)
	defer models.DB.Close()

	models.DB.AutoMigrate(&models.Book{}, &models.Author{}, &models.Country{})

	return nil
}
