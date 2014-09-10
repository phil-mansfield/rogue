package scratch

import (
	"fmt"
	"os"
	"time"

	"phil-mansfield/github.com/rogue/config"
	"phil-mansfield/github.com/rogue/error"
	"phil-mansfield/github.com/rogue/mvc"
)

func main() {
	configInfo, configErr := config.New(config.Path())

	// configInfo is still valid even if configErr != nil.
	model, view, controller, mvcErr := mvc.New(configInfo)

	if mvcErr != nil {
		error.Report(mvcErr)
		return
	}

	tick := time.Tick(time.Millesecond * configInfo.FrameTimeMillseconds)
	lock := make(chan int)

	for !model.GameOver() {

		select {
		case <-tick:

			lock <- 1

			err := model.PauseTasks()
			if err != nil { model.RespondError(err) }

			keysPressed, err := controller.KeysPressed()
			if err != nil { model.RespondError(err) }

			keysDown, err := controller.KeysDown()
			if err != nil { model.RespondError(err) }

			events, err := model.RespondKeys(keysPressed, keysDown)
			if err != nil { model.RespondError(err) }

			// I have no idea what could cause an error here, but we have to
			// prepare for it:

			err := view.DrawFrame(events, model.CurrentMap(), model.Player())
			if err != nil {
				model.RespondError(err)
				events = model.RespondKeys([]mvc.Key, []mvc.Key)
				if innerErr != nil {
					// Now is the time to panic.

				}
			}

			err := model.UnpauseTasks()

			// Note that if UnpauseTasks works sequentially, this will hang
			<-lock

		}
	}

	model.Close()
	controller.Close()
	view.Close()
}