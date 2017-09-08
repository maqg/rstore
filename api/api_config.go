package api

import "octlink/rstore/modules/systemconfig"

// ShowSystemConfig to add image by API
func ShowSystemConfig(paras *Paras) *Response {
	return &Response{
		Data: systemconfig.GetSystemConfig(),
	}
}
