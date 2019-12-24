/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha4

// Cluster contains kind cluster configuration
type Cluster struct {
	TypeMeta `yaml:",inline"`

	// Nodes contains the list of nodes defined in the `kind` Cluster
	// If unset this will default to a single control-plane node
	// Note that if more than one control plane is specified, an external
	// control plane load balancer will be provisioned implicitly
	Nodes []Node `yaml:"nodes,omitempty"`

	/* Advanced fields */

	// Networking contains cluster wide network settings
	Networking Networking `yaml:"networking,omitempty"`

	// KubeadmConfigPatches are applied to the generated kubeadm config as
	// merge patches. The `kind` field must match the target object, and
	// if `apiVersion` is specified it will only be applied to matching objects.
	//
	// This should be an inline yaml blob-string
	//
	// https://tools.ietf.org/html/rfc7386
	//
	// The cluster-level patches are appied before the node-level patches.
	KubeadmConfigPatches []string `yaml:"kubeadmConfigPatches,omitempty"`

	// KubeadmConfigPatchesJSON6902 are applied to the generated kubeadm config
	// as JSON 6902 patches. The `kind` field must match the target object, and
	// if group or version are specified it will only be objects matching the
	// apiVersion: group+"/"+version
	//
	// Name and Namespace are now ignored, but the fields continue to exist for
	// backwards compatibility of parsing the config. The name of the generated
	// config was/is always fixed as as is the namespace so these fields have
	// always been a no-op.
	//
	// https://tools.ietf.org/html/rfc6902
	//
	// The cluster-level patches are appied before the node-level patches.
	KubeadmConfigPatchesJSON6902 []PatchJSON6902 `yaml:"kubeadmConfigPatchesJSON6902,omitempty"`

	// ContainerdConfigPatches are applied to every node's containerd config
	// in the order listed.
	// These should be toml stringsto be applied as merge patches
	ContainerdConfigPatches []string `yaml:"containerdConfigPatches,omitempty"`

	// ContainerdConfigPatchesJSON6902 are applied to every node's containerd config
	// in the order listed.
	// These should be YAML or JSON formatting RFC 6902 JSON patches
	ContainerdConfigPatchesJSON6902 []string `yaml:"containerdConfigPatchesJSON6902,omitempty"`
}

// TypeMeta partially copies apimachinery/pkg/apis/meta/v1.TypeMeta
// No need for a direct dependence; the fields are stable.
type TypeMeta struct {
	Kind       string `json:"kind,omitempty" yaml:"kind,omitempty"`
	APIVersion string `json:"apiVersion,omitempty" yaml:"apiVersion,omitempty"`
}

// Node contains settings for a node in the `kind` Cluster.
// A node in kind config represent a container that will be provisioned with all the components
// required for the assigned role in the Kubernetes cluster
type Node struct {
	// Role defines the role of the node in the in the Kubernetes cluster
	// created by kind
	//
	// Defaults to "control-plane"
	Role NodeRole `yaml:"role,omitempty"`

	// Image is the node image to use when creating this node
	// If unset a default image will be used, see defaults.Image
	Image string `yaml:"image,omitempty"`

	/* Advanced fields */

	// TODO: cri-like types should be inline instead
	// ExtraMounts describes additional mount points for the node container
	// These may be used to bind a hostPath
	ExtraMounts []Mount `yaml:"extraMounts,omitempty"`

	// ExtraPortMappings describes additional port mappings for the node container
	// binded to a host Port
	ExtraPortMappings []PortMapping `yaml:"extraPortMappings,omitempty"`

	ExtraDockerOptions []string `yaml:"extraDockerOptions,omitempty"`

	// KubeadmConfigPatches are applied to the generated kubeadm config as
	// merge patches. The `kind` field must match the target object, and
	// if `apiVersion` is specified it will only be applied to matching objects.
	//
	// This should be an inline yaml blob-string
	//
	// https://tools.ietf.org/html/rfc7386
	//
	// The node-level patches will be applied after the cluster-level patches
	// have been applied. (See Cluster.KubeadmConfigPatches)
	KubeadmConfigPatches []string `yaml:"kubeadmConfigPatches,omitempty"`

	// KubeadmConfigPatchesJSON6902 are applied to the generated kubeadm config
	// as JSON 6902 patches. The `kind` field must match the target object, and
	// if group or version are specified it will only be objects matching the
	// apiVersion: group+"/"+version
	//
	// Name and Namespace are now ignored, but the fields continue to exist for
	// backwards compatibility of parsing the config. The name of the generated
	// config was/is always fixed as as is the namespace so these fields have
	// always been a no-op.
	//
	// https://tools.ietf.org/html/rfc6902
	//
	// The node-level patches will be applied after the cluster-level patches
	// have been applied. (See Cluster.KubeadmConfigPatchesJSON6902)
	KubeadmConfigPatchesJSON6902 []PatchJSON6902 `yaml:"kubeadmConfigPatchesJSON6902,omitempty"`
}

// NodeRole defines possible role for nodes in a Kubernetes cluster managed by `kind`
type NodeRole string

const (
	// ControlPlaneRole identifies a node that hosts a Kubernetes control-plane.
	// NOTE: in single node clusters, control-plane nodes act also as a worker
	// nodes, in which case the taint will be removed. see:
	// https://kubernetes.io/docs/setup/independent/create-cluster-kubeadm/#control-plane-node-isolation
	ControlPlaneRole NodeRole = "control-plane"
	// WorkerRole identifies a node that hosts a Kubernetes worker
	WorkerRole NodeRole = "worker"
)

