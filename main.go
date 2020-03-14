package main

import (
	"ShortLink/app"
)

func main() {
	a := app.App{}
	a.Initialize(app.GetEnv())
	a.Run(":8000")
}
