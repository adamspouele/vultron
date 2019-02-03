package console

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/adamspouele/vultron/container"
	"github.com/adamspouele/vultron/reader"

	"github.com/urfave/cli"
)

func Handle() {
	app := cli.NewApp()
	app.Name = "Vultron"
	app.Usage = "Vultron is a container orchestrator based on docker."

	myFlags := []cli.Flag{
		cli.StringFlag{
			Name:  "host",
			Value: "meshectares.com",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "images",
			Usage: "List all images of the current host",
			Action: func(c *cli.Context) error {
				images := container.ListAllImages()

				for _, image := range images {
					if len(image.RepoTags) > 0 {
						fmt.Println("Tag: " + image.RepoTags[0] + " - ID: " + image.ID)
					} else {
						fmt.Println("No Tag, ID: " + image.ID)
					}
				}

				return nil
			},
		},
		{
			Name:  "blueprint",
			Usage: "config file. Default : blueprint.yml",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "file",
					Value: "blueprint.yml",
				},
			},
			Action: func(c *cli.Context) error {
				ip, err := net.LookupIP(c.String("host"))
				if err != nil {
					fmt.Println(err)
					return err
				}
				for i := 0; i < len(ip); i++ {
					fmt.Println(ip[i])
				}

				//cfg, err := parseConfigFile(*config)

				return nil
			},
		},
		{
			Name:  "pidPath",
			Usage: "Looks up the MX records for a Particular Host",
			Flags: myFlags,
			Action: func(c *cli.Context) error {
				mx, err := net.LookupMX(c.String("host"))
				if err != nil {
					fmt.Println(err)
					return err
				}
				for i := 0; i < len(mx); i++ {
					fmt.Println(mx[i].Host, mx[i].Pref)
				}
				return nil
			},
		},
		{
			Name:  "service",
			Usage: "Looks up the MX records for a Particular Host",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "file",
					Value: "blueprint.yml",
				},
			},
			Action: func(c *cli.Context) error {
				fmt.Println(reader.ReadBlueprintFromPath(c.String("file")))
				return nil
			},
		},
		{
			Name:  "start",
			Usage: "Looks up the MX records for a Particular Host",
			Flags: myFlags,
			Action: func(c *cli.Context) error {

				return nil
			},
		},
		{
			Name:  "stop",
			Usage: "Looks up the MX records for a Particular Host",
			Flags: myFlags,
			Action: func(c *cli.Context) error {
				mx, err := net.LookupMX(c.String("host"))
				if err != nil {
					fmt.Println(err)
					return err
				}
				for i := 0; i < len(mx); i++ {
					fmt.Println(mx[i].Host, mx[i].Pref)
				}
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
