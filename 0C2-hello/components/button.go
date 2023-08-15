package components

import (
	"fmt"
	"time"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type myButton struct {
	app.Compo
	name string
}

func (c *myButton) Render() app.UI {
	return app.Button().OnClick(c.onClick).Text(c.name)
}

func (c *myButton) OnMount(ctx app.Context) {
	fmt.Println("component mounted")
	c.name = "click here"
}

func (c *myButton) onClick(ctx app.Context, e app.Event) {
	fmt.Println("onClick is called")
	c.name = "clicked"
	ctx.After(2*time.Second, c.revertText)
}

func (c *myButton) revertText(ctx app.Context) {
	c.name = "click here"
}
