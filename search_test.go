package goh4

import (
	"encoding/json"
	"testing"
)

func TestUnMarshalSimple(t *testing.T) {
	var err error

	t.Log("Starting TestUnMarshalSimple Test")

	c, err := loadJSONFromFile("testdata/queryFilterIsSimple.json")
	if err != nil {
		t.Errorf("Error in loadJSONFromFile: %s", err.Error())
		return
	}
	var qf QueryFilter
	err = json.Unmarshal(c, &qf)
	if err != nil {
		t.Errorf("Error in Unmarshal: %s", err.Error())
		return
	}

	if qf.Type != "and" {
		t.Errorf("Invalid type, expecting and, got %s\n", qf.Type)
	}

	if qf.IsFilter {
		t.Errorf("This is a filter?!\n")
	}

	seenVRFID := false
	seenIP := false

	for _, v := range qf.Filters {
		f := v.(*Filter)
		if f.Field == "vrf_id" {
			seenVRFID = true
			if f.Type != "eq" {
				t.Errorf("Invalid type for vrf_id, expecting eq, got %s\n", f.Type)
			} else if f.Value.(float64) != 1 {
				t.Errorf("Invalid value for vrf_id, expecting 1, got %d\n", f.Value)
			}
		}
		if f.Field == "ip" {
			seenIP = true
			if f.Type != "subnet" {
				t.Errorf("Invalid type for ip, expecting subnet, got %s\n", f.Type)
			} else if f.Value.(string) != "10.0.0.0/16" {
				t.Errorf("Invalid value for ip, expecting 10.0.0.0/16, got %s\n", f.Value)
			}
		}
	}

	if !seenVRFID {
		t.Errorf("Didn't find a vrf_id")
	}

	if !seenIP {
		t.Errorf("Didn't find an ip")
	}
}

func TestUnMarshalNested(t *testing.T) {
	var err error

	t.Log("Starting TestUnMarshalNested Test")

	c, err := loadJSONFromFile("testdata/queryFilterIsNested.json")
	if err != nil {
		t.Errorf("Error in loadJSONFromFile: %s", err.Error())
		return
	}
	var qf QueryFilter
	err = json.Unmarshal(c, &qf)
	if err != nil {
		t.Errorf("Error in Unmarshal: %s", err.Error())
		return
	}
	if qf.Type != "and" {
		t.Errorf("Invalid type, expecting and, got %s\n", qf.Type)
	}

	if qf.IsFilter {
		t.Errorf("This is a filter?!\n")
	}

	seenFilterVRFID := false
	seenFilterIP := false
	seenNestedIP1 := false
	seenNestedIP2 := false

	for _, v := range qf.Filters {
		switch v.(type) {
		case *Filter:
			f := v.(*Filter)
			if f.Field == "vrf_id" {
				seenFilterVRFID = true
				if f.Type != "eq" {
					t.Errorf("Invalid type for vrf_id, expecting eq, got %s\n", f.Type)
				} else if f.Value.(float64) != 1 {
					t.Errorf("Invalid value for vrf_id, expecting 1, got %d\n", f.Value)
				}
			}
			if f.Field == "ip" {
				seenFilterIP = true
				if f.Type != "subnet" {
					t.Errorf("Invalid type for ip, expecting subnet, got %s\n", f.Type)
				} else if f.Value.(string) != "172.21.0.0/16" {
					t.Errorf("Invalid value for ip, expecting 172.21.0.0/16, got %s\n", f.Value)
				}
			}

		case *QueryFilter:
			f := v.(*QueryFilter)
			if f.Type != "or" {
				t.Errorf("Invalid nested QF type, expecting or, got %s\n", f.Type)
			}
			if qf.IsFilter {
				t.Errorf("This Nested QF is a filter?!\n")
			}
			for _, nv := range f.Filters {
				nf := nv.(*Filter)
				if nf.Field == "ip" && nf.Value.(string) == "172.20.0.0/16" {
					seenNestedIP1 = true
				} else if nf.Field == "ip" && nf.Value.(string) == "172.21.0.0/16" {
					seenNestedIP2 = true
				}
			}
		}
	}

	if !seenFilterVRFID {
		t.Errorf("Didn't find a vrf_id Filter")
	}

	if !seenFilterIP {
		t.Errorf("Didn't find an ip Filter")
	}

	if !seenNestedIP1 {
		t.Errorf("Didn't find an ip Filter for 172.20.0.0/16 in nested QF")
	}

	if !seenNestedIP2 {
		t.Errorf("Didn't find an ip Filter for 172.21.0.0/16 in nested QF")
	}
}

func TestUnMarshalFilter(t *testing.T) {
	var err error

	t.Log("Starting TestUnMarshalFilter Test")

	c, err := loadJSONFromFile("testdata/queryFilterIsFilter.json")
	if err != nil {
		t.Errorf("Error in loadJSONFromFile: %s", err.Error())
		return
	}
	var qf QueryFilter
	err = json.Unmarshal(c, &qf)
	if err != nil {
		t.Errorf("Error in Unmarshal: %s", err.Error())
		return
	}
}

func TestUnMarshalInvalid(t *testing.T) {
	var err error

	t.Log("Starting TestUnMarshalInvalid Test")

	c, err := loadJSONFromFile("testdata/queryFilterIsInvalid.json")
	if err != nil {
		t.Errorf("Error in loadJSONFromFile: %s", err.Error())
		return
	}
	var qf QueryFilter
	err = json.Unmarshal(c, &qf)
	if err == nil {
		t.Errorf("Unmarshal did not error, that's bad")
		return
	}
}

func TestUnMarshalInvalidFilter(t *testing.T) {
	var err error

	t.Log("Starting TestUnMarshalInvalidFilter Test")

	c, err := loadJSONFromFile("testdata/queryFilterIsInvalidFilter.json")
	if err != nil {
		t.Errorf("Error in loadJSONFromFile: %s", err.Error())
		return
	}
	var qf QueryFilter
	err = json.Unmarshal(c, &qf)
	if err == nil {
		t.Errorf("Unmarshal did not error, that's bad")
		return
	}
}

func TestUnMarshalInvalidNestedQF(t *testing.T) {
	var err error

	t.Log("Starting TestUnMarshalInvalidNestedQF Test")

	c, err := loadJSONFromFile("testdata/queryFilterIsInvalidNestedQF.json")
	if err != nil {
		t.Errorf("Error in loadJSONFromFile: %s", err.Error())
		return
	}
	var qf QueryFilter
	err = json.Unmarshal(c, &qf)
	if err == nil {
		t.Errorf("Unmarshal did not error, that's bad")
		return
	}
}

func TestUnMarshalInvalidSimpleFilter(t *testing.T) {
	var err error

	t.Log("Starting TestUnMarshalInvalidSimpleFilter Test")

	c, err := loadJSONFromFile("testdata/queryFilterIsInvalidSimpleFilter.json")
	if err != nil {
		t.Errorf("Error in loadJSONFromFile: %s", err.Error())
		return
	}
	var qf QueryFilter
	err = json.Unmarshal(c, &qf)
	if err == nil {
		t.Errorf("Unmarshal did not error, that's bad")
		return
	}
}
