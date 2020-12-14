package ldapclient

import (
	"fmt"

	"github.com/go-ldap/ldap"
)

func Search(ldapServer string, bindDn string, bindPw string, baseDn string, groupFilter string, username string) (*ldap.Entry, error) {
	connLdap, err := ldap.Dial("tcp", ldapServer)
	if err != nil {
		return nil, err
	}
	defer connLdap.Close()

	err = connLdap.Bind(bindDn, bindPw)
	if err != nil {
		return nil, err
	}

	filter := fmt.Sprintf("(&(uid=%s)%s)", ldap.EscapeFilter(username), groupFilter)

	searchReq := ldap.NewSearchRequest(baseDn, ldap.ScopeWholeSubtree, 0, 0, 0, false, filter, nil, []ldap.Control{})

	result, err := connLdap.Search(searchReq)
	if err != nil {
		return nil, err
	}
	if len(result.Entries) == 0 {
		return nil, fmt.Errorf("Invalid user '%s'", username)
	}

	return result.Entries[0], nil
}

func IsValidUser(ldapServer string, userDN string, password string) bool {
	connLdap, err := ldap.Dial("tcp", ldapServer)
	if err != nil {
		return false
	}
	defer connLdap.Close()

	err = connLdap.Bind(userDN, password)
	if err != nil {
		return false
	}

	return true
}
