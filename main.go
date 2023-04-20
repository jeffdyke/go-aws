package main

import (
	"fmt"
	"log"
	"strconv"
	s "strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type PortForward struct {
	form   map[string]*widget.Entry
	window fyne.Window
}

func (pf *PortForward) loadForm() []*widget.FormItem {
	return []*widget.FormItem{
		{Text: "Host", Widget: pf.form["Host"]},
		{Text: "Local Port", Widget: pf.form["LocalPort"], HintText: "1515"},
		{Text: "Remote Port", Widget: pf.form["RemotePort"], HintText: "3306"},
	}
}

func (pf *PortForward) submitForm() {

	var region string
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

	go func() {
		portForward(lpf)
	}()

	connText := fmt.Sprintf("Host Name: %s\nLocal Port: %d\nRemote Port: %d\n", pf.form["Host"].Text, rp, lp)

	connected := widget.NewLabel(connText)
	cButton := &widget.Button{Icon: theme.ConfirmIcon(), Text: "Close Window and Session", Importance: widget.HighImportance, OnTapped: pf.window.Close}
	cButton.ExtendBaseWidget(cButton)

	c := container.New(layout.NewAdaptiveGridLayout(2), connected, layout.NewSpacer(), cButton, layout.NewSpacer())
	pf.window.SetContent(c)
}

// func displayForm
func main() {
	hostEV := &widget.Entry{Validator: validation.NewRegexp("[a-z0-3-]+", "Validation for Host fails")}
	hostEV.ExtendBaseWidget(hostEV)
	lpEV := &widget.Entry{Validator: validation.NewRegexp("[0-9]+", "Validation for LocalPort fails, number required"),
		Text: "1515", PlaceHolder: "1515"}
	lpEV.ExtendBaseWidget(lpEV)
	rpEV := &widget.Entry{Validator: validation.NewRegexp("[0-9]", "Validation for Remote fails, number required"),
		Text: "3306", PlaceHolder: "3306"}
	rpEV.ExtendBaseWidget(rpEV)
	a := app.New()

	var frm = map[string]*widget.Entry{
		"Host":       hostEV,
		"LocalPort":  lpEV,
		"RemotePort": rpEV,
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
