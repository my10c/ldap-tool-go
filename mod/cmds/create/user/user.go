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
	"regexp"
	"strconv"
	"strings"

	"my10c.ldap/ldap"
	"my10c.ldap/vars"
)

var (
	// not required for create a new user : cn, gidNumber, displayName, gecos
	// homeDirectory, shadowLastChange, shadowLastChange
	// groups is handled seperat

	fields = []string{"uid", "givenName", "sn",
		"uidNumber", "departmentNumber",
		"mail", "loginShell", "userPassword",
		"shadowWarning", "shadowMax",
		"sshPublicKey"}

	// construct base on FirstName + LastName
	userFullname = []string{"cn", "displayName", "gecos"}

	// given field value
	email       string
	passWord    string
	shells      string
	departments string
	nextUID     int
	shadowMax   int
)

func joinGroup(conn *ldap.Connection, funcs *vars.Funcs) int {
	var errored int = 0
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
		if !conn.AddToGroup() {
			errored++
		}
	}
	return errored
}

func createUserRecord(conn *ldap.Connection, funcs *vars.Funcs) bool {
	usersName := conn.GetAllUsers()
	usersNameUid := conn.GetUsersUID()
	reader := bufio.NewReader(os.Stdin)

	for _, fieldName := range fields {
		// these will be valid once the field was filled since they depends
		// on some of the fields value
		switch fieldName {
		case "uid":
			funcs.P.PrintPurple(
				fmt.Sprintf("\tThe userid / login name is case sensitive, it will be made all lowercase\n"))

		case "uidNumber":
			nextUID = conn.GetNextUID()
			funcs.P.PrintPurple(fmt.Sprintf("\t\tOptional set user's UID, press enter to use the next UID: %d\n",
				nextUID))

		case "departmentNumber":
			funcs.P.PrintYellow(fmt.Sprintf("\t\tValid departments: %s\n",
				strings.Join(conn.Config.GroupValues.Groups[:], ", ")))

		case "mail":
			email = fmt.Sprintf("%s.%s@%s",
				strings.ToLower(vars.WorkRecord.Fields["givenName"]),
				strings.ToLower(vars.WorkRecord.Fields["sn"]),
				conn.Config.ServerValues.EmailDomain)
			funcs.P.PrintCyan(fmt.Sprintf("\tDefault email: %s\n", email))

		case "loginShell":
			funcs.P.PrintYellow(fmt.Sprintf("\t\tValid shells: %s\n",
				strings.Join(conn.Config.DefaultValues.ValidShells[:], ", ")))

		case "userPassword":
			passWord = funcs.R.Generate()
			funcs.P.PrintPurple("\t\tPress Enter to accept the suggested password\n")
			funcs.P.PrintCyan(fmt.Sprintf("\tSuggested password: %s\n", passWord))

		case "shadowMax":
			funcs.P.PrintYellow(
				fmt.Sprintf("\t\tMin %d days and max %d days\n",
					conn.Config.DefaultValues.ShadowMin,
					conn.Config.DefaultValues.ShadowMax))
		}

		if vars.Template[fieldName].Value != "" {
			funcs.P.PrintPurple(fmt.Sprintf("\t ** Default to: %s **\n", vars.Template[fieldName].Value))
		}

		fmt.Printf("\t%s: ", vars.Template[fieldName].Prompt)

		reader = bufio.NewReader(os.Stdin)
		valueEntered, _ := reader.ReadString('\n')
		valueEntered = strings.TrimSuffix(valueEntered, "\n")

		switch fieldName {
		case "uid":
			if funcs.I.IsInList(usersName, valueEntered) {
				funcs.P.PrintRed(fmt.Sprintf("\n\tGiven user %s already exist, aborting...\n\n", valueEntered))
				return false
			}
			fmt.Printf("\t%s\n", funcs.P.PrintLine(vars.Purple, 50))
			vars.WorkRecord.Fields["uid"] = valueEntered
			vars.WorkRecord.ID = valueEntered
			funcs.P.PrintPurple(fmt.Sprintf("\tUsing user: %s\n", valueEntered))

		case "givenName", "sn":
			valueEntered = strings.Title(valueEntered)

		case "uidNumber":
			if len(valueEntered) > 0 {
				if userName, found := usersNameUid[valueEntered]; found {
					funcs.P.PrintRed(fmt.Sprintf("\n\tGiven uid id %s already use by the user %s , aborting...\n",
						valueEntered, userName))
					return false
				}
				valueEntered = valueEntered
			} else {
				valueEntered = strconv.Itoa(nextUID)
			}

		case "departmentNumber":
			if len(valueEntered) > 0 {
				if !funcs.I.IsInList(conn.Config.GroupValues.Groups, valueEntered) {
					funcs.P.PrintRed(fmt.Sprintf("\n\tGiven departmentNumber %s is not valid, aborting...\n\n",
						valueEntered))
					return false
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
			if len(valueEntered) == 0 {
				valueEntered = strings.ToUpper(conn.Config.DefaultValues.GroupName)
				vars.WorkRecord.Fields["gidNumber"] = strconv.Itoa(conn.Config.DefaultValues.GroupId)
			}

		case "mail":
			if len(valueEntered) == 0 {
				valueEntered = email
			}

		case "loginShell":
			if len(valueEntered) > 0 {
				if !funcs.I.IsInList(conn.Config.DefaultValues.ValidShells, valueEntered) {
					funcs.P.PrintRed(fmt.Sprintf("\n\tGiven shell %s is not valid, aborting...\n\n", valueEntered))
					return false
				}
				valueEntered = "/bin/" + valueEntered
			}
			if len(valueEntered) == 0 {
				valueEntered = "/bin/" + conn.Config.DefaultValues.Shell
			}

		case "userPassword":
			if len(valueEntered) > 0 {
				vars.WorkRecord.Fields["userPassword"] = valueEntered
			}
			if len(valueEntered) == 0 {
				valueEntered = passWord
			}

		case "shadowWarning":
			if len(valueEntered) == 0 {
				valueEntered = strconv.Itoa(conn.Config.DefaultValues.ShadowWarning)
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
			}
			if len(valueEntered) == 0 {
				valueEntered = strconv.Itoa(conn.Config.DefaultValues.ShadowAge)
			}

		default:
			if len(valueEntered) == 0 {
				valueEntered = vars.Template[fieldName].Value
			}

		}

		if len(valueEntered) == 0 && vars.Template[fieldName].NoEmpty == true {
			funcs.P.PrintRed("\tNo value was entered aborting...\n\n")
			return false
		}
		// set the default values
		if len(valueEntered) == 0 {
			vars.WorkRecord.Fields[fieldName] = vars.Template[fieldName].Value
		}
		// update the user record so it can be submitted
		vars.WorkRecord.Fields[fieldName] = valueEntered
	}

	for idx, _ := range userFullname {
		vars.WorkRecord.Fields[userFullname[idx]] =
			vars.WorkRecord.Fields["givenName"] + " " + vars.WorkRecord.Fields["sn"]
	}

	// dn is create base on given uid and user DN
	vars.WorkRecord.DN = fmt.Sprintf("uid=%s,%s", vars.WorkRecord.Fields["uid"], conn.Config.ServerValues.UserDN)

	// this is always /home + userlogin
	vars.WorkRecord.Fields["homeDirectory"] = "/home/" + vars.WorkRecord.Fields["uid"]

	// initialized to be today's epoch days
	vars.WorkRecord.Fields["shadowExpire"] = vars.Template["shadowExpire"].Value
	vars.WorkRecord.Fields["shadowLastChange"] = vars.Template["shadowLastChange"].Value

	// setup the groups for the user
	fmt.Printf("\n\t%s\n", funcs.P.PrintLine(vars.Purple, 50))
	fmt.Printf("\n")
	reg, _ := regexp.Compile("^cn=|,ou=groups,.*")
	availableGroups := conn.GetAllGroups()
	for _, joinGroup := range availableGroups {
		groupToJoin := reg.ReplaceAllString(joinGroup, "")
		if vars.WorkRecord.Fields["departmentNumber"] != strings.ToUpper(groupToJoin) {
			fmt.Printf("\t%sJoin%s the group %s%s%s? default not to join group, [Y/N]: ",
				vars.Green, vars.Off, vars.Green, groupToJoin, vars.Off)
			valueEntered, _ := reader.ReadString('\n')
			valueEntered = strings.ToLower(strings.TrimSuffix(valueEntered, "\n"))
			switch valueEntered {
			case "y", "yes":
				conn.Record.GroupAddList = append(conn.Record.GroupAddList, joinGroup)
			}
		} else {
			// need to append the group dn
			conn.Record.GroupAddList = append(conn.Record.GroupAddList,
				fmt.Sprintf("cn=%s,%s", groupToJoin, conn.Config.ServerValues.GroupDN))
		}
	}
	return conn.CreateUser()
}

func Create(conn *ldap.Connection, funcs *vars.Funcs) {
	fmt.Printf("\t%s\n", funcs.P.PrintHeader(vars.Blue, vars.Purple, "Create User", 18, true))
	if !createUserRecord(conn, funcs) {
		funcs.P.PrintRed(fmt.Sprintf("\n\tFailed adding the user %s, check the log file\n", vars.WorkRecord.Fields["uid"]))
	} else {
		funcs.P.PrintGreen(fmt.Sprintf("\n\tUser %s added successfully\n", vars.WorkRecord.Fields["uid"]))
		if len(conn.Record.GroupAddList) > 0 {
			if joinGroup(conn, funcs) != 0 {
				funcs.P.PrintRed(fmt.Sprintf("\n\tFailed adding the user %s to groups, check the log file\n",
					vars.WorkRecord.Fields["uid"]))
			}
		}
	}
	fmt.Printf("\t%s\n", funcs.P.PrintLine(vars.Purple, 50))
}
