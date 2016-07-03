package clibot

import "github.com/RyanPrintup/nimbus"

type Handle struct {
	listener nimbus.Listener
	done     chan bool
}

func (h *Handle) Run(msg *nimbus.Message) {
	h.listener(msg)
	h.done <- true
}

func (c *CliClient) AddListener(event string, l nimbus.Listener) {
	c.listeners[event] = append(c.listeners[event], l)
}

func (c *CliClient) GetListeners(event string) []nimbus.Listener {
	return c.listeners[event]
}

func (c *CliClient) Emit(event string, msg *nimbus.Message) {
	for _, listener := range c.listeners[event] {
		h := Handle{listener, make(chan bool)}
		go h.Run(msg)
		<-h.done
	}
}