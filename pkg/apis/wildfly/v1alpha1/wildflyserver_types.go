package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// WildFlyServerSpec defines the desired state of WildFlyServer
// +k8s:openapi-gen=true
type WildFlyServerSpec struct {
	// ApplicationImage is the name of the application image to be deployed
	ApplicationImage string `json:"applicationImage"`
	// BootableJar specifies whether the application image is using S2I Builder/Runtime images or Bootable Jar.
	// If omitted, it defaults to false (application image is expected to use S2I Builder/Runtime images)
	BootableJar bool `json:"bootableJar,omitempty"`
	// Replicas is the desired number of replicas for the application
	// +kubebuilder:validation:Minimum=0
	Replicas int32 `json:"replicas"`
	// SessionAffinity defines if connections from the same client ip are passed to the same WildFlyServer instance/pod each time (false if omitted)
	SessionAffinity bool `json:"sessionAffinity,omitempty"`
	// DisableHTTPRoute disables the creation a route to the HTTP port of the application service (false if omitted)
	DisableHTTPRoute bool `json:"disableHTTPRoute,omitempty"`
	// DeactivateTransactionRecovery disables the process of recoverying transactions (false if omitted)
	DeactivateTransactionRecovery bool                     `json:deactivateTransactionRecovery,omitempty`
	StandaloneConfigMap           *StandaloneConfigMapSpec `json:"standaloneConfigMap,omitempty"`
	// StorageSpec defines specific storage required for the server own data directory. If omitted, an EmptyDir is used (that will not
	// persist data across pod restart).
	Storage            *StorageSpec `json:"storage,omitempty"`
	ServiceAccountName string       `json:"serviceAccountName,omitempty"`
	// EnvFrom contains environment variables from a source such as a ConfigMap or a Secret
	// +kubebuilder:validation:MinItems=1
	// +listType=atomic
	EnvFrom []corev1.EnvFromSource `json:"envFrom,omitempty,list_type=corev1.EnvFromSource"`
	// Env contains environment variables for the containers running the WildFlyServer application
	// +kubebuilder:validation:MinItems=1
	// +listType=atomic
	Env []corev1.EnvVar `json:"env,omitempty"`
	// Secrets is a list of Secrets in the same namespace as the WildFlyServer
	// object, which shall be mounted into the WildFlyServer Pods.
	// The Secrets are mounted into /etc/secrets/<secret-name>.
	// +kubebuilder:validation:MinItems=1
	// +listType=set
	Secrets []string `json:"secrets,omitempty"`
	// ConfigMaps is a list of ConfigMaps in the same namespace as the WildFlyServer
	// object, which shall be mounted into the WildFlyServer Pods.
	// The ConfigMaps are mounted into /etc/configmaps/<configmap-name>.
	// +kubebuilder:validation:MinItems=1
	// +listType=set
	ConfigMaps []string `json:"configMaps,omitempty"`
	// ResourcesSpec defines the resources used by the WildFlyServer, ie CPU and memory, use limits and requests.
	// More info: https://pkg.go.dev/k8s.io/api@v0.18.14/core/v1#ResourceRequirements
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`
}

// StandaloneConfigMapSpec defines the desired configMap configuration to obtain the standalone configuration for WildFlyServer
// +k8s:openapi-gen=true
type StandaloneConfigMapSpec struct {
	Name string `json:"name"`
	// Key of the config map whose value is the standalone XML configuration file ("standalone.xml" if omitted)
	Key string `json:"key,omitempty"`
}

// StorageSpec defines the desired storage for WildFlyServer
// +k8s:openapi-gen=true
type StorageSpec struct {
	EmptyDir *corev1.EmptyDirVolumeSource `json:"emptyDir,omitempty"`
	// VolumeClaimTemplate defines the template to store WildFlyServer standalone data directory.
	// The name of the template is derived from the WildFlyServer name.
	//  The corresponding volume will be mounted in ReadWriteOnce access mode.
	// This template should be used to specify specific Resources requirements in the template spec.
	VolumeClaimTemplate corev1.PersistentVolumeClaim `json:"volumeClaimTemplate,omitempty"`
}

// WildFlyServerStatus defines the observed state of WildFlyServer
// +k8s:openapi-gen=true
type WildFlyServerStatus struct {
	// Replicas is the actual number of replicas for the application
	Replicas int32 `json:"replicas"`
	// +listType=atomic
	Pods []PodStatus `json:"pods,omitempty"`
	// +listType=set
	Hosts []string `json:"hosts,omitempty"`
	// Represents the number of pods which are in scaledown process
	// what particular pod is scaling down can be verified by PodStatus
	//
	// Read-only.
	ScalingdownPods int32 `json:"scalingdownPods"`
	// selector for pods, used by HorizontalPodAutoscaler
	Selector string `json:"selector"`
}

const (
	// (PodStatus.State) PodStateActive represents an active pod that is connected to the load balancer Service
	// and that can serve requests
	PodStateActive = "ACTIVE"
	// (PodStatus.State) PodStateScalingDownRecoveryInvestigation represents a pod that is under investigation
	// to find out if there are transactions to be recovered. A pod in this state will be updated to one of
	// the following states eventually
	PodStateScalingDownRecoveryInvestigation = "SCALING_DOWN_RECOVERY_INVESTIGATION"
	// (PodStatus.State) PodStateScalingDownRecoveryProcessing represents a pod that has transactions to be completed.
	// The Operator will wait until all transactions are processed
	PodStateScalingDownRecoveryProcessing = "SCALING_DOWN_RECOVERY_PROCESSING"
	// (PodStatus.State) PodStateScalingDownRecoveryHeuristic represents a pod that has heuristic transactions.
	// The Operator will wait until all heuristic transactions are manually solved
	PodStateScalingDownRecoveryHeuristic = "SCALING_DOWN_RECOVERY_HEURISTICS"
	// (PodStatus.State) PodStateScalingDownClean represents a pod that is ready to be scaled down
	PodStateScalingDownClean = "SCALING_DOWN_CLEAN"
)

// PodStatus defines the observed state of pods running the WildFlyServer application
// +k8s:openapi-gen=true
type PodStatus struct {
	Name  string `json:"name"`
	PodIP string `json:"podIP"`
	// Represent the state of the Pod, it is used especially during scale down.
	// +kubebuilder:validation:Enum=ACTIVE;SCALING_DOWN_RECOVERY_INVESTIGATION;SCALING_DOWN_RECOVERY_PROCESSING;SCALING_DOWN_RECOVERY_HEURISTICS;SCALING_DOWN_CLEAN
	State string `json:"state"`
	// Counts the recovery attempts when there are in-doubt txns
	RecoveryCounter int32 `json:recoveryCounter`
}

// WildFlyServer is the Schema for the wildflyservers API
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:subresource:status
// +kubebuilder:subresource:scale:specpath=.spec.replicas,statuspath=.status.replicas,selectorpath=.status.selector
// +kubebuilder:printcolumn:name="Replicas",type="integer",JSONPath=".spec.replicas"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:resource:shortName=wfly
type WildFlyServer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   WildFlyServerSpec   `json:"spec,omitempty"`
	Status WildFlyServerStatus `json:"status,omitempty"`
}

// WildFlyServerList contains a list of WildFlyServer
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type WildFlyServerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []WildFlyServer `json:"items"`
}

func init() {
	SchemeBuilder.Register(&WildFlyServer{}, &WildFlyServerList{})
}
