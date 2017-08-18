package api

var apiDescriptors = []ApiModule{
	ImageDescriptors,
}

func loadModules(module ApiModule) {

	if GApiConfig.Modules == nil {
		GApiConfig.Modules = make(map[string]ApiModule, 50)
	}

	GApiConfig.Modules[module.Name] = module
}

func init() {

	GApiServices = make(map[string]*ApiService, 10000)

	for _, descriptor := range apiDescriptors {
		for key, proto := range descriptor.Protos {
			service := new(ApiService)
			service.Name = proto.Name
			service.Handler = proto.handler
			proto.Key = API_PREFIX_CENTER + "." + descriptor.Name + "." + key
			descriptor.Protos[key] = proto
			GApiServices[API_PREFIX_CENTER+"."+descriptor.Name+"."+key] = service
		}

		loadModules(descriptor)
	}
}
