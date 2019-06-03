package main

import (
	"os"
	"time"

	"github.com/rickar/cal"
	"github.com/robfig/cron"
	"github.com/wm/go-flowdock/flowdock"
)

func main() {
	c := cron.New()
	_ = c.AddFunc("0 45 8,10,12,14 * * *", run)
	c.Start()

	select {}
}

func run() {
	c := workCalendar()
	if !c.IsWorkday(time.Now()) {
		return
	}

	edges, count, err := reviewRequests()

	if err != nil {
		panic(err)
	}

	if count <= 0 {
		return
	}

	str, err := format(edges, count)

	if err != nil {
		panic(err)
	}

	err = flowpost(str)

	if err != nil {
		panic(err)
	}
}

func workCalendar() *cal.Calendar {
	c := cal.NewCalendar()

	cal.AddDanishHolidays(c)
	c.AddHoliday(
		cal.DKJuleaften,
		cal.DKNytaarsaften,
	)

	return c
}

func flowpost(str string) error {
	flowname := os.Getenv("FLOWDOCK_FLOW")

	client := flowdock.NewClientWithToken(nil, os.Getenv("FLOWDOCK_API_TOKEN"))

	flows, _, err := client.Flows.List(true, &flowdock.FlowsListOptions{User: false})

	if err != nil {
		panic(err)
	}

	flowID := ""
	for _, f := range flows {
		if *f.ParameterizedName == flowname {
			flowID = *f.Id
		}
	}

	if flowID == "" {
		panic("Could not find flow.")
	}

	_, _, err = client.Messages.Create(&flowdock.MessagesCreateOptions{
		Event:            "message",
		FlowID:           flowID,
		Tags:             []string{"@team"},
		Content:          str,
		ExternalUserName: "ReviewBot",
	})

	return err
}
