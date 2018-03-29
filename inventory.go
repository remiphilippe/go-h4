package goh4

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/golang/glog"
)

// Inventory Inventory structure holding a single inventory item
// tags_scope_name can be a string
// E0329 11:34:51.834252   10294 inventory.go:46] Error unmarshalling json json: cannot unmarshal string into Go struct field Inventory.tags_scope_name of type []string / content: {"results":[{"address_type":"IPV4","host_name":"app01","host_uuid":"f22e3c96884084b5d4334ba977e9ca8a7564c1f1","iface_mac":"fa:16:3e:58:84:24","iface_name":"eth0","ip":"192.168.17.151","netmask":"255.255.255.0","os":"CentOS","os_version":"7.4","tags_is_internal":"true","tags_scope_name":"Default","user_aci-fab":null,"user_anp":null,"user_des":null,"user_desc":null,"user_dns":null,"user_epg":null,"user_pod":null,"user_subnet":null,"user_support-group":"cips_group","user_tenant":null,"user_type":null,"user_vrf":null,"user_zone":null,"uuid_src":"SENSOR","vrf_id":"1","vrf_name":"Default"}]}
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
		glog.Errorf("Error unmarshalling json %s / content: %s", err, postResp)
		return nil, err
	}

	return jsonResp["results"], nil
}
