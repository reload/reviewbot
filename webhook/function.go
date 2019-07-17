package function

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/wm/go-flowdock/flowdock"
	"gopkg.in/go-playground/webhooks.v5/github"
)

// Handle is the entrypoint for the Google Cloud Function.
func Handle(w http.ResponseWriter, r *http.Request) {
	// We log to Google Cloud Functions and don't need a timestamp
	// since it will be present in the log anyway. On the other
	// hand a reference to file and line number would be nice.
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
	payload, err := hook.Parse(r, github.PullRequestEvent)

	if err == github.ErrMissingHubSignatureHeader {
		http.Error(w, fmt.Sprintf("%s: %s", http.StatusText(http.StatusUnauthorized), err), http.StatusUnauthorized)

		return
	}

	if err != nil {
		http.Error(w, fmt.Sprintf("%s: %s", http.StatusText(http.StatusBadRequest), err), http.StatusBadRequest)

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

	str := fmt.Sprintf(
		"Review requested by `%s`:\n\n * %s#%d: [**%s**](%s) by `%s`",
		pullRequest.Sender.Login,
		pullRequest.Repository.FullName,
		pullRequest.Number,
		pullRequest.PullRequest.Title,
		pullRequest.PullRequest.HTMLURL,
		pullRequest.PullRequest.User.Login,
	)

	err = flowpost(str)

	if err != nil {
		log.Printf("Could not post to Flowdock: %s", err)
		http.Error(w, fmt.Sprintf("Could not post to Flowdock: %s", err), http.StatusInternalServerError)

		return
	}

	http.Error(w, str, http.StatusOK)

	return
}

func flowpost(msg string) error {
	flowname, ok := os.LookupEnv("FLOWDOCK_FLOW_NAME")

	if !ok {
		return fmt.Errorf("No Flowdock flow name configured (environment variable FLOWDOCK_FLOW_NAME)")
	}

	flowdockToken, ok := os.LookupEnv("FLOWDOCK_TOKEN")

	if !ok {
		return fmt.Errorf("No Flowdock token configured (environment variable FLOWDOCK_TOKEN)")
	}

	client := flowdock.NewClientWithToken(nil, flowdockToken)

	flows, _, err := client.Flows.List(true, &flowdock.FlowsListOptions{User: false})

	if err != nil {
		return err
	}

	flowID := ""
	for _, f := range flows {
		if *f.ParameterizedName == flowname {
			flowID = *f.Id
		}
	}

	if flowID == "" {
		return fmt.Errorf("Could not find flow named %s", flowname)
	}

	_, _, err = client.Messages.Create(&flowdock.MessagesCreateOptions{
		Event:            "message",
		FlowID:           flowID,
		Content:          msg,
		ExternalUserName: "ReviewBot",
	})

	return err
}
