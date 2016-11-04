package main

import (
	"encoding/json"
	"os"

	"github.com/peteclark-io/forecast/ledger"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "forecast"
	app.Usage = "Forecasting for ledger"

	/*flags := []cli.Flag{
		cli.StringFlag{
			Name:  "config",
			Value: "./config.yml",
			Usage: "Path to the YAML config file.",
		},
	}*/

	//app.Flags = flags
	app.Version = "v0.0.1"

	app.Action = func(ctx *cli.Context) error {
		parser := ledger.NewParser(os.Stdin)
		postings, err := parser.Parse()

		if err != nil {
			return err
		}

		encoder := json.NewEncoder(os.Stdout)
		encoder.Encode(postings)

		return nil
	}

	app.Run(os.Args)
}
