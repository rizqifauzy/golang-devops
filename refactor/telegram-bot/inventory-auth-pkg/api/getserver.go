package api

import (
	"database/sql"
	"fmt"
)

type ServerInfo struct {
	VM           string
	IPOAM        string
	IPService    string
	Powerstate   string
	Datacenter   string
	OS           string
	AppsName     string
	AppsPriority string
	AppsCustody  string
	ProdNonProd  string
	Environment  string
	Site         string
	ManagedBy    string
	SupportLevel string
	Notes        string
}

func getServerInfo(db *sql.DB, query string) (*ServerInfo, error) {
	var info ServerInfo
	row := db.QueryRow(`SELECT vm_name, ip_oam, ip_service, powerstate, datacenter, os_configuration, apps_name,   
		apps_priority, apps_custody, prod_non_prod, environment, site, managed_by,   
		support_level, notes FROM servers WHERE vm_name = $1 OR ip_oam = $1`, query)

	err := row.Scan(&info.VM, &info.IPOAM, &info.IPService, &info.Powerstate, &info.Datacenter, &info.OS,
		&info.AppsName, &info.AppsPriority, &info.AppsCustody, &info.ProdNonProd, &info.Environment,
		&info.Site, &info.ManagedBy, &info.SupportLevel, &info.Notes)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &info, nil
}

func formatServerInfo(info *ServerInfo) string {
	return fmt.Sprintf(`Informasi Server:  
VM: %s  
IPOAM: %s  
IPService: %s  
Powerstate: %s  
Datacenter: %s  
OS according to the configuration file: %s  
Apps Name: %s  
Apps Priority: %s  
Apps Custody (Email Address): %s  
Prod/Non Prod: %s  
Environment: %s  
Site: %s  
Managed By: %s  
Support Level: %s  
Notes: %s`,
		info.VM, info.IPOAM, info.IPService, info.Powerstate, info.Datacenter, info.OS,
		info.AppsName, info.AppsPriority, info.AppsCustody, info.ProdNonProd, info.Environment,
		info.Site, info.ManagedBy, info.SupportLevel, info.Notes)
}
