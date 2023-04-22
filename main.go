package main

import (
	"fyne.io/fyne/v2/app"
)

type WrappedError struct {
	StatusCode int
	Err        error
}

type AwsFilters struct {
	TagName         string
	PrivateIpFilter string
}

func main() {
	a := app.New()
	pf := new(PortForward)
	pf.window = a.NewWindow("BL Port Forwarding")
	pf.portForwardForm()
	pf.window.ShowAndRun()
}
