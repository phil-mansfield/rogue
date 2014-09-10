package mvc

import (
	"fmt"

	"github.com/phil-mansfield/rogue/actor"
	"github.com/phil-mansfield/rogue/error"
	"github.com/phil-mansfield/rogue/event"
	"github.com/phil-mansfield/rogue/world"
)

type DudModel struct {
	frames int
}

type DudController struct {
}

type DudView struct {
}

// DudModel methods

func (model *DudModel) PauseTasks() *error.Error {
	return nil
}

func (model *DudModel) ResumeTasks() *error.Error {
	return nil
}

func (model *DudModel) Respond([]Key) ([]event.Event, *error.Error) {
	model.frames += 1
	return []event.Event{event.Message{fmt.Sprintf("Frame # = %d", model.frames)}}, nil
}

func (model *DudModel) RespondError(
	err *error.Error,
) ([]event.Event, *error.Error) {
	
	return []event.Event{}, nil
}

func (model *DudModel) Map() world.Map {
	return nil
}

func (model *DudModel) Player() actor.Actor {
	return nil
}

func (model *DudModel) GameOver() bool {
	return model.frames == 40
}

func (model *DudModel) Close() {
	fmt.Println("Goodbye from DudModel!")
}


// DudView methods


func (view *DudView) Draw(
	currMap world.Map,
	player actor.Actor,
	events []event.Event,
) *error.Error {

	for _, ev := range events {
		switch ev := ev.(type) {
		case event.Message:
			fmt.Println(ev.String())
		default:
			return error.New(error.Sanity, "Unknown event type.")
		}
	}
	return nil
}
	
func (view *DudView) Respond([]Key) ([]Key, *error.Error) {
	return []Key{}, nil
}

func (view *DudView) Close() { 
	fmt.Println("Goodbye from DudView!")
}


// DudController methods


func (controller *DudController) KeysPressed() ([]Key, *error.Error) {
	return []Key{}, nil
}

func (controller *DudController) Close() {
	fmt.Println("Goodbye from DudController!")
}