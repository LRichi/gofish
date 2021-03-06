//
// SPDX-License-Identifier: BSD-3-Clause
//

package redfish

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/LRichi/WBfish/common"
)

var managerBody = `{
		"@Redfish.Copyright": "Copyright 2014-2019 DMTF. All rights reserved.",
		"@odata.context": "/redfish/v1/$metadata#Manager.Manager",
		"@odata.id": "/redfish/v1/Managers/BMC-1",
		"@odata.type": "#Manager.v1_1_0.Manager",
		"Id": "BMC-1",
		"Name": "Manager",
		"ManagerType": "BMC",
		"Description": "BMC",
		"AutoDSTEnabled": true,
		"ServiceEntryPointUUID": "92384634-2938-2342-8820-489239905423",
		"UUID": "00000000-0000-0000-0000-000000000000",
		"Model": "Joo Janta 200",
		"DateTime": "2015-03-13T04:14:33+06:00",
		"DateTimeLocalOffset": "+06:00",
		"PowerState": "On",
		"Status": {
			"State": "Enabled",
			"Health": "OK"
		},
		"GraphicalConsole": {
			"ServiceEnabled": true,
			"MaxConcurrentSessions": 2,
			"ConnectTypesSupported": [
				"KVMIP"
			]
		},
		"SerialConsole": {
			"ServiceEnabled": true,
			"MaxConcurrentSessions": 1,
			"ConnectTypesSupported": [
				"Telnet",
				"SSH",
				"IPMI"
			]
		},
		"CommandShell": {
			"ServiceEnabled": true,
			"MaxConcurrentSessions": 4,
			"ConnectTypesSupported": [
				"Telnet",
				"SSH"
			]
		},
		"FirmwareVersion": "1.00",
		"RemoteAccountService": {
			"@odata.id": "/redfish/v1/Managers/AccountService"
		},
		"RemoteRedfishServiceUri": "http://example.com/",
		"NetworkProtocol": {
			"@odata.id": "/redfish/v1/Managers/BMC-1/NetworkProtocol"
		},
		"HostInterfaces": {
			"@odata.id": "/redfish/v1/Managers/BMC-1/HostInterfaces"
		},
		"EthernetInterfaces": {
			"@odata.id": "/redfish/v1/Managers/BMC-1/EthernetInterfaces"
		},
		"SerialInterfaces": {
			"@odata.id": "/redfish/v1/Managers/BMC-1/SerialInterfaces"
		},
		"LogServices": {
			"@odata.id": "/redfish/v1/Managers/BMC-1/LogServices"
		},
		"VirtualMedia": {
			"@odata.id": "/redfish/v1/Managers/BMC-1/VM1"
		},
		"Links": {
			"ManagerForServers": [
				{
					"@odata.id": "/redfish/v1/Systems/System-1"
				}
			],
			"ManagerForChassis": [
				{
					"@odata.id": "/redfish/v1/Chassis/Chassis-1"
				}
			],
			"ManagerInChassis": {
				"@odata.id": "/redfish/v1/Chassis/Chassis-1"
			}
		},
		"Actions": {
			"#Manager.Reset": {
				"target": "/redfish/v1/Managers/BMC-1/Actions/Manager.Reset",
				"ResetType@Redfish.AllowableValues": [
					"ForceRestart",
					"GracefulRestart"
				]
			}
		}
	}`

// TestManager tests the parsing of Manager objects.
func TestManager(t *testing.T) {
	var result Manager
	err := json.NewDecoder(strings.NewReader(managerBody)).Decode(&result)

	if err != nil {
		t.Errorf("Error decoding JSON: %s", err)
	}

	if result.ID != "BMC-1" {
		t.Errorf("Received invalid ID: %s", result.ID)
	}

	if result.Name != "Manager" {
		t.Errorf("Received invalid name: %s", result.Name)
	}

	if result.ManagerType != BMCManagerType {
		t.Errorf("Received manager type: %s", result.ManagerType)
	}

	if result.PowerState != OnPowerState {
		t.Errorf("Received power state: %s", result.PowerState)
	}

	if !result.GraphicalConsole.ServiceEnabled {
		t.Error("Graphical console service state should be enabled")
	}

	if len(result.SerialConsole.ConnectTypesSupported) != 3 {
		t.Errorf("Serial console should have 3 connect types, got %d",
			len(result.SerialConsole.ConnectTypesSupported))
	}

	if result.managerForServers[0] != "/redfish/v1/Systems/System-1" {
		t.Errorf("Received manager for servers: %s", result.managerForServers)
	}

	if result.resetTarget != "/redfish/v1/Managers/BMC-1/Actions/Manager.Reset" {
		t.Errorf("Invalid Reset target: %s", result.resetTarget)
	}
}

// TestManagerUpdate tests the Update call.
func TestManagerUpdate(t *testing.T) {
	var result Manager
	err := json.NewDecoder(strings.NewReader(managerBody)).Decode(&result)

	if err != nil {
		t.Errorf("Error decoding JSON: %s", err)
	}

	testClient := &common.TestClient{}
	result.SetClient(testClient)

	result.AutoDSTEnabled = false
	result.DateTimeLocalOffset = "+05:00"
	err = result.Update()

	if err != nil {
		t.Errorf("Error making Update call: %s", err)
	}

	calls := testClient.CapturedCalls()

	if !strings.Contains(calls[0].Payload, "AutoDSTEnabled:false") {
		t.Errorf("Unexpected AutoDSTEnabled update payload: %s", calls[0].Payload)
	}

	if !strings.Contains(calls[0].Payload, "DateTimeLocalOffset:+05:00") {
		t.Errorf("Unexpected DateTimeLocalOffset update payload: %s", calls[0].Payload)
	}
}
