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
func (h *H4) GetFlows(q *FlowQuery, offset string) (*FlowResults, error) {
	var err error
	q.Offset = offset
	j, err := json.Marshal(q)
	if err != nil {
		return nil, err
	}

	r, err := h.Post("/flowsearch", string(j))
	if err != nil {
		return nil, err
	}

	var res *FlowResults
	err = json.Unmarshal(r, &res)
	if err != nil {
		return nil, err
	}

	// Some debug stuff we should remove later
	if false {
		fmt.Printf("offset: %s, count: %d\n", res.Offset, len(res.Results))
	}

	if res.Offset != "" {
		f, _ := h.GetFlows(q, res.Offset)
		res.Results = append(res.Results, f.Results...)
	}

	return res, nil
}
