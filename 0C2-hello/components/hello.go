package components

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type AppControl struct {
	app.Compo

	name             string
	isAppInstallable bool
}

func (uc *AppControl) OnMount(ctx app.Context) {
	uc.isAppInstallable = ctx.IsAppInstallable()
}

func (uc *AppControl) onInstallButtonClicked(ctx app.Context, e app.Event) {
	ctx.ShowAppInstallPrompt()
}

func (uc *AppControl) Render() app.UI {
	return app.Div().
		Body(
			app.H1().Body(
				app.Text("AppControl, "),
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
					OnChange(uc.ValueTo(&uc.name)),
			),
			app.Div().Style("margin", "1rem").
				Body(&myButton{}),
			app.If(uc.isAppInstallable,
				app.Button().
					Text("Install App").
					OnClick(uc.onInstallButtonClicked),
			),
		)
}
