package api

var ImageDescriptors = ApiModule{
	Name: "images",
	Protos: map[string]ApiProto{
		"APIAddImage": {
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
			Name:    "删除镜像（根据账号）",
			handler: APIDeleteImageByAccount,
			Paras: []ProtoPara{
				{
					Name:    "accountId",
					Type:    PARAM_TYPE_STRING,
					Desc:    "Account Id",
					Default: PARAM_NOT_NULL,
				},
			},
		},

		"APIShowImage": {
			Name:    "查看单个镜像",
			handler: APIShowImage,
			Paras: []ProtoPara{
				{
					Name:    "id",
					Type:    PARAM_TYPE_STRING,
					Desc:    "Image Id",
					Default: PARAM_NOT_NULL,
				},
			},
		},

		"APIShowAllImages": {
			Name:    "获取所有镜像",
			handler: APIShowAllImages,
			Paras: []ProtoPara{
				{
					Name:    "start",
					Type:    PARAM_TYPE_INT,
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
			handler: APIShowAccountList,
			Paras:   []ProtoPara{},
		},

		"APIDeleteImage": {
			Name:    "删除镜像",
			handler: APIDeleteImage,
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
}
