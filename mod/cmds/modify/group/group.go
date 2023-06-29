//
// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//

package modify

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"my10c.ldap/cmds/common"
	"my10c.ldap/ldap"
	"my10c.ldap/vars"
	"github.com/my10c/packages-go/print"
	ldapv3 "gopkg.in/ldap.v2"
)

var (
	valueEntered string
	modCount     int = 0
	p                = print.New()
)

// remove the user from groups its belong to if this was set during the template input
func deleteGroupEntries(conn *ldap.Connection, groupName string, funcs *vars.Funcs) {
	for _, userID := range vars.WorkRecord.GroupDelList {
		vars.WorkRecord.ID = userID
		if !conn.RemoveFromGroups() {
			funcs.P.PrintRed(fmt.Sprintf("\tFailed to remove the user %s from the group %s, check the log file\n",
				userID, vars.WorkRecord.DN))
		} else {
			funcs.P.PrintGreen(fmt.Sprintf("\tUser %s removed from group %s\n", userID, groupName))
		}
	}
}

// add the user to groups if this was set during the template input
func addGroupEntries(conn *ldap.Connection, groupName string, funcs *vars.Funcs) {
	for _, userID := range vars.WorkRecord.GroupAddList {
		vars.WorkRecord.ID = userID
		if !conn.AddToGroup() {
			funcs.P.PrintRed(fmt.Sprintf("\n\tFailed to add the user %s to the group %s, check the log file\n",
				userID, vars.WorkRecord.DN))
		} else {
			funcs.P.PrintGreen(fmt.Sprintf("\tUser %s added to group %s\n", userID, groupName))
		}
	}
}

// group modigy template
func modifyGroup(conn *ldap.Connection, records *ldapv3.SearchResult, funcs *vars.Funcs) bool {
	orgGroup := vars.WorkRecord.ID
	funcs.P.PrintPurple(fmt.Sprintf("\tUsing group: %s\n", orgGroup))
	funcs.P.PrintYellow(fmt.Sprintf("\tPress enter to leave the value unchanged\n"))
	fmt.Printf("\t%s\n", funcs.P.PrintLine(vars.Purple, 50))
	vars.WorkRecord.DN = fmt.Sprintf("cn=%s,%s", vars.WorkRecord.ID, conn.Config.ServerValues.GroupDN)
	for idx, entry := range records.Entries {
		if len(records.Entries[idx].GetAttributeValue("gidNumber")) != 0 {
			vars.WorkRecord.MemberType = "memberUid"
			for _, member := range entry.GetAttributeValues("memberUid") {
				funcs.P.PrintCyan(fmt.Sprintf("\tmember: %s\n", member))
			}
		} else {
			vars.WorkRecord.MemberType = "member"
			for _, member := range entry.GetAttributeValues("member") {
				funcs.P.PrintCyan(fmt.Sprintf("\tmember: %s\n", member))
			}
		}
	}

	funcs.P.PrintRed("\n\tEnter the user(s) to be deleted, select from the list above, (default to skip)\n")
	reader := bufio.NewReader(os.Stdin)
	for true {
		fmt.Printf("\tUser : ")
		valueEntered, _ = reader.ReadString('\n')
		valueEntered = strings.TrimSuffix(valueEntered, "\n")
		if valueEntered != "" {
			if vars.WorkRecord.MemberType == "member" {
				valueEntered = fmt.Sprintf("uid=%s,%s", valueEntered, conn.Config.ServerValues.UserDN)
			}
			vars.WorkRecord.GroupDelList = append(vars.WorkRecord.GroupDelList, valueEntered)
		} else {
			break
		}
	}

	funcs.P.PrintGreen("\n\tEnter the user(s) to be added, (default to skip)\n")
	for true {
		fmt.Printf("\tUser : ")
		valueEntered, _ := reader.ReadString('\n')
		valueEntered = strings.TrimSuffix(valueEntered, "\n")
		if valueEntered != "" {
			if vars.WorkRecord.MemberType == "member" {
				valueEntered = fmt.Sprintf("uid=%s,%s", valueEntered, conn.Config.ServerValues.UserDN)
			}
			vars.WorkRecord.GroupAddList = append(vars.WorkRecord.GroupAddList, valueEntered)
		} else {
			break
		}
	}

	modCount = len(vars.WorkRecord.GroupDelList) + len(vars.WorkRecord.GroupAddList)
	if modCount == 0 {
		funcs.P.PrintBlue(fmt.Sprintf("\n\tNo change, no modification made to group %s\n", orgGroup))
		return false
	}

	if len(vars.WorkRecord.GroupDelList) > 0 {
		deleteGroupEntries(conn, orgGroup, funcs)
	}
	if len(vars.WorkRecord.GroupAddList) > 0 {
		addGroupEntries(conn, orgGroup, funcs)
	}
	return true
}

func Modify(conn *ldap.Connection, funcs *vars.Funcs) {
	fmt.Printf("\t%s\n", funcs.P.PrintHeader(vars.Blue, vars.Purple, "Modify Group", 18, true))
	vars.SearchResultData.WildCardSearchBase = vars.GroupWildCardSearchBase
	vars.SearchResultData.RecordSearchbase = vars.GroupWildCardSearchBase
	vars.SearchResultData.DisplayFieldID = vars.GroupDisplayFieldID
	if common.GetObjectRecord(conn, true, "group", funcs) {
		modifyGroup(conn, vars.SearchResultData.SearchResult, funcs)
	}
	fmt.Printf("\t%s\n", funcs.P.PrintLine(vars.Purple, 50))
}
