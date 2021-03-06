package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/codegangsta/cli"
	"github.com/crackcomm/go-clitable"
	"github.com/mdb/seaweed"
)

func main() {
	app := cli.NewApp()
	app.Name = "seaweed-cli"
	app.Usage = "Should I go surfing?"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "apiKey",
			Usage:  "Magic Seaweed API key",
			EnvVar: "MAGIC_SEAWEED_API_KEY",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:        "forecast",
			Usage:       "forcast <spotId>",
			Description: "View the forecast for a spot",
			Action:      forecast,
		},
		{
			Name:        "today",
			Usage:       "today <spotId>",
			Description: "View today's forecast for a spot",
			Action:      today,
		},
		{
			Name:        "tomorrow",
			Usage:       "tomorrow <spotId>",
			Description: "View tomorrow's forecast for a spot",
			Action:      tomorrow,
		},
	}
	app.Run(os.Args)
}

func forecast(c *cli.Context) {
	client := client(c)
	spot := c.Args().First()
	forecast, err := client.Forecast(spot)
	if err != nil {
		panic(err)
	}

	printForecasts(spot, forecast)
}

func today(c *cli.Context) {
	client := client(c)
	spot := c.Args().First()
	forecast, err := client.Today(spot)
	if err != nil {
		panic(err)
	}

	printForecasts(spot, forecast)
}

func tomorrow(c *cli.Context) {
	client := client(c)
	spot := c.Args().First()
	forecast, err := client.Tomorrow(spot)
	if err != nil {
		panic(err)
	}

	printForecasts(spot, forecast)
}

func printForecasts(spot string, forecasts []seaweed.Forecast) {
	s := []map[string]interface{}{}
	for _, each := range forecasts {
		m := map[string]interface{}{}
		m["Date"] = time.Unix(each.LocalTimestamp, 0).Format("Mon 01/02 03:04 pm")
		m["Solid Rating"] = each.SolidRating
		m["Faded Rating"] = each.FadedRating
		m["Primary Swell Height"] = concat([]string{strconv.FormatFloat(each.Swell.Components.Primary.Height, 'f', 2, 64), each.Swell.Unit})
		m["Wind"] = concat([]string{strconv.Itoa(each.Wind.Speed), " ", each.Wind.Unit, " ", each.Wind.CompassDirection})
		s = append(s, m)
	}

	if len(s) != 0 {
		clitable.PrintTable([]string{"Date", "Primary Swell Height", "Wind", "Solid Rating", "Faded Rating"}, s)
	} else {
		fmt.Printf("No forecast found for spot: %s\n", spot)
	}
}

func client(c *cli.Context) *seaweed.Client {
	return seaweed.NewClient(os.Getenv("MAGIC_SEAWEED_API_KEY"))
}

func concat(arr []string) string {
	var buff bytes.Buffer

	for _, elem := range arr {
		buff.WriteString(elem)
	}

	return buff.String()
}
