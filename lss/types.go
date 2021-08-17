package lss;

import(
	"encoding/hex"
	"fmt"
	"crypto/sha1"
	"sync"
)

//
type ProblemInfo struct{
	Problem
	impr_token	string
	solutions		map[string]*Solution
};

//
func NewProblemInfo(pp *Problem) *ProblemInfo {
	var pi ProblemInfo
	pi.Problem= *pp
	pi.solutions= make(map[string]*Solution)
	return &pi
}

// server is used to implement helloworld.GreeterServer.
type LSSServer struct {
	UnimplementedLotsizeSolutionStoreServer
	mutex			sync.Mutex
	problems	map[string]*ProblemInfo
}

//
func NewLSSServer() *LSSServer{
	return &LSSServer{problems: make(map[string]*ProblemInfo)}
}

// type LSPProduct struct{
// 	Id 					string	`json:"id"`
// 	Name 				string	`json:"name"`
// 	InvUnitCost float64	`json:"inv_unit_cost"`
// }
//
// // func GetExProduct
//
// type LSPDemand struct{
// 	Id 					string	`json:"id"`
// 	ProductId		string	`json:"product"`
// 	PeriodIdx		uint		`json:"period"`
// 	Quantity		float64	`json:"quantity"`
// 	UnitRevenue float64	`json:"unit_revenue"`
// }
//
// type LSPMachine struct {
// 	Id 							string							`json:"id"`
// 	FixedCost 			float64							`json:"fixed_period_cost"`
// 	ProductionRates	map[string]float64	`json:"production_rates"`
// 	CHODsTable			map[string]float64	`json:"changeover_durations"`
// }
//
// type LSProblem struct{
// 	Name 					string				`json:"name"`
// 	Periods 			uint32				`json:"periods"`
// 	PeriodLength	uint32				`json:"period_length"`
// 	Products			[]*LSPProduct	`json:"products"`
// 	Demands				[]*LSPDemand	`json:"demands"`
// 	Machines			[]*LSPMachine	`json:"machines"`
// }
//
type LSColumn struct{
	Column
	hash				string
}

func (c *LSColumn) GetHash() string  {
	if c.hash != "" {
		return c.hash
	}

	fingerprint_str:= fmt.Sprintf("%s;%d;{", c.MachineId, c.BucketIdx)
	for d_id, d_ff := range c.Column.GetDemandsFf() {
		fingerprint_str+= fmt.Sprintf("%s;%f.2", d_id, d_ff)
	}
	fingerprint_str+= "};{"
	for _, p := range c.Column.GetProdsSeq() {
		fingerprint_str+= fmt.Sprintf("%s;%d;%d", p.GetProductId(), p.GetProdStart(), p.GetProdMinutes())
	}
	fingerprint_str+= "};"

	fingerprint_bytes := sha1.Sum([]byte(fingerprint_str))
	c.hash = hex.EncodeToString(fingerprint_bytes[:])
	return c.hash
}
