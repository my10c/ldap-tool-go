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
	"regexp"
	"strconv"
	"strings"

	"my10c.ldap/cmds/common"
	"my10c.ldap/ldap"
	"my10c.ldap/vars"
	ldapv3 "gopkg.in/ldap.v2"
)

var (
	fields = []string{"uidNumber", "givenName", "sn", "departmentNumber",
		"mail", "loginShell", "userPassword",
		"shadowMax", "shadowExpire", "sshPublicKey"}

	// input
	valueEntered string

	// user's groups
	userGroupList      []string
	availableGroupList []string
	// need to strip the full dn
	displayUserGroupList      []string
	displayAvailableGroupList []string

	// keep track if password was changed
	shadowMaxChanged bool = false
)

func leaveGroup(conn *ldap.Connection, funcs *vars.Funcs) {
	groupList := conn.GetGroupType()
	for _, groupName := range conn.Record.GroupDelList {
		if funcs.I.IsInList(groupList["posixGroup"], groupName) {
			vars.WorkRecord.MemberType = "posixGroup"
			vars.WorkRecord.MemberType = "memberUid"
			vars.WorkRecord.ID = vars.WorkRecord.ID
		}
		if funcs.I.IsInList(groupList["groupOfNames"], groupName) {
			vars.WorkRecord.MemberType = "groupOfNames"
			vars.WorkRecord.MemberType = "member"
			vars.WorkRecord.ID = fmt.Sprintf("uid=%s,%s", vars.WorkRecord.ID, conn.Config.ServerValues.UserDN)
		}
		vars.WorkRecord.DN = groupName
		conn.RemoveFromGroups()
	}
}

func joinGroup(conn *ldap.Connection, funcs *vars.Funcs) {
	groupList := conn.GetGroupType()
	for _, groupName := range conn.Record.GroupAddList {
		if funcs.I.IsInList(groupList["posixGroup"], groupName) {
			vars.WorkRecord.MemberType = "posixGroup"
			vars.WorkRecord.MemberType = "memberUid"
			vars.WorkRecord.ID = vars.WorkRecord.ID
		}
		if funcs.I.IsInList(groupList["groupOfNames"], groupName) {
			vars.WorkRecord.MemberType = "groupOfNames"
			vars.WorkRecord.MemberType = "member"
			vars.WorkRecord.ID = fmt.Sprintf("uid=%s,%s", vars.WorkRecord.ID, conn.Config.ServerValues.UserDN)
		}
		vars.WorkRecord.DN = groupName
		conn.AddToGroup()
	}
}

