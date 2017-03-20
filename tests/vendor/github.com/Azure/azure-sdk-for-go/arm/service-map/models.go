package servicemap

// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Code generated by Microsoft (R) AutoRest Code Generator 1.0.1.0
// Changes may cause incorrect behavior and will be lost if the code is
// regenerated.

import (
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/Azure/go-autorest/autorest/to"
	"net/http"
)

// Accuracy enumerates the values for accuracy.
type Accuracy string

const (
	// Actual specifies the actual state for accuracy.
	Actual Accuracy = "actual"
	// Estimated specifies the estimated state for accuracy.
	Estimated Accuracy = "estimated"
)

// Bitness enumerates the values for bitness.
type Bitness string

const (
	// SixFourbit specifies the six fourbit state for bitness.
	SixFourbit Bitness = "64bit"
	// ThreeTwobit specifies the three twobit state for bitness.
	ThreeTwobit Bitness = "32bit"
)

// ConnectionFailureState enumerates the values for connection failure state.
type ConnectionFailureState string

const (
	// Failed specifies the failed state for connection failure state.
	Failed ConnectionFailureState = "failed"
	// Mixed specifies the mixed state for connection failure state.
	Mixed ConnectionFailureState = "mixed"
	// Ok specifies the ok state for connection failure state.
	Ok ConnectionFailureState = "ok"
)

// HypervisorType enumerates the values for hypervisor type.
type HypervisorType string

const (
	// Hyperv specifies the hyperv state for hypervisor type.
	Hyperv HypervisorType = "hyperv"
	// Unknown specifies the unknown state for hypervisor type.
	Unknown HypervisorType = "unknown"
)

// MachineRebootStatus enumerates the values for machine reboot status.
type MachineRebootStatus string

const (
	// MachineRebootStatusNotRebooted specifies the machine reboot status not
	// rebooted state for machine reboot status.
	MachineRebootStatusNotRebooted MachineRebootStatus = "notRebooted"
	// MachineRebootStatusRebooted specifies the machine reboot status rebooted
	// state for machine reboot status.
	MachineRebootStatusRebooted MachineRebootStatus = "rebooted"
	// MachineRebootStatusUnknown specifies the machine reboot status unknown
	// state for machine reboot status.
	MachineRebootStatusUnknown MachineRebootStatus = "unknown"
)

// MonitoringState enumerates the values for monitoring state.
type MonitoringState string

const (
	// Discovered specifies the discovered state for monitoring state.
	Discovered MonitoringState = "discovered"
	// Monitored specifies the monitored state for monitoring state.
	Monitored MonitoringState = "monitored"
)

// OperatingSystemFamily enumerates the values for operating system family.
type OperatingSystemFamily string

const (
	// OperatingSystemFamilyAix specifies the operating system family aix state
	// for operating system family.
	OperatingSystemFamilyAix OperatingSystemFamily = "aix"
	// OperatingSystemFamilyLinux specifies the operating system family linux
	// state for operating system family.
	OperatingSystemFamilyLinux OperatingSystemFamily = "linux"
	// OperatingSystemFamilySolaris specifies the operating system family
	// solaris state for operating system family.
	OperatingSystemFamilySolaris OperatingSystemFamily = "solaris"
	// OperatingSystemFamilyUnknown specifies the operating system family
	// unknown state for operating system family.
	OperatingSystemFamilyUnknown OperatingSystemFamily = "unknown"
	// OperatingSystemFamilyWindows specifies the operating system family
	// windows state for operating system family.
	OperatingSystemFamilyWindows OperatingSystemFamily = "windows"
)

// VirtualizationState enumerates the values for virtualization state.
type VirtualizationState string

const (
	// VirtualizationStateHypervisor specifies the virtualization state
	// hypervisor state for virtualization state.
	VirtualizationStateHypervisor VirtualizationState = "hypervisor"
	// VirtualizationStatePhysical specifies the virtualization state physical
	// state for virtualization state.
	VirtualizationStatePhysical VirtualizationState = "physical"
	// VirtualizationStateUnknown specifies the virtualization state unknown
	// state for virtualization state.
	VirtualizationStateUnknown VirtualizationState = "unknown"
	// VirtualizationStateVirtual specifies the virtualization state virtual
	// state for virtualization state.
	VirtualizationStateVirtual VirtualizationState = "virtual"
)

