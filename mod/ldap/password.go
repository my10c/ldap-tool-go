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

// set an user's ldap passwod
func (conn *Connection) SetPassword() bool {
	// once the record is create we need to hash the password
	passwordReq := ldapv3.NewPasswordModifyRequest(
		vars.WorkRecord.DN, "", vars.WorkRecord.Fields["userPassword"])
	if _, err := conn.Conn.PasswordModify(passwordReq); err != nil {
		msg = fmt.Sprintf("Failed setting password for the user %s, error %s", vars.WorkRecord.ID, err.Error())
		logs.Log(msg, "ERROR")
		return false
	}
	msg = fmt.Sprintf("Successfully setting the password for for user %s to %s",
		vars.WorkRecord.ID, vars.WorkRecord.Fields["userPassword"])
	logs.Log(msg, "INFO")
	return true
}
