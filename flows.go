package goh4

import (
	"encoding/json"
	"fmt"

	"github.com/golang/glog"
)

// FlowResults results of a flow search
type FlowResults struct {
	Offset  string                   `json:"offset"`
	Results []map[string]interface{} `json:"results"`
}

// GetFlows Get Flows Matching
func (h *H4) GetFlows(q *FlowQuery) (*FlowResults, error) {
	// q.Limit max requested limit
	maxLimit := q.Limit
	if maxLimit > TetrationFlowQueryLimit || maxLimit == 0 {
		q.Limit = TetrationFlowQueryLimit
	}

	limit := 0
	offset := ""

	res := make(chan *FlowResults, 1)
	defer close(res)
	done := make(chan bool)
	defer close(done)

	fres := new(FlowResults)

	go func() {
		for r := range res {
			fres.Results = append(fres.Results, r.Results...)
			offset = r.Offset
			limit = limit + len(r.Results)
			fmt.Printf("offset: %s, limit: %d\n", offset, limit)

			if (maxLimit == 0 || limit < maxLimit) && offset != "" {
				go h.getFlows(q, offset, res)
			} else {
				done <- true
			}
		}
	}()

	h.getFlows(q, offset, res)
	<-done

	return fres, nil
}

func (h *H4) getFlows(q *FlowQuery, offset string, fres chan *FlowResults) {
	var err error
	q.Offset = offset
	j, err := json.Marshal(q)
	if err != nil {
		glog.Errorf("Error: %s\n", err)
	}

	r, err := h.Post("/flowsearch", string(j))
	if err != nil {
		glog.Errorf("Error: %s\n", err)
	}

	var res *FlowResults
	err = json.Unmarshal(r, &res)
	if err != nil {
		glog.Errorf("Error: %s\n", err)
	}

	// Some debug stuff we should remove later
	if false {
		fmt.Printf("offset: %s, count: %d\n", res.Offset, len(res.Results))
	}
	fres <- res
}
