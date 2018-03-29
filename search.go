package goh4

import (
	"encoding/json"
)

// QueryFilter defines a set of filter for a search
type QueryFilter struct {
	Type     string        `json:"type"`
	Filters  []interface{} `json:"filters"`
	IsFilter bool          `json:"-"`
}

// UnmarshalJSON Unmarshal a query filter
func (q *QueryFilter) UnmarshalJSON(data []byte) error {
	// Register a temporary QueryFilter without Filters, but with a RawMessage
	type tempQF struct {
		Type    string          `json:"type,omitempty"`
		Filters json.RawMessage `json:"filters,omitempty"`
	}

	// Register our variables
	var tqf tempQF
	var qf *QueryFilter
	var fi *Filter
	var err error

	// Unmarshal the current data in the temporary QueryFilter
	err = json.Unmarshal(data, &tqf)
	if err != nil {
		return err
	}

	// we don't like unmarshalling empty strings / nil
	if tqf.Filters != nil {
		// Start populating our real QueryFilter
		q.Type = tqf.Type
		q.IsFilter = false

		var rawArray []json.RawMessage
		err = json.Unmarshal(tqf.Filters, &rawArray)
		if err != nil {
			return err
		}

		for i := range rawArray {
			// We need to know what we're dealing with, at this point we look at the content of the raw message to decide
			m := make(map[string]interface{})
			err = json.Unmarshal(rawArray[i], &m)
			if err != nil {
				return err
			}

			// if there is no value and no field it's most likely a QF
			if m["value"] == nil && m["field"] == nil {
				// This is a nested struct, here we use recursion to get the content
				// Recursion is automatic via the Unmarshal
				qf = new(QueryFilter)
				err = json.Unmarshal(rawArray[i], &qf)
				if err != nil {
					return err
				}
				qf.IsFilter = false
				q.Filters = append(q.Filters, qf)

			} else {
				// This is a filter
				fi = new(Filter)
				err = json.Unmarshal(rawArray[i], &fi)
				if err != nil {
					return err
				}
				q.Filters = append(q.Filters, fi)
			}
		}
	} else {
		fi = new(Filter)
		err = json.Unmarshal(data, &fi)
		if err != nil {
			return err
		}
		q.IsFilter = true
		q.Filters = append(q.Filters, fi)
	}
	return nil
}

// Filter defines a filter
type Filter struct {
	Field string      `json:"field"`
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

// InventoryQuery defines a search query
type InventoryQuery struct {
	Filter *QueryFilter `json:"filter"`
	Scope  string       `json:"scopeName,omitempty"`
	Limit  int          `json:"limit"`
	Offset int          `json:"offset,omitempty"`
}
