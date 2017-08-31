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
	Listen(eventHandler func(*Event), systemLogHandler func(interface{}))
}

func Listen(listeners []Listener, eventLogger EventLogger, systemLogger *SystemLogger) {
	for _, listener := range listeners {
		go listener.Listen(eventHandler(eventLogger), sysLogHandler(systemLogger))
	}
	select {} //block forever
}

func eventHandler(eventLogger EventLogger) func(event *Event) {
	return func(e *Event) {
		eventLogger.Log(e)
	}
}

func sysLogHandler(systemLogger *SystemLogger) func(interface{}) {
	return func(i interface{}) {
		switch i := i.(type) {
		case error:
			systemLogger.Error(i)
		default:
			systemLogger.Info(i)
		}
	}
}

type FatalError struct {
	Err error
}

func (f FatalError) Error() string {
	return f.Err.Error()
}