// VirtualMachineType enumerates the values for virtual machine type.
type VirtualMachineType string

const (
	// VirtualMachineTypeHyperv specifies the virtual machine type hyperv state
	// for virtual machine type.
	VirtualMachineTypeHyperv VirtualMachineType = "hyperv"
	// VirtualMachineTypeLdom specifies the virtual machine type ldom state for
	// virtual machine type.
	VirtualMachineTypeLdom VirtualMachineType = "ldom"
	// VirtualMachineTypeLpar specifies the virtual machine type lpar state for
	// virtual machine type.
	VirtualMachineTypeLpar VirtualMachineType = "lpar"
	// VirtualMachineTypeUnknown specifies the virtual machine type unknown
	// state for virtual machine type.
	VirtualMachineTypeUnknown VirtualMachineType = "unknown"
	// VirtualMachineTypeVirtualPc specifies the virtual machine type virtual
	// pc state for virtual machine type.
	VirtualMachineTypeVirtualPc VirtualMachineType = "virtualPc"
	// VirtualMachineTypeVmware specifies the virtual machine type vmware state
	// for virtual machine type.
	VirtualMachineTypeVmware VirtualMachineType = "vmware"
	// VirtualMachineTypeXen specifies the virtual machine type xen state for
	// virtual machine type.
	VirtualMachineTypeXen VirtualMachineType = "xen"
)

// Acceptor is a process accepting on a port.
type Acceptor struct {
	ID                  *string `json:"id,omitempty"`
	Type                *string `json:"type,omitempty"`
	Name                *string `json:"name,omitempty"`
	*AcceptorProperties `json:"properties,omitempty"`
}

// AcceptorProperties is properties for an acceptor relationship.
type AcceptorProperties struct {
	Source      *PortReference    `json:"source,omitempty"`
	Destination *ProcessReference `json:"destination,omitempty"`
	StartTime   *date.Time        `json:"startTime,omitempty"`
	EndTime     *date.Time        `json:"endTime,omitempty"`
}

// AgentConfiguration is describes the configuration of the Dependency Agent
// installed on a machine.
type AgentConfiguration struct {
	AgentID                 *string             `json:"agentId,omitempty"`
	DependencyAgentID       *string             `json:"dependencyAgentId,omitempty"`
	DependencyAgentVersion  *string             `json:"dependencyAgentVersion,omitempty"`
	DependencyAgentRevision *string             `json:"dependencyAgentRevision,omitempty"`
	RebootStatus            MachineRebootStatus `json:"rebootStatus,omitempty"`
	ClockGranularity        *int32              `json:"clockGranularity,omitempty"`
}

// ClientGroup is represents a collection of clients of a resource. A client
// group can represent the clients of a port, process, or a machine.
type ClientGroup struct {
	autorest.Response      `json:"-"`
	ID                     *string `json:"id,omitempty"`
	Type                   *string `json:"type,omitempty"`
	Name                   *string `json:"name,omitempty"`
	Etag                   *string `json:"etag,omitempty"`
	*ClientGroupProperties `json:"properties,omitempty"`
}

// ClientGroupProperties is resource properties.
type ClientGroupProperties struct {
	ClientsOf *ResourceReference `json:"clientsOf,omitempty"`
}

// ClientGroupMember is represents a member of a client group
type ClientGroupMember struct {
	ID                           *string `json:"id,omitempty"`
	Type                         *string `json:"type,omitempty"`
	Name                         *string `json:"name,omitempty"`
	*ClientGroupMemberProperties `json:"properties,omitempty"`
}

// ClientGroupMemberProperties is resource properties.
type ClientGroupMemberProperties struct {
	IPAddress *string             `json:"ipAddress,omitempty"`
	Port      *PortReference      `json:"port,omitempty"`
	Processes *[]ProcessReference `json:"processes,omitempty"`
}

// ClientGroupMembersCollection is collection of ClientGroupMember resources.
type ClientGroupMembersCollection struct {
	autorest.Response `json:"-"`
	Value             *[]ClientGroupMember `json:"value,omitempty"`
	NextLink          *string              `json:"nextLink,omitempty"`
}

