package main

import (
	"fmt"
	"image/color"
	"log"
	"strconv"
	s "strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
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

func (pf *PortForward) updateWindow(lpf LocalPortForward) {
	connText := fmt.Sprintf("Host Name: %s\n\nLocal Port: %d\n\nRemote Port: %d\n\n", pf.form["Host"].Text, lpf.LocalPort, lpf.RemotePort)
	for k, v := range pf.form {
		log.Printf("Key %s Value %s", k, v.Text)
	}
	connected := widget.NewLabel(connText)
	cButton := &widget.Button{Icon: theme.ConfirmIcon(), Text: "Close Window and Session", Importance: widget.HighImportance, OnTapped: pf.window.Close}
	cButton.ExtendBaseWidget(cButton)
	c := &fyne.Container{Layout: layout.NewGridLayoutWithRows(1), Objects: []fyne.CanvasObject{connected, cButton}}

	pf.window.SetContent(c)
	pf.window.Resize(pf.window.Content().MinSize())
}

func (pf *PortForward) errorWindow(msg string, err error) {
	cRed := color.RGBA{207, 0, 15, 1}
	msgL := canvas.NewText(msg, cRed)
	errL := canvas.NewText(err.Error(), cRed)
	cButton := &widget.Button{Icon: theme.ErrorIcon(), Text: "Return to Form", Importance: widget.HighImportance, OnTapped: func() { defer pf.loadForm() }}
	pf.window.SetContent(container.New(layout.NewCenterLayout(), msgL, errL, canvas.NewLine(color.Black), cButton))

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
		dialog.ShowError(err, pf.window)
		pf.loadPFForm()
	}
	lp, err := strconv.Atoi(pf.form["LocalPort"].Text)
	if err != nil {
		dialog.ShowError(err, pf.window)
		pf.loadPFForm()
	}
	lpf := LocalPortForward{TargetConfig: tc, RemotePort: rp, LocalPort: lp}

	go func() {
		portForward(lpf)
	}()

	pf.updateWindow(lpf)

}

func (pf *PortForward) loadPFForm() {
	portForwardForm()
}

func portForwardForm() {
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
	c := container.New(layout.NewAdaptiveGridLayout(1), form)
	pf.window.Resize(fyne.NewSize(360, 240))

	pf.window.SetContent(c)
	pf.window.ShowAndRun()
}
