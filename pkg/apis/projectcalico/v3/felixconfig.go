// Copyright (c) 2017-2021 Tigera, Inc. All rights reserved.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v3

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/tigera/api/pkg/lib/numorstring"
)

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// FelixConfigurationList contains a list of FelixConfiguration object.
type FelixConfigurationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Items []FelixConfiguration `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type FelixConfiguration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec FelixConfigurationSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
}

type IptablesBackend string

const (
	KindFelixConfiguration     = "FelixConfiguration"
	KindFelixConfigurationList = "FelixConfigurationList"
	IptablesBackendLegacy      = "Legacy"
	IptablesBackendNFTables    = "NFT"
)

// +kubebuilder:validation:Enum=DoNothing;Enable;Disable
type AWSSrcDstCheckOption string

const (
	AWSSrcDstCheckOptionDoNothing AWSSrcDstCheckOption = "DoNothing"
	AWSSrcDstCheckOptionEnable                         = "Enable"
	AWSSrcDstCheckOptionDisable                        = "Disable"
)

// FelixConfigurationSpec contains the values of the Felix configuration.
type FelixConfigurationSpec struct {
	UseInternalDataplaneDriver *bool  `json:"useInternalDataplaneDriver,omitempty"`
	DataplaneDriver            string `json:"dataplaneDriver,omitempty"`

	IPv6Support *bool `json:"ipv6Support,omitempty" confignamev1:"Ipv6Support"`

	// RouteRefreshInterval is the period at which Felix re-checks the routes
	// in the dataplane to ensure that no other process has accidentally broken Calico's rules.
	// Set to 0 to disable route refresh. [Default: 90s]
	RouteRefreshInterval *metav1.Duration `json:"routeRefreshInterval,omitempty" configv1timescale:"seconds"`
	// InterfaceRefreshInterval is the period at which Felix rescans local interfaces to verify their state.
	// The rescan can be disabled by setting the interval to 0.
	InterfaceRefreshInterval *metav1.Duration `json:"interfaceRefreshInterval,omitempty" configv1timescale:"seconds"`
	// IptablesRefreshInterval is the period at which Felix re-checks the IP sets
	// in the dataplane to ensure that no other process has accidentally broken Calico's rules.
	// Set to 0 to disable IP sets refresh. Note: the default for this value is lower than the
	// other refresh intervals as a workaround for a Linux kernel bug that was fixed in kernel
	// version 4.11. If you are using v4.11 or greater you may want to set this to, a higher value
	// to reduce Felix CPU usage. [Default: 10s]
	IptablesRefreshInterval *metav1.Duration `json:"iptablesRefreshInterval,omitempty" configv1timescale:"seconds"`
	// IptablesPostWriteCheckInterval is the period after Felix has done a write
	// to the dataplane that it schedules an extra read back in order to check the write was not
	// clobbered by another process. This should only occur if another application on the system
	// doesn't respect the iptables lock. [Default: 1s]
	IptablesPostWriteCheckInterval *metav1.Duration `json:"iptablesPostWriteCheckInterval,omitempty" configv1timescale:"seconds" confignamev1:"IptablesPostWriteCheckIntervalSecs"`
	// IptablesLockFilePath is the location of the iptables lock file. You may need to change this
	// if the lock file is not in its standard location (for example if you have mapped it into Felix's
	// container at a different path). [Default: /run/xtables.lock]
	IptablesLockFilePath string `json:"iptablesLockFilePath,omitempty"`
	// IptablesLockTimeout is the time that Felix will wait for the iptables lock,
	// or 0, to disable. To use this feature, Felix must share the iptables lock file with all other
	// processes that also take the lock. When running Felix inside a container, this requires the
	// /run directory of the host to be mounted into the calico/node or calico/felix container.
	// [Default: 0s disabled]
	IptablesLockTimeout *metav1.Duration `json:"iptablesLockTimeout,omitempty" configv1timescale:"seconds" confignamev1:"IptablesLockTimeoutSecs"`
	// IptablesLockProbeInterval is the time that Felix will wait between
	// attempts to acquire the iptables lock if it is not available. Lower values make Felix more
	// responsive when the lock is contended, but use more CPU. [Default: 50ms]
	IptablesLockProbeInterval *metav1.Duration `json:"iptablesLockProbeInterval,omitempty" configv1timescale:"milliseconds" confignamev1:"IptablesLockProbeIntervalMillis"`
	// FeatureDetectOverride is used to override the feature detection.
	// Values are specified in a comma separated list with no spaces, example;
	// "SNATFullyRandom=true,MASQFullyRandom=false,RestoreSupportsLock=".
	// "true" or "false" will force the feature, empty or omitted values are
	// auto-detected.
	FeatureDetectOverride string `json:"featureDetectOverride,omitempty" validate:"omitempty,keyValueList"`
	// IpsetsRefreshInterval is the period at which Felix re-checks all iptables
	// state to ensure that no other process has accidentally broken Calico's rules. Set to 0 to
	// disable iptables refresh. [Default: 90s]
	IpsetsRefreshInterval *metav1.Duration `json:"ipsetsRefreshInterval,omitempty" configv1timescale:"seconds"`
	MaxIpsetSize          *int             `json:"maxIpsetSize,omitempty"`
	// IptablesBackend specifies which backend of iptables will be used. The default is legacy.
	IptablesBackend *IptablesBackend `json:"iptablesBackend,omitempty" validate:"omitempty,iptablesBackend"`

	// XDPRefreshInterval is the period at which Felix re-checks all XDP state to ensure that no
	// other process has accidentally broken Calico's BPF maps or attached programs. Set to 0 to
	// disable XDP refresh. [Default: 90s]
	XDPRefreshInterval *metav1.Duration `json:"xdpRefreshInterval,omitempty" configv1timescale:"seconds"`

	NetlinkTimeout *metav1.Duration `json:"netlinkTimeout,omitempty" configv1timescale:"seconds" confignamev1:"NetlinkTimeoutSecs"`

	// MetadataAddr is the IP address or domain name of the server that can answer VM queries for
	// cloud-init metadata. In OpenStack, this corresponds to the machine running nova-api (or in
	// Ubuntu, nova-api-metadata). A value of none (case insensitive) means that Felix should not
	// set up any NAT rule for the metadata path. [Default: 127.0.0.1]
	MetadataAddr string `json:"metadataAddr,omitempty"`
	// MetadataPort is the port of the metadata server. This, combined with global.MetadataAddr (if
	// not 'None'), is used to set up a NAT rule, from 169.254.169.254:80 to MetadataAddr:MetadataPort.
	// In most cases this should not need to be changed [Default: 8775].
	MetadataPort *int `json:"metadataPort,omitempty"`

	// OpenstackRegion is the name of the region that a particular Felix belongs to. In a multi-region
	// Calico/OpenStack deployment, this must be configured somehow for each Felix (here in the datamodel,
	// or in felix.cfg or the environment on each compute node), and must match the [calico]
	// openstack_region value configured in neutron.conf on each node. [Default: Empty]
	OpenstackRegion string `json:"openstackRegion,omitempty"`

	// InterfacePrefix is the interface name prefix that identifies workload endpoints and so distinguishes
	// them from host endpoint interfaces. Note: in environments other than bare metal, the orchestrators
	// configure this appropriately. For example our Kubernetes and Docker integrations set the 'cali' value,
	// and our OpenStack integration sets the 'tap' value. [Default: cali]
	InterfacePrefix string `json:"interfacePrefix,omitempty"`
	// InterfaceExclude is a comma-separated list of interfaces that Felix should exclude when monitoring for host
	// endpoints. The default value ensures that Felix ignores Kubernetes' IPVS dummy interface, which is used
	// internally by kube-proxy. If you want to exclude multiple interface names using a single value, the list
	// supports regular expressions. For regular expressions you must wrap the value with '/'. For example
	// having values '/^kube/,veth1' will exclude all interfaces that begin with 'kube' and also the interface
	// 'veth1'. [Default: kube-ipvs0]
	InterfaceExclude string `json:"interfaceExclude,omitempty"`

	// ChainInsertMode controls whether Felix hooks the kernel's top-level iptables chains by inserting a rule
	// at the top of the chain or by appending a rule at the bottom. insert is the safe default since it prevents
	// Calico's rules from being bypassed. If you switch to append mode, be sure that the other rules in the chains
	// signal acceptance by falling through to the Calico rules, otherwise the Calico policy will be bypassed.
	// [Default: insert]
	ChainInsertMode string `json:"chainInsertMode,omitempty"`
	// DefaultEndpointToHostAction controls what happens to traffic that goes from a workload endpoint to the host
	// itself (after the traffic hits the endpoint egress policy). By default Calico blocks traffic from workload
	// endpoints to the host itself with an iptables "DROP" action. If you want to allow some or all traffic from
	// endpoint to host, set this parameter to RETURN or ACCEPT. Use RETURN if you have your own rules in the iptables
	// "INPUT" chain; Calico will insert its rules at the top of that chain, then "RETURN" packets to the "INPUT" chain
	// once it has completed processing workload endpoint egress policy. Use ACCEPT to unconditionally accept packets
	// from workloads after processing workload endpoint egress policy. [Default: Drop]
	DefaultEndpointToHostAction string `json:"defaultEndpointToHostAction,omitempty" validate:"omitempty,dropAcceptReturn"`
	IptablesFilterAllowAction   string `json:"iptablesFilterAllowAction,omitempty" validate:"omitempty,acceptReturn"`
	IptablesMangleAllowAction   string `json:"iptablesMangleAllowAction,omitempty" validate:"omitempty,acceptReturn"`
	// LogPrefix is the log prefix that Felix uses when rendering LOG rules. [Default: calico-packet]
	LogPrefix string `json:"logPrefix,omitempty"`

	// LogDropActionOverride specifies whether or not to include the DropActionOverride in the logs when it is triggered.
	LogDropActionOverride *bool `json:"logDropActionOverride,omitempty"`

	// LogFilePath is the full path to the Felix log. Set to none to disable file logging. [Default: /var/log/calico/felix.log]
	LogFilePath string `json:"logFilePath,omitempty"`

	// LogSeverityFile is the log severity above which logs are sent to the log file. [Default: Info]
	LogSeverityFile string `json:"logSeverityFile,omitempty" validate:"omitempty,logLevel"`
	// LogSeverityScreen is the log severity above which logs are sent to the stdout. [Default: Info]
	LogSeverityScreen string `json:"logSeverityScreen,omitempty" validate:"omitempty,logLevel"`
	// LogSeveritySys is the log severity above which logs are sent to the syslog. Set to None for no logging to syslog.
	// [Default: Info]
	LogSeveritySys string `json:"logSeveritySys,omitempty" validate:"omitempty,logLevel"`

	IPIPEnabled *bool `json:"ipipEnabled,omitempty" confignamev1:"IpInIpEnabled"`
	// IPIPMTU is the MTU to set on the tunnel device. See Configuring MTU [Default: 1440]
	IPIPMTU *int `json:"ipipMTU,omitempty" confignamev1:"IpInIpMtu"`

	VXLANEnabled *bool `json:"vxlanEnabled,omitempty"`
	// VXLANMTU is the MTU to set on the tunnel device. See Configuring MTU [Default: 1440]
	VXLANMTU  *int `json:"vxlanMTU,omitempty"`
	VXLANPort *int `json:"vxlanPort,omitempty"`
	VXLANVNI  *int `json:"vxlanVNI,omitempty"`

	// AllowVXLANPacketsFromWorkloads controls whether Felix will add a rule to drop VXLAN encapsulated traffic
	// from workloads [Default: false]
	// +optional
	AllowVXLANPacketsFromWorkloads *bool `json:"allowVXLANPacketsFromWorkloads,omitempty"`
	// AllowIPIPPacketsFromWorkloads controls whether Felix will add a rule to drop IPIP encapsulated traffic
	// from workloads [Default: false]
	// +optional
	AllowIPIPPacketsFromWorkloads *bool `json:"allowIPIPPacketsFromWorkloads,omitempty"`

	// ReportingInterval is the interval at which Felix reports its status into the datastore or 0 to disable.
	// Must be non-zero in OpenStack deployments. [Default: 30s]
	ReportingInterval *metav1.Duration `json:"reportingInterval,omitempty" configv1timescale:"seconds" confignamev1:"ReportingIntervalSecs"`
	// ReportingTTL is the time-to-live setting for process-wide status reports. [Default: 90s]
	ReportingTTL *metav1.Duration `json:"reportingTTL,omitempty" configv1timescale:"seconds" confignamev1:"ReportingTTLSecs"`

	EndpointReportingEnabled *bool            `json:"endpointReportingEnabled,omitempty"`
	EndpointReportingDelay   *metav1.Duration `json:"endpointReportingDelay,omitempty" configv1timescale:"seconds" confignamev1:"EndpointReportingDelaySecs"`

	// IptablesMarkMask is the mask that Felix selects its IPTables Mark bits from. Should be a 32 bit hexadecimal
	// number with at least 8 bits set, none of which clash with any other mark bits in use on the system.
	// [Default: 0xff000000]
	IptablesMarkMask *uint32 `json:"iptablesMarkMask,omitempty"`

	DisableConntrackInvalidCheck *bool `json:"disableConntrackInvalidCheck,omitempty"`

	HealthEnabled *bool   `json:"healthEnabled,omitempty"`
	HealthHost    *string `json:"healthHost,omitempty"`
	HealthPort    *int    `json:"healthPort,omitempty"`

	// PrometheusMetricsEnabled enables the Prometheus metrics server in Felix if set to true. [Default: false]
	PrometheusMetricsEnabled *bool `json:"prometheusMetricsEnabled,omitempty"`
	// PrometheusMetricsHost is the host that the Prometheus metrics server should bind to. [Default: empty]
	PrometheusMetricsHost string `json:"prometheusMetricsHost,omitempty" validate:"omitempty,prometheusHost"`
	// PrometheusMetricsPort is the TCP port that the Prometheus metrics server should bind to. [Default: 9091]
	PrometheusMetricsPort *int `json:"prometheusMetricsPort,omitempty"`
	// PrometheusGoMetricsEnabled disables Go runtime metrics collection, which the Prometheus client does by default, when
	// set to false. This reduces the number of metrics reported, reducing Prometheus load. [Default: true]
	PrometheusGoMetricsEnabled *bool `json:"prometheusGoMetricsEnabled,omitempty"`
	// PrometheusProcessMetricsEnabled disables process metrics collection, which the Prometheus client does by default, when
	// set to false. This reduces the number of metrics reported, reducing Prometheus load. [Default: true]
	PrometheusProcessMetricsEnabled *bool `json:"prometheusProcessMetricsEnabled,omitempty"`
	// PrometheusWireGuardMetricsEnabled disables wireguard metrics collection, which the Prometheus client does by default, when
	// set to false. This reduces the number of metrics reported, reducing Prometheus load. [Default: true]
	PrometheusWireGuardMetricsEnabled *bool `json:"prometheusWireGuardMetricsEnabled,omitempty"`
	// TLS credentials for this port.
	PrometheusMetricsCertFile string `json:"prometheusMetricsCertFile,omitempty"`
	PrometheusMetricsKeyFile  string `json:"prometheusMetricsKeyFile,omitempty"`
	PrometheusMetricsCAFile   string `json:"prometheusMetricsCAFile,omitempty"`

	// FailsafeInboundHostPorts is a list of UDP/TCP ports and CIDRs that Felix will allow incoming traffic to host endpoints
	// on irrespective of the security policy. This is useful to avoid accidentally cutting off a host with incorrect configuration.
	// For back-compatibility, if the protocol is not specified, it defaults to "tcp". If a CIDR is not specified, it will allow
	// traffic from all addresses. To disable all inbound host ports, use the value none. The default value allows ssh access
	// and DHCP.
	// [Default: tcp:22, udp:68, tcp:179, tcp:2379, tcp:2380, tcp:6443, tcp:6666, tcp:6667]
	FailsafeInboundHostPorts *[]ProtoPort `json:"failsafeInboundHostPorts,omitempty"`
	// FailsafeOutboundHostPorts is a list of UDP/TCP ports and CIDRs that Felix will allow outgoing traffic from host endpoints
	// to irrespective of the security policy. This is useful to avoid accidentally cutting off a host with incorrect configuration.
	// For back-compatibility, if the protocol is not specified, it defaults to "tcp". If a CIDR is not specified, it will allow
	// traffic from all addresses. To disable all outbound host ports, use the value none. The default value opens etcd's standard
	// ports to ensure that Felix does not get cut off from etcd as well as allowing DHCP and DNS.
	// [Default: tcp:179, tcp:2379, tcp:2380, tcp:6443, tcp:6666, tcp:6667, udp:53, udp:67]
	FailsafeOutboundHostPorts *[]ProtoPort `json:"failsafeOutboundHostPorts,omitempty"`

	// KubeMasqueradeBit should be set to the same value as --iptables-masquerade-bit of kube-proxy
	// when TPROXY is used. The default is the same as kube-proxy default thus only needs a change
	// if kube-proxy is using a non-standard setting. Must be within the range of 0-31.  [Default: 14]
	KubeMasqueradeBit *int `json:"kubeMasqueradeBit,omitempty" validate:"omitempty,gte=0,lte=31"`

	// KubeNodePortRanges holds list of port ranges used for service node ports. Only used if felix detects kube-proxy running in ipvs mode.
	// Felix uses these ranges to separate host and workload traffic. [Default: 30000:32767].
	KubeNodePortRanges *[]numorstring.Port `json:"kubeNodePortRanges,omitempty" validate:"omitempty,dive"`

	// PolicySyncPathPrefix is used to by Felix to communicate policy changes to external services,
	// like Application layer policy. [Default: Empty]
	PolicySyncPathPrefix string `json:"policySyncPathPrefix,omitempty"`

	// UsageReportingEnabled reports anonymous Calico version number and cluster size to projectcalico.org. Logs warnings returned by the usage
	// server. For example, if a significant security vulnerability has been discovered in the version of Calico being used. [Default: true]
	UsageReportingEnabled *bool `json:"usageReportingEnabled,omitempty"`

	// UsageReportingInitialDelay controls the minimum delay before Felix makes a report. [Default: 300s]
	UsageReportingInitialDelay *metav1.Duration `json:"usageReportingInitialDelay,omitempty" configv1timescale:"seconds" confignamev1:"UsageReportingInitialDelaySecs"`
	// UsageReportingInterval controls the interval at which Felix makes reports. [Default: 86400s]
	UsageReportingInterval *metav1.Duration `json:"usageReportingInterval,omitempty" configv1timescale:"seconds" confignamev1:"UsageReportingIntervalSecs"`

	// NATPortRange specifies the range of ports that is used for port mapping when doing outgoing NAT. When unset the default behavior of the
	// network stack is used.
	NATPortRange *numorstring.Port `json:"natPortRange,omitempty"`

	// NATOutgoingAddress specifies an address to use when performing source NAT for traffic in a natOutgoing pool that
	// is leaving the network. By default the address used is an address on the interface the traffic is leaving on
	// (ie it uses the iptables MASQUERADE target)
	NATOutgoingAddress string `json:"natOutgoingAddress,omitempty"`

	// This is the source address to use on programmed device routes. By default the source address is left blank,
	// leaving the kernel to choose the source address used.
	DeviceRouteSourceAddress string `json:"deviceRouteSourceAddress,omitempty"`

	// This defines the route protocol added to programmed device routes, by default this will be RTPROT_BOOT
	// when left blank.
	DeviceRouteProtocol *int `json:"deviceRouteProtocol,omitempty"`
	// Whether or not to remove device routes that have not been programmed by Felix. Disabling this will allow external
	// applications to also add device routes. This is enabled by default which means we will remove externally added routes.
	RemoveExternalRoutes *bool `json:"removeExternalRoutes,omitempty"`

	// ExternalNodesCIDRList is a list of CIDR's of external-non-calico-nodes which may source tunnel traffic and have
	// the tunneled traffic be accepted at calico nodes.
	ExternalNodesCIDRList *[]string `json:"externalNodesList,omitempty"`

	NfNetlinkBufSize  string `json:"nfNetlinkBufSize,omitempty"`
	StatsDumpFilePath string `json:"statsDumpFilePath,omitempty"`

	// Felix Denied Packet Metrics configuration parameters.
	PrometheusReporterEnabled   *bool  `json:"prometheusReporterEnabled,omitempty"`
	PrometheusReporterPort      *int   `json:"prometheusReporterPort,omitempty"`
	PrometheusReporterCertFile  string `json:"prometheusReporterCertFile,omitempty"`
	PrometheusReporterKeyFile   string `json:"prometheusReporterKeyFile,omitempty"`
	PrometheusReporterCAFile    string `json:"prometheusReporterCAFile,omitempty"`
	DeletedMetricsRetentionSecs *int   `json:"deletedMetricsRetentionSecs,omitempty"`

	// DropActionOverride overrides the Drop action in Felix, optionally changing the behavior to Accept, and optionally adding Log.
	// Possible values are Drop, LogAndDrop, Accept, LogAndAccept. [Default: Drop]
	DropActionOverride string `json:"dropActionOverride,omitempty" validate:"omitempty,dropActionOverride"`

	DebugMemoryProfilePath          string           `json:"debugMemoryProfilePath,omitempty"`
	DebugDisableLogDropping         *bool            `json:"debugDisableLogDropping,omitempty"`
	DebugSimulateCalcGraphHangAfter *metav1.Duration `json:"debugSimulateCalcGraphHangAfter,omitempty" configv1timescale:"seconds"`
	DebugSimulateDataplaneHangAfter *metav1.Duration `json:"debugSimulateDataplaneHangAfter,omitempty" configv1timescale:"seconds"`

	IptablesNATOutgoingInterfaceFilter string `json:"iptablesNATOutgoingInterfaceFilter,omitempty" validate:"omitempty,ifaceFilter"`

	// SidecarAccelerationEnabled enables experimental sidecar acceleration [Default: false]
	SidecarAccelerationEnabled *bool `json:"sidecarAccelerationEnabled,omitempty"`

	// XDPEnabled enables XDP acceleration for suitable untracked incoming deny rules. [Default: true]
	XDPEnabled *bool `json:"xdpEnabled,omitempty" confignamev1:"XDPEnabled"`

	// GenericXDPEnabled enables Generic XDP so network cards that don't support XDP offload or driver
	// modes can use XDP. This is not recommended since it doesn't provide better performance than
	// iptables. [Default: false]
	GenericXDPEnabled *bool `json:"genericXDPEnabled,omitempty" confignamev1:"GenericXDPEnabled"`

	// BPFEnabled, if enabled Felix will use the BPF dataplane. [Default: false]
	BPFEnabled *bool `json:"bpfEnabled,omitempty" validate:"omitempty"`
	// BPFDisableUnprivileged, if enabled, Felix sets the kernel.unprivileged_bpf_disabled sysctl to disable
	// unprivileged use of BPF.  This ensures that unprivileged users cannot access Calico's BPF maps and
	// cannot insert their own BPF programs to interfere with Calico's. [Default: true]
	BPFDisableUnprivileged *bool `json:"bpfDisableUnprivileged,omitempty" validate:"omitempty"`
	// BPFLogLevel controls the log level of the BPF programs when in BPF dataplane mode.  One of "Off", "Info", or
	// "Debug".  The logs are emitted to the BPF trace pipe, accessible with the command `tc exec bpf debug`.
	// [Default: Off].
	BPFLogLevel string `json:"bpfLogLevel,omitempty" validate:"omitempty,bpfLogLevel"`
	// BPFDataIfacePattern is a regular expression that controls which interfaces Felix should attach BPF programs to
	// in order to catch traffic to/from the network.  This needs to match the interfaces that Calico workload traffic
	// flows over as well as any interfaces that handle incoming traffic to nodeports and services from outside the
	// cluster.  It should not match the workload interfaces (usually named cali...).
	BPFDataIfacePattern string `json:"bpfDataIfacePattern,omitempty" validate:"omitempty,regexp"`
	// BPFConnectTimeLoadBalancingEnabled when in BPF mode, controls whether Felix installs the connection-time load
	// balancer.  The connect-time load balancer is required for the host to be able to reach Kubernetes services
	// and it improves the performance of pod-to-service connections.  The only reason to disable it is for debugging
	// purposes.  [Default: true]
	BPFConnectTimeLoadBalancingEnabled *bool `json:"bpfConnectTimeLoadBalancingEnabled,omitempty" validate:"omitempty"`
	// BPFExternalServiceMode in BPF mode, controls how connections from outside the cluster to services (node ports
	// and cluster IPs) are forwarded to remote workloads.  If set to "Tunnel" then both request and response traffic
	// is tunneled to the remote node.  If set to "DSR", the request traffic is tunneled but the response traffic
	// is sent directly from the remote node.  In "DSR" mode, the remote node appears to use the IP of the ingress
	// node; this requires a permissive L2 network.  [Default: Tunnel]
	BPFExternalServiceMode string `json:"bpfExternalServiceMode,omitempty" validate:"omitempty,bpfServiceMode"`
	// BPFExtToServiceConnmark in BPF mode, control a 32bit mark that is set on connections from an
	// external client to a local service. This mark allows us to control how packets of that
	// connection are routed within the host and how is routing intepreted by RPF check. [Default: 0]
	BPFExtToServiceConnmark *int `json:"bpfExtToServiceConnmark,omitempty" validate:"omitempty,gte=0,lte=4294967295"`
	// BPFKubeProxyIptablesCleanupEnabled, if enabled in BPF mode, Felix will proactively clean up the upstream
	// Kubernetes kube-proxy's iptables chains.  Should only be enabled if kube-proxy is not running.  [Default: true]
	BPFKubeProxyIptablesCleanupEnabled *bool `json:"bpfKubeProxyIptablesCleanupEnabled,omitempty" validate:"omitempty"`
	// BPFKubeProxyMinSyncPeriod, in BPF mode, controls the minimum time between updates to the dataplane for Felix's
	// embedded kube-proxy.  Lower values give reduced set-up latency.  Higher values reduce Felix CPU usage by
	// batching up more work.  [Default: 1s]
	BPFKubeProxyMinSyncPeriod *metav1.Duration `json:"bpfKubeProxyMinSyncPeriod,omitempty" validate:"omitempty" configv1timescale:"seconds"`
	// BPFKubeProxyEndpointSlicesEnabled in BPF mode, controls whether Felix's
	// embedded kube-proxy accepts EndpointSlices or not.
	BPFKubeProxyEndpointSlicesEnabled *bool `json:"bpfKubeProxyEndpointSlicesEnabled,omitempty" validate:"omitempty"`

	SyslogReporterNetwork string `json:"syslogReporterNetwork,omitempty"`
	SyslogReporterAddress string `json:"syslogReporterAddress,omitempty"`

	// IPSecMode controls which mode IPSec is operating on.
	// Default value means IPSec is not enabled. [Default: ""]
	IPSecMode string `json:"ipsecMode,omitempty" validate:"omitempty,ipsecMode"`
	// IPSecAllowUnsecuredTraffic controls whether non-IPsec traffic is allowed in addition to IPsec traffic. Enabling this
	// negates the anti-spoofing protections of IPsec but it is useful when migrating to/from IPsec. [Default: false]
	IPSecAllowUnsecuredTraffic *bool `json:"ipsecAllowUnsecuredTraffic,omitempty"`
	// IPSecIKEAlgorithm sets IPSec IKE algorithm. Default is NIST suite B recommendation. [Default: aes128gcm16-prfsha256-ecp256]
	IPSecIKEAlgorithm string `json:"ipsecIKEAlgorithm,omitempty"`
	// IPSecESAlgorithm sets IPSec ESP algorithm. Default is NIST suite B recommendation. [Default: aes128gcm16-ecp256]
	IPSecESPAlgorithm string `json:"ipsecESPAlgorithm,omitempty"`
	// IPSecLogLevel controls log level for IPSec components. Set to None for no logging.
	// A generic log level terminology is used [None, Notice, Info, Debug, Verbose].
	// [Default: Info]
	IPSecLogLevel string `json:"ipsecLogLevel,omitempty" validate:"omitempty,ipsecLogLevel"`
	// IPSecPolicyRefreshInterval is the interval at which Felix will check the kernel's IPsec policy tables and
	// repair any inconsistencies. [Default: 600s]
	IPSecPolicyRefreshInterval *metav1.Duration `json:"ipsecPolicyRefreshInterval,omitempty" configv1timescale:"seconds"`

	// FlowLogsFlushInterval configures the interval at which Felix exports flow logs.
	FlowLogsFlushInterval *metav1.Duration `json:"flowLogsFlushInterval,omitempty" configv1timescale:"seconds"`
	// FlowLogsEnableHostEndpoint enables Flow logs reporting for HostEndpoints.
	FlowLogsEnableHostEndpoint *bool `json:"flowLogsEnableHostEndpoint,omitempty"`
	// FlowLogsEnableNetworkSets enables Flow logs reporting for GlobalNetworkSets.
	FlowLogsEnableNetworkSets *bool `json:"flowLogsEnableNetworkSets,omitempty"`
	// FlowLogsMaxOriginalIPsIncluded specifies the number of unique IP addresses (if relevant) that should be included in Flow logs.
	FlowLogsMaxOriginalIPsIncluded *int `json:"flowLogsMaxOriginalIPsIncluded,omitempty"`
	// FlowLogsCollectProcessInfo, if enabled Felix will load the kprobe BPF programs to collect process info. [Default: false]
	FlowLogsCollectProcessInfo *bool `json:"flowLogsCollectProcessInfo,omitempty" validate:"omitempty"`
	// FlowLogsCollectTcpStats enables flow logs reporting TCP socket stats
	FlowLogsCollectTcpStats *bool `json:"flowLogsCollectTcpStats,omitempty" validate:"omitempty"`
	// When FlowLogsCollectProcessPath and FlowLogsCollectProcessInfo are
	// both enabled, each flow log will include information about the process
	// that is sending or receiving the packets in that flow: the
	// `process_name` field will contain the full path of the process
	// executable, and the `process_args` field will have the arguments with
	// which the executable was invoked.  Process information will not be
	// reported for connections which use raw sockets.
	FlowLogsCollectProcessPath *bool `json:"flowLogsCollectProcessPath,omitempty" validate:"omitempty"`

	// FlowLogsFileEnabled when set to true, enables logging flow logs to a file. If false no flow logging to file will occur.
	FlowLogsFileEnabled *bool `json:"flowLogsFileEnabled,omitempty"`
	// FlowLogsFileMaxFiles sets the number of log files to keep.
	FlowLogsFileMaxFiles *int `json:"flowLogsFileMaxFiles,omitempty"`
	// FlowLogsFileMaxFileSizeMB sets the max size in MB of flow logs files before rotation.
	FlowLogsFileMaxFileSizeMB *int `json:"flowLogsFileMaxFileSizeMB,omitempty"`
	// FlowLogsFileDirectory sets the directory where flow logs files are stored.
	FlowLogsFileDirectory *string `json:"flowLogsFileDirectory,omitempty"`
	// FlowLogsFileIncludeLabels is used to configure if endpoint labels are included in a Flow log entry written to file.
	FlowLogsFileIncludeLabels *bool `json:"flowLogsFileIncludeLabels,omitempty"`
	// FlowLogsFileIncludePolicies is used to configure if policy information are included in a Flow log entry written to file.
	FlowLogsFileIncludePolicies *bool `json:"flowLogsFileIncludePolicies,omitempty"`
	// FlowLogsFileIncludeService is used to configure if the destination service is included in a Flow log entry written to file.
	// The service information can only be included if the flow was explicitly determined to be directed at the service (e.g.
	// when the pre-DNAT destination corresponds to the service ClusterIP and port).
	FlowLogsFileIncludeService *bool `json:"flowLogsFileIncludeService,omitempty"`
	// FlowLogsFileAggregationKindForAllowed is used to choose the type of aggregation for flow log entries created for
	// allowed connections. [Default: 2 - pod prefix name based aggregation].
	// Accepted values are 0, 1 and 2.
	// 0 - No aggregation
	// 1 - Source port based aggregation
	// 2 - Pod prefix name based aggreagation.
	FlowLogsFileAggregationKindForAllowed *int `json:"flowLogsFileAggregationKindForAllowed,omitempty" validate:"omitempty,flowLogAggregationKind"`
	// FlowLogsFileAggregationKindForDenied is used to choose the type of aggregation for flow log entries created for
	// denied connections. [Default: 1 - source port based aggregation].
	// Accepted values are 0, 1 and 2.
	// 0 - No aggregation
	// 1 - Source port based aggregation
	// 2 - Pod prefix name based aggregation.
	// 3 - No destination ports based aggregation
	FlowLogsFileAggregationKindForDenied *int `json:"flowLogsFileAggregationKindForDenied,omitempty" validate:"omitempty,flowLogAggregationKind"`
	// FlowLogsFileEnabledForAllowed is used to enable/disable flow logs entries created for allowed connections. Default is true.
	// This parameter only takes effect when FlowLogsFileReporterEnabled is set to true.
	FlowLogsFileEnabledForAllowed *bool `json:"flowLogsFileEnabledForAllowed,omitempty"`
	// FlowLogsFileEnabledForDenied is used to enable/disable flow logs entries created for denied flows. Default is true.
	// This parameter only takes effect when FlowLogsFileReporterEnabled is set to true.
	FlowLogsFileEnabledForDenied *bool `json:"flowLogsFileEnabledForDenied,omitempty"`
	// FlowLogsDynamicAggregationEnabled is used to enable/disable dynamically changing aggregation levels. Default is true.
	FlowLogsDynamicAggregationEnabled *bool `json:"flowLogsDynamicAggregationEnabled,omitempty"`
	// FlowLogsPositionFilePath is used specify the position of the external pipeline that reads flow logs. Default is /var/log/calico/flows.log.pos.
	// This parameter only takes effect when FlowLogsDynamicAggregationEnabled is set to true.
	FlowLogsPositionFilePath *string `json:"flowLogsPositionFilePath,omitempty"`
	// FlowLogsAggregationThresholdBytes is used specify how far behind the external pipeline that reads flow logs can be. Default is 8192 bytes.
	// This parameter only takes effect when FlowLogsDynamicAggregationEnabled is set to true.
	FlowLogsAggregationThresholdBytes *int `json:"flowLogsAggregationThresholdBytes,omitempty"`
	// FlowLogsFilePerFlowProcessLimit, is used to specify the maximum number of flow log entries with distinct process information
	// beyond which process information will be aggregated. [Default: 2]
	FlowLogsFilePerFlowProcessLimit *int `json:"flowLogsFilePerFlowProcessLimit,omitempty" validate:"omitempty"`

	// WindowsFlowLogsFileDirectory sets the directory where flow logs files are stored on Windows nodes. [Default: "c:\\TigeraCalico\\flowlogs"].
	WindowsFlowLogsFileDirectory string `json:"windowsFlowLogsFileDirectory,omitempty"`
	// WindowsFlowLogsPositionFilePath is used to specify the position of the external pipeline that reads flow logs on Windows nodes.
	// [Default: "c:\\TigeraCalico\\flowlogs\\flows.log.pos"].
	// This parameter only takes effect when FlowLogsDynamicAggregationEnabled is set to true.
	WindowsFlowLogsPositionFilePath string `json:"windowsFlowLogsPositionFilePath,omitempty"`
	// WindowsStatsDumpFilePath is used to specify the path of the stats dump file on Windows nodes. [Default: "c:\\TigeraCalico\\stats\\dump"]
	WindowsStatsDumpFilePath string `json:"windowsStatsDumpFilePath,omitempty"`
	// The name of the file that Felix uses to preserve learnt DNS information when restarting. [Default:
	// "c:\\TigeraCalico\\felix-dns-cache.txt"].
	WindowsDNSCacheFile string `json:"windowsDnsCacheFile,omitempty"`
	// Extra time to keep IPs and alias names that are learnt from DNS, in addition to each name
	// or IP's advertised TTL. The default value is 120s which is same as the default value of
	// ServicePointManager.DnsRefreshTimeout on .net framework. [Default: 120s].
	WindowsDNSExtraTTL *metav1.Duration `json:"windowsDnsExtraTTL,omitempty" configv1timescale:"seconds"`

	// The DNS servers that Felix should trust. Each entry here must be `<ip>[:<port>]` - indicating an
	// explicit DNS server IP - or `k8s-service:[<namespace>/]<name>[:port]` - indicating a Kubernetes DNS
	// service. `<port>` defaults to the first service port, or 53 for an IP, and `<namespace>` to
	// `kube-system`. An IPv6 address with a port must use the square brackets convention, for example
	// `[fd00:83a6::12]:5353`.Note that Felix (calico-node) will need RBAC permission to read the details of
	// each service specified by a `k8s-service:...` form. [Default: "k8s-service:kube-dns"].
	DNSTrustedServers *[]string `json:"dnsTrustedServers,omitempty" validate:"omitempty,dive,ipOrK8sService"`
	// The name of the file that Felix uses to preserve learnt DNS information when restarting. [Default:
	// "/var/run/calico/felix-dns-cache.txt"].
	DNSCacheFile string `json:"dnsCacheFile,omitempty"`
	// The periodic interval at which Felix saves learnt DNS information to the cache file. [Default:
	// 60s].
	DNSCacheSaveInterval *metav1.Duration `json:"dnsCacheSaveInterval,omitempty" configv1timescale:"seconds"`
	// An arbitrary number that can be changed, at runtime, to tell Felix to discard all its
	// learnt DNS information. [Default: 0].
	DNSCacheEpoch *int `json:"dnsCacheEpoch,omitempty"`
	// Extra time to keep IPs and alias names that are learnt from DNS, in addition to each name
	// or IP's advertised TTL. [Default: 0s].
	DNSExtraTTL *metav1.Duration `json:"dnsExtraTTL,omitempty" configv1timescale:"seconds"`

	// DNSLogsFlushInterval configures the interval at which Felix exports DNS logs.
	// [Default: 300s]
	DNSLogsFlushInterval *metav1.Duration `json:"dnsLogsFlushInterval,omitempty" configv1timescale:"seconds"`
	// DNSLogsFileEnabled controls logging DNS logs to a file. If false no DNS logging to file will occur.
	// [Default: false]
	DNSLogsFileEnabled *bool `json:"dnsLogsFileEnabled,omitempty"`
	// DNSLogsFileMaxFiles sets the number of DNS log files to keep.
	// [Default: 5]
	DNSLogsFileMaxFiles *int `json:"dnsLogsFileMaxFiles,omitempty"`
	// DNSLogsFileMaxFileSizeMB sets the max size in MB of DNS log files before rotation.
	// [Default: 100]
	DNSLogsFileMaxFileSizeMB *int `json:"dnsLogsFileMaxFileSizeMB,omitempty"`
	// DNSLogsFileDirectory sets the directory where DNS log files are stored.
	// [Default: /var/log/calico/dnslogs]
	DNSLogsFileDirectory *string `json:"dnsLogsFileDirectory,omitempty"`
	// DNSLogsFileIncludeLabels is used to configure if endpoint labels are included in a DNS log entry written to file.
	// [Default: true]
	DNSLogsFileIncludeLabels *bool `json:"dnsLogsFileIncludeLabels,omitempty"`
	// DNSLogsFileAggregationKind is used to choose the type of aggregation for DNS log entries.
	// [Default: 1 - client name prefix aggregation].
	// Accepted values are 0 and 1.
	// 0 - No aggregation
	// 1 - Aggregate over clients with the same name prefix
	DNSLogsFileAggregationKind *int `json:"dnsLogsFileAggregationKind,omitempty" validate:"omitempty,dnsAggregationKind"`
	// Limit on the number of DNS logs that can be emitted within each flush interval.  When
	// this limit has been reached, Felix counts the number of unloggable DNS responses within
	// the flush interval, and emits a WARNING log with that count at the same time as it
	// flushes the buffered DNS logs.  [Default: 0, meaning no limit]
	DNSLogsFilePerNodeLimit *int `json:"dnsLogsFilePerNodeLimit,omitempty"`
	// DNSLogsLatency indicates to include measurements of DNS request/response latency in each DNS log.
	// [Default: true]
	DNSLogsLatency *bool `json:"dnsLogsLatency,omitempty"`

	// L7LogsFlushInterval configures the interval at which Felix exports L7 logs.
	// [Default: 300s]
	L7LogsFlushInterval *metav1.Duration `json:"l7LogsFlushInterval,omitempty" configv1timescale:"seconds"`
	// L7LogsFileEnabled controls logging L7 logs to a file. If false no L7 logging to file will occur.
	// [Default: true]
	L7LogsFileEnabled *bool `json:"l7LogsFileEnabled,omitempty"`
	// L7LogsFileMaxFiles sets the number of L7 log files to keep.
	// [Default: 5]
	L7LogsFileMaxFiles *int `json:"l7LogsFileMaxFiles,omitempty"`
	// L7LogsFileMaxFileSizeMB sets the max size in MB of L7 log files before rotation.
	// [Default: 100]
	L7LogsFileMaxFileSizeMB *int `json:"l7LogsFileMaxFileSizeMB,omitempty"`
	// L7LogsFileDirectory sets the directory where L7 log files are stored.
	// [Default: /var/log/calico/l7logs]
	L7LogsFileDirectory *string `json:"l7LogsFileDirectory,omitempty"`
	// L7LogsFileAggregationHTTPHeaderInfo is used to choose the type of aggregation for HTTP header data on L7 log entries.
	// [Default: ExcludeL7HTTPHeaderInfo - http header info removal].
	// Accepted values are IncludeL7HTTPHeaderInfo and ExcludeL7HTTPHeaderInfo.
	// IncludeL7HTTPHeaderInfo - Include HTTP header data in the logs.
	// ExcludeL7HTTPHeaderInfo - Aggregate over all other fields ignoring the user agent and log type.
	L7LogsFileAggregationHTTPHeaderInfo *string `json:"l7LogsFileAggregationHTTPHeaderInfo,omitempty" validate:"omitempty,l7HTTPHeaderAggregation"`
	// L7LogsFileAggregationHTTPMethod is used to choose the type of aggregation for the HTTP request method on L7 log entries.
	// [Default: IncludeL7HTTPMethod - include the HTTP method].
	// Accepted values are IncludeL7HTTPMethod and ExcludeL7HTTPMethod.
	// IncludeL7HTTPMethod - Include HTTP method in the logs.
	// ExcludeL7HTTPMethod - Aggregate over all other fields ignoring the HTTP method.
	L7LogsFileAggregationHTTPMethod *string `json:"l7LogsFileAggregationHTTPMethod,omitempty" validate:"omitempty,l7HTTPMethodAggregation"`
	// L7LogsFileAggregationServiceInfo is used to choose the type of aggregation for the service data on L7 log entries.
	// [Default: IncludeL7ServiceInfo - include service data].
	// Accepted values are IncludeL7ServiceInfo and ExcludeL7ServiceInfo.
	// IncludeL7ServiceInfo - Include service data in the logs.
	// ExcludeL7ServiceInfo - Aggregate over all other fields ignoring the service name, namespace, and port.
	L7LogsFileAggregationServiceInfo *string `json:"l7LogsFileAggregationServiceInfo,omitempty" validate:"omitempty,l7ServiceAggregation"`
	// L7LogsFileAggregationDestinationInfo is used to choose the type of aggregation for the destination metadata on L7 log entries.
	// [Default: IncludeL7DestinationInfo - include destination metadata].
	// Accepted values are IncludeL7DestinationInfo and ExcludeL7DestinationInfo.
	// IncludeL7DestinationInfo - Include destination metadata in the logs.
	// ExcludeL7DestinationInfo - Aggregate over all other fields ignoring the destination aggregated name, namespace, and type.
	L7LogsFileAggregationDestinationInfo *string `json:"l7LogsFileAggregationDestinationInfo,omitempty" validate:"omitempty,l7DestinationAggregation"`
	// L7LogsFileAggregationExcludeSourceInfo is used to choose the type of aggregation for the source metadata on L7 log entries.
	// [Default: IncludeL7SourceInfoNoPort - include all source metadata except for the source port].
	// Accepted values are IncludeL7SourceInfo, IncludeL7SourceInfoNoPort, and ExcludeL7SourceInfo.
	// IncludeL7SourceInfo - Include source metadata in the logs.
	// IncludeL7SourceInfoNoPort - Include source metadata in the logs excluding the source port.
	// ExcludeL7SourceInfo - Aggregate over all other fields ignoring the source aggregated name, namespace, and type.
	L7LogsFileAggregationSourceInfo *string `json:"l7LogsFileAggregationSourceInfo,omitempty" validate:"omitempty,l7SourceAggregation"`
	// L7LogsFileAggregationResponseCode is used to choose the type of aggregation for the response code on L7 log entries.
	// [Default: IncludeL7ResponseCode - include the response code].
	// Accepted values are IncludeL7ResponseCode and ExcludeL7ResponseCode.
	// IncludeL7ResponseCode - Include the response code in the logs.
	// ExcludeL7ResponseCode - Aggregate over all other fields ignoring the response code.
	L7LogsFileAggregationResponseCode *string `json:"l7LogsFileAggregationResponseCode,omitempty" validate:"omitempty,l7ResponseCodeAggregation"`
	// L7LogsFileAggregationTrimURL is used to choose the type of aggregation for the url on L7 log entries.
	// [Default: IncludeL7FullURL - include the full URL up to however many path components are allowed by L7LogsFileAggregationNumURLPath].
	// Accepted values are 0 and 1.
	// IncludeL7FullURL - Include the full URL up to however many path components are allowed by L7LogsFileAggregationNumURLPath.
	// TrimURLQuery - Aggregate over all other fields ignoring the query parameters on the URL.
	// TrimURLQueryAndPath - Aggregate over all other fields and the base URL only.
	// ExcludeL7URL - Aggregate over all other fields ignoring the URL entirely.
	L7LogsFileAggregationTrimURL *string `json:"l7LogsFileAggregationTrimURL,omitempty" validate:"omitempty,l7URLAggregation"`
	// L7LogsFileAggregationNumURLPath is used to choose the number of components in the url path to display.
	// This allows for the url to be truncated in case parts of the path provide no value. Setting this value
	// to negative will allow all parts of the path to be displayed.
	// [Default: 5].
	L7LogsFileAggregationNumURLPath *int `json:"l7LogsFileAggregationNumURLPath,omitempty"`
	// Limit on the length of the URL collected in L7 logs. When a URL length reaches this limit
	// it is sliced off, and the sliced URL is sent to log storage. [Default: 250]
	L7LogsFileAggregationURLCharLimit *int `json:"l7LogsFileAggregationURLCharLimit,omitempty"`
	// Limit on the number of L7 logs that can be emitted within each flush interval.  When
	// this limit has been reached, Felix counts the number of unloggable L7 responses within
	// the flush interval, and emits a WARNING log with that count at the same time as it
	// flushes the buffered L7 logs. A value of 0 means no limit. [Default: 1500]
	L7LogsFilePerNodeLimit *int `json:"l7LogsFilePerNodeLimit,omitempty"`

	// WindowsNetworkName specifies which Windows HNS networks Felix should operate on.  The default is to match
	// networks that start with "calico".  Supports regular expression syntax.
	WindowsNetworkName *string `json:"windowsNetworkName,omitempty"`

	// RouteSource configures where Felix gets its routing information.
	// - WorkloadIPs: use workload endpoints to construct routes.
	// - CalicoIPAM: the default - use IPAM data to construct routes.
	RouteSource string `json:"routeSource,omitempty" validate:"omitempty,routeSource"`

	// Calico programs additional Linux route tables for various purposes.  RouteTableRange
	// specifies the indices of the route tables that Calico should use.
	RouteTableRange *RouteTableRange `json:"routeTableRange,omitempty" validate:"omitempty"`

	// EgressIPSupport defines three different support modes for egress IP function. [Default: Disabled]
	// - Disabled:                    Egress IP function is disabled.
	// - EnabledPerNamespace:         Egress IP function is enabled and can be configured on a per-namespace basis;
	//                                per-pod egress annotations are ignored.
	// - EnabledPerNamespaceOrPerPod: Egress IP function is enabled and can be configured per-namespace or per-pod,
	//                                with per-pod egress annotations overriding namespace annotations.
	EgressIPSupport string `json:"egressIPSupport,omitempty" validate:"omitempty,oneof=Disabled EnabledPerNamespace EnabledPerNamespaceOrPerPod"`
	// EgressIPVXLANPort is the port number of vxlan tunnel device for egress traffic. [Default: 4790]
	EgressIPVXLANPort *int `json:"egressIPVXLANPort,omitempty"`
	// EgressIPVXLANVNI is the VNI ID of vxlan tunnel device for egress traffic. [Default: 4097]
	EgressIPVXLANVNI *int `json:"egressIPVXLANVNI,omitempty"`
	// EgressIPRoutingRulePriority controls the priority value to use for the egress IP routing rule. [Default: 100]
	EgressIPRoutingRulePriority *int `json:"egressIPRoutingRulePriority,omitempty" validate:"omitempty,gt=0,lt=32766"`

	// WireguardEnabled controls whether Wireguard is enabled. [Default: false]
	WireguardEnabled *bool `json:"wireguardEnabled,omitempty"`
	// WireguardListeningPort controls the listening port used by Wireguard. [Default: 51820]
	WireguardListeningPort *int `json:"wireguardListeningPort,omitempty" validate:"omitempty,gt=0,lte=65535"`
	// WireguardRoutingRulePriority controls the priority value to use for the Wireguard routing rule. [Default: 99]
	WireguardRoutingRulePriority *int `json:"wireguardRoutingRulePriority,omitempty" validate:"omitempty,gt=0,lt=32766"`
	// WireguardInterfaceName specifies the name to use for the Wireguard interface. [Default: wg.calico]
	WireguardInterfaceName string `json:"wireguardInterfaceName,omitempty" validate:"omitempty,interface"`
	// WireguardMTU controls the MTU on the Wireguard interface. See Configuring MTU [Default: 1420]
	WireguardMTU *int `json:"wireguardMTU,omitempty"`
	// WireguardHostEncryptionEnabled controls whether Wireguard host-to-host encryption is enabled. [Default: false]
	WireguardHostEncryptionEnabled *bool `json:"wireguardHostEncryptionEnabled,omitempty"`

	// +kubebuilder:validation:MinLength=1
	// CaptureDir controls directory to store file capture. [Default: /var/log/calico/pcap]
	CaptureDir *string `json:"captureDir,omitempty" validate:"omitempty,gt=0"`

	// +kubebuilder:validation:Minimum=1
	// CaptureMaxSizeBytes controls the max size of a file capture. [Default: 10000000]
	CaptureMaxSizeBytes *int `json:"captureMaxSizeBytes,omitempty" validate:"omitempty,gt=0"`

	// +kubebuilder:validation:Minimum=1
	// CaptureRotationSeconds controls the time rotation of a packet capture. [Default: 3600]
	CaptureRotationSeconds *int `json:"captureRotationSeconds,omitempty" validate:"omitempty,gt=0"`

	// +kubebuilder:validation:Minimum=1
	// CaptureMaxFiles controls number of rotated capture file to keep. [Default: 2]
	CaptureMaxFiles *int `json:"captureMaxFiles,omitempty" validate:"omitempty,gt=0"`

	// Set source-destination-check on AWS EC2 instances. Accepted value must be one of "DoNothing", "Enabled" or "Disabled".
	// [Default: DoNothing]
	AWSSrcDstCheck *AWSSrcDstCheckOption `json:"awsSrcDstCheck,omitempty" validate:"omitempty,oneof=DoNothing Enable Disable"`

	// When service IP advertisement is enabled, prevent routing loops to service IPs that are
	// not in use, by dropping or rejecting packets that do not get DNAT'd by kube-proxy.
	// Unless set to "Disabled", in which case such routing loops continue to be allowed.
	// [Default: Drop]
	ServiceLoopPrevention string `json:"serviceLoopPrevention,omitempty" validate:"omitempty,oneof=Drop Reject Disabled"`

	// MTUIfacePattern is a regular expression that controls which interfaces Felix should scan in order
	// to calculate the host's MTU.
	// This should not match workload interfaces (usually named cali...).
	// +optional
	MTUIfacePattern string `json:"mtuIfacePattern,omitempty" validate:"omitempty,regexp"`

	// TPROXYMode sets whether traffic is directed through a transparent proxy
	// for further processing or not and how is the proxying done.
	// [Default: Disabled]
	TPROXYMode string `json:"tproxyMode,omitempty" validate:"omitempty,oneof=Disabled Enabled EnabledAllServices"`
	// TPROXYPort sets to which port proxied traffic should be redirected.
	// [Default: 16001]
	TPROXYPort *int `json:"tproxyPort,omitempty" validate:"omitempty,gt=0,lte=65535"`
}

type RouteTableRange struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

// ProtoPort is combination of protocol, port, and CIDR. Protocol and port must be specified.
type ProtoPort struct {
	Protocol string `json:"protocol"`
	Port     uint16 `json:"port"`
	// +optional
	Net string `json:"net"`
}

// New FelixConfiguration creates a new (zeroed) FelixConfiguration struct with the TypeMetadata
// initialized to the current version.
func NewFelixConfiguration() *FelixConfiguration {
	return &FelixConfiguration{
		TypeMeta: metav1.TypeMeta{
			Kind:       KindFelixConfiguration,
			APIVersion: GroupVersionCurrent,
		},
	}
}
