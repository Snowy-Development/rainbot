// Copyright 2015-2017, Rodolfo Castillo-Valladares. All rights reserved.
//
// This file is governed by the Modified BSD License. You should have
// received a copy of the license (LICENSE.md) with this file's program.
// You may find the program here: https://github.com/raindevteam/rain
//
// Contact Rodolfo at rcvallada@gmail.com for any inquiries of this file.

package handler

import (
	"github.com/bwmarrin/discordgo"
	"github.com/raindevteam/rain/internal/dapi"
	"github.com/raindevteam/rain/internal/hail"
)

const (
	// NoDroplet is the string used to check whether a droplet listener has an
	// owner assigned or not.
	NoDroplet = "__NO_DROPLET__"

	// Internal is the string constant identifier for bot listeners.
	Internal = "__INTERNAL__"
)

// Listener is the interface for both internal listeners and droplet
// listeners. Use the Owner() function to discern if a listener is internal. If
// it is, it will return the string constant "__INTERNAL__".
type Listener interface {
	Run() error
	SetActionHandler(a Action)
	SetEvent(v interface{})
	Owner() string
}

// Command is an interface for both internal and droplets.
type Command interface {
	Invoke() error
	SetArguments(string)
	SetEvent(*discordgo.MessageCreate)
	SetName(string)
	GetName() string
	Owner() string
}

// InternalListener is the struct for internal listeners. It does not have a
// droplet field like DropletListener.
type InternalListener struct {
	Event   interface{}
	enabled bool
	action  Action
}

// Run will execute the listener's Action.
func (il *InternalListener) Run() error {
	if il.action == nil {
		return hail.Err(hail.Fhandler,
			"listener tried to run without action set")
	}
	if il.Event == nil {
		return hail.Err(hail.Fhandler,
			"listener tried to run without event set")
	}
	il.action.Do(il.Event)
	return nil
}

// SetActionHandler will set the listener's action. It is set during a call to
// AddInternalListener().
func (il *InternalListener) SetActionHandler(a Action) {
	il.action = a
}

// SetEvent will set the listener's event for future use when Do() is called on
// the listener's action.
func (il *InternalListener) SetEvent(v interface{}) {
	il.Event = v
}

// Owner for InternalListener returns the string constant "__INTERNAL__".
func (il *InternalListener) Owner() string {
	return Internal
}

// The DropletListener struct is similar to InternalListener, but is used for
// droplet listeners. The most notable difference is how action is handled.
type DropletListener struct {
	Event  interface{}
	drop   *dapi.Droplet
	action Action
}

// Run will run the listener's action that being an RPC call.
func (l *DropletListener) Run() error {
	return nil
}

// SetActionHandler will set the listener's action.
func (l *DropletListener) SetActionHandler(a Action) {
	l.action = a
}

// SetEvent will set a listener's event for future use when calling action.Do().
func (l *DropletListener) SetEvent(v interface{}) {
	l.Event = v
}

// Owner will return the name of the droplet to which this listener belongs to.
func (l *DropletListener) Owner() string {
	if l.drop == nil {
		return NoDroplet
	}
	return NoDroplet
}

// The InternalCommand struct is a wrapper for internal commands.
type InternalCommand struct {
	Event       *discordgo.MessageCreate
	CommandFunc func(string, *discordgo.MessageCreate) error
	Arguments   string
	name        string
}

// Invoke runs the commandFunc for the internal command.
func (ic *InternalCommand) Invoke() error {
	err := ic.CommandFunc(ic.Arguments, ic.Event)
	if err != nil {
		return err
	}
	return nil
}

// SetArguments sets the internal command's argument string.
func (ic *InternalCommand) SetArguments(arguments string) {
	ic.Arguments = arguments
}

// SetEvent sets the internal command's MessageCreate event. It will always be
// the one corresponding to the message that specified the command.
func (ic *InternalCommand) SetEvent(e *discordgo.MessageCreate) {
	ic.Event = e
}

// SetName will set the name for the command. This name is used to call the
// command from Discord.
func (ic *InternalCommand) SetName(name string) {
	ic.name = name
}

// GetName will return the name for this command.
func (ic *InternalCommand) GetName() string {
	return ic.name
}

// Owner returns the internal command's owner. Always handler.Internal.
func (ic *InternalCommand) Owner() string {
	return Internal
}
