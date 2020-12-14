package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/zerodotfive/ldaptokenauth/src/authserver"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	pflag.String("listen", "127.0.0.1:8090", "Listen address")
	pflag.String("secret", "", "Signing secret key. Minimum is 24 symbols")
	pflag.String("ldap-server", "", "LDAP server address, exmaple: ldap.domain.com:389")
	pflag.String("bind-dn", "", "LDAP Bind DB, example: uid=admin,cn=users,cn=accounts,dc=domain,dc=com")
	pflag.String("bind-pw", "", "LDAP Bind password")
	pflag.String("base-dn", "", "LDAP base DN cn=users,cn=accounts,dc=domain,dc=com")
	pflag.String("group-filter", "", "LDAP filter (memberOf=cn=mygroup,cn=groups,cn=accounts,dc=domain,dc=com)")
	pflag.Duration("ttl", 24*time.Hour, "Token TTL")

	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	s := authserver.AuthServer{
		Listen:      viper.GetString("listen"),
		Secret:      viper.GetString("secret"),
		LdapServer:  viper.GetString("ldap-server"),
		BindDn:      viper.GetString("bind-dn"),
		BindPw:      viper.GetString("bind-pw"),
		BaseDN:      viper.GetString("base-dn"),
		GroupFilter: viper.GetString("group-filter"),
		TTL:         viper.GetDuration("ttl"),
	}

	varsValid := true
	if len(s.Secret) < 24 {
		varsValid = false
		fmt.Printf("'--secret' should be at least 24 symbols. Also can be set via 'SECRET' environment variable.\n\n")
	}

	if len(s.LdapServer) == 0 {
		varsValid = false
		fmt.Printf("'--ldap-server' should be set. Also can be set via 'LDAP_SERVER' environment variable.\n\n")
	}

	if len(s.BindDn) == 0 {
		varsValid = false
		fmt.Printf("'--bind-dn' should be set. Also can be set via 'BIND_DN' environment variable.\n\n")
	}

	if len(s.BindPw) == 0 {
		varsValid = false
		fmt.Printf("'--bind-pw' should be set. Also can be set via 'BIND_PW' environment variable.\n\n")
	}

	if len(s.BaseDN) == 0 {
		varsValid = false
		fmt.Printf("'--base-dn' should be set. Also can be set via 'BASE_DN' environment variable.\n\n")
	}

	if len(s.GroupFilter) == 0 {
		varsValid = false
		fmt.Printf("'--group-filter' should be set. Also can be set via 'GROUP_FILTER' environment variable.\n\n")
	}

	if !varsValid {
		pflag.Usage()
		os.Exit(1)
	}

	s.Run()
}