// ClientGroupMembersCollectionPreparer prepares a request to retrieve the next set of results. It returns
// nil if no more results exist.
func (client ClientGroupMembersCollection) ClientGroupMembersCollectionPreparer() (*http.Request, error) {
	if client.NextLink == nil || len(to.String(client.NextLink)) <= 0 {
		return nil, nil
	}
	return autorest.Prepare(&http.Request{},
		autorest.AsJSON(),
		autorest.AsGet(),
		autorest.WithBaseURL(to.String(client.NextLink)))
}

// ClientGroupMembersCount is specifies the number of members in a client
// group.
type ClientGroupMembersCount struct {
	autorest.Response `json:"-"`
	StartTime         *date.Time `json:"startTime,omitempty"`
	EndTime           *date.Time `json:"endTime,omitempty"`
	GroupID           *string    `json:"groupId,omitempty"`
	Count             *int32     `json:"count,omitempty"`
	Accuracy          Accuracy   `json:"accuracy,omitempty"`
}

// Connection is a network connection.
type Connection struct {
	ID                    *string `json:"id,omitempty"`
	Type                  *string `json:"type,omitempty"`
	Name                  *string `json:"name,omitempty"`
	*ConnectionProperties `json:"properties,omitempty"`
}

// ConnectionCollection is collection of Connection resources.
type ConnectionCollection struct {
	autorest.Response `json:"-"`
	Value             *[]Connection `json:"value,omitempty"`
	NextLink          *string       `json:"nextLink,omitempty"`
}

// ConnectionCollectionPreparer prepares a request to retrieve the next set of results. It returns
// nil if no more results exist.
func (client ConnectionCollection) ConnectionCollectionPreparer() (*http.Request, error) {
	if client.NextLink == nil || len(to.String(client.NextLink)) <= 0 {
		return nil, nil
	}
	return autorest.Prepare(&http.Request{},
		autorest.AsJSON(),
		autorest.AsGet(),
		autorest.WithBaseURL(to.String(client.NextLink)))
}

// ConnectionProperties is properties for a connection resource.
type ConnectionProperties struct {
	Source       *ResourceReference     `json:"source,omitempty"`
	Destination  *ResourceReference     `json:"destination,omitempty"`
	StartTime    *date.Time             `json:"startTime,omitempty"`
	EndTime      *date.Time             `json:"endTime,omitempty"`
	ServerPort   *PortReference         `json:"serverPort,omitempty"`
	FailureState ConnectionFailureState `json:"failureState,omitempty"`
}

// CoreResource is marker resource for the core Service Map resources
type CoreResource struct {
	ID   *string `json:"id,omitempty"`
	Type *string `json:"type,omitempty"`
	Name *string `json:"name,omitempty"`
	Etag *string `json:"etag,omitempty"`
}

// Error is error details.
type Error struct {
	Code    *string `json:"code,omitempty"`
	Message *string `json:"message,omitempty"`
}

// ErrorResponse is an error response from the API.
type ErrorResponse struct {
	Error *Error `json:"error,omitempty"`
}

// HypervisorConfiguration is describes the hypervisor configuration of a
// machine.
type HypervisorConfiguration struct {
	HypervisorType      HypervisorType `json:"hypervisorType,omitempty"`
	NativeHostMachineID *string        `json:"nativeHostMachineId,omitempty"`
}

// Ipv4NetworkInterface is describes an IPv4 network interface.
type Ipv4NetworkInterface struct {
	IPAddress  *string `json:"ipAddress,omitempty"`
	SubnetMask *string `json:"subnetMask,omitempty"`
}

// Ipv6NetworkInterface is describes an IPv6 network interface.
type Ipv6NetworkInterface struct {
	IPAddress *string `json:"ipAddress,omitempty"`
}

// Liveness is specifies the contents of a check liveness response.
type Liveness struct {
	autorest.Response `json:"-"`
	StartTime         *date.Time `json:"startTime,omitempty"`
	EndTime           *date.Time `json:"endTime,omitempty"`
	Live              *bool      `json:"live,omitempty"`
}

