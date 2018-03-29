package goh4

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/golang/glog"
)

// Inventory Inventory structure holding a single inventory item
type Inventory struct {
	Hostname    string   `json:"host_name"`
	AddressType string   `json:"address_type"`
	HostUUID    string   `json:"host_uuid"`
	IfaceMac    string   `json:"iface_mac"`
	IfaceName   string   `json:"iface_name"`
	IP          net.IP   `json:"ip,string"`
	Netmask     string   `json:"netmask"`
	OS          string   `json:"os"`
	OSVersion   string   `json:"os_version"`
	Internal    bool     `json:"tags_is_internal,string"`
	Scopes      []string `json:"tags_scope_name"`
	UUIDSource  string   `json:"uuid_src"`
	VRFID       int      `json:"vrf_id,string"`
	VRFName     string   `json:"vrf_name"`
}

// InventorySearch Perform a search on the inventory
func (h *H4) InventorySearch(s *InventoryQuery) ([]Inventory, error) {
	jsonStr, err := json.Marshal(&s)
	if err != nil {
		glog.Errorf("Error Marshalling search %s", err)
		return nil, err
	}

	postResp, err := h.Post("/inventory/search", fmt.Sprintf("%s", jsonStr))
	if err != nil {
		glog.Errorf("POST error %s", err)
		return nil, err
	}

	jsonResp := make(map[string][]Inventory)
	err = json.Unmarshal(postResp, &jsonResp)
	if err != nil {
		glog.Errorf("Error unmarshalling json %s", err)
		return nil, err
	}

	return jsonResp["results"], nil
}
