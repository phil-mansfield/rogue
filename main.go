package main

import (
	"flag"
	"os"
	"time"

	"github.com/phil-mansfield/rogue/config"
	"github.com/phil-mansfield/rogue/error"
	"github.com/phil-mansfield/rogue/mvc"
)

var (
	configPtr = flag.String("config", "", "Configuration file location. " + 
		"Empty string indicates that the default values will be used.")
)

func main() {

	// Setup

	flag.Parse()

	info := getConfigInfo()

	model, view, controller, mvcErr := mvc.New(info)
	if mvcErr != nil {
		error.Report(mvcErr)
		os.Exit(1)
	}

	// Draw starting screen

	events, err := model.Respond([]mvc.Key{})
	if err != nil { drawError(model, view, err) }

	err = view.Draw(model.Map(), model.Player(), events)
	if err != nil { drawError(model, view, err) }

	// mainloop

	mainloop(model, view, controller, info)

	model.Close()
	view.Close()
	controller.Close()
}

// getConfigInfo returns the the config.Info instance corresponding to the
// user-provided config file location. If a fatal error occurs, it is reported
// using the low level tools provided in the error package and the program
// terminates.
//
// TODO: Consider propogating Configuration errors to the the point where they
// can be reported with the standard mvc apparatus.
func getConfigInfo() *config.Info {
	configPath := *configPtr

	var (
		info *config.Info
		err *error.Error
	)

	if configPath == "" {
		info, err = config.Default()
	} else {
		info, err = config.Parse(configPath)
	}

	if err != nil {
		if err.Code == error.Sanity || err.Code == error.Library {
			error.Report(err)
			os.Exit(1)
		} else {
			error.Report(err)
			os.Exit(1)
		}
	}

	return info
}

// drawError handles the  boiler plate code associated with drawing the 
// specified error. If an error occurs during this step, the process is
// considered a lost cause and terminates after using low-level error reporting
// tools.
func drawError(model mvc.Model, view mvc.View, err *error.Error) {
	events, err := model.RespondError(err)
	if err != nil {
		error.Report(err)
		os.Exit(1)
	}

	err = view.Draw(model.Map(), model.Player(), events)
	if err != nil {
		error.Report(err)
		os.Exit(1)
	}
}

// mainloop draws a frame at the user specified rate until the game ends.
func mainloop(
	model mvc.Model,
	view mvc.View, 
	controller mvc.Controller, 
	info *config.Info,
) {

	ms := int(1000.0 / float64(info.FramesPerSecond))
	tick := time.Tick(time.Millisecond * time.Duration(ms))

	for !model.GameOver() {
		select {
		case <-tick:
			err := model.PauseTasks()
			if err != nil { drawError(model, view, err) }

			keys, err := controller.KeysPressed()
			if err != nil { drawError(model, view, err) }

			keys, err = view.Respond(keys)
			if err != nil { drawError(model, view, err) }

			events, err := model.Respond(keys)
			if err != nil { drawError(model, view, err) }

			err = view.Draw(model.Map(), model.Player(), events)
			if err != nil { drawError(model, view, err) }

			err = model.ResumeTasks()
			if err != nil { drawError(model, view, err) }
		}
	}	
}