//
// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//

package ldap

import (
	"fmt"

	"my10c.ldap/logs"
	"my10c.ldap/vars"
	ldapv3 "gopkg.in/ldap.v2"
)

// delete a ldap record
func (conn *Connection) Delete(recordId, recordType string) bool {
	delReq := ldapv3.NewDelRequest(vars.WorkRecord.DN, []ldapv3.Control{})
	if err := conn.Conn.Del(delReq); err != nil {
		msg := fmt.Sprintf("Error deleting the %s %s, error %s", recordType, recordId, err.Error())
		logs.Log(msg, "ERROR")
		return false
	}
	msg := fmt.Sprintf("The %s %s has been deleted", recordType, recordId)
	logs.Log(msg, "INFO")
	return true
}
