package jaal

type Event struct {
	UnixTime       int64
	Timestamp      string
	CorrelationID  string
	Source         string
	SourceHostName string
	Type           string
	Summary        string
	Data           interface{}
}

type Listener interface {
	Listen(eventHandler func(*Event), errHandler func(error))
}

func Listen(listeners []Listener, eventLogger EventLogger, errLogger *ErrLogger) {
	for _, listener := range listeners {
		go listener.Listen(eventHandler(eventLogger), errHandler(errLogger))
	}
}

func eventHandler(eventLogger EventLogger) func(event *Event) {
	return func(e *Event) {
		eventLogger.Log(e)
	}
}

func errHandler(errLogger *ErrLogger) func(error) {
	return func(e error) {
		errLogger.Log(e)
	}
}

type FatalError struct {
	Err error
}

func (f FatalError) Error() string {
	return f.Err.Error()
}
