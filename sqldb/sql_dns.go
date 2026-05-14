package sqldb

import "fmt"

type SqlDNS struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

func (dns SqlDNS) getPostgresDNS() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", dns.User, dns.Password, dns.Host, dns.Port, dns.Database)
}
