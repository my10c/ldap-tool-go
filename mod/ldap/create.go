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

// create a ldap user
func (conn *Connection) CreateUser() bool {
	newUserReq := ldapv3.NewAddRequest(vars.WorkRecord.DN)
	newUserReq.Attribute("objectClass", vars.UserObjectClass)
	for _, fieldName := range vars.UserFields {
		if fieldName != "userPassword" {
			newUserReq.Attribute(fieldName, []string{vars.WorkRecord.Fields[fieldName]})
		}
	}
	if err := conn.Conn.Add(newUserReq); err != nil {
		msg = fmt.Sprintf("Error creating the user %s error %s",
			vars.WorkRecord.Fields["uid"], err.Error())
		logs.Log(msg, "ERROR")
		return false
	}
	msg = fmt.Sprintf("The user %s has been created", vars.WorkRecord.Fields["uid"])
	logs.Log(msg, "INFO")

	//set password
	return conn.SetPassword()
}

// create a ldap group
func (conn *Connection) CreateGroup() bool {
	newGroupReq := ldapv3.NewAddRequest(vars.WorkRecord.Fields["dn"])
	newGroupReq.Attribute("objectClass", []string{vars.WorkRecord.Fields["objectClass"]})
	newGroupReq.Attribute("cn", []string{vars.WorkRecord.Fields["cn"]})
	if vars.WorkRecord.Fields["objectClass"] == "posixGroup" {
		newGroupReq.Attribute("gidNumber", []string{vars.WorkRecord.Fields["gidNumber"]})
	}
	if vars.WorkRecord.Fields["objectClass"] == "groupOfNames" {
		newGroupReq.Attribute("member", []string{vars.WorkRecord.Fields["member"]})
	}
	if err := conn.Conn.Add(newGroupReq); err != nil {
		msg = fmt.Sprintf("Error creating the group %s error %s",
			vars.WorkRecord.Fields["cn"], err.Error())
		logs.Log(msg, "ERROR")
		return false
	}
	msg = fmt.Sprintf("The group %s has been created", vars.WorkRecord.Fields["cn"])
	logs.Log(msg, "INFO")
	return true
}

// create a ldap sudo rule
func (conn *Connection) CreateSudoRule() bool {
	newSudoRuleReq := ldapv3.NewAddRequest(vars.WorkRecord.Fields["dn"])
	newSudoRuleReq.Attribute("objectClass", []string{vars.WorkRecord.Fields["objectClass"]})
	for _, field := range vars.SudoFields {
		if len(vars.WorkRecord.Fields[field]) > 0 {
			newSudoRuleReq.Attribute(field, []string{vars.WorkRecord.Fields[field]})
		}
	}
	if err := conn.Conn.Add(newSudoRuleReq); err != nil {
		msg = fmt.Sprintf("Error creating the sudo rule %s error %s",
			vars.WorkRecord.Fields["cn"], err.Error())
		logs.Log(msg, "ERROR")
		return false
	}
	msg = fmt.Sprintf("The sudo rule %s has been created", vars.WorkRecord.Fields["cn"])
	logs.Log(msg, "INFO")
	return true
}
