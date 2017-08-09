package v1

var RouteDescriptorsMap map[string]RouteDescriptor

// RouteDescriptor describes a route specified by name.
type RouteDescriptor struct {
	// Name is the name of the route, as specified in RouteNameXXX exports.
	// These names a should be considered a unique reference for a route. If
	// the route is registered with gorilla, this is the name that will be
	// used.
	Name string

	// Path is a gorilla/mux-compatible regexp that can be used to match the
	// route. For any incoming method and path, only one route descriptor
	// should match.
	Path string
}

var routeDescriptors = []RouteDescriptor{
	{
		Name: "v1base",
		Path: "/v1/",
	},
}

func init() {
	RouteDescriptorsMap = make(map[string]RouteDescriptor, len(routeDescriptors))
	for _, descriptor := range routeDescriptors {
		RouteDescriptorsMap[descriptor.Name] = descriptor
	}
}
