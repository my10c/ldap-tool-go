//
// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//

package delete

import (
	"fmt"

	"my10c.ldap/cmds/common"
	"my10c.ldap/ldap"
	"my10c.ldap/vars"
)

// once an user has been deleted, we need to make sure
// it are removed from all the group its belong too
func removeUserFromGroups(conn *ldap.Connection, funcs *vars.Funcs) {
	var groupsList []string
	userUID := vars.WorkRecord.ID
	conn.SearchInfo.SearchBase = "(&(objectClass=posixGroup))"
	conn.SearchInfo.SearchAttribute = []string{"cn", "memberUid"}

	records, _ := conn.Search()
	for idx, entry := range records.Entries {
		for _, member := range entry.GetAttributeValues("memberUid") {
			if member == userUID {
				groupsList = append(groupsList, records.Entries[idx].GetAttributeValue("cn"))
			}
		}
	}
	if len(groupsList) > 0 {
		for _, groupName := range groupsList {
			vars.WorkRecord.DN = fmt.Sprintf("cn=%s,%s", groupName, conn.Config.ServerValues.GroupDN)
			if !conn.RemoveFromGroups() {
				funcs.P.PrintRed(fmt.Sprintf("User % was not remobe from the group %s, check the log...\n",
					userUID, groupName))
			}
		}
	}
}

func Delete(conn *ldap.Connection, funcs *vars.Funcs) {
	fmt.Printf("\t%s\n", funcs.P.PrintHeader(vars.Blue, vars.Purple, "Delete User", 18, true))
	vars.SearchResultData.WildCardSearchBase = vars.UserWildCardSearchBase
	vars.SearchResultData.RecordSearchbase = vars.UserWildCardSearchBase
	vars.SearchResultData.DisplayFieldID = vars.UserDisplayFieldID
	// we only handle posix group
	vars.WorkRecord.GroupType = "posix"
	vars.WorkRecord.MemberType = "memberUid"
	if common.GetObjectRecord(conn, true, "user", funcs) {
		common.DeleteObjectRecord(conn, vars.SearchResultData.SearchResult, "user", funcs)
		removeUserFromGroups(conn, funcs)
	}
	fmt.Printf("\t%s\n", funcs.P.PrintLine(vars.Purple, 50))
}
