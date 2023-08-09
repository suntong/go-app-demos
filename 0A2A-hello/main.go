package main

import (
	"log"
	"net/http"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// hello is a component that displays a simple "Hello World!". A component is a
// customizable, independent, and reusable UI element. It is created by
// embedding app.Compo into a struct.
type hello struct {
	app.Compo
	name string
}

// The Render method is where the component appearance is defined. Here, a
// "Hello World!" is displayed as a heading.
func (h *hello) Render() app.UI {
	return app.Div().Body(
		app.H1().Body(
			app.Text("Hello, "),
			app.If(h.name != "",
				app.Text(h.name),
			).Else(
				app.Text("World!"),
			),
		),
		app.P().Body(
			app.Input().
				Type("text").
				Value(h.name).
				// OK: Name("name").ID("frmNameA").Attr("autocomplete", "name").
				// OK: Name("email").ID("frmEmailA").Attr("autocomplete", "email").
				// OK: Name("tel").ID("frmTelA").Attr("autocomplete", "tel").
				// OK: Name("ra").Attr("autocomplete", "billing street-address").
				// NOK: Name("ra").Attr("autocomplete", "some other string").
				// NOK: Name("ra").Attr("autocomplete", "ContactID").
				// NOK: Name("ra").Attr("autocomplete", "section-blue billing nickname").
				// NOK: Name("ra").Attr("autocomplete", "billing nickname").
				// NOK: Name("ra").Attr("autocomplete", "billing address-level4").
				// NOK: Name("ra").Attr("autocomplete", "billing cc-csc").
				// NOK: Name("ra").Attr("autocomplete", "one-time-code").
				Name("ra").Attr("autocomplete", "nickname"). // NOK
				Placeholder("What is your name?").
				AutoFocus(true).
				OnChange(h.ValueTo(&h.name)),
		),
	)
}

// The main function is the entry point where the app is configured and started.
// It is executed in 2 different environments: A client (the web browser) and a
// server.
func main() {
	// The first thing to do is to associate the hello component with a path.
	//
	// This is done by calling the Route() function,  which tells go-app what
	// component to display for a given path, on both client and server-side.
	app.Route("/", &hello{})
	app.Route("/hello", &hello{})

	// Once the routes set up, the next thing to do is to either launch the app
	// or the server that serves the app.
	//
	// When executed on the server-side, RunWhenOnBrowser() does nothing, which
	// lets room for server implementation without the need for precompiling
	// instructions.
	app.RunWhenOnBrowser()

	// Finally, launching the server that serves the app is done by using the Go
	// standard HTTP package.
	//
	// The Handler is an HTTP handler that serves the client and all its
	// required resources to make it work into a web browser. Here it is
	// configured to handle requests with a path that starts with "/".
	http.Handle("/", &app.Handler{
		Name:        "Hello",
		Description: "An Hello World! example",
	})

	log.Println("Listening on http://:8000")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}