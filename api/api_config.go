package api

import "octlink/rstore/modules/systemconfig"

// ShowSystemConfig to add image by API
func ShowSystemConfig(paras *Paras) *Response {
	resp := new(Response)
	resp.Data = systemconfig.GetSystemConfig()
	return resp
}
