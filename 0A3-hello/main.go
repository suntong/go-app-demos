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

// == Lifecycle Events https://go-app.dev/components

// PRERENDER
// A component is prerendered when it is used on the server-side to generate HTML markup that is included in a requested HTML page, allowing search engines to index contents created with go-app.
func (h *hello) OnPreRender(ctx app.Context) {
	log.Println("component prerendered")
}

// OnInit is called before the component gets mounted
// This is before Render was called for the first time
func (h *hello) OnInit() {
	app.Log("OnInit")
	log.Println("component initiated")
}

// MOUNT
// A component is mounted when it is inserted into the webpage DOM.
func (h *hello) OnMount(ctx app.Context) {
	app.Log("OnMount")
	log.Println("component mounted")
}

// NAV
// A component is navigated when a page is loaded, reloaded, or navigated from an anchor link or an HREF change. It can occur multiple times during a component life.
func (h *hello) OnNav(ctx app.Context) {
	app.Log("OnNav")
	log.Println("component navigated:", h, ctx)
}

// DISMOUNT
// A component is dismounted when it is removed from the webpage DOM.
func (h *hello) OnDismount() {
	app.Log("OnDismount")
	log.Println("component dismounted")
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
