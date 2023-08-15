package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

////////////////////////////////////////
// Credit:
// https://github.com/maxence-charriere/go-app/issues/859#issuecomment-1677198131
// https://github.com/maxence-charriere/go-app/issues/872#issuecomment-1677725579
////////////////////////////////////////

// Define component, a customizable, independent, and reusable UI
// element. It is created by embedding app.Compo into a struct.
type codeBlockModel struct {
	app.Compo
	code []string
}

func (m *codeBlockModel) OnInit() {
	m.code = []string{`
                    // Code block 1
                    function helloWorld() {
                        console.log("Hello, world!");
                    }
`,
		`
                    // Code block 2
                    for (let i = 0; i < 5; i++) {
                        console.log(i);
                    }
`}
}

// The Render method is where the component appearance is defined.
func (m *codeBlockModel) Render() app.UI {
	return app.Div().Body(
		app.H1().Text("H1"),
		app.H4().Text("H4"),

		app.Div().Class("code-container").Body(
			app.Range(m.code).Slice(func(i int) app.UI {
				id := len(m.code) - 1 - i
				//id = i
				return &CodeBlock{code: m.code[id], id: fmt.Sprintf("codeBlock%02d", id)}
			}),
		),
	)
}

type CodeBlock struct {
	app.Compo
	id   string
	code string
}

func (m *CodeBlock) Render() app.UI {
	return app.Div().Class("code-block").Body(
		copySVG(),
		&CopyButton{text: "Copy code", from: m},
		app.Pre().Body(
			app.Code().ID(m.id).Text(m.code),
		),
	)
}

func copySVG() app.UI {
	return app.Raw(`<svg stroke="currentColor" fill="none" stroke-width="2" viewBox="0 0 24 24" stroke-linecap="round" stroke-linejoin="round" class="copy-svg h-4 w-4" height="1em" width="1em" xmlns="http://www.w3.org/2000/svg"><path d="M16 4h2a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2h2"></path><rect x="8" y="2" width="8" height="4" rx="1" ry="1"></rect></svg>`)
}

type CopyButton struct {
	app.Compo
	from *CodeBlock
	text string
}

func (cb *CopyButton) Render() app.UI {
	return app.Button().Class("copy-button").Text(cb.text).OnClick(cb.onClick)
}

func (cb *CopyButton) onClick(ctx app.Context, e app.Event) {
	// using go
	log.Println("copy via go\n", cb.from.code)

	// using element id
	val := app.Window().GetElementByID(cb.from.id).Get("innerHTML")
	log.Println("copy via dom\n", val.String())

	isSecure := app.Window().Get("isSecureContext")
	log.Println("isSecureContext:", isSecure)
	if isSecure.Bool() {
		app.Window().Get("navigator").Get("clipboard").Call("writeText", val)
	} else {
		copyToClipboard(cb.from.code)
	}
	cb.text = "Copied"
	ctx.After(2*time.Second, cb.revertText)
}

func (cb *CopyButton) revertText(ctx app.Context) {
	cb.text = "Copy code"
}

func copyToClipboard(text string) {
	//app.Log("Copying to clipboard: %q", text)
	app.Window().Call("copyToClipboard", text)
}

// The main function is the entry point where the app is configured and started.
// It is executed in 2 different environments: A client (the web browser) and a
// server.
func main() {
	// The first thing to do is to associate the component with a path.
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
		RawHeaders: []string{
			`<script>
    function copyToClipboard(text) {
        if (window.clipboardData && window.clipboardData.setData) {
            // Internet Explorer-specific code path to prevent textarea being shown while dialog is visible.
            return window.clipboardData.setData("Text", text);

        } else if (document.queryCommandSupported && document.queryCommandSupported("copy")) {
            // fallback for old browsers (probably not needed)
            var textarea = document.createElement("textarea");
            textarea.textContent = text;
            textarea.style.position = "fixed";  // Prevent scrolling to bottom of page in Microsoft Edge.
            document.body.appendChild(textarea);
            textarea.select();
            try {
                return document.execCommand("copy");  // Security exception may be thrown by some browsers.
            } catch (ex) {
                console.warn("Copy to clipboard failed.", ex);
                return prompt("Copy to clipboard: Ctrl+C, Enter", text);
            } finally {
                document.body.removeChild(textarea);
            }
        }
    }
</script>
`},
	})

	log.Println("Listening on http://:8000")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
