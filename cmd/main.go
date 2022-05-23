package main

import (
	"fmt"
	"log"
	"math"
	"net"
	"orders/src/server"
	"orders/src/storage"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "Orders",
		Usage: "",
		Commands: []*cli.Command{
			cmdServer(),
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Println("ERROR:", err)
		os.Exit(1)
	}
}

func cmdServer() *cli.Command {
	var (
		restPort        uint
		dbConnectionURL string
	)
	return &cli.Command{
		Name:     "server",
		Category: "server",
		Usage:    "Start http server to get/add/update/delete users/products/orders",
		Flags: []cli.Flag{
			&cli.UintFlag{
				Name:        "port",
				Usage:       "Defines port for REST API",
				DefaultText: "8080",
				Destination: &restPort,
			},
			&cli.StringFlag{
				Name:        "pg",
				Destination: &dbConnectionURL,
				Usage:       "Defines PostgreSQL connection URL",
			},
		},
		Action: func(c *cli.Context) error {
			if restPort != 0 && restPort >= math.MaxUint16 {
				return fmt.Errorf("REST port number is too big %v", restPort)
			}
			if restPort == 0 {
				restPort = 8080
			}
			ln, err := net.Listen("tcp", fmt.Sprintf(":%v", restPort))
			if err != nil {
				return fmt.Errorf("port %v is busy, set another rest api port", restPort)
			}
			ln.Close()
			if len(dbConnectionURL) == 0 {
				log.Fatal("set db connection URL")
			}
			err = storage.InitSQL(dbConnectionURL)
			if err != nil {
				log.Fatalf("Failed to initialize server storage system %v", err)
			}
			server.Start(restPort)
			return nil
		},
	}
}
