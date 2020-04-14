package v1

import (
	"fmt"
	"sort"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	OperatorGroupAnnotationKey             = "olm.operatorGroup"
	OperatorGroupNamespaceAnnotationKey    = "olm.operatorNamespace"
	OperatorGroupTargetsAnnotationKey      = "olm.targetNamespaces"
	OperatorGroupProvidedAPIsAnnotationKey = "olm.providedAPIs"

	OperatorGroupKind = "OperatorGroup"

	operatorGroupLabelPrefix   = "olm.operatorgroup/"
	operatorGroupLabelTemplate = operatorGroupLabelPrefix + "%s.%s"
)

// OperatorGroupSpec is the spec for an OperatorGroup resource.
type OperatorGroupSpec struct {
	// Selector selects the OperatorGroup's target namespaces.
	// +optional
	Selector *metav1.LabelSelector `json:"selector,omitempty"`

	// TargetNamespaces is an explicit set of namespaces to target.
	// If it is set, Selector is ignored.
	// +optional
	// +listType=set
	TargetNamespaces []string `json:"targetNamespaces,omitempty"`

	// ServiceAccountName is the admin specified service account which will be
	// used to deploy operator(s) in this operator group.
	ServiceAccountName string `json:"serviceAccountName,omitempty"`

	// Static tells OLM not to update the OperatorGroup's providedAPIs annotation
	// +optional
	StaticProvidedAPIs bool `json:"staticProvidedAPIs,omitempty"`
}

// OperatorGroupStatus is the status for an OperatorGroupResource.
type OperatorGroupStatus struct {
	// Namespaces is the set of target namespaces for the OperatorGroup.
	// +listType=set
	Namespaces []string `json:"namespaces,omitempty"`

	// ServiceAccountRef references the service account object specified.
	ServiceAccountRef *corev1.ObjectReference `json:"serviceAccountRef,omitempty"`

	// LastUpdated is a timestamp of the last time the OperatorGroup's status was Updated.
	LastUpdated *metav1.Time `json:"lastUpdated"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +genclient

// OperatorGroup is the unit of multitenancy for OLM managed operators.
// It constrains the installation of operators in its namespace to a specified set of target namespaces.
type OperatorGroup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	// +optional
	Spec   OperatorGroupSpec   `json:"spec"`
	Status OperatorGroupStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// OperatorGroupList is a list of OperatorGroup resources.
type OperatorGroupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	// +listType=set
	Items []OperatorGroup `json:"items"`
}

func (o *OperatorGroup) BuildTargetNamespaces() string {
	ns := make([]string, len(o.Status.Namespaces))
	copy(ns, o.Status.Namespaces)
	sort.Strings(ns)
	return strings.Join(ns, ",")
}

// IsServiceAccountSpecified returns true if the spec has a service account name specified.
func (o *OperatorGroup) IsServiceAccountSpecified() bool {
	if o.Spec.ServiceAccountName == "" {
		return false
	}

	return true
}

// HasServiceAccountSynced returns true if the service account specified has been synced.
func (o *OperatorGroup) HasServiceAccountSynced() bool {
	if o.IsServiceAccountSpecified() && o.Status.ServiceAccountRef != nil {
		return true
	}

	return false
}

// OGLabelKeyAndValue returns a key and value that should be applied to namespaces listed in the OperatorGroup.
func (o *OperatorGroup) OGLabelKeyAndValue() (string, string) {
	return fmt.Sprintf(operatorGroupLabelTemplate, o.GetNamespace(), o.GetName()), ""
}

// NamespaceLabelSelector provides a selector that can be used to filter namespaces that belong to the OperatorGroup.
func (o *OperatorGroup) NamespaceLabelSelector() *metav1.LabelSelector {
	if len(o.Spec.TargetNamespaces) == 0 {
		// If no target namespaces are set, check if a selector exists.
		if o.Spec.Selector != nil {
			return o.Spec.Selector
		}
		// No selector exists, return nil which should be used to select EVERYTHING.
		return nil
	}
	// Return a label that should be present on all namespaces defined in the OperatorGroup.Spec.TargetNamespaces field.
	ogKey, ogValue := o.OGLabelKeyAndValue()
	return &metav1.LabelSelector{
		MatchLabels: map[string]string{
			ogKey: ogValue,
		},
	}
}

// IsOperatorGroupLabel returns true if the label is an OperatorGroup label.
func IsOperatorGroupLabel(label string) bool {
	return strings.HasPrefix(label, operatorGroupLabelPrefix)
}
