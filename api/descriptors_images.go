package api

// ImageDescriptors for image management by API
var ImageDescriptors = Module{
	Name: "image",
	Protos: map[string]Proto{
		"APIAddImage": {
			Name:    "添加镜像",
			handler: AddImage,
			Paras: []ProtoPara{
				{
					Name:    "name",
					Type:    ParamTypeString,
					Desc:    "Image Name",
					Default: ParamNotNull,
				},
				{
					Name:    "id",
					Type:    ParamTypeString,
					Desc:    "UUID of Image",
					Default: ParamNotNull,
				},
				{
					Name:    "accountId",
					Type:    ParamTypeString,
					Desc:    "Account Id",
					Default: "",
				},
				{
					Name:    "url1",
					Type:    ParamTypeString,
					Desc:    "Url1 of Source Image",
					Default: "",
				},
				{
					Name:    "url2",
					Type:    ParamTypeString,
					Desc:    "Url2 of Source Image",
					Default: "",
				},
				{
					Name:    "url3",
					Type:    ParamTypeString,
					Desc:    "Url3 of Source Image",
					Default: "",
				},
				{
					Name:    "format",
					Type:    ParamTypeString,
					Desc:    "Image Format, qcow2,raw,iso",
					Default: "qcow2",
				},
				{
					Name:    "arch",
					Type:    ParamTypeString,
					Desc:    "Arch of Image, amd64 or x86",
					Default: "amd64",
				},
				{
					Name:    "mediaType",
					Type:    ParamTypeString,
					Desc:    "镜像的媒体类型，RootVolumeTemplate，DataValumeTemplate, ISO",
					Default: "RootVolumeTemplate",
				},
				{
					Name:    "platform",
					Type:    ParamTypeString,
					Desc:    "平台类型，如Linux，Windows， Other",
					Default: "Linux",
				},
				{
					Name:    "guestOsType",
					Type:    ParamTypeString,
					Desc:    "客户操作系统类型，如CentOS，Debian",
					Default: "Debian7",
				},
				{
					Name:    "isSystem",
					Type:    ParamTypeBoolean,
					Desc:    "是否为系统镜像",
					Default: false,
				},
				{
					Name:    "username",
					Type:    ParamTypeString,
					Desc:    "FTP or HTTP username",
					Default: "",
				},
				{
					Name:    "password",
					Type:    ParamTypeString,
					Desc:    "FTP or HTTP password",
					Default: "",
				},
			},
		},

		"APIUpdateImage": {
			Name:    "编辑镜像",
			handler: UpdateImage,
			Paras: []ProtoPara{
				{
					Name:    "id",
					Type:    ParamTypeString,
					Desc:    "UUID of Image",
					Default: ParamNotNull,
				},
				{
					Name:    "name",
					Type:    ParamTypeString,
					Desc:    "Image Name",
					Default: ParamNotNull,
				},
				{
					Name:    "arch",
					Type:    ParamTypeString,
					Desc:    "Arch of Image, amd64 or x86",
					Default: "amd64",
				},
				{
					Name:    "platform",
					Type:    ParamTypeString,
					Desc:    "平台类型，如Linux，Windows， Other",
					Default: "Linux",
				},
				{
					Name:    "guestOsType",
					Type:    ParamTypeString,
					Desc:    "客户操作系统类型，如CentOS，Debian",
					Default: "Debian7",
				},
			},
		},

		"APIRemoveImageByAccount": {
			Name:    "删除镜像（根据账号）",
			handler: DeleteImageByAccount,
			Paras: []ProtoPara{
				{
					Name:    "accountId",
					Type:    ParamTypeString,
					Desc:    "Account Id",
					Default: ParamNotNull,
				},
			},
		},

		"APIShowImage": {
			Name:    "查看单个镜像",
			handler: ShowImage,
			Paras: []ProtoPara{
				{
					Name:    "id",
					Type:    ParamTypeString,
					Desc:    "Image Id",
					Default: ParamNotNull,
				},
			},
		},

		"APIShowAllImage": {
			Name:    "获取所有镜像",
			handler: ShowAllImages,
			Paras: []ProtoPara{
				{
					Name:    "start",
					Type:    ParamTypeInt,
					Desc:    "开始位置",
					Default: 0,
				},
				{
					Name:    "limit",
					Type:    "int",
					Desc:    "获取条目",
					Default: 15,
				},
				{
					Name:    "accountId",
					Type:    "string",
					Desc:    "账号ID",
					Default: "",
				},
				{
					Name:    "mediaType",
					Type:    "string",
					Desc:    "镜像类型，RootVolumeTemplate，DataValumeTemplate, ISO",
					Default: "",
				},
				{
					Name:    "keyword",
					Type:    "string",
					Desc:    "Key word for image name or Id",
					Default: "",
				},
			},
		},

		"APIShowAccountList": {
			Name:    "查看账号列表",
			handler: ShowAccountList,
			Paras:   []ProtoPara{},
		},

		"APIRemoveImage": {
			Name:    "删除镜像",
			handler: DeleteImage,
			Paras: []ProtoPara{
				{
					Name:    "id",
					Type:    ParamTypeString,
					Desc:    "UUID of Image",
					Default: ParamNotNull,
				},
				{
					Name:    "mediaType",
					Type:    ParamTypeString,
					Desc:    "镜像的媒体类型，RootVolumeTemplate，DataValumeTemplate, ISO",
					Default: "RootVolumeTemplate",
				},
			},
		},
	},
}
