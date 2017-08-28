package v1

import (
	"octlink/rstore/reference"
)

// RouteDescriptorsMap for router map
var RouteDescriptorsMap map[string]RouteDescriptor

// RouteDescriptor describes a route specified by name.
type RouteDescriptor struct {
	// Name is the name of the route, as specified in RouteNameXXX exports.
	// These names a should be considered a unique reference for a route. If
	// the route is registered with gorilla, this is the name that will be
	// used.
	Name string `json:"name"`

	// path is a gorilla/mux-compatible regexp that can be used to match the
	// route. For any incoming method and path, only one route descriptor
	// should match.
	path       string
	PathSimple string `json:"pathSimple"`

	// Description for this router.
	Description string `json:"description"`

	// Methods should describe the various HTTP methods that may be used on
	// this route, including request and response formats.
	Methods []MethodDescriptor `json:"methods"`
}

// MethodDescriptor provides a description of the requests that may be
// conducted with the target method.
type MethodDescriptor struct {

	// Method is an HTTP method, such as GET, PUT or POST.
	Method string `json:"method"`

	// Description should provide an overview of the functionality provided by
	// the covered method, suitable for use in documentation. Use of markdown
	// here is encouraged.
	Description string `json:"description"`

	// Requests is a slice of request descriptors enumerating how this
	// endpoint may be used.
	Requests []RequestDescriptor `json:"request"`
}

// RequestDescriptor covers a particular set of headers and parameters that
// can be carried out with the parent method. Its most helpful to have one
// RequestDescriptor per API use case.
type RequestDescriptor struct {
	// Name provides a short identifier for the request, usable as a title or
	// to provide quick context for the particular request.
	Name string `json:"name"`

	// Description should cover the requests purpose, covering any details for
	// this particular use case.
	Description string `json:"description"`

	// Headers describes headers that must be used with the HTTP request.
	Headers []ParameterDescriptor `json:"headers"`

	// PathParameters enumerate the parameterized path components for the
	// given request, as defined in the route's regular expression.
	PathParameters []ParameterDescriptor `json:"pathParameterDescriptor"`

	// QueryParameters provides a list of query parameters for the given
	// request.
	QueryParameters []ParameterDescriptor `json:"parameterDescriptor"`
}

// ParameterDescriptor describes the format of a request parameter, which may
// be a header, path parameter or query parameter.
type ParameterDescriptor struct {
	// Name is the name of the parameter, either of the path component or
	// query parameter.
	Name string `json:"name"`

	// Type specifies the type of the parameter, such as string, integer, etc.
	Type string `json:"type"`

	// format is a specifying the string format accepted by this parameter.
	format string

	// Description provides a human-readable description of the parameter.
	Description string `json:"description"`

	// Required or not
	Required bool `json:"required"`

	// Examples provides multiple examples for the values that might be valid
	// for this parameter.
	Examples []string `json:"examples"`
}

var (
	nameParameterDescriptor = ParameterDescriptor{
		Name:        "name",
		Type:        "string",
		format:      reference.NameRegexp.String(),
		Required:    true,
		Description: `Name of the target repository.`,
	}

	referenceParameterDescriptor = ParameterDescriptor{
		Name:        "reference",
		Type:        "string",
		format:      reference.TagRegexp.String(),
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
		format:      reference.DigestRegexp.String(),
		Description: `Digest of desired blob.`,
	}

	contentLengthZeroHeader = ParameterDescriptor{
		Name:        "Content-Length",
		Description: "The `Content-Length` header must be zero and the body must be empty.",
		Type:        "integer",
		format:      "0",
	}

	digestHeader = ParameterDescriptor{
		Name:        "Docker-Content-Digest",
		Description: "Digest of the targeted content for the request.",
		Type:        "digest",
		format:      "<digest>",
	}
)

