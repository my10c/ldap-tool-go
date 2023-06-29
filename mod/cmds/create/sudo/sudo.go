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

func createSudoRecord(conn *ldap.Connection, funcs *vars.Funcs) bool {
	sudoRules := conn.GetAllSudoRules()

	for _, fieldName := range vars.SudoFields {
		if vars.Template[fieldName].Value != "" {
			fmt.Printf("\t%sDefault to:%s %s%s%s\n",
				vars.Purple, vars.Off, vars.Cyan, vars.Template[fieldName].Value, vars.Off)
		}

		fmt.Printf("\t%s: ", vars.Template[fieldName].Prompt)

		reader := bufio.NewReader(os.Stdin)
		valueEntered, _ := reader.ReadString('\n')
		valueEntered = strings.TrimSuffix(valueEntered, "\n")

		// make sure any combination of `all` is made uppercase
		if strings.ToLower(valueEntered) == "all" {
			valueEntered = "ALL"
		}

		switch fieldName {
		case "cn":
			if funcs.I.IsInList(sudoRules, valueEntered) {
				funcs.P.PrintRed(fmt.Sprintf("\n\tGiven cn %s already exist, aborting...\n\n", valueEntered))
				return false
			}
			fmt.Printf("\t%s\n", funcs.P.PrintLine(vars.Purple, 50))
			funcs.P.PrintPurple(fmt.Sprintf("\tUsing Sudo Rule: %s\n\n", valueEntered))
			vars.WorkRecord.Fields["cn"] = valueEntered
			vars.WorkRecord.Fields["dn"] = fmt.Sprintf("cn=%s,%s", valueEntered, conn.Config.SudoValues.SudoersBase)
			vars.WorkRecord.Fields["objectClass"] = "sudoRole"

		case "sudoCommand", "sudoHost", "sudoRunAsUser":
			if len(valueEntered) > 0 {
				vars.WorkRecord.Fields[fieldName] = valueEntered
			} else {
				vars.WorkRecord.Fields[fieldName] = vars.Template[fieldName].Value
			}
		case "sudoOption":
			if len(valueEntered) > 0 {
				vars.WorkRecord.Fields[fieldName] = valueEntered
			}
		case "sudoOrder":
			vars.WorkRecord.Fields[fieldName] = vars.Template[fieldName].Value
			if len(valueEntered) > 0 {
				value, _ := strconv.Atoi(valueEntered)
				if value < 3 || value > 10 {
					funcs.P.PrintRed(fmt.Sprintf("%s\tGiven order %s is not allowed, set to the default %s\n",
						vars.OneLineUP, valueEntered, vars.Template[fieldName].Value))
					valueEntered = vars.Template[fieldName].Value
				} else {
					vars.WorkRecord.Fields[fieldName] = valueEntered
				}
			}
		}
		if len(valueEntered) == 0 && vars.Template[fieldName].NoEmpty == true {
			funcs.P.PrintRed("\tNo value was entered aborting...\n\n")
			return false
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\n")
	return conn.CreateSudoRule()
}

func Create(conn *ldap.Connection, funcs *vars.Funcs) {
	fmt.Printf("\t%s\n", funcs.P.PrintHeader(vars.Blue, vars.Purple, "Create Group", 18, true))
	if createSudoRecord(conn, funcs) {
		funcs.P.PrintGreen(fmt.Sprintf("\tSudo rule %s created\n", vars.WorkRecord.Fields["cn"]))
	} else {
		funcs.P.PrintRed(fmt.Sprintf("\tFailed to create the sudo rule %s, check the log file\n",
			vars.WorkRecord.Fields["cn"]))
	}
	fmt.Printf("\t%s\n", funcs.P.PrintLine(vars.Purple, 50))
}
