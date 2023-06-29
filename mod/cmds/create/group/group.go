//
// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//

package create

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"my10c.ldap/ldap"
	"my10c.ldap/vars"
)

var (
	// group fields
	// required: groupName and groupType
	// required if posix: gidNumber
	// autofilled: objectClass, cn
	// autofilled if not posix: member

	valueEntered string
	nextGID      int
	fields       = []string{"groupName", "groupType"}
	validTypes   = []string{"posix", "groupOfNames"}
)

func createGroup(conn *ldap.Connection, funcs *vars.Funcs) bool {
	allGroupDN := conn.GetAllGroups()

	for _, fieldName := range fields {
		if vars.Template[fieldName].Value != "" {
			fmt.Printf("\t%sDefault to:%s %s%s%s\n",
				vars.Purple, vars.Off, vars.Cyan, vars.Template[fieldName].Value, vars.Off)
		}

		fmt.Printf("\t%s: ", vars.Template[fieldName].Prompt)

		reader := bufio.NewReader(os.Stdin)
		valueEntered, _ = reader.ReadString('\n')
		valueEntered = strings.TrimSuffix(valueEntered, "\n")

		switch fieldName {
		case "groupName":
			groupDN := fmt.Sprintf("cn=%s,%s", valueEntered, conn.Config.ServerValues.GroupDN)
			if funcs.I.IsInList(allGroupDN, groupDN) {
				funcs.P.PrintRed(fmt.Sprintf("\n\tGiven group %s already exist, aborting...\n\n", valueEntered))
				return false
			}
			fmt.Printf("\t%s\n", funcs.P.PrintLine(vars.Purple, 50))
			funcs.P.PrintPurple(fmt.Sprintf("\tUsing Group: %s\n\n", valueEntered))
			vars.WorkRecord.Fields["cn"] = valueEntered
			vars.WorkRecord.Fields["dn"] = groupDN

		case "groupType":
			switch valueEntered {
			case "", "p", "posix":
				vars.WorkRecord.Fields["objectClass"] = "posixGroup"
				valueEntered = "posix"
			case "g", "groupOfNames":
				vars.WorkRecord.Fields["objectClass"] = "groupOfNames"
				// hard coded, groupOfNames must have at least 1 member
				vars.WorkRecord.Fields["member"] = fmt.Sprintf("uid=initial-user,%s", conn.Config.ServerValues.GroupDN)
				valueEntered = "groupOfNames"
			}
			if !funcs.I.IsInList(validTypes, valueEntered) {
				funcs.P.PrintRed(fmt.Sprintf("\tWrong group type (%s) aborting...\n\n", valueEntered))
				return false
			}
		}

		if (len(valueEntered) == 0) && (vars.Template[fieldName].NoEmpty == true) {
			funcs.P.PrintRed("\tNo value was entered aborting..%s.\n\n")
			return false
		}
		fmt.Printf("\n")
	}

	if vars.WorkRecord.Fields["objectClass"] == "posixGroup" {
		vars.WorkRecord.Fields["gidNumber"] = strconv.Itoa(conn.GetNextGID())
		funcs.P.PrintPurple(fmt.Sprintf("\tOptional set groups's GID, press enter to use the next GID: %s\n",
			vars.WorkRecord.Fields["gidNumber"]))
		fmt.Printf("\t%s: ", vars.Template["gidNumber"].Prompt)
		reader := bufio.NewReader(os.Stdin)
		valueEntered, _ := reader.ReadString('\n')
		valueEntered = strings.TrimSuffix(valueEntered, "\n")
		if len(valueEntered) > 0 {
			gitNumberList := conn.GetAllGroupsGID()
			if groupname, found := gitNumberList[valueEntered]; found {
				funcs.P.PrintRed(fmt.Sprintf("\n\tGiven group id %s already use by the group %s , aborting...\n",
					valueEntered, groupname))
				return false
			}
			vars.WorkRecord.Fields["gidNumber"] = valueEntered
		}
	}
	fmt.Printf("\n")
	return conn.CreateGroup()
}

func Create(conn *ldap.Connection, funcs *vars.Funcs) {
	fmt.Printf("\t%s\n", funcs.P.PrintHeader(vars.Blue, vars.Purple, "Create Group", 18, true))
	if createGroup(conn, funcs) {
		funcs.P.PrintGreen(fmt.Sprintf("\tGroup %s created\n", vars.WorkRecord.Fields["cn"]))
	} else {
		funcs.P.PrintRed(fmt.Sprintf("\tFailed to create the group %s, check the log file\n",
			vars.WorkRecord.Fields["cn"]))
	}
	fmt.Printf("\t%s\n", funcs.P.PrintLine(vars.Purple, 50))
}
