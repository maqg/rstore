package api

import (
	"octlink/mirage/src/modules/session"
	"octlink/mirage/src/utils/httpresponse"
	"octlink/mirage/src/utils/merrors"
	"octlink/mirage/src/utils/octlog"
	"octlink/mirage/src/utils/octmysql"

	"github.com/gin-gonic/gin"
)

type ApiResponse struct {
	Error    int         `json:"error"`
	ErrorLog string      `json:"errorLog"`
	Data     interface{} `json:"data"`
}

/*
{
	"module": "octlink.mirage.center.host.APIAddHost",
	"paras": {
		"ip": "kk",
		"account": "root",
		"password": ""
	},
	"async": false,
	"session": {
		"uuid": "00000000000000000000000000000000",
		"skey": "00000000000000000000000000000000"
	}
}
*/
type InputParas struct {
	Module  string
	Api     string
	Paras   map[string]interface{}
	Async   bool
	Session map[string]interface{}
}

type ApiParas struct {
	Proto   *ApiProto
	Session *session.Session
	InParas *InputParas
	Db      *octmysql.OctMysql
}

func (api *Api) ApiTest(c *gin.Context) {
	httpresponse.Ok(c, "Api Server is Running")
}

var GApiServices map[string]*ApiService

type ApiService struct {
	Name    string                       `json:"name"`
	Handler func(*ApiParas) *ApiResponse `json:"handler"`
}

func GetApiService(key string) *ApiService {
	service, ok := GApiServices[key]
	if !ok {
		octlog.Error("no service for %s found\n", key)
		return nil
	}

	return service
}

var SessionExceptions = []string{
	"octlink.mirage.center.account.APILoginByAccount",
	"octlink.mirage.center.user.APILoginByUser",
}

func NeedSessionCheck(api string) bool {
	for _, tmp_api := range SessionExceptions {
		if api == tmp_api {
			return false
		}
	}

	return true
}

func getApiParas(c *gin.Context) (*ApiParas, int) {

	var sid string
	var apiParas *ApiParas = new(ApiParas)

	c.BindJSON(&apiParas.InParas)

	octlog.Debug("got api %s\n", apiParas.InParas.Api)

	if apiParas.InParas.Api == "" {
		octlog.Error("got null api\n")
		return nil, merrors.ERR_NO_SUCH_API
	}

	proto := FindApiProto(apiParas.InParas.Api)
	if proto == nil {
		octlog.Error("no api proto found for %s\n",
			apiParas.InParas.Api)
		return nil, merrors.ERR_NO_SUCH_API
	}

	apiParas.Proto = proto
	apiParas.Db = new(octmysql.OctMysql)

	if NeedSessionCheck(apiParas.InParas.Api) {
		sid = apiParas.InParas.Session["uuid"].(string)
	} else {
		sid = session.SESSION_DEFAULT_ID
	}
	octlog.Debug("found session id " + sid)

	session := session.FindSession(apiParas.Db, sid)
	if session == nil {
		octlog.Error("not session found for this id" + sid)
		return nil, merrors.ERR_USER_NOT_LOGIN
	}

	apiParas.Session = session

	return apiParas, 0
}

func checkParas(apiParas *ApiParas) (int, string) {

	protoParas := apiParas.Proto.Paras

	for i := 0; i < len(protoParas); i++ {

		protoParam := protoParas[i]
		inParam := apiParas.InParas.Paras[protoParam.Name]

		// if paras have default value and no input sepecified, set a default value
		if protoParam.Default != PARAM_NOT_NULL && inParam == nil {
			apiParas.InParas.Paras[protoParam.Name] = protoParam.Default
		}

		octlog.Debug("param:%s, default:%s, value:%s\n", protoParam.Name,
			protoParam.Default, inParam)

		if protoParam.Default == PARAM_NOT_NULL && inParam.(string) == "" {
			errorMsg := "paras \"" + protoParam.Name + "\" must be specified"
			return merrors.ERR_NOT_ENOUGH_PARAS, errorMsg
		}
	}

	return merrors.ERR_OCT_SUCCESS, ""
}

func (api *Api) ApiDispatch(c *gin.Context) {

	octlog.Error("got api request\n")

	paras, err := getApiParas(c)
	if paras == nil {
		octlog.Error("No match proto found\n")
		httpresponse.Error(c, err, nil)
		return
	}

	service := GetApiService(paras.InParas.Api)
	if service == nil {
		octlog.Error("No match service found\n")
		httpresponse.Error(c, merrors.ERR_NO_SUCH_API, paras.InParas.Api)
		return
	}

	ret, msg := checkParas(paras)
	if ret != merrors.ERR_OCT_SUCCESS {
		octlog.Error("Not Enough Paras\n")
		httpresponse.Error(c, merrors.ERR_NOT_ENOUGH_PARAS, msg)
		return
	}

	resp := service.Handler(paras)
	defer paras.Db.Close()

	if resp.Error == 0 {
		httpresponse.Ok(c, resp.Data)
	} else {
		httpresponse.Error(c, resp.Error, resp.ErrorLog)
	}
}
