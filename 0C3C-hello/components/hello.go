package components

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type Hello struct {
	app.Compo
	updateAvailable  bool
	isAppInstallable bool
}

func (h *Hello) OnAppUpdate(ctx app.Context) {
	h.updateAvailable = ctx.AppUpdateAvailable() // Reports that an app update is available.
}

func (h *Hello) OnMount(ctx app.Context) {
	h.isAppInstallable = ctx.IsAppInstallable()
}

func (h *Hello) onInstallButtonClicked(ctx app.Context, e app.Event) {
	ctx.ShowAppInstallPrompt()
}

func (h *Hello) Render() app.UI {
	if app.IsServer {
		// this gets called on the server before the page is delivered
		return app.Div().Text("app is loading")
	}
	return app.Div().
		Body(
			func() app.UI {
				if h.updateAvailable {
					return app.H1().
						Style("text-align", "center").
						Text("Update available, please reload.")
				}
				return app.Div()
			}(),
			app.If(h.isAppInstallable,
				app.Button().
					Text("Install App").
					OnClick(h.onInstallButtonClicked),
			),
			app.Div().Body(&HelloUI{}),
		)
}

type HelloUI struct {
	app.Compo
	Name string
}

func (h *HelloUI) Render() app.UI {
	return app.Div().
		Body(
			app.If(h.Name == "",
				app.P().Body(
					app.Input().
						Type("text").
						Value(h.Name).
						Placeholder("What is your name?").
						AutoFocus(true).
						OnChange(h.ValueTo(&h.Name)),
				),
			).Else(
				app.H1().Body(
					app.Text("Hello, "),
					app.Text(h.Name),
				),
			))
}
