package mvc

import (
	"github.com/phil-mansfield/rogue/actor"
	"github.com/phil-mansfield/rogue/config"
	"github.com/phil-mansfield/rogue/error"
	"github.com/phil-mansfield/rogue/event"
	"github.com/phil-mansfield/rogue/world"
)

type Model interface {
	PauseTasks() *error.Error
	ResumeTasks() *error.Error
	Respond([]Key) ([]event.Event, *error.Error)
	RespondError(*error.Error) ([]event.Event, *error.Error)

	Map() world.Map
	Player() actor.Actor

	GameOver() bool

	Close()
}

type View interface {
	Draw(world.Map, actor.Actor, []event.Event) *error.Error
	Respond([]Key) ([]Key, *error.Error)

	Close()
}

type Controller interface {
	KeysPressed() ([]Key, *error.Error)

	Close()
}

type Key struct {
}

func New(info *config.Info) (Model, View, Controller, *error.Error) {
	return &DudModel{}, &DudView{}, &DudController{}, nil
}