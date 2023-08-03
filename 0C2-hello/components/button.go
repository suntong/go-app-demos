package components

import (
	"fmt"

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
	//c.Update()
}
