package graph

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.
import svc "github.com/tsirysndr/mada/interfaces"

type Resolver struct {
	CommuneService   svc.CommuneSvc
	DistrictService  svc.DistrictSvc
	FokontanyService svc.FokontanySvc
	RegionService    svc.RegionSvc
	SearchService    svc.SearchSvc
}
