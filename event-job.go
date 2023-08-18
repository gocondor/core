package core

type EventJob func(event *Event, requestContext *Context)
