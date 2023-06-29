//
// BSD 3-Clause License
//
// Copyright (c) 2022, © Badassops LLC / Luc Suryo
// All rights reserved.
//

package ldap

import (
	"github.com/my10c/packages-go/exit"
	"github.com/my10c/packages-go/lock"
	ldapv3 "gopkg.in/ldap.v2"
)

// search the ldap database
func (conn *Connection) Search() (*ldapv3.SearchResult, int) {
	searchRecords := ldapv3.NewSearchRequest(
		conn.Config.ServerValues.BaseDN,
		ldapv3.ScopeWholeSubtree,
		ldapv3.NeverDerefAliases, 0, 0, false,
		conn.SearchInfo.SearchBase,
		conn.SearchInfo.SearchAttribute,
		nil,
	)
	searchResult, err := conn.Conn.Search(searchRecords)
	if err != nil {
		conn.Conn.Close()
		// create a new lock, so we can remove the lock file
		// without the need to pass it as an argument
		// so with exit
		lockPtr := lock.New(conn.Config.DefaultValues.LockFile)
		exitPtr := exit.New("ldap search", 1)
		lockPtr.LockRelease()
		exitPtr.ExitError(err)
	}
	return searchResult, len(searchResult.Entries)
}
