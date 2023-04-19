package main

import (
	"log"
	"strconv"
	s "strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type PortForward struct {
	form   map[string]*widget.Entry
	window fyne.Window
}

func (pf *PortForward) loadForm() []*widget.FormItem {
	return []*widget.FormItem{
		{Text: "Host", Widget: pf.form["Host"]},
		{Text: "Local Port", Widget: pf.form["LocalPort"]},
		{Text: "Remote Port", Widget: pf.form["RemotePort"]},
	}
}

func (pf *PortForward) submitForm() {
	ui := widget.NewTextGrid()
	var region string
	ui.SetText("Initializing")
	ui.SetText("Starting Session")
	if s.Contains(pf.form["Host"].Text, "west") {
		region = "us-west-2"
	} else {
		region = "us-east-1"
	}
	tc := buildConfig(pf.form["Host"].Text, region)
	rp, err := strconv.Atoi(pf.form["RemotePort"].Text)
	if err != nil {
		log.Fatal("RemotePort could not be converted to int", pf.form["RemotePort"].Text, err)
	}
	lp, err := strconv.Atoi(pf.form["LocalPort"].Text)
	if err != nil {
		log.Fatal("LocalPort could not be converted to int", pf.form["LocalPort"].Text, err)
	}
	lpf := LocalPortForward{TargetConfig: tc, RemotePort: rp, LocalPort: lp}
	ui.SetText("Starting Forwarding session")
	mysql(lpf)

}

// func displayForm
func main() {

	a := app.New()

	var frm = map[string]*widget.Entry{
		"Host":       widget.NewEntry(),
		"LocalPort":  widget.NewEntry(),
		"RemotePort": widget.NewEntry(),
	}
	pf := PortForward{form: frm}
	pf.window = a.NewWindow("BL Port Forwarding")

	form := &widget.Form{
		Items:    pf.loadForm(),
		OnSubmit: func() { defer pf.submitForm() },
		OnCancel: func() {
			pf.window.Close()
		},
	}
	grid := container.New(layout.NewCenterLayout(), form)

	pf.window.Resize(fyne.NewSize(240, 240))

	pf.window.SetContent(grid)
	pf.window.ShowAndRun()
}
