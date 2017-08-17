package api

var apiDescriptors = []ApiModule{
	{
		Name: "images",
		Protos: map[string]ApiProto{
			"APIAddImage": {
				Key:     "APIAddImage",
				Name:    "添加镜像",
				handler: APIAddImage,
				Paras: []ProtoPara{
					{
						Name:    "name",
						Type:    PARAM_TYPE_STRING,
						Desc:    "Image Name",
						Default: PARAM_NOT_NULL,
					},
					{
						Name:    "id",
						Type:    PARAM_TYPE_STRING,
						Desc:    "UUID of Image",
						Default: PARAM_NOT_NULL,
					},
					{
						Name:    "accountId",
						Type:    PARAM_TYPE_STRING,
						Desc:    "Account Id",
						Default: "",
					},
					{
						Name:    "url1",
						Type:    PARAM_TYPE_STRING,
						Desc:    "Url1 of Source Image",
						Default: PARAM_NOT_NULL,
					},
					{
						Name:    "url2",
						Type:    PARAM_TYPE_STRING,
						Desc:    "Url2 of Source Image",
						Default: "",
					},
					{
						Name:    "url3",
						Type:    PARAM_TYPE_STRING,
						Desc:    "Url3 of Source Image",
						Default: "",
					},
					{
						Name:    "format",
						Type:    PARAM_TYPE_STRING,
						Desc:    "Image Format, qcow2,raw,iso",
						Default: "qcow2",
					},
					{
						Name:    "arch",
						Type:    PARAM_TYPE_STRING,
						Desc:    "Arch of Image, amd64 or x86",
						Default: "amd64",
					},
					{
						Name:    "mediaType",
						Type:    PARAM_TYPE_STRING,
						Desc:    "镜像的媒体类型，RootVolumeTemplate，DataValumeTemplate, ISO",
						Default: "RootVolumeTemplate",
					},
					{
						Name:    "platform",
						Type:    PARAM_TYPE_STRING,
						Desc:    "平台类型，如Linux，Windows， Other",
						Default: "Linux",
					},
					{
						Name:    "guestOsType",
						Type:    PARAM_TYPE_STRING,
						Desc:    "客户操作系统类型，如CentOS，Debian",
						Default: "Debian7",
					},
					{
						Name:    "isSystem",
						Type:    PARAM_TYPE_BOOLEAN,
						Desc:    "是否为系统镜像",
						Default: false,
					},
					{
						Name:    "username",
						Type:    PARAM_TYPE_STRING,
						Desc:    "FTP or HTTP username",
						Default: "",
					},
					{
						Name:    "password",
						Type:    PARAM_TYPE_STRING,
						Desc:    "FTP or HTTP password",
						Default: "",
					},
				},
			},

			"APIDeleteImageByAccount": {
				Key:     "APIDeleteImageByAccount",
				Name:    "删除镜像",
				handler: APIDeleteImage,
				Paras: []ProtoPara{
					{
						Name:    "accountId",
						Type:    PARAM_TYPE_STRING,
						Desc:    "Account Id",
						Default: PARAM_NOT_NULL,
					},
				},
			},

			"APIDeleteImage": {
				Key:     "APIDeleteImage",
				Name:    "删除镜像",
				handler: APIDeleteImageByAccount,
				Paras: []ProtoPara{
					{
						Name:    "id",
						Type:    PARAM_TYPE_STRING,
						Desc:    "UUID of Image",
						Default: PARAM_NOT_NULL,
					},
					{
						Name:    "mediaType",
						Type:    PARAM_TYPE_STRING,
						Desc:    "镜像的媒体类型，RootVolumeTemplate，DataValumeTemplate, ISO",
						Default: "RootVolumeTemplate",
					},
				},
			},
		},
	},
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
		for _, proto := range descriptor.Protos {
			service := new(ApiService)
			service.Name = proto.Name
			service.Handler = proto.handler
			GApiServices["octlink.rstore.center."+descriptor.Name+"."+proto.Key] = service
		}
		loadModules(descriptor)
	}
}
