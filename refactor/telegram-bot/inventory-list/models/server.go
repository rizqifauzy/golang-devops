package models

type Server struct {
	VMName          string `gorm:"vm_name"`
	IPOAM           string `gorm:"ipoam"`
	IPService       string `gorm:"ip_service"`
	Powerstate      string `gorm:"powerstate"`
	Datacenter      string `gorm:"datacenter"`
	OSConfiguration string `gorm:"os_configuration"`
	AppsName        string `gorm:"apps_name"`
	AppsPriority    string `gorm:"apps_priority"`
	AppsCustody     string `gorm:"apps_custody"`
	ProdNonProd     string `gorm:"prod_non_prod"`
	Environment     string `gorm:"environment"`
	Site            string `gorm:"site"`
	ManagedBy       string `gorm:"managed_by"`
	SupportLevel    string `gorm:"support_level"`
	Notes           string `gorm:"notes"`
}
