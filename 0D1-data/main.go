package main

import (
	"log"
	"net/http"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

const dataStorageKey = "0D1-data.Name"

// appControl is a component that displays a simple "Hello World!". A component is a
// customizable, independent, and reusable UI element. It is created by
// embedding app.Compo into a struct.
type appControl struct {
	app.Compo
	name string

	removeEventListeners []func()
}

func (uc *appControl) OnMount(ctx app.Context) {
	uc.removeEventListeners = []func(){
		app.Window().AddEventListener("storage", func(ctx app.Context, e app.Event) { // This event only fires in other tabs; it does not lead to local race conditions with c.writeKeysToLocalStorage
			uc.readFromLocalStorage()
			uc.Update()
		}),
	}
}

func (uc *appControl) OnDismount() {
	if uc.removeEventListeners != nil {
		for _, clearListener := range uc.removeEventListeners {
			clearListener()
		}
	}
}

// The Render method is where the component appearance is defined. Here, a
// "Hello World!" is displayed as a heading.
func (uc *appControl) Render() app.UI {
	if uc.name == "" {
		uc.readFromLocalStorage()
		if uc.name == "<null>" {
			uc.name = ""
		}
	}
	return app.Div().Body(
		app.H1().Body(
			app.Text("Hello, "),
			app.If(uc.name != "",
				app.Text(uc.name),
			).Else(
				app.Text("World!"),
			),
		),
		app.P().Body(
			app.Input().
				Type("text").
				Value(uc.name).
				Placeholder("What is your name?").
				AutoFocus(true).
				//OnChange(uc.ValueTo(&uc.name)),
				OnChange(uc.OnChange),
		),
	)
}

func (uc *appControl) OnChange(ctx app.Context, e app.Event) {
	uc.name = ctx.JSSrc().Get("value").String()
	app.Window().Get("localStorage").Call("setItem", dataStorageKey, uc.name)
	uc.name = ""
	uc.readFromLocalStorage()
}

func (uc *appControl) readFromLocalStorage() {
	uc.name = app.Window().Get("localStorage").Call("getItem", dataStorageKey).String()
	log.Println("readFromLocalStorage:", uc.name)
}

// The main function is the entry point where the app is configured and started.
// It is executed in 2 different environments: A client (the web browser) and a
// server.
func main() {
	// The first thing to do is to associate the appControl component with a path.
	//
	// This is done by calling the Route() function,  which tells go-app what
	// component to display for a given path, on both client and server-side.
	app.Route("/", &appControl{})
	app.Route("/appControl", &appControl{})

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
