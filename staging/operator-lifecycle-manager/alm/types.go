package alm

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/openapi"
)

/////////////////
//  App Types  //
/////////////////

const (
	ALMGroup             = "app.coreos.com"
	AppTypeCRDName       = "apptype-v1s.app.coreos.com"
	AppTypeCRDAPIVersion = "apiextensions.k8s.io/v1beta1" // API version w/ CRD support
)

// AppType defines an Application that can be installed
type AppType struct {
	DisplayName string       `json:"displayName"`
	Description string       `json:"description"`
	Keywords    []string     `json:"keywords"`
	Maintainers []Maintainer `json:"maintainers"`
	Links       []AppLink    `json:"links"`
	Icon        Icon         `json:"iconURL"`
}

type Maintainer struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type AppLink struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Icon struct {
	Data      string `json:"base64data"`
	MediaType string `json:"mediatype"`
}

// Custom Resource of type "AppType" (AppType CRD created by ALM)
type AppTypeResource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	Spec   *AppType      `json:"spec"`
	Status metav1.Status `json:"status"`
}

func NewAppTypeResource(app *AppType) *AppTypeResource {
	resource := AppTypeResource{}
	resource.Kind = AppTypeCRDName
	resource.APIVersion = AppTypeCRDAPIVersion
	resource.Spec = app
	return &resource
}

type AppTypeList struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	Items []*AppType `json:"items"`
}

/////////////////////////////
//  Application Instances  //
/////////////////////////////

// CRD's representing the Apps that will be controlled by their OperatorVersionSpec-installed operator
type AppCRD struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	Spec   AppCRDSpec    `json:"spec"`
	Status metav1.Status `json:"status"`
}

// CRD's must correspond to this schema to be recognized by the ALM
type AppCRDSpec struct {
	metav1.GroupVersionForDiscovery `json:",inline"`

	Scope      string                    `json:"scope"`
	Validation openapi.OpenAPIDefinition `json:"validation"`
	Outputs    []AppOutput               `json:"outputs"`
	Names      ResourceNames             `json:"names"`
}

type AppOutput struct {
	Name         string   `json:"string"`
	Capabilities []string `json:"x-alm-capabilities,omitempty"`
	Description  string   `json:"description"`
}

type ResourceNames struct {
	Plural   string `json:"plural"`
	Singular string `json:"singular"`
	Kind     string `json:"kind"`
}
