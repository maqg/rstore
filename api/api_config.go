package api

import "octlink/rstore/modules/config"

// ShowSystemConfig to add image by API
func ShowSystemConfig(paras *Paras) *Response {
	resp := new(Response)
	resp.Data = config.GetSystemConfig()
	return resp
}