// Machine is a machine resource represents a discovered computer system. It
// can be *monitored*, i.e., a Dependency Agent is running on it, or
// *discovered*, i.e., its existence was inferred by observing the data stream
// from monitored machines. As machines change, prior versions of the machine
// resource are preserved and available for access. A machine is live during an
// interval of time, if either its Dependency Agent has reported data during
// (parts) of that interval, or a Dependency agent running on other machines
// has reported activity associated with the machine.
type Machine struct {
	autorest.Response  `json:"-"`
	ID                 *string `json:"id,omitempty"`
	Type               *string `json:"type,omitempty"`
	Name               *string `json:"name,omitempty"`
	Etag               *string `json:"etag,omitempty"`
	*MachineProperties `json:"properties,omitempty"`
}

// MachineProperties is resource properties.
type MachineProperties struct {
	Timestamp                *date.Time                     `json:"timestamp,omitempty"`
	MonitoringState          MonitoringState                `json:"monitoringState,omitempty"`
	VirtualizationState      VirtualizationState            `json:"virtualizationState,omitempty"`
	DisplayName              *string                        `json:"displayName,omitempty"`
	ComputerName             *string                        `json:"computerName,omitempty"`
	FullyQualifiedDomainName *string                        `json:"fullyQualifiedDomainName,omitempty"`
	BootTime                 *date.Time                     `json:"bootTime,omitempty"`
	Timezone                 *Timezone                      `json:"timezone,omitempty"`
	Agent                    *AgentConfiguration            `json:"agent,omitempty"`
	Resources                *MachineResourcesConfiguration `json:"resources,omitempty"`
	Networking               *NetworkConfiguration          `json:"networking,omitempty"`
	OperatingSystem          *OperatingSystemConfiguration  `json:"operatingSystem,omitempty"`
	VirtualMachine           *VirtualMachineConfiguration   `json:"virtualMachine,omitempty"`
	Hypervisor               *HypervisorConfiguration       `json:"hypervisor,omitempty"`
}

// MachineCollection is collection of Machine resources.
type MachineCollection struct {
	autorest.Response `json:"-"`
	Value             *[]Machine `json:"value,omitempty"`
	NextLink          *string    `json:"nextLink,omitempty"`
}

// MachineCollectionPreparer prepares a request to retrieve the next set of results. It returns
// nil if no more results exist.
func (client MachineCollection) MachineCollectionPreparer() (*http.Request, error) {
	if client.NextLink == nil || len(to.String(client.NextLink)) <= 0 {
		return nil, nil
	}
	return autorest.Prepare(&http.Request{},
		autorest.AsJSON(),
		autorest.AsGet(),
		autorest.WithBaseURL(to.String(client.NextLink)))
}

// MachineCountsByOperatingSystem is machines by operating system.
type MachineCountsByOperatingSystem struct {
	Windows *int32 `json:"windows,omitempty"`
	Linux   *int32 `json:"linux,omitempty"`
}

// MachineGroup is a user-defined logical grouping of machines.
type MachineGroup struct {
	autorest.Response       `json:"-"`
	ID                      *string `json:"id,omitempty"`
	Type                    *string `json:"type,omitempty"`
	Name                    *string `json:"name,omitempty"`
	Etag                    *string `json:"etag,omitempty"`
	*MachineGroupProperties `json:"properties,omitempty"`
}

// MachineGroupProperties is resource properties.
type MachineGroupProperties struct {
	DisplayName *string                      `json:"displayName,omitempty"`
	Machines    *[]MachineReferenceWithHints `json:"machines,omitempty"`
}

// MachineGroupCollection is collection of Machine Group resources.
type MachineGroupCollection struct {
	autorest.Response `json:"-"`
	Value             *[]MachineGroup `json:"value,omitempty"`
	NextLink          *string         `json:"nextLink,omitempty"`
}

// MachineGroupCollectionPreparer prepares a request to retrieve the next set of results. It returns
// nil if no more results exist.
func (client MachineGroupCollection) MachineGroupCollectionPreparer() (*http.Request, error) {
	if client.NextLink == nil || len(to.String(client.NextLink)) <= 0 {
		return nil, nil
	}
	return autorest.Prepare(&http.Request{},
		autorest.AsJSON(),
		autorest.AsGet(),
		autorest.WithBaseURL(to.String(client.NextLink)))
}

// MachineGroupMapRequest is specifies the computation of a machine group
// dependency map. A machine group dependency map includes all direct
// dependencies of a group of machines.
type MachineGroupMapRequest struct {
	StartTime       *date.Time `json:"startTime,omitempty"`
	EndTime         *date.Time `json:"endTime,omitempty"`
	MachineGroupID  *string    `json:"machineGroupId,omitempty"`
	FilterProcesses *bool      `json:"filterProcesses,omitempty"`
}

