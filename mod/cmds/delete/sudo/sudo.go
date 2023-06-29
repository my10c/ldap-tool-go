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

func Delete(conn *ldap.Connection, funcs *vars.Funcs) {
	fmt.Printf("\t%s\n", funcs.P.PrintHeader(vars.Blue, vars.Purple, "Delete Sudo rules", 18, true))
	vars.SearchResultData.WildCardSearchBase = vars.SudoWildCardSearchBase
	vars.SearchResultData.RecordSearchbase = vars.SudoWildCardSearchBase
	vars.SearchResultData.DisplayFieldID = vars.SudoDisplayFieldID
	if common.GetObjectRecord(conn, true, "sudo rules", funcs) {
		common.DeleteObjectRecord(conn, vars.SearchResultData.SearchResult, "sudo rules", funcs)
	}
	fmt.Printf("\t%s\n", funcs.P.PrintLine(vars.Purple, 50))
}
