package function

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/containrrr/shoutrrr"
	"github.com/containrrr/shoutrrr/pkg/types"
	"github.com/rickar/cal"
	"gopkg.in/go-playground/webhooks.v5/github"
)

// Handle is the entrypoint for the Google Cloud Function.
func Handle(w http.ResponseWriter, r *http.Request) {
	log.SetFlags(log.Lshortfile)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	githubSecret, ok := os.LookupEnv("GITHUB_SECRET")
	if !ok {
		log.Printf("No GitHub secret defined (environment variable GITHUB_SECRET)")
		http.Error(w, "No GitHub secret defined (environment variable GITHUB_SECRET)", http.StatusInternalServerError)

		return
	}

	teamSlug, ok := os.LookupEnv("GITHUB_TEAM_SLUG")
	if !ok {
		log.Printf("No GitHub team ID defined (environment variable GITHUB_TEAM_ID)")
		http.Error(w, "No GitHub team ID defined (environment variable GITHUB_TEAM_ID)", http.StatusInternalServerError)

		return
	}

	hook, _ := github.New(github.Options.Secret(githubSecret))
	payload, err := hook.Parse(r, github.PullRequestEvent, github.PingEvent)

	if errors.Is(err, github.ErrMissingHubSignatureHeader) {
		http.Error(w, fmt.Sprintf("%s: %s", http.StatusText(http.StatusUnauthorized), err), http.StatusUnauthorized)

		return
	}

	if err != nil {
		http.Error(w, fmt.Sprintf("%s: %s", http.StatusText(http.StatusBadRequest), err), http.StatusBadRequest)

		return
	}

	if _, ok := payload.(github.PingPayload); ok {
		http.Error(w, "Pong", http.StatusOK)

		return
	}

	pullRequest, ok := payload.(github.PullRequestPayload)
	if !ok {
		http.Error(w, fmt.Sprintf("Could not parse as pull request payload: %#v", payload), http.StatusBadRequest)

		return
	}

	if pullRequest.Action != "review_requested" {
		http.Error(w, "Not a review request", http.StatusOK)

		return
	}

	if pullRequest.RequestedTeam.Slug != teamSlug {
		http.Error(w, fmt.Sprintf("Not requesting review from %s", teamSlug), http.StatusOK)

		return
	}

	sendMessage(w, pullRequest)
}

func sendMessage(w http.ResponseWriter, pullRequest github.PullRequestPayload) {
	str := fmt.Sprintf(
		"Review requested by `%s`:\n\n * %s#%d: **[%s](%s)** by `%s`",
		pullRequest.Sender.Login,
		pullRequest.Repository.FullName,
		pullRequest.Number,
		pullRequest.PullRequest.Title,
		pullRequest.PullRequest.HTMLURL,
		pullRequest.PullRequest.User.Login,
	)

	if withinWorkingHours() {
		str = fmt.Sprintf("@**all**, %s", str)
	}

	err := send(str, pullRequest.Repository.FullName)
	if err != nil {
		log.Printf("Could not post message: %s", err)
		http.Error(w, fmt.Sprintf("Could not post message: %s", err), http.StatusInternalServerError)

		return
	}

	http.Error(w, str, http.StatusOK)
}

func send(message string, topic string) error {
	services := os.Getenv("NOTIFY")
	notify, err := shoutrrr.CreateSender(strings.Split(services, ",")...)

	if err != nil {
		return fmt.Errorf("error creating notification sender(s): %w", err)
	}

	params := types.Params{
		"topic": topic,
	}

	errs := notify.Send(message, &params)

	if len(errs) > 0 {
		return fmt.Errorf("error creating notification sender(s): %v", errs) //nolint:goerr113
	}

	return nil
}

func withinWorkingHours() bool {
	c := workCalendar()
	now := time.Now()

	if !c.IsWorkday(now) {
		return false
	}

	if (now.Hour() < 9) || (now.Hour() > 16) {
		return false
	}

	return true
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
