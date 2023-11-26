package lib

type IncrementalRouteManager struct {
	Routes map[string]Route
}

func (receiver *IncrementalRouteManager) HasRoute(route Route) bool {
	_, ok := receiver.Routes[route.RoutePath]
	return ok
}

func (receiver *IncrementalRouteManager) IncludeRoute(route Route) {
	receiver.Routes[route.RoutePath] = route
}

func (receiver *IncrementalRouteManager) ExcludeRoute(route Route) {
	delete(receiver.Routes, route.RoutePath)
}
