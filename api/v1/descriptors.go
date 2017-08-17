package v1

import (
	"octlink/rstore/digest"
	"octlink/rstore/reference"
)

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

	// Description for this router.
	Description string

	// Methods should describe the various HTTP methods that may be used on
	// this route, including request and response formats.
	Methods []MethodDescriptor
}

// MethodDescriptor provides a description of the requests that may be
// conducted with the target method.
type MethodDescriptor struct {

	// Method is an HTTP method, such as GET, PUT or POST.
	Method string

	// Description should provide an overview of the functionality provided by
	// the covered method, suitable for use in documentation. Use of markdown
	// here is encouraged.
	Description string

	// Requests is a slice of request descriptors enumerating how this
	// endpoint may be used.
	Requests []RequestDescriptor
}

// RequestDescriptor covers a particular set of headers and parameters that
// can be carried out with the parent method. Its most helpful to have one
// RequestDescriptor per API use case.
type RequestDescriptor struct {
	// Name provides a short identifier for the request, usable as a title or
	// to provide quick context for the particular request.
	Name string

	// Description should cover the requests purpose, covering any details for
	// this particular use case.
	Description string

	// Headers describes headers that must be used with the HTTP request.
	Headers []ParameterDescriptor

	// PathParameters enumerate the parameterized path components for the
	// given request, as defined in the route's regular expression.
	PathParameters []ParameterDescriptor

	// QueryParameters provides a list of query parameters for the given
	// request.
	QueryParameters []ParameterDescriptor
}

// ParameterDescriptor describes the format of a request parameter, which may
// be a header, path parameter or query parameter.
type ParameterDescriptor struct {
	// Name is the name of the parameter, either of the path component or
	// query parameter.
	Name string

	// Type specifies the type of the parameter, such as string, integer, etc.
	Type string

	// Format is a specifying the string format accepted by this parameter.
	Format string

	// Description provides a human-readable description of the parameter.
	Description string

	// Required or not
	Required bool

	// Examples provides multiple examples for the values that might be valid
	// for this parameter.
	Examples []string
}

var (
	nameParameterDescriptor = ParameterDescriptor{
		Name:        "name",
		Type:        "string",
		Format:      reference.NameRegexp.String(),
		Required:    true,
		Description: `Name of the target repository.`,
	}

	referenceParameterDescriptor = ParameterDescriptor{
		Name:        "reference",
		Type:        "string",
		Format:      reference.TagRegexp.String(),
		Required:    true,
		Description: `Tag or digest of the target manifest.`,
	}

	uuidParameterDescriptor = ParameterDescriptor{
		Name:        "uuid",
		Type:        "opaque",
		Required:    true,
		Description: "A uuid identifying the upload. This field can accept characters that match `[a-zA-Z0-9-_.=]+`.",
	}

	digestPathParameter = ParameterDescriptor{
		Name:        "digest",
		Type:        "path",
		Required:    true,
		Format:      digest.DigestRegexp.String(),
		Description: `Digest of desired blob.`,
	}

	hostHeader = ParameterDescriptor{
		Name:        "Host",
		Type:        "string",
		Description: "Standard HTTP Host Header. Should be set to the registry host.",
		Format:      "<registry host>",
		Examples:    []string{"registry-1.docker.io"},
	}

	authHeader = ParameterDescriptor{
		Name:        "Authorization",
		Type:        "string",
		Description: "An RFC7235 compliant authorization header.",
		Format:      "<scheme> <token>",
		Examples:    []string{"Bearer dGhpcyBpcyBhIGZha2UgYmVhcmVyIHRva2VuIQ=="},
	}

	authChallengeHeader = ParameterDescriptor{
		Name:        "WWW-Authenticate",
		Type:        "string",
		Description: "An RFC7235 compliant authentication challenge header.",
		Format:      `<scheme> realm="<realm>", ..."`,
		Examples: []string{
			`Bearer realm="https://auth.docker.com/", service="registry.docker.com", scopes="repository:library/ubuntu:pull"`,
		},
	}

	contentLengthZeroHeader = ParameterDescriptor{
		Name:        "Content-Length",
		Description: "The `Content-Length` header must be zero and the body must be empty.",
		Type:        "integer",
		Format:      "0",
	}

	dockerUploadUUIDHeader = ParameterDescriptor{
		Name:        "Docker-Upload-UUID",
		Description: "Identifies the docker upload uuid for the current request.",
		Type:        "uuid",
		Format:      "<uuid>",
	}

	digestHeader = ParameterDescriptor{
		Name:        "Docker-Content-Digest",
		Description: "Digest of the targeted content for the request.",
		Type:        "digest",
		Format:      "<digest>",
	}

	linkHeader = ParameterDescriptor{
		Name:        "Link",
		Type:        "link",
		Description: "RFC5988 compliant rel='next' with URL to next result set, if available",
		Format:      `<<url>?n=<last n value>&last=<last entry from response>>; rel="next"`,
	}

	paginationParameters = []ParameterDescriptor{
		{
			Name:        "n",
			Type:        "integer",
			Description: "Limit the number of entries in each response. It not present, all entries will be returned.",
			Format:      "<integer>",
			Required:    false,
		},
		{
			Name:        "last",
			Type:        "string",
			Description: "Result set will include values lexically after last.",
			Format:      "<integer>",
			Required:    false,
		},
	}
)

