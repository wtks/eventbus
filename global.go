package eventbus

var g = New()

func Subscribe(name string, fn Handler) {
	g.Subscribe(name, fn)
}

func Unsubscribe(name string, fn Handler) bool {
	return g.Unsubscribe(name, fn)
}

func UnsubscribeAll(name string) {
	g.UnsubscribeAll(name)
}

func Publish(name string, args ...interface{}) {
	g.Publish(name, args...)
}

func WaitAll() {
	g.WaitAll()
}
