package goh4

import (
	"encoding/json"
	"fmt"
)

// Application struct representing and application workspace
type Application struct {
	ID                 string                `json:"id,omitempty"`
	Scope              *Scope                `json:"app_scope_id"`
	Name               string                `json:"name"`
	Description        string                `json:"description"`
	Author             string                `json:"author"`
	Primary            bool                  `json:"primary"`
	AlternateQueryMode bool                  `json:"alternate_query_mode"`
	LatestVersion      int                   `json:"latest_adm_version"`
	EnforcementEnabled bool                  `json:"enforcement_enabled"`
	EnforcementVersion int                   `json:"enforced_version"`
	Clusters           []*ApplicationCluster `json:"clusters"`
	Filters            []*InventoryFilter    `json:"inventory_filters"`
	CatchAll           string                `json:"catch_all_action"`
	AbsolutePolicies   []*ApplicationPolicy  `json:"absolute_policies"`
	DefaultPolicies    []*ApplicationPolicy  `json:"default_policies"`
}

// MarshalJSON Converts Struct to JSON
func (a *Application) MarshalJSON() ([]byte, error) {
	type Alias Application
	return json.Marshal(&struct {
		Scope string `json:"app_scope_id"`
		*Alias
	}{
		Scope: a.Scope.ID,
		Alias: (*Alias)(a),
	})
}

// UnmarshalJSON Converts JSON to struct
func (a *Application) UnmarshalJSON(data []byte) error {
	var err error
	type Alias Application

	aux := &struct {
		Scope string `json:"app_scope_id"`
		*Alias
	}{
		Alias: (*Alias)(a),
	}

	if err = json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if h4 != nil {
		var scope *Scope

		if aux.Scope != "" {
			scope, err = h4.GetScope(aux.Scope)
			if err != nil {
				return err
			}
		} else {
			scope = nil
		}

		a.Scope = scope
	} else {
		return fmt.Errorf("H4 is not defined")
	}
	return nil
}

// L4Params struct representing L4 information (TCP / UDP / Port...)
type L4Params struct {
	Proto int   `json:"proto"`
	Port  []int `json:"port"`
}

// ApplicationPolicy Policy defined within an Application workspace
type ApplicationPolicy struct {
	ConsumerFilter string      `json:"consumer_filter_id"`
	ProviderFilter string      `json:"provider_filter_id"`
	Action         string      `json:"action"`
	L4Params       []*L4Params `json:"l4_params"`
}

// // MarshalJSON Converts Struct to JSON
// func (ap *ApplicationPolicy) MarshalJSON() ([]byte, error) {
// 	type Alias ApplicationPolicy
// 	return json.Marshal(&struct {
// 		ConsumerFilter string `json:"consumer_filter_id"`
// 		ProviderFilter string `json:"provider_filter_id"`
// 		*Alias
// 	}{
// 		ConsumerFilter: ap.ConsumerFilter.ID,
// 		ProviderFilter: ap.ProviderFilter.ID,
// 		Alias:          (*Alias)(ap),
// 	})
// }
//
// // UnmarshalJSON Converts JSON to struct
// func (ap *ApplicationPolicy) UnmarshalJSON(data []byte) error {
// 	var err error
// 	type Alias ApplicationPolicy
//
// 	aux := &struct {
// 		ConsumerFilter string `json:"consumer_filter_id"`
// 		ProviderFilter string `json:"provider_filter_id"`
// 		*Alias
// 	}{
// 		Alias: (*Alias)(ap),
// 	}
//
// 	if err = json.Unmarshal(data, &aux); err != nil {
// 		return err
// 	}
// 	if h4 != nil {
// 		var consumerFilter *InventoryFilter
// 		var providerFilter *InventoryFilter
//
// 		if aux.ConsumerFilter != "" {
// 			consumerFilter, err = h4.GetFilter(aux.ConsumerFilter)
// 			if err != nil {
// 				return fmt.Errorf("Error getting consumerFilter %s - Error: %s", aux.ConsumerFilter, err.Error())
// 			}
// 		} else {
// 			consumerFilter = nil
// 		}
//
// 		if aux.ProviderFilter != "" {
// 			providerFilter, err = h4.GetFilter(aux.ProviderFilter)
// 			if err != nil {
// 				return fmt.Errorf("Error getting providerFilter %s - Error: %s", aux.ProviderFilter, err.Error())
// 			}
// 		} else {
// 			providerFilter = nil
// 		}
//
// 		ap.ConsumerFilter = consumerFilter
// 		ap.ProviderFilter = providerFilter
// 	} else {
// 		return fmt.Errorf("H4 is not defined")
// 	}
// 	return nil
// }

// ApplicationNode Node that is part of a cluster
type ApplicationNode struct {
	IP   string `json:"ip"`
	Name string `json:"name"`
}

// ApplicationCluster Group of nodes
type ApplicationCluster struct {
	ID             string            `json:"id,omitempty"`
	Name           string            `json:"name"`
	Nodes          []ApplicationNode `json:"nodes"`
	ConsistentUUID string            `json:"consistent_uuid"`
}

// GetApplication Get a single application by ID
func (h *H4) GetApplication(id string) (*Application, error) {
	getResp, err := h.Get(fmt.Sprintf("/applications/%s/details", id))
	if err != nil {
		return nil, fmt.Errorf("GET error: %s", err.Error())
	}

	var jsonResp *Application
	//fmt.Printf("%s\n", getResp)

	err = json.Unmarshal(getResp, &jsonResp)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling JSON: %s / JSON: %s", err.Error(), getResp)
	}

	return jsonResp, nil
}

// GetAllApplications Get all applications
func (h *H4) GetAllApplications(details bool) ([]*Application, error) {
	getResp, err := h.Get(fmt.Sprintf("/applications"))
	if err != nil {
		return nil, fmt.Errorf("GET error: %s", err.Error())
	}
	//fmt.Printf("%s", getResp)

	var jsonResp []*Application
	err = json.Unmarshal(getResp, &jsonResp)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling JSON: %s / JSON: %s", err.Error(), getResp)
	}

	if details {
		for i := range jsonResp {
			app, err := h.GetApplication(jsonResp[i].ID)
			if err != nil {
				return nil, fmt.Errorf("Error unmarshalling ID %s / %s", jsonResp[i].ID, err.Error())
			}
			jsonResp[i] = app
		}
	}

	return jsonResp, nil
}
