package api

import (
	"fmt"
	"octlink/mirage/src/modules/account"
	"octlink/mirage/src/modules/session"
	"octlink/mirage/src/utils"
	"octlink/mirage/src/utils/merrors"
	"octlink/mirage/src/utils/octlog"
	"octlink/mirage/src/utils/uuid"
)

func APIAddAccount(paras *ApiParas) *ApiResponse {
	resp := new(ApiResponse)

	newAccount := account.FindAccountByName(paras.Db, paras.InParas.Paras["account"].(string))
	if newAccount != nil {
		logger.Errorf("account %s already exist\n", newAccount.Name)
		resp.Error = merrors.ERR_SEGMENT_ALREADY_EXIST
		return resp
	}

	newAccount = new(account.Account)
	newAccount.Id = uuid.Generate().Simple()
	newAccount.Name = paras.InParas.Paras["account"].(string)
	newAccount.Type = utils.ParasInt(paras.InParas.Paras["type"])
	newAccount.Email = paras.InParas.Paras["email"].(string)
	newAccount.PhoneNumber = paras.InParas.Paras["phoneNumber"].(string)
	newAccount.Password = paras.InParas.Paras["password"].(string)
	newAccount.Desc = paras.InParas.Paras["desc"].(string)

	resp.Error = newAccount.Add(paras.Db)

	return resp
}

func APILoginByAccount(paras *ApiParas) *ApiResponse {

	resp := new(ApiResponse)

	user := paras.InParas.Paras["account"].(string)
	password := paras.InParas.Paras["password"].(string)

	logger.Debugf("Login %s:%s", user, password)

	account := account.FindAccountByName(paras.Db, user)
	if account == nil {
		logger.Errorf("account %s already exist\n", user)
		resp.Error = merrors.ERR_USER_NOT_EXIST
		return resp
	}

	session := account.Login(paras.Db, password)
	if session == nil {
		resp.Error = merrors.ERR_PASSWORD_DONT_MATCH
		return resp
	}

	resp.Data = session

	return resp
}

func APIShowAccount(paras *ApiParas) *ApiResponse {
	resp := new(ApiResponse)

	accountId := paras.InParas.Paras["id"].(string)
	temp := account.FindAccount(paras.Db, accountId)

	if temp == nil {
		resp.Error = merrors.ERR_SEGMENT_NOT_EXIST
		resp.ErrorLog = fmt.Sprintf("user %s not found", accountId)
		return resp
	}

	resp.Data = temp

	octlog.Debug("found User %s", temp.Name)

	return resp
}

func APIUpdateAccount(paras *ApiParas) *ApiResponse {
	resp := new(ApiResponse)

	id := paras.InParas.Paras["id"].(string)

	ac := account.FindAccount(paras.Db, id)
	if ac == nil {
		resp.Error = merrors.ERR_USER_NOT_EXIST
		resp.ErrorLog = "User " + id + "Not Exist"
		return resp
	}

	ac.Email = paras.InParas.Paras["email"].(string)
	ac.PhoneNumber = paras.InParas.Paras["phoneNumber"].(string)
	ac.Desc = paras.InParas.Paras["desc"].(string)

	ret := ac.Update(paras.Db)
	if ret != 0 {
		resp.Error = ret
		return resp
	}

	return resp
}

func APIShowAccountList(paras *ApiParas) *ApiResponse {

	resp := new(ApiResponse)

	octlog.Debug("running in APIShowAccountList\n")

	rows, err := paras.Db.Query("SELECT ID,U_Name FROM tb_account")
	if err != nil {
		logger.Errorf("query account list error %s\n", err.Error())
	}
	defer rows.Close()

	accountList := make([]map[string]string, 0)

	for rows.Next() {
		var account account.Account
		err = rows.Scan(&account.Id, &account.Name)
		if err == nil {
			logger.Debugf("query result: %s:%s\n", account.Id, account.Name)
		} else {
			logger.Errorf("query account list error %s\n", err.Error())
		}
		accountList = append(accountList, account.Brief())
	}

	resp.Data = accountList

	return resp
}

func APIDeleteAccount(paras *ApiParas) *ApiResponse {

	octlog.Debug("running in APIDeleteAccount\n")

	resp := new(ApiResponse)

	account := account.FindAccount(paras.Db,
		paras.InParas.Paras["id"].(string))
	if account == nil {
		resp.Error = merrors.ERR_SEGMENT_NOT_EXIST
		return resp
	}

	resp.Error = account.Delete(paras.Db)

	return resp
}

func APIShowAllAccount(paras *ApiParas) *ApiResponse {
	resp := new(ApiResponse)

	octlog.Debug("running in APIShowAllAccount\n")

	offset := utils.ParasInt(paras.InParas.Paras["start"])
	limit := utils.ParasInt(paras.InParas.Paras["limit"])

	rows, err := paras.Db.Query("SELECT ID,U_Name,U_State,U_Type,U_Email,U_PhoneNumber,"+
		"U_Description,U_CreateTime,U_LastLogin,U_LastSync "+
		"FROM tb_account LIMIT ?,?", offset, limit)
	if err != nil {
		logger.Errorf("query account list error %s\n", err.Error())
		resp.Error = merrors.ERR_DB_ERR
		return resp
	}
	defer rows.Close()

	accountList := make([]account.Account, 0)

	for rows.Next() {
		var account account.Account
		err = rows.Scan(&account.Id, &account.Name, &account.State,
			&account.Type, &account.Email, &account.PhoneNumber, &account.Desc,
			&account.CreateTime, &account.LastLogin, &account.LastSync)
		if err == nil {
			logger.Debugf("query result: %s:%s\n", account.Id,
				account.Name, account.State, account.Type)
		} else {
			logger.Errorf("query account list error %s\n", err.Error())
		}
		accountList = append(accountList, account)
	}

	count := account.GetAccountCount(paras.Db)

	result := make(map[string]interface{}, 3)
	result["total"] = count
	result["count"] = len(accountList)
	result["data"] = accountList

	resp.Data = result

	return resp
}

func APIResetAccountPassword(paras *ApiParas) *ApiResponse {

	resp := new(ApiResponse)

	id := paras.InParas.Paras["id"].(string)
	account := account.FindAccount(paras.Db, id)
	if account == nil {
		resp.Error = merrors.ERR_SEGMENT_NOT_EXIST
		return resp
	}

	password := paras.InParas.Paras["password"].(string)

	ret := account.ResetPassword(paras.Db, password)
	if ret != 0 {
		resp.Error = ret
		return resp
	}

	return resp
}

func APIUpdateAccountPassword(paras *ApiParas) *ApiResponse {

	resp := new(ApiResponse)

	id := paras.InParas.Paras["id"].(string)
	account := account.FindAccount(paras.Db, id)
	if account == nil {
		resp.Error = merrors.ERR_USER_NOT_EXIST
		return resp
	}

	oldPassword := paras.InParas.Paras["oldPassword"].(string)
	newPassword := paras.InParas.Paras["newPassword"].(string)

	ret := account.UpdatePassword(paras.Db, oldPassword, newPassword)
	if ret != 0 {
		octlog.Warn("update password for %s error", account.Name)
		resp.Error = ret
		return resp
	}

	return resp
}

func APIAccountLogOut(paras *ApiParas) *ApiResponse {

	resp := new(ApiResponse)

	sessionId := paras.InParas.Paras["sessionId"].(string)

	session := session.FindSession(paras.Db, sessionId)
	if session != nil {
		session.Delete(paras.Db)
		octlog.Warn("session %s:%s cleared", session.Id, session.UserName)
	}

	return resp
}