func createModifyUserRecord(conn *ldap.Connection, records *ldapv3.SearchResult, funcs *vars.Funcs) int {
	reader := bufio.NewReader(os.Stdin)

	vars.WorkRecord.DN = fmt.Sprintf("uid=%s,%s", vars.WorkRecord.ID, conn.Config.ServerValues.UserDN)

	funcs.P.PrintPurple(fmt.Sprintf("\tUsing user: %s\n", vars.WorkRecord.ID))
	funcs.P.PrintYellow(fmt.Sprintf("\tPress enter to leave the value unchanged\n"))
	fmt.Printf("\t%s\n", funcs.P.PrintLine(vars.Purple, 50))

	for _, fieldName := range fields {
		// these will be valid once the field was filled since they depends
		// on some of the fields value
		switch fieldName {
		case "uidNumber":
			fmt.Printf("\t%s\n", vars.DangerZone)
			funcs.P.PrintCyan(fmt.Sprintf("\tCurrent value: %s%s%s\n",
				vars.Green, records.Entries[0].GetAttributeValue(fieldName), vars.Off))

		case "givenName":
			funcs.P.PrintCyan(fmt.Sprintf("\tCurrent value: %s%s%s\n",
				vars.Green, records.Entries[0].GetAttributeValue(fieldName), vars.Off))

		case "sn":
			funcs.P.PrintCyan(fmt.Sprintf("\tCurrent value: %s%s%s\n",
				vars.Green, records.Entries[0].GetAttributeValue(fieldName), vars.Off))

		case "mail":
			funcs.P.PrintCyan(fmt.Sprintf("\tCurrent value: %s%s%s\n",
				vars.Green, records.Entries[0].GetAttributeValue(fieldName), vars.Off))

		case "departmentNumber":
			funcs.P.PrintCyan(fmt.Sprintf("\tCurrent value: %s%s%s\n",
				vars.Green, records.Entries[0].GetAttributeValue(fieldName), vars.Off))
			funcs.P.PrintYellow(fmt.Sprintf("\t\tValid departments: %s\n",
				strings.Join(conn.Config.GroupValues.Groups[:], ", ")))

		case "loginShell":
			funcs.P.PrintCyan(fmt.Sprintf("\tCurrent value: %s%s%s\n",
				vars.Green, records.Entries[0].GetAttributeValue(fieldName), vars.Off))
			funcs.P.PrintYellow(fmt.Sprintf("\t\tValid shells: %s\n",
				strings.Join(conn.Config.DefaultValues.ValidShells[:], ", ")))

		case "userPassword":
			passWord := funcs.R.Generate()
			funcs.P.PrintCyan(fmt.Sprintf("\tCurrent value (encrypted!): %s%s%s\n",
				vars.Green, records.Entries[0].GetAttributeValue(fieldName), vars.Off))
			funcs.P.PrintYellow(fmt.Sprintf("\t\tsuggested password: %s\n", passWord))

		case "shadowMax":
			funcs.P.PrintCyan(fmt.Sprintf("\tCurrent max password age: %s%s%s\n",
				vars.Green, records.Entries[0].GetAttributeValue(fieldName), vars.Off))
			funcs.P.PrintYellow(
				fmt.Sprintf("\t\tMin %d days and max %d days\n",
					conn.Config.DefaultValues.ShadowMin,
					conn.Config.DefaultValues.ShadowMax))

		case "shadowExpire":
			value, _ := strconv.ParseInt(records.Entries[0].GetAttributeValue(fieldName), 10, 64)
			passExpired := funcs.E.ReadableEpoch(value * 86400)
			funcs.P.PrintCyan(fmt.Sprintf("\tCurrent password will expire on: %s%s%s\n",
				vars.Green, passExpired, vars.Off))

		case "sshPublicKey":
			funcs.P.PrintCyan(fmt.Sprintf("\tCurrent value: %s%s%s\n",
				vars.Green, records.Entries[0].GetAttributeValue(fieldName), vars.Off))
		}

		funcs.P.PrintPurple(fmt.Sprintf("\t%s: ", vars.Template[fieldName].Prompt))

		valueEntered, _ = reader.ReadString('\n')
		valueEntered = strings.TrimSuffix(valueEntered, "\n")
		switch fieldName {
		case "givenName", "sn":
			valueEntered = strings.Title(valueEntered)

		case "mail":
			valueEntered = strings.ToLower(valueEntered)
			valueEntered = strings.ToLower(valueEntered)

		case "loginShell":
			if len(valueEntered) > 0 {
				if !funcs.I.IsInList(conn.Config.DefaultValues.ValidShells, valueEntered) {
					funcs.P.PrintRed("\t\tInvalid shell was given, it will be ignored\n")
					valueEntered = ""
				} else {
					valueEntered = "/bin/" + valueEntered
				}
			}

		case "departmentNumber":
			if len(valueEntered) != 0 {
				if !funcs.I.IsInList(conn.Config.GroupValues.Groups, valueEntered) {
					funcs.P.PrintRed("\t\tInvalid departments was given, it will be ignored\n")
					valueEntered = ""
				} else {
					for _, mapValues := range conn.Config.GroupValues.GroupsMap {
						if mapValues.Name == valueEntered {
							vars.WorkRecord.Fields["gidNumber"] = strconv.Itoa(mapValues.Gid)
							break
						}
					}
					valueEntered = strings.ToUpper(valueEntered)
				}
			}

		case "shadowMax":
			if len(valueEntered) != 0 {
				shadowMax, _ := strconv.Atoi(valueEntered)
				if shadowMax < conn.Config.DefaultValues.ShadowMin ||
					shadowMax > conn.Config.DefaultValues.ShadowMax {
					funcs.P.PrintRed(fmt.Sprintf("\t\tGiven value %d, is out or range, is set to %d\n",
						shadowMax, conn.Config.DefaultValues.ShadowAge))
					valueEntered = strconv.Itoa(conn.Config.DefaultValues.ShadowAge)
				}
				shadowMaxChanged = true
			}

		case "shadowExpire":
			if len(valueEntered) == 0 {
				funcs.P.PrintCyan(fmt.Sprintf("\tPassword expiration date will not be changed\n"))
			} else {
				// calculate when it will be expired based on default value if shadowMax
				// otherwise it will be today + new shadowMax value
				currShadowMax := records.Entries[0].GetAttributeValue("shadowMax")
				if shadowMaxChanged == true {
					currShadowMax = vars.WorkRecord.Fields["shadowMax"]
					funcs.P.PrintYellow(fmt.Sprintf("\t\tCalculate from new given max password age\n"))
				}

				// set last changed to now
				vars.WorkRecord.Fields["shadowLastChange"] = vars.Template["shadowLastChange"].Value
				// calculate the new shadowExpire
				shadowLastChange, _ := strconv.ParseInt(vars.WorkRecord.Fields["shadowLastChange"], 10, 64)
				shadowMax, _ := strconv.ParseInt(currShadowMax, 10, 64)
				passExpired := funcs.E.ReadableEpoch((shadowLastChange + shadowMax) * 86400)
				funcs.P.PrintCyan(fmt.Sprintf("\tCurrent password will now expire on: %s\n", passExpired))
				// replace the 'Y' with the correct value
				valueEntered = strconv.FormatInt((shadowLastChange + shadowMax), 10)
			}
		}

		if len(valueEntered) != 0 {
			vars.WorkRecord.Fields[fieldName] = valueEntered
		}
	}

	fmt.Printf("\n")
	userID := vars.WorkRecord.ID
	userDN := records.Entries[0].DN
	conn.GetUserGroups(userID, userDN)
	conn.GetAvailableGroups(userID, userDN)

	fmt.Printf("\n\t%s\n", funcs.P.PrintLine(vars.Purple, 50))
	reg, _ := regexp.Compile("^cn=|,ou=groups,.*")

	for _, userGroup := range conn.Record.UserGroups {
		userGroupList = append(userGroupList, userGroup)
		displayUserGroupList = append(displayUserGroupList, reg.ReplaceAllString(userGroup, " "))
	}
	for _, availableGroup := range conn.Record.AvailableGroups {
		availableGroupList = append(availableGroupList, availableGroup)
		displayAvailableGroupList = append(displayAvailableGroupList, reg.ReplaceAllString(availableGroup, " "))
	}

	funcs.P.PrintPurple(fmt.Sprintf("\tUser %s groups: %s\n", vars.WorkRecord.ID,
		strings.Join(displayUserGroupList[:], "")))

	funcs.P.PrintGreen(fmt.Sprintf("\tAvailable groups: %s\n",
		strings.Join(displayAvailableGroupList[:], " ")))

	fmt.Printf("\n\t%s\n", funcs.P.PrintLine(vars.Purple, 50))
	for _, leaveGroup := range userGroupList {
		fmt.Printf("\t%sRemove%s the group %s%s%s? default to not remove group, [Y/N]: ",
			vars.Red, vars.Off, vars.Red, reg.ReplaceAllString(leaveGroup, ""), vars.Off)
		valueEntered, _ = reader.ReadString('\n')
		valueEntered = strings.ToLower(strings.TrimSuffix(valueEntered, "\n"))
		switch valueEntered {
		case "y", "yes", "d", "del":
			conn.Record.GroupDelList = append(conn.Record.GroupDelList, leaveGroup)
		}
	}

	fmt.Printf("\n")
	for _, joinGroup := range availableGroupList {
		fmt.Printf("\t%sJoin%s the group %s%s%s? default not to join group, [Y/N]: ",
			vars.Green, vars.Off, vars.Green, reg.ReplaceAllString(joinGroup, ""), vars.Off)
		valueEntered, _ = reader.ReadString('\n')
		valueEntered = strings.ToLower(strings.TrimSuffix(valueEntered, "\n"))
		switch valueEntered {
		case "y", "yes":
			conn.Record.GroupAddList = append(conn.Record.GroupAddList, joinGroup)
		}
	}

	if len(conn.Record.GroupDelList) == 0 && len(conn.Record.GroupAddList) == 0 && len(vars.WorkRecord.Fields) == 0 {
		fmt.Printf("\n\t%s\n", funcs.P.PrintLine(vars.Purple, 50))
		funcs.P.PrintBlue(fmt.Sprintf("\n\tNo field were changed, no modification was made for the user %s\n",
			vars.WorkRecord.ID))
		return 0
	}
	return 1
}

