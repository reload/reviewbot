module github.com/reload/reviewbot/webhook

// Don't bump above 1.21 - it's unsupported by Google Cloud Functions gen 1
go 1.21.13

require (
	github.com/containrrr/shoutrrr v0.8.0
	github.com/rickar/cal/v2 v2.1.15
	gopkg.in/go-playground/webhooks.v5 v5.17.0
)

require (
	github.com/fatih/color v1.16.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/stretchr/testify v1.8.4 // indirect
	golang.org/x/sys v0.17.0 // indirect
)
