package transformer

import (
	"github.com/wso2/apk/common-go-libs/apis/dp/v1alpha1"
	dpv1alpha2 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha2"
	dpv1alpha3 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha3"
	dpv1alpha4 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha4"
	corev1 "k8s.io/api/core/v1"
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"
)

// K8sArtifacts k8s artifact representation of API
type K8sArtifacts struct {
	API                 dpv1alpha3.API
	HTTPRoutes          map[string]*gwapiv1.HTTPRoute
	GQLRoutes           map[string]*dpv1alpha2.GQLRoute
	Backends            map[string]*dpv1alpha2.Backend
	Scopes              map[string]*v1alpha1.Scope
	Authentication      map[string]*dpv1alpha2.Authentication
	APIPolicies         map[string]*dpv1alpha4.APIPolicy
	InterceptorServices map[string]*v1alpha1.InterceptorService
	ConfigMaps          map[string]*corev1.ConfigMap
	Secrets             map[string]*corev1.Secret
	BackendJWT          *v1alpha1.BackendJWT
	RateLimitPolicies   map[string]*v1alpha1.RateLimitPolicy
	AIRateLimitPolicies map[string]*dpv1alpha3.AIRateLimitPolicy
}
