
package lss

import (
	"io"
	"encoding/json"
)


func ReadProblem(r io.Reader) (*Problem, error){
	var p = new(Problem)
	dec:= json.NewDecoder(r)
	if err:= dec.Decode(&p); err != nil {
		return nil, err;
	}
	return p, nil
}
