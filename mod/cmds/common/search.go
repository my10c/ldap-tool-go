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

var (
	noMatch bool = false
	enterData string
)

func GetObjectRecord(c *ldap.Connection, firstTime bool, displayName string, funcs *vars.Funcs) bool {
	var records *ldapv3.SearchResult
	var recordCount int
	var displayFieldID = vars.SearchResultData.DisplayFieldID
	var wildCardSearchBase = vars.SearchResultData.WildCardSearchBase
	var recordSearchbase = vars.SearchResultData.RecordSearchbase

	if (noMatch == true ) {
		funcs.P.PrintYellow(fmt.Sprintf("\n\tNo record matching wildcard %s was found, aborting...\n", enterData))
		return false
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("\tEnter %s name to be use: ", displayName)
	enterData, _ = reader.ReadString('\n')
	enterData = strings.TrimSuffix(enterData, "\n")

	if enterData == "" {
		funcs.P.PrintRed(fmt.Sprintf("\n\tNo %s name was given aborting...\n", displayName))
		return false
	}

	if firstTime {
		fmt.Printf("\tUse wildcard (default to N)? [y/n]: ")
		wildCard, _ := reader.ReadString('\n')
		wildCard = strings.TrimSuffix(wildCard, "\n")
		if readinput.ReadYN(wildCard, false) == true {
			enterData = "*" + enterData + "*"
			wildCardSearchBase = strings.ReplaceAll(wildCardSearchBase, "VALUE", enterData)
			c.SearchInfo.SearchBase = wildCardSearchBase
			c.SearchInfo.SearchAttribute = []string{displayFieldID}
			records, recordCount = c.Search()
			if recordCount == 0 {
				noMatch = true
				return GetObjectRecord(c, false, displayName, funcs)
			}
			for idx, _ := range records.Entries {
				funcs.P.PrintBlue(fmt.Sprintf("\t%s: %s\n",
					displayFieldID,
					records.Entries[idx].GetAttributeValue(displayFieldID)))
			}
			fmt.Printf("\n\tSelect the %s from the above list:\n", displayName)
			return GetObjectRecord(c, false, displayName, funcs)
		}
	}

	recordSearchbase = strings.ReplaceAll(recordSearchbase, "VALUE", enterData)
	c.SearchInfo.SearchBase = recordSearchbase
	c.SearchInfo.SearchAttribute = []string{}
	records, recordCount = c.Search()

	if recordCount == 0 {
		funcs.P.PrintRed(fmt.Sprintf("\n\t%s %s was not found, aborting...\n", strings.Title(displayName), enterData))
		return false
	}
	vars.SearchResultData.RecordCount = recordCount
	vars.SearchResultData.SearchResult = records
	vars.WorkRecord.ID = enterData
	return true
}
