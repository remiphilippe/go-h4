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
	Author             string                `json:"author,omitempty"`
	StrictValidation   bool                  `json:"strict_validation"`
	Primary            bool                  `json:"primary"`
	AlternateQueryMode bool                  `json:"alternate_query_mode"`
	LatestVersion      int                   `json:"latest_adm_version,omitempty"`
	EnforcementEnabled bool                  `json:"enforcement_enabled,omitempty"`
	EnforcementVersion int                   `json:"enforced_version,omitempty"`
	Clusters           []*ApplicationCluster `json:"clusters,omitempty"`
	Filters            []*InventoryFilter    `json:"inventory_filters,omitempty"`
	CatchAll           string                `json:"catch_all_action,omitempty"`
	AbsolutePolicies   []*ApplicationPolicy  `json:"absolute_policies,omitempty"`
	DefaultPolicies    []*ApplicationPolicy  `json:"default_policies,omitempty"`
	h4                 *H4
}

// MarshalJSON Converts Struct to JSON
func (a *Application) MarshalJSON() ([]byte, error) {
	if a.Scope == nil {
		return nil, fmt.Errorf("scope is not defined")
	}

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
	a              *Application
	ID             string           `json:"id,omitempty"`
	ConsumerFilter *InventoryFilter `json:"consumer_filter_id"`
	ProviderFilter *InventoryFilter `json:"provider_filter_id"`
	Action         string           `json:"action"`
	L4Params       []*L4Params      `json:"l4_params"`
	Version        string           `json:"version,omitempty"`
	Priority       int              `json:"priority,omitempty"`
	Rank           string           `json:"rank,omitempty"`
}

// MarshalJSON Converts Struct to JSON
func (ap *ApplicationPolicy) MarshalJSON() ([]byte, error) {
	type Alias ApplicationPolicy
	return json.Marshal(&struct {
		ConsumerFilter string `json:"consumer_filter_id"`
		ProviderFilter string `json:"provider_filter_id"`
		*Alias
	}{
		ConsumerFilter: ap.ConsumerFilter.ID,
		ProviderFilter: ap.ProviderFilter.ID,
		Alias:          (*Alias)(ap),
	})
}

// UnmarshalJSON Converts JSON to struct
func (ap *ApplicationPolicy) UnmarshalJSON(data []byte) error {
	var err error
	type Alias ApplicationPolicy

	aux := &struct {
		ConsumerFilter string `json:"consumer_filter_id"`
		ProviderFilter string `json:"provider_filter_id"`
		*Alias
	}{
		Alias: (*Alias)(ap),
	}

	if err = json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if h4 != nil {
		var providerFilter *InventoryFilter
		var consumerFilter *InventoryFilter

		if aux.ProviderFilter != "" {
			providerFilter, err = h4.GetFilter(aux.ProviderFilter)
			if err != nil {
				return err
			}
		} else {
			providerFilter = nil
		}

		ap.ProviderFilter = providerFilter

		if aux.ConsumerFilter != "" {
			consumerFilter, err = h4.GetFilter(aux.ConsumerFilter)
			if err != nil {
				return err
			}
		} else {
			consumerFilter = nil
		}

		ap.ConsumerFilter = consumerFilter
	} else {
		return fmt.Errorf("H4 is not defined")
	}
	return nil
}

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

	jsonResp.h4 = h

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

	for k := range jsonResp {
		if jsonResp[k] != nil {
			jsonResp[k].h4 = h
		}
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

// AddApplication Add a new Application Workspace to Tetration
func (h *H4) AddApplication(a *Application) error {
	jsonStr, err := json.Marshal(&a)
	if err != nil {
		return fmt.Errorf("Error Marshalling application workspace %s", err)
	}
	postResp, err := h.Post("/applications", fmt.Sprintf("%s", jsonStr))
	if err != nil {
		return fmt.Errorf("POST error: %s / POST: %s", err.Error(), jsonStr)
	}

	err = json.Unmarshal(postResp, &a)
	if err != nil {
		return fmt.Errorf("Error unmarshalling JSON: %s / JSON: %s", err.Error(), postResp)
	}

	a.h4 = h

	return nil
}

// DeleteApplication Delete an Application Workspace
func (h *H4) DeleteApplication(applicationID string) error {
	err := h.Delete(fmt.Sprintf("/applications/%s", applicationID), "")
	if err != nil {
		return fmt.Errorf("Error deleting application workspace %s: %s", applicationID, err)
	}

	return nil
}

// UpdateApplication Update an Application Workspace
func (h *H4) UpdateApplication(a *Application) error {
	if a.ID == "" {
		return fmt.Errorf("error ID is nil")
	}

	jsonStr, err := json.Marshal(&a)
	if err != nil {
		return fmt.Errorf("Error Marshalling application workspace %s", err)
	}

	putURL := fmt.Sprintf("/applications/%s", a.ID)

	postResp, err := h.Put(putURL, fmt.Sprintf("%s", jsonStr))
	if err != nil {
		return fmt.Errorf("POST error: %s / POST: %s", err.Error(), jsonStr)
	}

	err = json.Unmarshal(postResp, &a)
	if err != nil {
		return fmt.Errorf("Error unmarshalling JSON: %s / JSON: %s", err.Error(), postResp)
	}

	a.h4 = h

	return nil
}

func (a *Application) addPolicy(p *ApplicationPolicy, polType string) error {
	var postURL string
	jsonStr, err := json.Marshal(&p)
	if err != nil {
		return fmt.Errorf("Error Marshalling application policy %s", err)
	}

	if polType == "default" || polType == "" {
		postURL = fmt.Sprintf("/applications/%s/default_policies", a.ID)
	} else {
		postURL = fmt.Sprintf("/applications/%s/absolute_policies", a.ID)
	}

	postResp, err := a.h4.Post(postURL, fmt.Sprintf("%s", jsonStr))
	if err != nil {
		return fmt.Errorf("POST error: %s / POST: %s", err.Error(), jsonStr)
	}

	err = json.Unmarshal(postResp, &p)
	if err != nil {
		return fmt.Errorf("Error unmarshalling JSON: %s / JSON: %s", err.Error(), postResp)
	}

	p.a = a

	return nil
}

// AddDefaultPolicy Adds a Default Policy
func (a *Application) AddDefaultPolicy(p *ApplicationPolicy) error {
	return a.addPolicy(p, "default")
}

// AddAbsolutePolicy Adds a Default Policy
func (a *Application) AddAbsolutePolicy(p *ApplicationPolicy) error {
	return a.addPolicy(p, "absolute")
}

// AddServicePort Adds a Service Port to a Policy
func (ap *ApplicationPolicy) AddServicePort(start, end, proto int) error {
	m := make(map[string]interface{})
	m["start_port"] = start
	m["end_port"] = end
	m["proto"] = proto
	m["version"] = fmt.Sprintf("v%d", ap.a.LatestVersion)

	jsonStr, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("Error Marshalling service port %s", err)
	}

	postURL := fmt.Sprintf("/policies/%s/l4_params", ap.ID)

	_, err = ap.a.h4.Post(postURL, fmt.Sprintf("%s", jsonStr))
	if err != nil {
		return fmt.Errorf("POST error: %s / POST: %s", err.Error(), jsonStr)
	}

	return nil
}

// GetApplicationByName Returns a scope based on it's name
func (h *H4) GetApplicationByName(name string) (*Application, error) {
	apps, err := h.GetAllApplications(false)
	if err != nil {
		return nil, err
	}

	for _, a := range apps {
		if a.Name == name {
			return h.GetApplication(a.ID)
		}
	}

	return nil, nil
}
