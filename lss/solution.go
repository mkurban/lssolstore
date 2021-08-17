package lss


type Solution struct{
	Id					string															`json:"id"`
	ProblemId		string															`json:"problem_id"`
	// map [machine_id] -> seq [ map[column_hash] -> column ]
	Columns			map[string][](map[string]*LSColumn)	`json:"columns"`
	imprToken		string
}

func NewSolution(id string, p_id string, p *Problem) *Solution{
	var s Solution
	s.Id = id
	s.ProblemId= p_id
	s.Columns= make(map[string][](map[string]*LSColumn))
	for _, m := range p.Machines {
		s.Columns[m.Id]= make([]map[string]*LSColumn, p.GetNumPeriods(), p.GetNumPeriods())
		var i uint32 = 0;
		for ; i < p.GetNumPeriods(); i++ {
			s.Columns[m.Id][i]= make(map[string]*LSColumn)
		}
	}
	return &s
}
