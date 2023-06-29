//
// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//

package ldap

import (
	"fmt"
	"strconv"
)

// get the next user UID from the ldap database
func (conn *Connection) GetNextUID() int {
	var uidValue int
	startUID := conn.Config.DefaultValues.UidStart
	conn.SearchInfo.SearchBase = "(objectClass=person)"
	conn.SearchInfo.SearchAttribute = []string{"uidNumber"}
	records, _ := conn.Search()
	for _, entry := range records.Entries {
		uidValue, _ = strconv.Atoi(entry.GetAttributeValue("uidNumber"))
		if uidValue > startUID {
			startUID = uidValue
		}
	}
	return startUID + 1
}

// get the next group GID from the ldap database
func (conn *Connection) GetNextGID() int {
	var uidValue int
	startGID := conn.Config.DefaultValues.GidStart
	conn.SearchInfo.SearchBase = "(objectClass=posixGroup)"
	conn.SearchInfo.SearchAttribute = []string{"gidNumber"}
	records, _ := conn.Search()
	for _, entry := range records.Entries {
		uidValue, _ = strconv.Atoi(entry.GetAttributeValue("gidNumber"))
		if uidValue == conn.Config.DefaultValues.GroupId {
			// we skip this special gid
			continue
		}
		if uidValue > startGID {
			startGID = uidValue
		}
	}
	return startGID + 1
}

// get all the user uid and UID
func (conn *Connection) GetUsersUID() map[string]string {
	userUIDList := make(map[string]string)
	conn.SearchInfo.SearchBase = "(&(objectClass=inetOrgPerson))"
	conn.SearchInfo.SearchAttribute = []string{"uidNumber", "uid"}
	records, _ := conn.Search()
	for _, uidNumber := range records.Entries {
		userUIDList[uidNumber.GetAttributeValue("uidNumber")] = uidNumber.GetAttributeValue("uid")
	}
	return userUIDList
}

// get all user uid
func (conn *Connection) GetAllUsers() []string {
	var usersList []string
	usersNameUid := conn.GetUsersUID()
	for userUID, _ := range usersNameUid {
		usersList = append(usersList, usersNameUid[userUID])
	}
	return usersList
}

// get the groups an user belong to
func (conn *Connection) GetUserGroups(userID, userDN string) int {
	conn.SearchInfo.SearchBase =
		fmt.Sprintf("(|(&(objectClass=posixGroup)(memberUid=%s))(&(objectClass=groupOfNames)(member=%s)))",
			userID, userDN)
	conn.SearchInfo.SearchAttribute = []string{"dn"}
	records, recordsCount := conn.Search()
	for _, entry := range records.Entries {
		conn.Record.UserGroups = append(conn.Record.UserGroups, entry.DN)
	}
	return recordsCount
}

// get the group of which a user does not belong to
func (conn *Connection) GetAvailableGroups(userID, userDN string) int {
	conn.SearchInfo.SearchBase =
		fmt.Sprintf("(|(&(objectClass=posixGroup)(!memberUid=%s))(&(objectClass=groupOfNames)(!member=%s)))",
			userID, userDN)
	conn.SearchInfo.SearchAttribute = []string{"dn"}
	records, recordsCount := conn.Search()
	for _, entry := range records.Entries {
		conn.Record.AvailableGroups = append(conn.Record.AvailableGroups, entry.DN)
	}
	return recordsCount
}

// get all group and their type: posix or groupOfNames
func (conn *Connection) GetGroupType() map[string][]string {
	result := make(map[string][]string)
	conn.SearchInfo.SearchBase = "(&(objectClass=posixGroup))"
	conn.SearchInfo.SearchAttribute = []string{"dn"}
	records, _ := conn.Search()
	for _, posix := range records.Entries {
		result["posixGroup"] = append(result["posixGroup"], posix.DN)
	}
	conn.SearchInfo.SearchBase = "(&(objectClass=groupOfNames))"
	conn.SearchInfo.SearchAttribute = []string{"dn"}
	records, _ = conn.Search()
	for _, groupOfNames := range records.Entries {
		result["groupOfNames"] = append(result["groupOfNames"], groupOfNames.DN)
	}
	return result
}

// get all group in the ldap database
func (conn *Connection) GetAllGroups() []string {
	groups := conn.GetGroupType()
	return append(groups["posixGroup"], groups["groupOfNames"]...)
}

// get all the posixGroup's group GID
func (conn *Connection) GetAllGroupsGID() map[string]string {
	gitNumberList := make(map[string]string)
	conn.SearchInfo.SearchBase = "(&(objectClass=posixGroup))"
	conn.SearchInfo.SearchAttribute = []string{"gidNumber", "cn"}
	records, _ := conn.Search()
	for _, gidNumber := range records.Entries {
		gitNumberList[gidNumber.GetAttributeValue("gidNumber")] = gidNumber.GetAttributeValue("cn")
	}
	return gitNumberList
}

// get all sudo rule in the ldap database
func (conn *Connection) GetAllSudoRules() []string {
	var sudoRuleList []string
	conn.SearchInfo.SearchBase = "(&(objectClass=sudoRole))"
	conn.SearchInfo.SearchAttribute = []string{"gidNumber", "cn"}
	records, _ := conn.Search()
	for _, sudoRule := range records.Entries {
		sudoRuleList = append(sudoRuleList, sudoRule.GetAttributeValue("cn"))
	}
	return sudoRuleList
}
