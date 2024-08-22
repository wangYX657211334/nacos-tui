package base

import "github.com/wangYX657211334/nacos-tui/pkg/event"

func Route(path string, params ...any) {
	p := []any{path}
	p = append(p, params...)
	event.Publish(event.RouteEvent, p...)
}

func BackRoute() {
	event.Publish(event.BackRouteEvent)
}
