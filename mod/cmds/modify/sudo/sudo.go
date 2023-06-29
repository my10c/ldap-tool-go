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
	allowedField     = []string{"sudoCommand", "sudoHost", "sudoOption",
		"sudoOrder", "sudoRunAsUser"}
)

func modifySudo(conn *ldap.Connection, records []*ldapv3.Entry, funcs *vars.Funcs) bool {
	vars.WorkRecord.DN = fmt.Sprintf("cn=%s,%s", vars.WorkRecord.ID, conn.Config.SudoValues.SudoersBase)
	fmt.Printf("\t%s\n", funcs.P.PrintLine(vars.Purple, 50))
	reader := bufio.NewReader(os.Stdin)
	for _, entry := range records {
		funcs.P.PrintBlue(fmt.Sprintf("\tdn: %s\n", entry.DN))
		for _, attributes := range entry.Attributes {
			for _, value := range attributes.Values {
				if (attributes.Name != "objectClass") && (attributes.Name != "cn") {
					funcs.P.PrintCyan(fmt.Sprintf("\tField: %s%s%s\n", vars.Red, attributes.Name, vars.Off))
					switch attributes.Name {
					case "sudoCommand":
						funcs.P.PrintCyan(fmt.Sprintf("\tCurrent value: %s\n", value))

					case "sudoHost":
						funcs.P.PrintCyan(fmt.Sprintf("\tCurrent value: %s\n", value))

					case "sudoOption":
						funcs.P.PrintCyan(fmt.Sprintf("\tCurrent value: %s\n", value))

					case "sudoOrder":
						funcs.P.PrintCyan(fmt.Sprintf("\tCurrent value: %s\n", value))

					case "sudoRunAsUser":
						funcs.P.PrintCyan(fmt.Sprintf("\tCurrent value: %s\n", value))
					}
					fmt.Printf("\tEnter %sdelete%s to delete or press enter to keep: ", vars.RedUnderline, vars.Off)
					valueEntered, _ = reader.ReadString('\n')
					valueEntered = strings.TrimSuffix(valueEntered, "\n")
					if valueEntered == "delete" {
						vars.WorkRecord.SudoDelList[attributes.Name] =
							append(vars.WorkRecord.SudoDelList[attributes.Name], value)
					} else {
						fmt.Printf("\n")
					}
				}
			}
		}
	}

	funcs.P.PrintCyan(fmt.Sprintf("\n\tEach field can have multiple entries\n"))
	funcs.P.PrintCyan(fmt.Sprintf("\tPress enter to skip, or enter value for field\n"))
	for _, fieldname := range allowedField {
		for true {
			funcs.P.PrintGreen(fmt.Sprintf("\tField %s%s%s: enter value: ", vars.Purple, fieldname, vars.Off))
			reader := bufio.NewReader(os.Stdin)
			valueEntered, _ := reader.ReadString('\n')
			valueEntered = strings.TrimSuffix(valueEntered, "\n")
			if len(valueEntered) != 0 {
				vars.WorkRecord.SudoAddList[fieldname] =
					append(vars.WorkRecord.SudoAddList[fieldname], valueEntered)
			} else {
				fmt.Printf("\n")
				break
			}
		}
	}

	modCount = len(vars.WorkRecord.SudoDelList) + len(vars.WorkRecord.SudoAddList)
	if modCount == 0 {
		funcs.P.PrintBlue(fmt.Sprintf("\n\tNo change, no modification made to sudo rule %s\n", vars.WorkRecord.ID))
		return false
	}

	if len(vars.WorkRecord.SudoDelList) > 0 {
		conn.DeleteSudoRule()
	}
	if len(vars.WorkRecord.SudoAddList) > 0 {
		conn.AddSudoRule()
	}
	return true
}

func Modify(conn *ldap.Connection, funcs *vars.Funcs) {
	fmt.Printf("\t%s\n", funcs.P.PrintHeader(vars.Blue, vars.Purple, "Sudo Rules", 18, true))
	vars.SearchResultData.WildCardSearchBase = vars.SudoWildCardSearchBase
	vars.SearchResultData.RecordSearchbase = vars.SudoWildCardSearchBase
	vars.SearchResultData.DisplayFieldID = vars.SudoDisplayFieldID
	if common.GetObjectRecord(conn, true, "sudo rule", funcs) {
		modifySudo(conn, vars.SearchResultData.SearchResult.Entries, funcs)
	}
	fmt.Printf("\t%s\n", funcs.P.PrintLine(vars.Purple, 50))
}
