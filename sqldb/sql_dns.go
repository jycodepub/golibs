package sqldb

import "fmt"

type SqlDNS struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	SSLMode  string
}

func (dns SqlDNS) getPostgresDNS() string {
	sslmode := dns.SSLMode
	if sslmode == "" {
		sslmode = "disable"
	}
	return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s", dns.User, dns.Password, dns.Host, dns.Port, dns.Database, sslmode)
}

