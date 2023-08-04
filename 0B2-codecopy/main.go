package main

import (
	"log"
	"net/http"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// Define component, a customizable, independent, and reusable UI
// element. It is created by embedding app.Compo into a struct.
type codeBlockModel struct {
	app.Compo
	code string
}

func (m *codeBlockModel) OnInit() {
	m.code = `
                    // Code block 1
                    function helloWorld() {
                        console.log("Hello, world!");
                    }
`
}

// The Render method is where the component appearance is defined.
func (m *codeBlockModel) Render() app.UI {
	return app.Div().Class("code-container").Body(
		app.Div().Class("code-block").Body(
			app.Raw(`<svg class="copy-svg" stroke="currentColor" fill="none" stroke-width="2" viewBox="0 0 24 24" stroke-linecap="round" stroke-linejoin="round" class="h-4 w-4" height="1em" width="1em" xmlns="http://www.w3.org/2000/svg"><path d="M16 4h2a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2h2"></path><rect x="8" y="2" width="8" height="4" rx="1" ry="1"></rect></svg><button class="copy-button">Copy code</button>`),
			app.Raw(`<button class="copy-button">Copy code</button>`),
			app.Pre().Body(
				app.Code().Text(m.code),
			),
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
	app.Route("/", &codeBlockModel{})

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
		Title:   "Code Copy Example",
		Author:  "Suntown Studio",
		Styles:  []string{"/web/styles.css"},
		Scripts: []string{"/web/script.js"},
		Icon: app.Icon{
			Default:    "/web/copy-icon.png",
			Large:      "/web/copy-icon.png",
			AppleTouch: "/web/copy-icon.png",
		},
	})

	log.Println("Listening on http://:8000")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
