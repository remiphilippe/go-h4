package goh4

import (
	"encoding/json"
	"fmt"
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
	err := make(chan error, 1)
	defer close(err)

	fres := new(FlowResults)

	go func() {
		for r := range res {
			fres.Results = append(fres.Results, r.Results...)
			offset = r.Offset
			limit = limit + len(r.Results)
			if false {
				fmt.Printf("offset: %s, limit: %d\n", offset, limit)
			}

			if (maxLimit == 0 || limit < maxLimit) && offset != "" {
				go h.getFlows(q, offset, res, err)
			} else {
				done <- true
			}
		}
	}()

	h.getFlows(q, offset, res, err)

	select {
	case <-done:
	case e := <-err:
		if e != nil {
			return nil, e
		}
	}

	<-done

	return fres, nil
}

func (h *H4) getFlows(q *FlowQuery, offset string, fres chan *FlowResults, err chan error) {
	var e error
	q.Offset = offset
	j, e := json.Marshal(q)
	if e != nil {
		err <- e
		return
	}

	r, e := h.Post("/flowsearch", string(j))
	if e != nil {
		err <- e
		return
	}

	var res *FlowResults
	e = json.Unmarshal(r, &res)
	if e != nil {
		err <- e
		return
	}

	// Some debug stuff we should remove later
	if false {
		fmt.Printf("offset: %s, count: %d\n", res.Offset, len(res.Results))
	}
	fres <- res
}
