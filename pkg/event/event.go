package event

var (
	QuitEvent               = "QuitEvent"
	ApplicationMessageEvent = "ApplicationMessageEvent"
	RefreshScreenEvent      = "RefreshScreenEvent"
	RouteEvent              = "RouteEvent"
	BackRouteEvent          = "BackRouteEvent"
)

var defaultPubSub = newPubSub()

type pubSub struct {
	subscribers map[any][]SubscriberCallback
}

type SubscriberCallback func(...any)

func newPubSub() *pubSub {
	ps := &pubSub{subscribers: make(map[any][]SubscriberCallback)}
	return ps
}

func (ps *pubSub) publish(e any, param ...any) {
	if callbacks, ok := ps.subscribers[e]; ok {
		for _, callback := range callbacks {
			callback(param...)
		}
	}
}

func (ps *pubSub) registerSubscribe(e any, callback SubscriberCallback) {
	if _, ok := ps.subscribers[e]; !ok {
		ps.subscribers[e] = []SubscriberCallback{}
	}
	ps.subscribers[e] = append(ps.subscribers[e], callback)
}

func Publish(e any, param ...any) {
	defaultPubSub.publish(e, param...)
}

func RegisterSubscribe(e any, callback SubscriberCallback) {
	defaultPubSub.registerSubscribe(e, callback)
}
