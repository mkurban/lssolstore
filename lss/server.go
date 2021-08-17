package lss;

import(
	"context"
	"io"
	"os"
	"log"
	"fmt"
	"encoding/hex"
	"strings"
	"crypto/sha1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//
func (s *LSSServer) ReadProblem(ctx context.Context, req *ReadProblemRequest) (*ProblemIdentifier, error) {
	problem_filename := req.GetProblemFile()
	reader, err := io.Reader(nil), error(nil)

	if problem_filename != "" {

		reader, err = os.Open("lss/testdata/"+problem_filename)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "Could not read JSON input 1: %v", err)
		}
	}else{
		reader	= strings.NewReader(req.GetProblemStr())
	}

	p, err := ReadProblem(reader)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Could not read JSON input 2: %v", err)
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	problem_name_as_bytes := sha1.Sum([]byte(p.Name))
	problem_id := hex.EncodeToString(problem_name_as_bytes[:])
	s.problems[problem_id]= NewProblemInfo(p)
	log.Printf("P_id == %s created. Total: %d problems!", problem_id, len(s.problems))

	return &ProblemIdentifier{ProblemId: problem_id}, nil
}

//
func (s *LSSServer) CreateSolution(_ context.Context, pi *ProblemIdentifier) (*ProblemEntityIdentifier, error) {
	problem_id := pi.GetProblemId()
	if problem_id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "LSS_BADPROBLEMID")
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()

	problem_info, ok := s.problems[problem_id]
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "LSS_BADPROBLEMID")
	}

	var fingerprint_str= fmt.Sprintf("%s;%d", problem_id, len(problem_info.solutions))
	var fingerprint_bytes= sha1.Sum([]byte(fingerprint_str))
	var sol_id= hex.EncodeToString(fingerprint_bytes[:])

	log.Printf("Creating Solution w/ID %s", sol_id)
	problem_info.solutions[sol_id]= NewSolution(sol_id, problem_id, &problem_info.Problem)

	return &ProblemEntityIdentifier{ProblemId: problem_id, EntityId: sol_id}, nil
}

//
func (s *LSSServer) GetProblem(_ context.Context, pi *ProblemIdentifier) (*Problem, error) {
	problem_id := pi.GetProblemId()
	if problem_id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "LSS_BADPROBLEMID")
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()

	problem_info, ok := s.problems[problem_id]
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "LSS_BADPROBLEMID")
	}

	log.Printf("Sending Problem w/ID %s", pi.GetProblemId())
	return &problem_info.Problem, nil
}

//
func (s *LSSServer) GetMachine(_ context.Context, pei *ProblemEntityIdentifier) (*Machine, error) {
	problem_id := pei.GetProblemId()
	if problem_id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "LSS_BADPROBLEMID")
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()

	problem_info, ok := s.problems[problem_id]
	if !ok {
		return nil, status.Errorf(codes.Unknown, "LSS_BADPROBLEMID")
	}

	for _, m := range problem_info.Problem.Machines {
		if pei.GetEntityId() == m.GetId() {
			log.Printf("Returning Machine w/id == %s", m.GetId())
			return m, nil
		}
	}

	return nil, status.Errorf(codes.NotFound, "")
}

//
func (s *LSSServer) GetDemandsForBucket(pbi *ProblemBucketIdentifier, demands_stream LotsizeSolutionStore_GetDemandsForBucketServer) error {
	problem_id := pbi.GetProblemId()
	if problem_id == "" {
		return status.Errorf(codes.InvalidArgument, "LSS_BADPROBLEMID")
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()

	problem_info, ok := s.problems[problem_id]
	if !ok {
		return status.Errorf(codes.Unknown, "LSS_BADPROBLEMID")
	}

	var target_machine = (*Machine)(nil)
	for _, m := range problem_info.Problem.Machines {
		if pbi.GetMachineId() == m.GetId() {
			log.Printf("Using Machine w/id == %s", m.GetId())
			target_machine= m
		}
	}
	if target_machine == nil {
		return status.Errorf(codes.NotFound, "")
	}

	log.Printf("Sending Demands for Bucket [%s,%d].", pbi.GetMachineId(), pbi.GetBucketIdx())
	bucket_idx := pbi.GetBucketIdx()
	for _, d := range problem_info.Problem.Demands {
		if bucket_idx <= d.GetPeriodIdx() {
			prod_rates := target_machine.GetProductionRates()
			if _, ok := prod_rates[d.GetProductId()]; ok {
				if err := demands_stream.Send(d); err != nil {
	        return err
	      }
			}
		}
	}

	return nil
}

//
func (s *LSSServer) GetSolutionsIDsList(_ context.Context, pi *ProblemIdentifier) (*IDsList, error) {
	problem_id := pi.GetProblemId()
	if problem_id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "LSS_BADPROBLEMID")
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()

	problem_info, ok := s.problems[problem_id]
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "LSS_BADPROBLEMID")
	}

	log.Printf("Sending Solution IDs for Problem w/id == %s", pi.GetProblemId())
	var ids = make([]string, 0, len(problem_info.solutions))
	for id := range problem_info.solutions {
		ids = append(ids, id)
	}

	return &IDsList{Ids:ids}, nil
}

