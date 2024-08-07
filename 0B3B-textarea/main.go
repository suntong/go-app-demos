package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// appControl is a component that displays a simple text area. A component is a
// customizable, independent, and reusable UI element. It is created by
// embedding app.Compo into a struct.
type appControl struct {
	app.Compo
	Content string
	Image   string
}

// PasteData is the type and data pasted into the clipboard
type PasteData struct {
	Type string
	Data string
}

func (uc *appControl) OnMount(ctx app.Context) {
	app.Log("network status: mount - online")
}

// The Render method is where the component appearance is defined. Here, a
// "Hello World!" is displayed as a heading.
func (uc *appControl) Render() app.UI {
	return app.Div().Body(
		app.Textarea().
			Placeholder("Paste text or image here...").
			AutoFocus(true).
			OnPaste(uc.OnPaste).
			Style("width", "100%").
			Style("height", "200px"),
		app.If(uc.Content != "", app.Div().Text(uc.Content)),
		app.If(uc.Image != "", app.Img().Src(uc.Image).Style("max-width", "100%")),
	)
}

func (uc *appControl) OnPaste(ctx app.Context, e app.Event) {
	e.PreventDefault()
	items := e.Get("clipboardData").Get("items")

	for i := 0; i < items.Length(); i++ {
		item := items.Call("item", i)
		clipboardItem, err := uc.readPasteData(item)
		if err == ProtectedData {
			continue
		}
		if err != nil {
			fmt.Println("error reading clipboard item", i, err)
			continue
		}
		ctx.NewActionWithValue("Clipboard:paste", clipboardItem, app.Tag{Name: "id", Value: c.ID})
	}
}

func (uc *appControl) OnTextPaste(ctx app.Context, value app.Value) {
	uc.Content = value.String()
}

var ProtectedData = errors.New("protected data")

func (uc *appControl) readPasteData(item app.Value) (result *PasteData, err error) {
	result = &PasteData{Type: item.Get("type").String()}
	switch item.Get("kind").String() {
	case "string":
		done := make(chan bool)
		item.Call("getAsString", app.FuncOf(func(this app.Value, args []app.Value) interface{} {
			result.Data = args[0].String()
			done <- true
			return nil
		}))
		<-done
	case "file":
		done := make(chan bool)
		reader := app.Window().Get("FileReader").New()
		reader.Set("onloadend", app.FuncOf(func(this app.Value, args []app.Value) interface{} {
			done <- true
			return nil
		}))
		reader.Call("readAsDataURL", item.Call("getAsFile"))
		<-done
		if reader.Get("error").Truthy() {
			// TODO: extract real error from reader.Get("error")
			return nil, errors.New("error reading clipboard data")
		}
		result.Data = reader.Get("result").String()
	default:
		return nil, ProtectedData
	}
	return result, nil
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
		Name:        "Go-App Paste Example",
		Description: "A simple app demonstrating paste functionality.",
	})

	log.Println("Listening on http://:8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