var routeDescriptors = []RouteDescriptor{
	{
		Name:        RouteNameBase,
		Path:        "/v1/",
		Description: "Base V1 API route",
		Methods: []MethodDescriptor{
			{
				Method:      "GET",
				Description: "Check implements API V1.",
			},
		},
	},
	{
		Name:        RouteNameTags,
		Path:        "/v1/{name:" + reference.NameRegexp.String() + "}/tags/list",
		Description: "Retrieve information about tags.",
		Methods: []MethodDescriptor{
			{
				Method:      "GET",
				Description: "Fetch the tags under the repository identified by `name`.",
				Requests: []RequestDescriptor{
					{
						Name:        "Tags",
						Description: "Return all tags for the repository",
						PathParameters: []ParameterDescriptor{
							nameParameterDescriptor,
						},
					},
				},
			},
		},
	},
	{
		Name:        RouteNameBlobUpload,
		Path:        "/v1/{name:" + reference.NameRegexp.String() + "}/blobs/uploads/",
		Description: "Initiate a blob upload.",
		Methods: []MethodDescriptor{
			{
				Method:      "POST",
				Description: "Initiate a resumable blob upload. On success, upload location will be provided. Optionally, if the `digest` parameter is present, the request body will be used to complete the upload in a single request.",
				Requests: []RequestDescriptor{
					{
						Name:        "Initiate Monolithic Blob Upload",
						Description: "Upload a blob identified by the `digest` parameter in single request. This upload will not beresumable unless a recoverable error is returned.",
						Headers: []ParameterDescriptor{
							{
								Name:   "Content-Length",
								Type:   "integer",
								Format: "<length of blob>",
							},
						},
						PathParameters: []ParameterDescriptor{
							nameParameterDescriptor,
						},
						QueryParameters: []ParameterDescriptor{
							{
								Name:   "digest",
								Type:   "query",
								Format: "<digest>",
								Description: `Digest of uploaded blob. If present, the upload will be completed, in a single reques
								t, with contents of the request body as the resulting blob.`,
							},
						},
					},
					{
						Name:        "Initiate Resumable Blob Upload",
						Description: "Initiate a resumable blob upload with an empty request body.",
						Headers: []ParameterDescriptor{
							contentLengthZeroHeader,
						},
						PathParameters: []ParameterDescriptor{
							nameParameterDescriptor,
						},
					},
					{
						Name:        "Mount Blob",
						Description: "Mount a blob identified by the `mount` parameter from another repository.",
						Headers: []ParameterDescriptor{
							contentLengthZeroHeader,
						},
						PathParameters: []ParameterDescriptor{
							nameParameterDescriptor,
						},
						QueryParameters: []ParameterDescriptor{
							{
								Name:        "mount",
								Type:        "query",
								Format:      "<digest>",
								Description: `Digest of blob to mount from the source repository.`,
							},
							{
								Name:        "from",
								Type:        "query",
								Format:      "<repository name>",
								Description: `Name of the source repository.`,
							},
						},
					},
				},
			},
		},
	},

	{
		Name:        RouteNameManifest,
		Path:        "/v1/{name:" + reference.NameRegexp.String() + "}/manifests/{reference:" + digest.DigestRegexp.String() + "}",
		Description: "Create, update, delete and retrieve manifests.",
		Methods: []MethodDescriptor{
			{
				Method:      "GET",
				Description: "Fetch the manifest identified by `name` and `reference`",
				Requests: []RequestDescriptor{
					{
						Headers: []ParameterDescriptor{
							hostHeader,
							authHeader,
						},
						PathParameters: []ParameterDescriptor{
							nameParameterDescriptor,
							referenceParameterDescriptor,
						},
					},
				},
			},
			{
				Method:      "PUT",
				Description: "Put the manifest identified by `name` and `reference`",
				Requests: []RequestDescriptor{
					{
						Headers: []ParameterDescriptor{
							hostHeader,
							authHeader,
						},
						PathParameters: []ParameterDescriptor{
							nameParameterDescriptor,
							referenceParameterDescriptor,
						},
					},
				},
			},
			{
				Method:      "DELETE",
				Description: "Delete the manifest identified by `name` and `reference`",
				Requests: []RequestDescriptor{
					{
						Headers: []ParameterDescriptor{
							hostHeader,
							authHeader,
						},
						PathParameters: []ParameterDescriptor{
							nameParameterDescriptor,
							referenceParameterDescriptor,
						},
					},
				},
			},
		},
	},
	{
		Name: RouteNameHelp,
		Path: "/v1/help/",
		Methods: []MethodDescriptor{
			{
				Method:      "GET",
				Description: "Print API Help Message for V1.",
			},
		},
	},
	{
		Name:    RouteNameBlob,
		Path:    "/v1/{name:" + reference.NameRegexp.String() + "}/blobs/",
		Methods: []MethodDescriptor{},
	},
}

func init() {
	RouteDescriptorsMap = make(map[string]RouteDescriptor, len(routeDescriptors))
	for _, descriptor := range routeDescriptors {
		RouteDescriptorsMap[descriptor.Name] = descriptor
	}
}
