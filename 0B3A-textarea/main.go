package main

import (
	"fmt"
	"image"
	"log"
	"net/http"

	"github.com/maxence-charriere/go-app/v9/pkg/app"

	"github.com/mlctrez/imgtofactbp/components/clipboard"
	"github.com/mlctrez/imgtofactbp/conversions"
)

const ImageRenderWidth = 300

// appControl is a component that displays a simple "Hello World!". A component is a
// customizable, independent, and reusable UI element. It is created by
// embedding app.Compo into a struct.
type appControl struct {
	app.Compo
	textStr   string
	clipboard *clipboard.Clipboard
	original  image.Image
	scaled    image.Image
	grayscale image.Image
	inverted  bool

	//threshold      *slider.Continuous
	thresholdValue uint32
}

func (uc *appControl) OnMount(ctx app.Context) {
	uc.clipboard.HandlePaste(ctx, "image/", uc.imagePaste)
}

// The Render method is where the component appearance is defined. Here, a
// "Hello World!" is displayed as a heading.
func (uc *appControl) Render() app.UI {
	if uc.clipboard == nil {
		uc.clipboard = &clipboard.Clipboard{ID: "clipboard"}
	}
	return app.Div().Body(
		uc.clipboard,
		app.If(uc.textStr != "",
			app.Textarea().Text(uc.textStr).Cols(80).ReadOnly(true),
		).Else(
			app.P().Body(
				app.Textarea().
					Text(uc.textStr).
					Spellcheck(true).
					Style("border", "solid 1px orange;").
					Placeholder("Paste your text").
					AutoFocus(true).
					OnChange(uc.ValueTo(&uc.textStr)),
			),
		),
		uc.imagesRow(),
	)
}

func (uc *appControl) imagesRow() app.HTMLDiv {
	return app.Div().Style("display", "flex").Body(
		app.Img().ID("uploadedImage").Src("/web/logo-512.png").Width(ImageRenderWidth).
			Style("cursor", "pointer"),
	)
}

func (uc *appControl) imagePaste(data *clipboard.PasteData) {
	pastedImage, _, err := conversions.Base64ToImage(data.Data)
	if err != nil {
		fmt.Println(err)
		return
	}
	uc.original = pastedImage
	uc.renderImages()
}

func (uc *appControl) renderImages() {
	// normalize image width to 400px
	//uc.scaled = resize.Resize(ImageRenderWidth, 0, uc.original, resize.Lanczos3)
	setImageSrc("uploadedImage", conversions.ImageToBase64(uc.original))
}

func setImageSrc(id string, src string) {
	app.Window().GetElementByID(id).Set("src", src)
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
		Name:        "Hello",
		Description: "An Hello World! example",
	})

	log.Println("Listening on http://:8000")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