//
func (s *LSSServer) GetSolution(pei *ProblemEntityIdentifier, columns_stream LotsizeSolutionStore_GetSolutionServer) error {
	problem_id := pei.GetProblemId()
	if problem_id == "" {
		return status.Errorf(codes.InvalidArgument, "LSS_BADPROBLEMID")
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()

	problem_info, ok := s.problems[problem_id]
	if !ok {
		return status.Errorf(codes.Unknown, "LSS_BADPROBLEMID")
	}

	var solution *Solution
	solution, ok = problem_info.solutions[pei.GetEntityId()]
	if !ok {
		return status.Errorf(codes.Unknown, "LSS_BADSOLUTIONID")
	}

	log.Printf("Sending Solution w/ID %s", pei.GetEntityId())
	for _, cols_machine := range solution.Columns {
		for _, cols_bucket := range cols_machine {
			for _, col := range cols_bucket {
				if err := columns_stream.Send(&col.Column); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

//
func (s *LSSServer) SolutionImprovementStart(_ context.Context, pei *ProblemEntityIdentifier) (*SolutionEntityIdentifier, error) {
	problem_id := pei.GetProblemId()
	if problem_id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "LSS_BADPROBLEMID")
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()

	problem_info, ok := s.problems[problem_id]
	if !ok {
		return nil, status.Errorf(codes.Unknown, "LSS_BADPROBLEMID")
	}

	var solution *Solution
	solution, ok = problem_info.solutions[pei.GetEntityId()]
	if !ok {
		return nil, status.Errorf(codes.Unknown, "LSS_BADSOLUTIONID")
	}

	if solution.imprToken != "" {
		return nil, status.Errorf(codes.Unavailable, "LSS_BADTOKEN")
	}

	var token_seed			= fmt.Sprintf("%s;%s;1", problem_id, solution.Id)
	var token_bytes			= sha1.Sum([]byte(token_seed))
	solution.imprToken	= hex.EncodeToString(token_bytes[:])

	return &SolutionEntityIdentifier{ProblemId: problem_id, SolutionId: solution.Id, ImpToken: solution.imprToken}, nil
}

func (s *LSSServer) SolutionImprovementColsAdded(_ context.Context, sei *SolutionEntityIdentifier) (*Int32Result, error) {
	problem_id := sei.GetProblemId()
	if problem_id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "LSS_BADPROBLEMID")
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()

	problem_info, ok := s.problems[problem_id]
	if !ok {
		return nil, status.Errorf(codes.Unknown, "LSS_BADPROBLEMID")
	}

	var solution *Solution
	solution, ok = problem_info.solutions[sei.GetSolutionId()]
	if !ok {
		return nil, status.Errorf(codes.Unknown, "LSS_BADSOLUTIONID")
	}

	if solution.imprToken != sei.ImpToken {
		return nil, status.Errorf(codes.Unavailable, "LSS_BADTOKEN")
	}

	// ToDo: Add receved cols

	return &Int32Result{Result: -2}, nil
}

// //
func (s *LSSServer) SolutionImprovementStop(_ context.Context, sei *SolutionEntityIdentifier) (*Int32Result, error) {
	problem_id := sei.GetProblemId()
	if problem_id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "LSS_BADPROBLEMID")
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()

	problem_info, ok := s.problems[problem_id]
	if !ok {
		return nil, status.Errorf(codes.Unknown, "LSS_BADPROBLEMID")
	}

	var solution *Solution
	solution, ok = problem_info.solutions[sei.GetSolutionId()]
	if !ok {
		return nil, status.Errorf(codes.Unknown, "LSS_BADSOLUTIONID")
	}

	if solution.imprToken != sei.ImpToken {
		return nil, status.Errorf(codes.Unavailable, "LSS_BADTOKEN")
	}

	// ToDo: Mark receved cols as seleted ones
	solution.imprToken = ""

	return &Int32Result{Result: -1}, nil
}

func (s *LSSServer)	mustEmbedUnimplementedLotsizeSolutionStoreServer(){}
