package main

import (
	"fmt"
	"os"
	"strings"
	"time"
	_ "time/tzdata"

	"github.com/containrrr/shoutrrr"
	"github.com/containrrr/shoutrrr/pkg/types"
	"github.com/rickar/cal/v2"
	"github.com/rickar/cal/v2/dk"
	"github.com/robfig/cron"
	_ "golang.org/x/crypto/x509roots/fallback"
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

	err = send(str)

	if err != nil {
		panic(err)
	}
}

func workCalendar() *cal.BusinessCalendar {
	c := cal.NewBusinessCalendar()

	c.AddHoliday(dk.Holidays...)

	//nolint:exhaustivestruct
	c.AddHoliday(&cal.Holiday{
		Month: time.December,
		Day:   24, //nolint:gomnd
		Func:  cal.CalcDayOfMonth,
	})

	//nolint:exhaustivestruct
	c.AddHoliday(&cal.Holiday{
		Month: time.December,
		Day:   31, //nolint:gomnd
		Func:  cal.CalcDayOfMonth,
	})

	return c
}

func send(message string) error {
	services := os.Getenv("NOTIFY")

	notify, err := shoutrrr.CreateSender(strings.Split(services, ",")...)
	if err != nil {
		return fmt.Errorf("error creating notification sender(s): %w", err)
	}

	t := time.Now()
	params := types.Params{
		"topic": fmt.Sprintf("%d-%02d-%02d", t.Year(), t.Month(), t.Day()),
	}

	errs := notify.Send(message, &params)

	if len(errs) > 0 {
		return fmt.Errorf("error creating notification sender(s): %v", errs) //nolint:goerr113
	}

	return nil
}
