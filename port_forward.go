package main

import (
	"image/color"
	"log"
	"strconv"
	s "strings"

	"fyne.io/fyne/v2"
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

func (pf *PortForward) successWindow(lpf LocalPortForward) {
	blue := color.RGBA{5, 55, 247, 1}
	// green := color.RGBA{45, 172, 31, 1}
	connText := []fyne.CanvasObject{
		canvas.NewText("Host Name", color.White), canvas.NewText(pf.form["Host"].Text, color.White),
		canvas.NewText("Local Port", color.White), canvas.NewText(pf.form["LocalPort"].Text, color.White),
		canvas.NewText("Remote Port", color.White), canvas.NewText(pf.form["RemotePort"].Text, color.White),
	}

	for k, v := range pf.form {
		log.Printf("Key %s Value %s", k, v.Text)
	}
	connected := &fyne.Container{Layout: layout.NewGridLayout(2), Objects: connText}
	cButton := &widget.Button{Icon: theme.ConfirmIcon(), Text: "Close Window and Session",
		Importance: widget.HighImportance, OnTapped: pf.window.Close}
	cButton.ExtendBaseWidget(cButton)
	c := &fyne.Container{Layout: layout.NewGridLayoutWithRows(1), Objects: []fyne.CanvasObject{connected, container.NewMax(canvas.NewRectangle(blue), cButton)}}

	pf.window.SetContent(c)
	pf.window.Resize(pf.window.Content().MinSize())
}

func (pf *PortForward) errorWindow(msg string, err error) {

	msgL := widget.NewLabelWithStyle(msg, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	errL := widget.NewLabelWithStyle(err.Error(), fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	cButton := widget.NewButton("Return to form", func() {
		pf.portForwardForm()
	})
	pf.window.SetContent(container.New(layout.NewGridLayout(2), widget.NewLabel("Message"), msgL, widget.NewLabel("Error"), errL, layout.NewSpacer(), cButton))
	pf.window.Resize(pf.window.Content().MinSize())
}
func (pf *PortForward) submitForm() {

	var region string
	if s.Contains(pf.form["Host"].Text, "west") {
		region = "us-west-2"
	} else {
		region = "us-east-1"
	}
	tc, err := instanceConfig(pf.form["Host"].Text, region)
	if err != nil {
		pf.errorWindow("Build Config Failed", err)
		return
	}

	rp, err := strconv.Atoi(pf.form["RemotePort"].Text)
	if err != nil {
		dialog.ShowError(err, pf.window)
	}
	lp, err := strconv.Atoi(pf.form["LocalPort"].Text)
	if err != nil {
		dialog.ShowError(err, pf.window)
	}

	lpf := LocalPortForward{TargetConfig: tc, RemotePort: rp, LocalPort: lp}

	go func() {
		portForward(lpf)
	}()

	pf.successWindow(lpf)

}

func (pf *PortForward) portForwardForm() {
	hostEV := &widget.Entry{Validator: validation.NewRegexp("[a-z0-3-]+", "Validation for Host fails")}
	hostEV.ExtendBaseWidget(hostEV)
	lpEV := &widget.Entry{Validator: validation.NewRegexp("[0-9]+", "Validation for LocalPort fails, number required"),
		Text: "1515", PlaceHolder: "1515"}
	lpEV.ExtendBaseWidget(lpEV)
	rpEV := &widget.Entry{Validator: validation.NewRegexp("[0-9]", "Validation for Remote fails, number required"),
		Text: "3306", PlaceHolder: "3306"}
	rpEV.ExtendBaseWidget(rpEV)

	var frm = map[string]*widget.Entry{
		"Host":       hostEV,
		"LocalPort":  lpEV,
		"RemotePort": rpEV,
	}
	pf.form = frm

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
}
