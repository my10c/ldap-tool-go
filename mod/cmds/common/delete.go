//
// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//

package common

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"my10c.ldap/ldap"
	"my10c.ldap/vars"
	"github.com/my10c/packages-go/readinput"
	ldapv3 "gopkg.in/ldap.v2"
)

func DeleteObjectRecord(conn *ldap.Connection, records *ldapv3.SearchResult, objectType string, funcs *vars.Funcs) {
	reader := bufio.NewReader(os.Stdin)

	switch objectType {
	case "user":
		vars.ObjectID = records.Entries[0].GetAttributeValue("uid")
		vars.ProtectedList = conn.Config.GroupValues.Groups
		vars.WorkRecord.DN = fmt.Sprintf("uid=%s,%s", vars.ObjectID, conn.Config.ServerValues.UserDN)
		vars.WorkRecord.ID = vars.ObjectID

	case "group":
		vars.ObjectID = records.Entries[0].GetAttributeValue("cn")
		vars.ProtectedList = append(conn.Config.GroupValues.Groups, conn.Config.GroupValues.SpecialGroups...)
		vars.WorkRecord.DN = fmt.Sprintf("cn=%s,%s", vars.ObjectID, conn.Config.ServerValues.GroupDN)
		vars.WorkRecord.ID = vars.ObjectID

	case "sudo rules":
		vars.ObjectID = records.Entries[0].GetAttributeValue("cn")
		vars.ProtectedList = conn.Config.SudoValues.ExcludeSudo
		vars.WorkRecord.DN = fmt.Sprintf("cn=%s,%s", vars.ObjectID, conn.Config.SudoValues.SudoersBase)
		vars.WorkRecord.ID = vars.ObjectID
	}

	if objectType != "user" {
		if funcs.I.IsInList(vars.ProtectedList, vars.ObjectID) {
			funcs.P.PrintRed(fmt.Sprintf("\n\tGiven %s %s is protected and can not be deleted, aborting...\n\n",
				objectType, vars.ObjectID))
			return
		}
	}

	funcs.P.PrintRed(fmt.Sprintf("\n\tGiven %s %s will be delete, this can not be undo!\n", objectType, vars.ObjectID))
	funcs.P.PrintYellow(fmt.Sprintf("\tContinue (default to N)? [y/n]: "))
	continueDelete, _ := reader.ReadString('\n')
	continueDelete = strings.ToLower(strings.TrimSuffix(continueDelete, "\n"))
	if readinput.ReadYN(continueDelete, false) == true {
		if !conn.Delete(vars.ObjectID, objectType) {
			funcs.P.PrintRed(fmt.Sprintf("\n\tFailed to delete the %s %s, check the log file\n",
				objectType, vars.ObjectID))
		} else {
			funcs.P.PrintGreen(fmt.Sprintf("\n\tGiven %s %s has been deleted\n", objectType, vars.ObjectID))
		}
	} else {
		funcs.P.PrintBlue(fmt.Sprintf("\n\tDeletion of the %s %s cancelled\n", objectType, vars.ObjectID))
	}
	return
}
