package eventbus

var g = New()

func Subscribe(event string, fn Handler) {
	g.Subscribe(event, fn)
}

func Unsubscribe(event string, fn Handler) bool {
	return g.Unsubscribe(event, fn)
}

func UnsubscribeAll(event string) {
	g.UnsubscribeAll(event)
}

func Publish(event string, args ...interface{}) {
	g.Publish(event, args...)
}

func WaitAll() {
	g.WaitAll()
}
