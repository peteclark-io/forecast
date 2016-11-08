package main

import (
	"encoding/json"
	"os"
	"regexp"

	"github.com/peteclark-io/forecast/forecasting"
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

	app.Commands = []cli.Command{
		{
			Name:    "balance",
			Aliases: []string{"b", "ba", "bal"},
			Usage:   "Balance an account",
			Action: func(c *cli.Context) error {
				parser := ledger.NewParser(os.Stdin)
				postings, err := parser.Parse()

				if err != nil {
					return err
				}

				val, err := forecasting.Balance("Assets:Current:HSBC", postings)
				if err != nil {
					return err
				}

				encoder := json.NewEncoder(os.Stdout)
				encoder.Encode(val)

				return nil
			},
		},
		{
			Name:    "forecast",
			Aliases: []string{"f", "fo", "for"},
			Usage:   "Forecast",
			Action: func(c *cli.Context) error {
				parser := ledger.NewParser(os.Stdin)
				postings, err := parser.Parse()

				if err != nil {
					return err
				}

				filter, err := regexp.Compile(c.Args().First())
				if err != nil {
					return err
				}

				data := forecasting.AverageByDay(filter, postings)
				w := json.NewEncoder(os.Stdout)
				w.SetIndent("", "  ")
				w.Encode(data)

				return nil
			},
		},
	}

	app.Run(os.Args)
}