// Networking contains cluster wide network settings
type Networking struct {
	// IPFamily is the network cluster model, currently it can be ipv4 or ipv6
	IPFamily ClusterIPFamily `yaml:"ipFamily,omitempty"`
	// APIServerPort is the listen port on the host for the Kubernetes API Server
	// Defaults to a random port on the host
	APIServerPort int32 `yaml:"apiServerPort,omitempty"`
	// APIServerAddress is the listen address on the host for the Kubernetes
	// API Server. This should be an IP address.
	//
	// Defaults to 127.0.0.1
	APIServerAddress string `yaml:"apiServerAddress,omitempty"`
	// PodSubnet is the CIDR used for pod IPs
	// kind will select a default if unspecified
	PodSubnet string `yaml:"podSubnet,omitempty"`
	// ServiceSubnet is the CIDR used for services VIPs
	// kind will select a default if unspecified for IPv6
	ServiceSubnet string `yaml:"serviceSubnet,omitempty"`
	// If DisableDefaultCNI is true, kind will not install the default CNI setup.
	// Instead the user should install their own CNI after creating the cluster.
	DisableDefaultCNI bool `yaml:"disableDefaultCNI,omitempty"`
}

// ClusterIPFamily defines cluster network IP family
type ClusterIPFamily string

const (
	// IPv4Family sets ClusterIPFamily to ipv4
	IPv4Family ClusterIPFamily = "ipv4"
	// IPv6Family sets ClusterIPFamily to ipv6
	IPv6Family ClusterIPFamily = "ipv6"
)

// PatchJSON6902 represents an inline kustomize json 6902 patch
// https://tools.ietf.org/html/rfc6902
type PatchJSON6902 struct {
	// these fields specify the patch target resource
	Group   string `yaml:"group"`
	Version string `yaml:"version"`
	Kind    string `yaml:"kind"`
	// Patch should contain the contents of the json patch as a string
	Patch string `yaml:"patch"`
}

/*
These types are from
https://github.com/kubernetes/kubernetes/blob/063e7ff358fdc8b0916e6f39beedc0d025734cb1/pkg/kubelet/apis/cri/runtime/v1alpha2/api.pb.go#L183
*/

// Mount specifies a host volume to mount into a container.
// This is a close copy of the upstream cri Mount type
// see: k8s.io/kubernetes/pkg/kubelet/apis/cri/runtime/v1alpha2
// It additionally serializes the "propagation" field with the string enum
// names on disk as opposed to the int32 values, and the serlialzed field names
// have been made closer to core/v1 VolumeMount field names
// In yaml this looks like:
//  containerPath: /foo
//  hostPath: /bar
//  readOnly: true
//  selinuxRelabel: false
//  propagation: None
// Propagation may be one of: None, HostToContainer, Bidirectional
type Mount struct {
	// Path of the mount within the container.
	ContainerPath string `yaml:"containerPath,omitempty"`
	// Path of the mount on the host. If the hostPath doesn't exist, then runtimes
	// should report error. If the hostpath is a symbolic link, runtimes should
	// follow the symlink and mount the real destination to container.
	HostPath string `yaml:"hostPath,omitempty"`
	// If set, the mount is read-only.
	Readonly bool `yaml:"readOnly,omitempty"`
	// If set, the mount needs SELinux relabeling.
	SelinuxRelabel bool `yaml:"selinuxRelabel,omitempty"`
	// Requested propagation mode.
	Propagation MountPropagation `yaml:"propagation,omitempty"`
}

// PortMapping specifies a host port mapped into a container port.
// In yaml this looks like:
//  containerPort: 80
//  hostPort: 8000
//  listenAddress: 127.0.0.1
//  protocol: TCP
type PortMapping struct {
	// Port within the container.
	ContainerPort int32 `yaml:"containerPort,omitempty"`
	// Port on the host.
	HostPort int32 `yaml:"hostPort,omitempty"`
	// TODO: add protocol (tcp/udp) and port-ranges
	ListenAddress string `yaml:"listenAddress,omitempty"`
	// Protocol (TCP/UDP)
	Protocol PortMappingProtocol `yaml:"protocol,omitempty"`
}

// MountPropagation represents an "enum" for mount propagation options,
// see also Mount.
type MountPropagation string

const (
	// MountPropagationNone specifies that no mount propagation
	// ("private" in Linux terminology).
	MountPropagationNone MountPropagation = "None"
	// MountPropagationHostToContainer specifies that mounts get propagated
	// from the host to the container ("rslave" in Linux).
	MountPropagationHostToContainer MountPropagation = "HostToContainer"
	// MountPropagationBidirectional specifies that mounts get propagated from
	// the host to the container and from the container to the host
	// ("rshared" in Linux).
	MountPropagationBidirectional MountPropagation = "Bidirectional"
)

// PortMappingProtocol represents an "enum" for port mapping protocol options,
// see also PortMapping.
type PortMappingProtocol string

const (
	// PortMappingProtocolTCP specifies TCP protocol
	PortMappingProtocolTCP PortMappingProtocol = "TCP"
	// PortMappingProtocolUDP specifies UDP protocol
	PortMappingProtocolUDP PortMappingProtocol = "UDP"
	// PortMappingProtocolSCTP specifies SCTP protocol
	PortMappingProtocolSCTP PortMappingProtocol = "SCTP"
)
