package main

import (
	"bufio"
	"fmt"
	"image/color"
	"os/exec"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type PortForward struct {
	form   map[string]*widget.Entry
	window fyne.Window
}

const LocalCommand = "/src/oddjob/ssm/port_forward.sh"

func (pf *PortForward) loadForm() []*widget.FormItem {
	return []*widget.FormItem{
		{Text: "Host", Widget: pf.form["Host"]},
		{Text: "Local Port", Widget: pf.form["LocalPort"]},
		{Text: "Remote Port", Widget: pf.form["RemotePort"]},
	}
}

func (pf *PortForward) submitForm() {
	ui := widget.NewTextGrid()
	ui.SetText("Starting Terminal")
	c := exec.Command("cd", "/src/oddjob/ssh", "/bin/bash", "+x", "/src/oddjob/ssm/port_forward.sh", pf.form["Host"].Text, pf.form["LocalPort"].Text, pf.form["RemotePort"].Text)
	stdout, err := c.StdoutPipe()
	if err != nil {
		panic(err)
	}

	err = c.Start()
	fmt.Println("The command is running")
	if err != nil {
		fmt.Println(err)
	}
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		m := scanner.Text()
		ui.SetText("Trying ...")
		ui.SetText(fmt.Sprintln(m))

	}
	c.Wait()
	text := scanner.Text()
	grid := container.New(layout.NewCenterLayout(), canvas.NewText(text, color.White))
	pf.window.SetContent(grid)

}

// func displayForm
func program() {

	// var formWidgets [3]FormPair

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