// MachineReference is reference to a machine.
type MachineReference struct {
	ID   *string `json:"id,omitempty"`
	Type *string `json:"type,omitempty"`
	Name *string `json:"name,omitempty"`
}

// MachineReferenceWithHints is a machine reference with a hint of the
// machine's name and operating system.
type MachineReferenceWithHints struct {
	ID                                   *string `json:"id,omitempty"`
	Type                                 *string `json:"type,omitempty"`
	Name                                 *string `json:"name,omitempty"`
	*MachineReferenceWithHintsProperties `json:"properties,omitempty"`
}

// MachineReferenceWithHintsProperties is machine reference with name and os
// hints.
type MachineReferenceWithHintsProperties struct {
	DisplayNameHint *string               `json:"displayNameHint,omitempty"`
	OsFamilyHint    OperatingSystemFamily `json:"osFamilyHint,omitempty"`
}

// MachineResourcesConfiguration is describes the resources of a machine.
type MachineResourcesConfiguration struct {
	PhysicalMemory   *int32   `json:"physicalMemory,omitempty"`
	Cpus             *int32   `json:"cpus,omitempty"`
	CPUSpeed         *int32   `json:"cpuSpeed,omitempty"`
	CPUSpeedAccuracy Accuracy `json:"cpuSpeedAccuracy,omitempty"`
}

// MachinesSummary is a summary of the machines in the workspace.
type MachinesSummary struct {
	autorest.Response          `json:"-"`
	ID                         *string `json:"id,omitempty"`
	Type                       *string `json:"type,omitempty"`
	Name                       *string `json:"name,omitempty"`
	*MachinesSummaryProperties `json:"properties,omitempty"`
}

// MachinesSummaryProperties is summarizes machines in the workspace.
type MachinesSummaryProperties struct {
	StartTime *date.Time                      `json:"startTime,omitempty"`
	EndTime   *date.Time                      `json:"endTime,omitempty"`
	Total     *int32                          `json:"total,omitempty"`
	Live      *int32                          `json:"live,omitempty"`
	Os        *MachineCountsByOperatingSystem `json:"os,omitempty"`
}

// Map is a map of resources and relationships between them.
type Map struct {
	Nodes *MapNodes `json:"nodes,omitempty"`
	Edges *MapEdges `json:"edges,omitempty"`
}

// MapEdges is the edges (relationships) of a map.
type MapEdges struct {
	Connections *[]Connection `json:"connections,omitempty"`
	Acceptors   *[]Acceptor   `json:"acceptors,omitempty"`
}

// MapNodes is the nodes (entities) of a map.
type MapNodes struct {
	Machines     *[]Machine     `json:"machines,omitempty"`
	Processes    *[]Process     `json:"processes,omitempty"`
	Ports        *[]Port        `json:"Ports,omitempty"`
	ClientGroups *[]ClientGroup `json:"ClientGroups,omitempty"`
}

// MapRequest is specifies the contents of request to generate a map.
type MapRequest struct {
	StartTime *date.Time `json:"startTime,omitempty"`
	EndTime   *date.Time `json:"endTime,omitempty"`
}

// MapResponse is specified the contents of a map response.
type MapResponse struct {
	autorest.Response `json:"-"`
	StartTime         *date.Time `json:"startTime,omitempty"`
	EndTime           *date.Time `json:"endTime,omitempty"`
	Map               *Map       `json:"map,omitempty"`
}

// NetworkConfiguration is describes the network configuration of a machine.
type NetworkConfiguration struct {
	Ipv4Interfaces      *[]Ipv4NetworkInterface `json:"ipv4Interfaces,omitempty"`
	Ipv6Interfaces      *[]Ipv6NetworkInterface `json:"ipv6Interfaces,omitempty"`
	DefaultIpv4Gateways *[]string               `json:"defaultIpv4Gateways,omitempty"`
	MacAddresses        *[]string               `json:"macAddresses,omitempty"`
	DNSNames            *[]string               `json:"dnsNames,omitempty"`
}

