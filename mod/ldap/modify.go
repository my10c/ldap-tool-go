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

var (
	msg string
)

// remove an user from a group
func (conn *Connection) RemoveFromGroups() bool {
	// posix uses user name
	// groupOfNames uses user's full dn
	removeReq := ldapv3.NewModifyRequest(vars.WorkRecord.DN)
	removeReq.Delete(vars.WorkRecord.MemberType, []string{vars.WorkRecord.ID})
	if err := conn.Conn.Modify(removeReq); err != nil {
		msg = fmt.Sprintf("Error removing the user %s from group %s, error %s",
			vars.WorkRecord.ID, vars.WorkRecord.DN, err.Error())
		logs.Log(msg, "ERROR")
		return false
	}
	msg = fmt.Sprintf("The %s %s has been modify", vars.WorkRecord.ID, vars.WorkRecord.DN)
	logs.Log(msg, "INFO")
	return true
}

// add an user to a group
func (conn *Connection) AddToGroup() bool {
	// posix uses user name
	// groupOfNames uses user's full dn
	addReq := ldapv3.NewModifyRequest(vars.WorkRecord.DN)
	addReq.Add(vars.WorkRecord.MemberType, []string{vars.WorkRecord.ID})
	if err := conn.Conn.Modify(addReq); err != nil {
		msg = fmt.Sprintf("Error adding the user %s to group %s, error %s",
			vars.WorkRecord.ID, vars.WorkRecord.DN, err.Error())
		logs.Log(msg, "ERROR")
		return false
	}
	msg = fmt.Sprintf("The %s %s has been modify", vars.WorkRecord.ID, vars.WorkRecord.DN)
	logs.Log(msg, "INFO")
	return true
}

// modify an user ldap record
func (conn *Connection) ModifyUser() bool {
	var passChanged bool = false
	modifyRecord := ldapv3.NewModifyRequest(vars.WorkRecord.DN)
	for fieldName, fieldValue := range vars.WorkRecord.Fields {
		if fieldName != "userPassword" {
			modifyRecord.Replace(fieldName, []string{fieldValue})
		}
		if fieldName == "userPassword" {
			passChanged = true
		}
	}
	if err := conn.Conn.Modify(modifyRecord); err != nil {
		msg = fmt.Sprintf("Error modifying the user %s, error %s",
			vars.WorkRecord.ID, err.Error())
		logs.Log(msg, "ERROR")
		return false
	}
	if passChanged {
		return conn.SetPassword()
	}
	return true
}

// delete a sudo rule
func (conn *Connection) DeleteSudoRule() bool {
	delSudoRule := ldapv3.NewModifyRequest(vars.WorkRecord.DN)
	for fieldName, _ := range vars.WorkRecord.SudoDelList {
		for _, value := range vars.WorkRecord.SudoDelList[fieldName] {
			delSudoRule.Delete(fieldName, []string{value})
		}
	}
	if err := conn.Conn.Modify(delSudoRule); err != nil {
		msg = fmt.Sprintf("Error deleting some of the entries of sudo rule %s, error %s",
			vars.WorkRecord.ID, err.Error())
		logs.Log(msg, "ERROR")
		return false
	}
	msg = fmt.Sprintf("The sudo rule %s entries has been modified", vars.WorkRecord.ID)
	logs.Log(msg, "INFO")
	return true
}

// add a sudo rule
func (conn *Connection) AddSudoRule() bool {
	addSudoRule := ldapv3.NewModifyRequest(vars.WorkRecord.DN)
	for fieldName, _ := range vars.WorkRecord.SudoAddList {
		for _, value := range vars.WorkRecord.SudoAddList[fieldName] {
			addSudoRule.Add(fieldName, []string{value})
		}
	}
	if err := conn.Conn.Modify(addSudoRule); err != nil {
		msg = fmt.Sprintf("Error adding some of the entries of sudo rule %s, error %s",
			vars.WorkRecord.ID, err.Error())
		logs.Log(msg, "ERROR")
		return false
	}
	msg = fmt.Sprintf("The sudo rule %s entries has been modified", vars.WorkRecord.ID)
	logs.Log(msg, "INFO")
	return true
}
