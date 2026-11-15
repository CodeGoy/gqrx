package main

import (
	"image/color"
	"log"

	"github.com/CodeGoy/gqrx"

	"github.com/golang-gui/goui/layout"
	"github.com/golang-gui/goui/theme/modern"
	"github.com/golang-gui/goui/ui"
)

type App struct {
	gqrx *gqrx.GQRX
}

func main() {
	a := App{}
	g := &gqrx.GQRX{
		Address: "127.0.0.1",
		Port:    "7356",
	}
	a.gqrx = g
	g.Init()
	a.gui()
}

func (a *App) gui() {
	var buttons []ui.View
	for i := 0; i < len(gqrx.Demods); i++ {
		nb := ui.Button(gqrx.Demods[i])
		nb.OnClick(func() {
			a.gqrx.SetDemod(gqrx.Demods[i])
		})
		buttons = append(buttons, nb)
	}
	err := ui.Run("com.example.Hello", func(app ui.App) ui.RootView {
		settings := app.Settings()
		stylesheet := modern.Sheet(modern.Options{
			Dark: true, AccentColor: color.NRGBA{R: 140, G: 70, B: 220, A: 255},
			FontFamily: settings.FontFamily(), FontSize: settings.FontSize(),
		})

		return ui.Root().
			StyleSheet(stylesheet). // Use the Modern theme
			Windows(
				ui.Window("main").
					Title("GQRX Demod remote").
					Size(850, 75).
					Content(
						ui.VBox(
							ui.HBox(
								buttons...,
							).Spacing(5).Padding(2),
						).MainAlign(layout.MainCenter).Padding(15),
					),
			)
	})
	if err != nil {
		log.Fatal(err)
	}
}