// OperatingSystemConfiguration is describes the configuration of the operating
// system of a machine.
type OperatingSystemConfiguration struct {
	Family   OperatingSystemFamily `json:"family,omitempty"`
	FullName *string               `json:"fullName,omitempty"`
	Bitness  Bitness               `json:"bitness,omitempty"`
}

// Port is a port resource represents a server port on a machine. The port may
// be actively *monitored*, i.e., a Dependency Agent is running on its machine,
// or *discovered*, i.e., its existence was inferred by observing the data
// stream from monitored machines. A port is live during an interval of time,
// if that port had associated activity during (parts) of that interval.
type Port struct {
	autorest.Response `json:"-"`
	ID                *string `json:"id,omitempty"`
	Type              *string `json:"type,omitempty"`
	Name              *string `json:"name,omitempty"`
	Etag              *string `json:"etag,omitempty"`
	*PortProperties   `json:"properties,omitempty"`
}

// PortProperties is resource properties.
type PortProperties struct {
	MonitoringState MonitoringState    `json:"monitoringState,omitempty"`
	Machine         *ResourceReference `json:"machine,omitempty"`
	DisplayName     *string            `json:"displayName,omitempty"`
	IPAddress       *string            `json:"ipAddress,omitempty"`
	PortNumber      *int32             `json:"portNumber,omitempty"`
}

// PortCollection is collection of Port resources.
type PortCollection struct {
	autorest.Response `json:"-"`
	Value             *[]Port `json:"value,omitempty"`
	NextLink          *string `json:"nextLink,omitempty"`
}

// PortCollectionPreparer prepares a request to retrieve the next set of results. It returns
// nil if no more results exist.
func (client PortCollection) PortCollectionPreparer() (*http.Request, error) {
	if client.NextLink == nil || len(to.String(client.NextLink)) <= 0 {
		return nil, nil
	}
	return autorest.Prepare(&http.Request{},
		autorest.AsJSON(),
		autorest.AsGet(),
		autorest.WithBaseURL(to.String(client.NextLink)))
}

// PortReference is reference to a port.
type PortReference struct {
	ID                       *string `json:"id,omitempty"`
	Type                     *string `json:"type,omitempty"`
	Name                     *string `json:"name,omitempty"`
	*PortReferenceProperties `json:"properties,omitempty"`
}

// PortReferenceProperties is resource properties.
type PortReferenceProperties struct {
	Machine    *MachineReference `json:"machine,omitempty"`
	IPAddress  *string           `json:"ipAddress,omitempty"`
	PortNumber *int32            `json:"portNumber,omitempty"`
}

// Process is a process resource represents a process running on a machine. The
// process may be actively *monitored*, i.e., a Dependency Agent is running on
// its machine, or *discovered*, i.e., its existence was inferred by observing
// the data stream from monitored machines. A process resource represents a
// pool of actual operating system resources that share command lines and
// metadata. As the process pool evolves over time, prior versions of the
// process resource are preserved and available for access. A process is live
// during an interval of time, if that process is executing during (parts) of
// that interval
type Process struct {
	autorest.Response  `json:"-"`
	ID                 *string `json:"id,omitempty"`
	Type               *string `json:"type,omitempty"`
	Name               *string `json:"name,omitempty"`
	Etag               *string `json:"etag,omitempty"`
	*ProcessProperties `json:"properties,omitempty"`
}

// ProcessProperties is resource properties.
type ProcessProperties struct {
	Timestamp       *date.Time         `json:"timestamp,omitempty"`
	MonitoringState MonitoringState    `json:"monitoringState,omitempty"`
	Machine         *ResourceReference `json:"machine,omitempty"`
	ExecutableName  *string            `json:"executableName,omitempty"`
	DisplayName     *string            `json:"displayName,omitempty"`
	StartTime       *date.Time         `json:"startTime,omitempty"`
	Details         *ProcessDetails    `json:"details,omitempty"`
	User            *ProcessUser       `json:"user,omitempty"`
	ClientOf        *ResourceReference `json:"clientOf,omitempty"`
	AcceptorOf      *ResourceReference `json:"acceptorOf,omitempty"`
}

// ProcessCollection is collection of Process resources.
type ProcessCollection struct {
	autorest.Response `json:"-"`
	Value             *[]Process `json:"value,omitempty"`
	NextLink          *string    `json:"nextLink,omitempty"`
}

