package components

import (
	"fmt"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type FormKind int

const (
	Login FormKind = iota
	Register
	Recover
)

type LoginForm struct {
	app.Compo
	username string
	password string
	fType    FormKind
}

func (l *LoginForm) setUsername(ctx app.Context, e app.Event) {
	l.username = ctx.JSSrc().Get("value").String()
}

func (l *LoginForm) setPassword(ctx app.Context, e app.Event) {
	l.password = ctx.JSSrc().Get("value").String()
}

func (l *LoginForm) btnClickAnimation(btn app.Value, e app.Event) {
	// simple click animation (ease-in-out)
	btn.Call("animate", []interface{}{
		map[string]interface{}{
			"transform": "scale(0.9)",
		},
		map[string]interface{}{
			"transform": "scale(1)",
		},
	}, map[string]interface{}{
		"duration": 200,
		"easing":   "ease-in-out",
	})
}

func (l *LoginForm) btnClick(ctx app.Context, e app.Event) {
	btn := e.Get("currentTarget").JSValue()
	ctx.Async(
		func() {
			l.btnClickAnimation(btn, e)
		},
	)
}

func (l *LoginForm) inputFocus(ctx app.Context, e app.Event) {
	// remove placeholder
	ctx.JSSrc().Set("placeholder", "")
}

func (l *LoginForm) inputBlur(ctx app.Context, e app.Event) {
	// if value is empty, set placeholder
	if ctx.JSSrc().Get("value").String() == "" {
		// capitalize first letter of placeholder
		placeholder := ctx.JSSrc().Get("name").String()
		placeholder = string(placeholder[0]-32) + placeholder[1:]
		ctx.JSSrc().Set("placeholder", placeholder)
	}
}

func (l *LoginForm) OnInit(ctx app.Context, e app.Event) {
	fmt.Println("type", l.fType)
}

func (l *LoginForm) OnMount(ctx app.Context, e app.Event) {
	fmt.Println("type", l.fType)
}

func (l *LoginForm) Render() app.UI {
	switch l.fType {
	case Login:
		return app.Div().Class("fill").Body(
			app.Div().Class("login-form").Body(
				app.H4().Class("login-title").Text("Login"),
				app.Input().Type("text").Placeholder("Username").Name("username").Value(l.username).OnChange(l.setUsername).OnFocus(l.inputFocus).OnBlur(l.inputBlur),
				app.Input().Type("password").Placeholder("Password").Name("password").Value(l.password).OnChange(l.setPassword).OnFocus(l.inputFocus).OnBlur(l.inputBlur),
				app.Div().Class("login-button-container").Body(
					app.Button().Class("login-button").Text("Login").OnClick(l.btnClick),
					app.Button().Class("login-button").Text("Register").OnClick(l.btnClick),
				),
			),
		)
	case Register:
		return app.Div().Class("fill").Body(
			app.Div().Class("login-form").Body(
				app.H4().Class("login-title").Text("Register"),
				app.Input().Type("text").Placeholder("Username").Name("username").Value(l.username).OnChange(l.setUsername).OnFocus(l.inputFocus).OnBlur(l.inputBlur),
				app.Input().Type("password").Placeholder("Password").Name("password").Value(l.password).OnChange(l.setPassword).OnFocus(l.inputFocus).OnBlur(l.inputBlur),
				app.Input().Type("password").Placeholder("Confirm Password").Name("confirm-password").Value(l.password).OnChange(l.setPassword).OnFocus(l.inputFocus).OnBlur(l.inputBlur),
				app.Div().Class("login-button-container").Body(
					app.Button().Class("login-button").Text("Register").OnClick(l.btnClick),
					app.Button().Class("login-button").Text("Login").OnClick(l.btnClick),
				),
			),
		)
	case Recover:
		return app.Div().Class("fill").Body(
			app.Div().Class("login-form").Body(
				app.H4().Class("login-title").Text("Recover"),
				app.Input().Type("text").Placeholder("Username").Name("username").Value(l.username).OnChange(l.setUsername).OnFocus(l.inputFocus).OnBlur(l.inputBlur),
				app.Div().Class("login-button-container").Body(
					app.Button().Class("login-button").Text("Recover").OnClick(l.btnClick),
					app.Button().Class("login-button").Text("Login").OnClick(l.btnClick),
				),
			),
		)
	default:
		return app.Div().Class("fill").Body(
			app.Span().Text("Invalid form kind"),
		)
	}
}

func NewLoginForm(k FormKind) *LoginForm {
	return &LoginForm{
		fType: k,
	}
}
