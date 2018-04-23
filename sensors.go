package goh4

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

// AgentInterface Interfaces of SW agents
type AgentInterface struct {
	Name    string `json:"name"`
	VRFID   int    `json:"vrf_id"`
	VRFName string `json:"vrf"`
	MAC     string `json:"mac"`
	IP      net.IP `json:"ip,string"`
	Family  string `json:"family_type"`
	NetMask string `json:"netmask"`
}

// SWAgent SWAgent representation
type SWAgent struct {
	UUID              string           `json:"uuid"`
	Hostname          string           `json:"host_name"`
	Platform          string           `json:"platform"`
	Type              string           `json:"agent_type"`
	CurrentVersion    string           `json:"current_sw_version"`
	DesiredVersion    string           `json:"desired_sw_version"`
	LastUpdate        time.Time        `json:"last_software_update_at"`
	DataPlaneDisabled bool             `json:"data_plane_disabled"`
	Forensics         bool             `json:"enable_forensics"`
	PIDLookup         bool             `json:"enable_pid_lookup"`
	AutoUpgradeOptOut bool             `json:"auto_upgrade_opt_out"`
	CPUQuotaMode      int              `json:"cpu_quota_mode"`
	CPUQuotaUSec      int              `json:"cpu_quota_usec"`
	ConfigFetch       time.Time        `json:"last_config_fetch_at"`
	Interfaces        []AgentInterface `json:"interfaces"`
}

// MarshalJSON Converts Struct to JSON
func (s *SWAgent) MarshalJSON() ([]byte, error) {
	type Alias SWAgent
	return json.Marshal(&struct {
		LastUpdate  int64 `json:"last_software_update_at"`
		ConfigFetch int64 `json:"last_config_fetch_at"`
		*Alias
	}{
		LastUpdate:  s.LastUpdate.Unix(),
		ConfigFetch: s.ConfigFetch.Unix(),
		Alias:       (*Alias)(s),
	})
}

// UnmarshalJSON Converts JSON to struct
func (s *SWAgent) UnmarshalJSON(data []byte) error {
	var err error
	type Alias SWAgent

	aux := &struct {
		LastUpdate  int64 `json:"last_software_update_at"`
		ConfigFetch int64 `json:"last_config_fetch_at"`
		*Alias
	}{
		Alias: (*Alias)(s),
	}

	if err = json.Unmarshal(data, &aux); err != nil {
		return err
	}

	s.LastUpdate = time.Unix(aux.LastUpdate, 0)
	s.ConfigFetch = time.Unix(aux.ConfigFetch, 0)
	return nil
}

type resultSWAgent struct {
	Results []SWAgent `json:"results"`
}

// GetSWAgents Get all Software Agents
func (h *H4) GetSWAgents() ([]SWAgent, error) {
	getResp, err := h.Get("/sensors")
	if err != nil {
		return nil, fmt.Errorf("GET error: %s / GET: %s", err.Error(), getResp)
	}

	var jsonResp resultSWAgent
	err = json.Unmarshal([]byte(getResp), &jsonResp)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling JSON: %s / JSON: %s", err.Error(), getResp)
	}

	return jsonResp.Results, nil
}
