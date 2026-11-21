package main

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/CodeGoy/gqrx"
)

func main() {
	g := &gqrx.GQRX{
		Address: "127.0.0.1",
		Port:    "7356",
	}
	g.Init()
	a := app.New()
	window := a.NewWindow("GQRX controller")
	var demodButtons []fyne.CanvasObject
	for i := 0; i < len(gqrx.Demods); i++ {
		nb := widget.NewButton(gqrx.Demods[i], func() {
			g.SetDemod(gqrx.Demods[i])
		})
		demodButtons = append(demodButtons, nb)
	}
	lnaLabel := widget.NewLabel("LNA: 0.0")
	lnaSlider := widget.NewSlider(0, 49.6)
	lnaSlider.Step = 0.2
	lnaSlider.OnChanged = func(value float64) {
		g.SetLnaGain(value)
		fyne.Do(func() {
			lnaLabel.SetText(fmt.Sprintf("LNA: %.1f", value))
		})
	}
	audioLabel := widget.NewLabel("LNA: 0.0")
	audioSlider := widget.NewSlider(-80, 49.6)
	audioSlider.Step = 0.2
	audioSlider.OnChanged = func(value float64) {
		g.SetAF(value)
		fyne.Do(func() {
			audioLabel.SetText(fmt.Sprintf("Audio: %.1f", value))
		})
	}
	window.SetContent(container.NewVBox(
		container.NewHBox(
			demodButtons...,
		),
		&layout.Spacer{
			FixHorizontal: true,
			FixVertical:   false,
		},
		widget.NewSeparator(),
		lnaLabel, lnaSlider,
		widget.NewSeparator(),
		audioLabel, audioSlider,
	))
	a.Lifecycle().SetOnStarted(func() {
		fyne.Do(func() {
			lnaGain, err := g.GetLnaGain()
			if err != nil {
				log.Fatal(err)
			}
			lnaSlider.SetValue(lnaGain)
			lnaLabel.SetText(fmt.Sprintf("LNA: %.1f", lnaGain))

			audioGain, err := g.GetAF()
			if err != nil {
				log.Fatal(err)
			}
			audioSlider.SetValue(audioGain)
			audioLabel.SetText(fmt.Sprintf("Audio: %.1f", audioGain))
		})
	})
	window.ShowAndRun()
}
