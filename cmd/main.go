package main

import (
	"context"

	"github.com/mttchrry/oxio-phone-lookup/pkg/app"
	httppkg "github.com/mttchrry/oxio-phone-lookup/pkg/http"
	"github.com/mttchrry/oxio-phone-lookup/pkg/phoneNumbers"
)

func main() {
	app.Start(appStart)
}

func appStart(ctx context.Context, a *app.App) ([]app.Listener, error) {
	p := phoneNumbers.New()

	h, err := httppkg.New(p, "8081")
	if err != nil {
		return nil, err
	}

	// Start listening for HTTP requests
	return []app.Listener{
		h,
	}, nil
}