// ProcessCollectionPreparer prepares a request to retrieve the next set of results. It returns
// nil if no more results exist.
func (client ProcessCollection) ProcessCollectionPreparer() (*http.Request, error) {
	if client.NextLink == nil || len(to.String(client.NextLink)) <= 0 {
		return nil, nil
	}
	return autorest.Prepare(&http.Request{},
		autorest.AsJSON(),
		autorest.AsGet(),
		autorest.WithBaseURL(to.String(client.NextLink)))
}

// ProcessDetails is describes process metadata.
type ProcessDetails struct {
	PersistentKey    *string `json:"persistentKey,omitempty"`
	PoolID           *int32  `json:"poolId,omitempty"`
	FirstPid         *int32  `json:"firstPid,omitempty"`
	Description      *string `json:"description,omitempty"`
	CompanyName      *string `json:"companyName,omitempty"`
	InternalName     *string `json:"internalName,omitempty"`
	ProductName      *string `json:"productName,omitempty"`
	ProductVersion   *string `json:"productVersion,omitempty"`
	FileVersion      *string `json:"fileVersion,omitempty"`
	CommandLine      *string `json:"commandLine,omitempty"`
	ExecutablePath   *string `json:"executablePath,omitempty"`
	WorkingDirectory *string `json:"workingDirectory,omitempty"`
}

// ProcessReference is reference to a process.
type ProcessReference struct {
	ID                          *string `json:"id,omitempty"`
	Type                        *string `json:"type,omitempty"`
	Name                        *string `json:"name,omitempty"`
	*ProcessReferenceProperties `json:"properties,omitempty"`
}

// ProcessReferenceProperties is resource properties.
type ProcessReferenceProperties struct {
	Machine *MachineReference `json:"machine,omitempty"`
}

// ProcessUser is describes the user under which a process is running.
type ProcessUser struct {
	UserName   *string `json:"userName,omitempty"`
	UserDomain *string `json:"userDomain,omitempty"`
}

// Relationship is a typed relationship between two entities.
type Relationship struct {
	ID   *string `json:"id,omitempty"`
	Type *string `json:"type,omitempty"`
	Name *string `json:"name,omitempty"`
}

// RelationshipProperties is relationship properties.
type RelationshipProperties struct {
	Source      *ResourceReference `json:"source,omitempty"`
	Destination *ResourceReference `json:"destination,omitempty"`
	StartTime   *date.Time         `json:"startTime,omitempty"`
	EndTime     *date.Time         `json:"endTime,omitempty"`
}

// Resource is resource model definition.
type Resource struct {
	ID   *string `json:"id,omitempty"`
	Type *string `json:"type,omitempty"`
	Name *string `json:"name,omitempty"`
}

// ResourceReference is represents a reference to another resource.
type ResourceReference struct {
	ID   *string `json:"id,omitempty"`
	Type *string `json:"type,omitempty"`
	Name *string `json:"name,omitempty"`
}

// SingleMachineDependencyMapRequest is specifies the computation of a single
// server dependency map. A single server dependency map includes all direct
// dependencies of a given machine.
type SingleMachineDependencyMapRequest struct {
	StartTime *date.Time `json:"startTime,omitempty"`
	EndTime   *date.Time `json:"endTime,omitempty"`
	MachineID *string    `json:"machineId,omitempty"`
}

// Summary is base for all resource summaries.
type Summary struct {
	ID   *string `json:"id,omitempty"`
	Type *string `json:"type,omitempty"`
	Name *string `json:"name,omitempty"`
}

// SummaryProperties is base for all summaries.
type SummaryProperties struct {
	StartTime *date.Time `json:"startTime,omitempty"`
	EndTime   *date.Time `json:"endTime,omitempty"`
}

// Timezone is describes a timezone.
type Timezone struct {
	FullName *string `json:"fullName,omitempty"`
}

// VirtualMachineConfiguration is describes the virtualizaton-related
// configuration of a machine.
type VirtualMachineConfiguration struct {
	VirtualMachineType  VirtualMachineType `json:"virtualMachineType,omitempty"`
	NativeMachineID     *string            `json:"nativeMachineId,omitempty"`
	VirtualMachineName  *string            `json:"virtualMachineName,omitempty"`
	NativeHostMachineID *string            `json:"nativeHostMachineId,omitempty"`
}