func Modify(conn *ldap.Connection, funcs *vars.Funcs) {
	fmt.Printf("\t%s\n", funcs.P.PrintHeader(vars.Blue, vars.Purple, "Modify user", 18, true))
	vars.SearchResultData.WildCardSearchBase = vars.UserWildCardSearchBase
	vars.SearchResultData.RecordSearchbase = vars.UserWildCardSearchBase
	vars.SearchResultData.DisplayFieldID = vars.UserDisplayFieldID
	if common.GetObjectRecord(conn, true, "user", funcs) {
		if createModifyUserRecord(conn, vars.SearchResultData.SearchResult, funcs) > 0 {
			if len(vars.WorkRecord.Fields) > 0 {
				if !conn.ModifyUser() {
					funcs.P.PrintRed(fmt.Sprintf("\n\tFailed modify the user %s, check the log file\n",
						vars.WorkRecord.ID))
				} else {
					funcs.P.PrintGreen(fmt.Sprintf("\n\tUser %s modified successfully\n", vars.WorkRecord.ID))
				}
			}
			if len(conn.Record.GroupDelList) > 0 {
				leaveGroup(conn, funcs)
			}
			if len(conn.Record.GroupAddList) > 0 {
				joinGroup(conn, funcs)
			}
		}
	}
	fmt.Printf("\t%s\n", funcs.P.PrintLine(vars.Purple, 50))
}