// RouteDescriptors for route descriptor list
var RouteDescriptors = []RouteDescriptor{
	{
		Name:        RouteNameBase,
		path:        "/v1/",
		PathSimple:  "/v1/",
		Description: "Base V1 API route",
		Methods: []MethodDescriptor{
			{
				Method:      "GET",
				Description: "Check implements API V1.",
			},
		},
	},
	{
		Name:        RouteNameBlobUpload,
		path:        "/v1/{name:" + reference.NameRegexp.String() + "}/blobs/uploads/{digest:" + reference.DigestRegexp.String() + "}",
		PathSimple:  "/v1/{name}/blobs/uploads/{digest}",
		Description: "Upload blob by name and digest.",
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
								format: "<length of blob>",
							},
						},
						PathParameters: []ParameterDescriptor{
							nameParameterDescriptor,
							digestPathParameter,
						},
					},
				},
			},
		},
	},

	{
		Name:        RouteNameManifest,
		path:        "/v1/{name:" + reference.NameRegexp.String() + "}/manifests/{digest:" + reference.DigestRegexp.String() + "}",
		PathSimple:  "/v1/{name}/manifests/{digest}/",
		Description: "Create, update, delete and retrieve manifests.",
		Methods: []MethodDescriptor{
			{
				Method:      "GET",
				Description: "Fetch the manifest identified by `name` and `digest`",
				Requests: []RequestDescriptor{
					{
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
		Name:       RouteNameHelp,
		path:       "/v1/help/",
		PathSimple: "/v1/help/",
		Methods: []MethodDescriptor{
			{
				Method:      "GET",
				Description: "Print API Help Message for V1.",
			},
		},
	},

	{
		Name:       RouteNameHelpModule,
		path:       "/v1/help/{module:" + reference.NameRegexp.String() + "}",
		PathSimple: "/v1/help/{module}",
		Methods: []MethodDescriptor{
			{
				Method:      "GET",
				Description: "Print API Help Message for V1.",
			},
		},
	},

	{
		Name:       RouteNameBlobsManifest,
		path:       "/v1/{name:" + reference.NameRegexp.String() + "}/blobsmanifest/{blobsum:" + reference.DigestRegexp.String() + "}",
		PathSimple: "/v1/{name}/blobsmanifest/{blobsum}",
		Methods: []MethodDescriptor{
			{
				Method:      "GET",
				Description: "Get blobsmanifest config for blobs pulling.",
			},
		},
	},

	{
		Name:        RouteNameBlob,
		path:        "/v1/{name:" + reference.NameRegexp.String() + "}/blobs/{digest:" + reference.NameRegexp.String() + "}",
		PathSimple:  "/v1/{name}/blobs/{digest}",
		Description: "Operations on blobs identified by `name` and `digest`. Used to fetch or delete layers by digest.",
		Methods: []MethodDescriptor{
			{
				Method:      "GET",
				Description: "Retrieve the blob from the registry identified by `digest`. A `HEAD` request can also be issued to this endpoint to obtain resource information without receiving all data.",
				Requests: []RequestDescriptor{
					{
						Name: "Fetch Blob",
						PathParameters: []ParameterDescriptor{
							nameParameterDescriptor,
							digestPathParameter,
						},
					},
					{
						Name:        "Fetch Blob Part",
						Description: "This endpoint may also support RFC7233 compliant range requests. Support can be detected by issuing a HEAD request. If the header `Accept-Range: bytes` is returned, range requests can be used to fetch partial content.",
						Headers: []ParameterDescriptor{
							{
								Name:        "Range",
								Type:        "string",
								Description: "HTTP Range header specifying blob chunk.",
							},
						},
						PathParameters: []ParameterDescriptor{
							nameParameterDescriptor,
							digestPathParameter,
						},
					},
				},
			},
			{
				Method:      "DELETE",
				Description: "Delete the blob identified by `name` and `digest`",
				Requests: []RequestDescriptor{
					{
						PathParameters: []ParameterDescriptor{
							nameParameterDescriptor,
							digestPathParameter,
						},
					},
				},
			},
		},
	},
}

func init() {
	RouteDescriptorsMap = make(map[string]RouteDescriptor, len(RouteDescriptors))
	for _, descriptor := range RouteDescriptors {
		RouteDescriptorsMap[descriptor.Name] = descriptor
	}
}
