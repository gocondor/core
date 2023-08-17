package core

type EventJob func(event *Event, c *Context)
