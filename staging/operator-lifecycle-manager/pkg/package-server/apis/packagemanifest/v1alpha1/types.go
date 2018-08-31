package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// PackageManifestList is a list of PackageManifest objects.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type PackageManifestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []PackageManifest `json:"items"`
}

// PackageManifest holds information about a package, which is a reference to one (or more)
// channels under a single package.
// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type PackageManifest struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PackageManifestSpec   `json:"spec,omitempty"`
	Status PackageManifestStatus `json:"status,omitempty"`
}

// PackageManifestSpec defines the desired state of PackageManifest
type PackageManifestSpec struct{}

// PackageManifestStatus represents the current status of the PackageManifest
type PackageManifestStatus struct {
	// CatalogSourceName is the name of the CatalogSource this package belongs to
	CatalogSourceName string `json:"catalogSource"`

	//  CatalogSourceNamespace is the namespace of the owning CatalogSource
	CatalogSourceNamespace string `json:"catalogSourceNamespace"`

	// PackageName is the name of the overall package, ala `etcd`.
	PackageName string `json:"packageName"`

	// Channels are the declared channels for the package, ala `stable` or `alpha`.
	Channels []PackageChannel `json:"channels"`

	// DefaultChannelName is, if specified, the name of the default channel for the package. The
	// default channel will be installed if no other channel is explicitly given. If the package
	// has a single channel, then that channel is implicitly the default.
	DefaultChannelName string `json:"defaultChannel"`
}

// GetDefaultChannel gets the default channel or returns the only one if there's only one. returns empty string if it
// can't determine the default
func (m PackageManifest) GetDefaultChannel() string {
	if m.Status.DefaultChannelName != "" {
		return m.Status.DefaultChannelName
	}
	if len(m.Status.Channels) == 1 {
		return m.Status.Channels[0].Name
	}
	return ""
}

// PackageChannel defines a single channel under a package, pointing to a version of that
// package.
type PackageChannel struct {
	// Name is the name of the channel, e.g. `alpha` or `stable`
	Name string `json:"name"`

	// CurrentCSVName defines a reference to the CSV holding the version of this package currently
	// for the channel.
	CurrentCSVName string `json:"currentCSV"`
}

// IsDefaultChannel returns true if the PackageChannel is the default for the PackageManifest
func (pc PackageChannel) IsDefaultChannel(pm PackageManifest) bool {
	return pc.Name == pm.Status.DefaultChannelName || len(pm.Status.Channels) == 1
}
