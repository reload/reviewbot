# Reviewbot

[![MIT Licensed](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go](https://img.shields.io/badge/language-Go-blue.svg)](https://go.dev/)

An experimental bot that posts pending team review requests to Zulip,
plus a Google Cloud Function that turns GitHub webhooks into Zulip
messages requesting reviews.

## Overview

**Reviewbot** bridges the gap between GitHub pull requests and your
team’s Zulip chat. Whenever code is ready for review, Reviewbot
notifies your Zulip stream or topic, making it easier for teams to
track PRs awaiting attention.

- **GitHub → Zulip Integration:** Converts webhooks into Zulip
  messages when PRs need review.
- **Cloud Function Support:** Deployable as a Google Cloud Function
  for scalable automation.
- **Team Notifications:** Ensures timely notifications for code review
  requests.

## Features

- Monitors GitHub repositories for new or pending PRs.
- Posts relevant review requests in designated Zulip streams/topics.
- Intended for team environments to reduce review friction.
- Written in Go for performance and cloud-native deployment.

## Requirements

- [Go](https://go.dev/) (for development, building, or local running)
- GitHub repository with webhook permissions
- Zulip account, stream, and API information
- [Google Cloud Functions](https://cloud.google.com/functions/)
  (optional; for serverless deployment)

## Installation

1. **Clone the Repo**

   ```bash
   git clone https://github.com/reload/reviewbot.git
   cd reviewbot
   ```

2. **Configure**
   - Set up your GitHub webhook to point to your Reviewbot endpoint.
   - Provide Zulip bot credentials and stream configuration.

3. **Deploy**
   - _As a Google Cloud Function_: Follow Google’s deployment
     instructions and provide the proper environment variables.
   - _Locally_: Run with `go run main.go` (additional configuration
     may be needed).

## Usage

- When a pull request is created or marked ready for review in your
  GitHub repository:
  - Reviewbot receives the webhook, parses it, and posts a structured
    message in your Zulip stream/topic requesting review from assigned
    users or teams.

## License

MIT © [Reload](https://github.com/reload)

---

_Experimental project. Contributions and feedback welcome!_
