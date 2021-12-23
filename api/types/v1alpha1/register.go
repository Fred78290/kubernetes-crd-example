package v1alpha1

import (
	"context"
	"reflect"

	apiextensionv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextension "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
)

const GroupName = "example.martin-helmich.de"
const GroupVersion = "v1alpha1"
const CRDPlural = "projects"
const CRDSingular = "project"
const CRDShortName = "pj"
const FullCRDName = CRDPlural + "." + GroupName

var SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: GroupVersion}

var (
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme   = SchemeBuilder.AddToScheme
)

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&Project{},
		&ProjectList{},
	)

	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}

func CreateCRD(c *rest.Config) error {
	var err error
	var clientset apiextension.Interface

	if clientset, err = apiextension.NewForConfig(c); err != nil {
		return err
	}

	crd := &apiextensionv1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{Name: FullCRDName},
		Spec: apiextensionv1.CustomResourceDefinitionSpec{
			Group: GroupName,
			Versions: []apiextensionv1.CustomResourceDefinitionVersion{
				{
					Name:    GroupVersion,
					Served:  true,
					Storage: true,
					Schema: &apiextensionv1.CustomResourceValidation{
						OpenAPIV3Schema: &apiextensionv1.JSONSchemaProps{
							Type: "object",
							Properties: map[string]apiextensionv1.JSONSchemaProps{
								"spec": {
									Type: "object",
									Properties: map[string]apiextensionv1.JSONSchemaProps{
										"replicas": {
											Type: "integer",
										},
									},
								},
							},
						},
					},
				},
			},
			Scope: apiextensionv1.NamespaceScoped,
			Names: apiextensionv1.CustomResourceDefinitionNames{
				Plural:     CRDPlural,
				Singular:   CRDSingular,
				ShortNames: []string{CRDShortName},
				Kind:       reflect.TypeOf(Project{}).Name(),
			},
		},
	}

	_, err = clientset.ApiextensionsV1().CustomResourceDefinitions().Create(context.TODO(), crd, metav1.CreateOptions{})

	if err != nil && apierrors.IsAlreadyExists(err) {
		return nil
	}

	return err
}
