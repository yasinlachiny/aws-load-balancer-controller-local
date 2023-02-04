package networking

import (
	"context"
	"fmt"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	networking "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/aws-load-balancer-controller/pkg/annotations"
	"sigs.k8s.io/aws-load-balancer-controller/pkg/config"
	"sigs.k8s.io/aws-load-balancer-controller/pkg/ingress"
	"sigs.k8s.io/aws-load-balancer-controller/pkg/webhook"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

const (
	apiPathValidateNetworkingIngress = "/validate-networking-v1-ingress"
	lbAttrsDeletionProtectionEnabled = "deletion_protection.enabled"
	schemDefault                     = "internal"
)

// NewIngressValidator returns a validator for Ingress API.
func NewIngressValidator(client client.Client, ingConfig config.IngressConfig, logger logr.Logger) *ingressValidator {
	return &ingressValidator{
		annotationParser:              annotations.NewSuffixAnnotationParser(annotations.AnnotationPrefixIngress),
		classAnnotationMatcher:        ingress.NewDefaultClassAnnotationMatcher(ingConfig.IngressClass),
		classLoader:                   ingress.NewDefaultClassLoader(client),
		disableIngressClassAnnotation: ingConfig.DisableIngressClassAnnotation,
		disableIngressGroupAnnotation: ingConfig.DisableIngressGroupNameAnnotation,
		logger:                        logger,
	}
}

var _ webhook.Validator = &ingressValidator{}

type ingressValidator struct {
	annotationParser              annotations.Parser
	classAnnotationMatcher        ingress.ClassAnnotationMatcher
	classLoader                   ingress.ClassLoader
	disableIngressClassAnnotation bool
	disableIngressGroupAnnotation bool
	logger                        logr.Logger
}

func (v *ingressValidator) Prototype(req admission.Request) (runtime.Object, error) {
	return &networking.Ingress{}, nil
}

func (v *ingressValidator) ValidateCreate(ctx context.Context, obj runtime.Object) error {
	ing := obj.(*networking.Ingress)
	if err := v.checkIngressClassAnnotationUsage(ing, nil); err != nil {
		return err
	}
	if err := v.checkGroupNameAnnotationUsage(ing, nil); err != nil {
		return err
	}
	if err := v.checkIngressClassUsage(ctx, ing, nil); err != nil {
		return err
	}
	return nil
}

func (v *ingressValidator) ValidateUpdate(ctx context.Context, obj runtime.Object, oldObj runtime.Object) error {
	ing := obj.(*networking.Ingress)
	oldIng := oldObj.(*networking.Ingress)
	if err := v.checkIngressClassAnnotationUsage(ing, oldIng); err != nil {
		return err
	}
	if err := v.checkGroupNameAnnotationUsage(ing, oldIng); err != nil {
		return err
	}
	if err := v.checkIngressClassUsage(ctx, ing, oldIng); err != nil {
		return err
	}
	if err := v.validateDeletionProtectionAnnotation(ctx, ing, oldIng); err != nil {
		return err
	}

	return nil
}

func (v *ingressValidator) validateDeletionProtectionAnnotation(ctx context.Context, ing *networking.Ingress, oldIng *networking.Ingress) error {
	// Get the values of the "alb.ingress.kubernetes.io/scheme" and "alb.ingress.kubernetes.io/ingress.class" annotations for the old and new Ingress objects
	defaultIngressClass, err := v.classLoader.GetDefaultIngressClass(ctx)
	if err != nil {
		return err
	}
	fmt.Println(defaultIngressClass)
	fmt.Println(defaultIngressClass)
	fmt.Println("test")

	ingClass := defaultIngressClass
	oldIngClass := defaultIngressClass
	fmt.Println(ingClass)
	fmt.Println(oldIngClass)
	if value, ok := ing.Annotations[annotations.IngressClass]; ok {
		ingClass = value
	}

	if value, ok := oldIng.Annotations[annotations.IngressClass]; ok {
		oldIngClass = value
	}

	rawSchemaold := ""
	rawSchema := ""

	if exists := v.annotationParser.ParseStringAnnotation(annotations.IngressSuffixScheme, &rawSchemaold, oldIng.Annotations); !exists {
		rawSchemaold = schemDefault
	}

	if exists := v.annotationParser.ParseStringAnnotation(annotations.IngressSuffixScheme, &rawSchema, ing.Annotations); !exists {
		rawSchema = schemDefault
	}

	// // Check if the scheme or type of the load balancer changed in the new Ingress object
	// if rawSchemaold != rawSchema || ingClass != oldIngClass {
	// 	// Check if the Ingress object had the deletion protection annotation enabled
	// 	enabled, err := v.getDeletionProtectionEnabled(ing)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	if enabled == "true" {
	// 		return errors.Errorf("cannot change the scheme or type of ingress %s/%s with deletion protection enabled", ing.Namespace, ing.Name)
	// 	}
	// }
	return nil
}

// getDeletionProtectionEnabled extracts the value of the "deletion_protection.enabled" attribute from the "alb.ingress.kubernetes.io/load-balancer-attributes" annotation of the given Ingress object.
// If the annotation or the attribute is not present, it returns an empty string.
func (v *ingressValidator) getDeletionProtectionEnabled(ing *networking.Ingress) (string, error) {
	var lbAttributes map[string]string

	_, err := v.annotationParser.ParseStringMapAnnotation(annotations.IngressSuffixLoadBalancerAttributes, &lbAttributes, ing.Annotations)
	if err != nil {
		return "", err
	}

	return lbAttributes[lbAttrsDeletionProtectionEnabled], nil

}

func (v *ingressValidator) ValidateDelete(ctx context.Context, obj runtime.Object) error {
	return nil
}

// checkIngressClassAnnotationUsage checks the usage of kubernetes.io/ingress.class annotation.
// kubernetes.io/ingress.class annotation cannot be set to the ingress class for this controller once disabled,
// so that we enforce users to use spec.ingressClassName in Ingress and IngressClass resource instead.
func (v *ingressValidator) checkIngressClassAnnotationUsage(ing *networking.Ingress, oldIng *networking.Ingress) error {
	if !v.disableIngressClassAnnotation {
		return nil
	}
	usedInNewIng := false
	usedInOldIng := false
	if ingClassAnnotation, exists := ing.Annotations[annotations.IngressClass]; exists {
		if v.classAnnotationMatcher.Matches(ingClassAnnotation) {
			usedInNewIng = true
		}
	}
	if oldIng != nil {
		if ingClassAnnotation, exists := oldIng.Annotations[annotations.IngressClass]; exists {
			if v.classAnnotationMatcher.Matches(ingClassAnnotation) {
				usedInOldIng = true
			}
		}
	}
	if !usedInOldIng && usedInNewIng {
		return errors.Errorf("new usage of `%s` annotation is forbidden", annotations.IngressClass)
	}
	return nil
}

// checkGroupNameAnnotationUsage checks the usage of "group.name" annotation.
// "group.name" annotation cannot be set once disabled,
// so that we enforce users to use spec.group in IngressClassParams resource instead.
func (v *ingressValidator) checkGroupNameAnnotationUsage(ing *networking.Ingress, oldIng *networking.Ingress) error {
	if !v.disableIngressGroupAnnotation {
		return nil
	}
	usedInNewIng := false
	usedInOldIng := false
	newGroupName := ""
	oldGroupName := ""
	if exists := v.annotationParser.ParseStringAnnotation(annotations.IngressSuffixGroupName, &newGroupName, ing.Annotations); exists {
		usedInNewIng = true
	}
	if oldIng != nil {
		if exists := v.annotationParser.ParseStringAnnotation(annotations.IngressSuffixGroupName, &oldGroupName, oldIng.Annotations); exists {
			usedInOldIng = true
		}
	}

	if usedInNewIng {
		if !usedInOldIng || (newGroupName != oldGroupName) {
			return errors.Errorf("new usage of `%s/%s` annotation is forbidden", annotations.AnnotationPrefixIngress, annotations.IngressSuffixGroupName)
		}
	}
	return nil
}

// checkIngressClassUsage checks the usage of "ingressClassName" field.
// if ingressClassName is mutated, it must refer to a existing & valid IngressClass.
func (v *ingressValidator) checkIngressClassUsage(ctx context.Context, ing *networking.Ingress, oldIng *networking.Ingress) error {
	usedInNewIng := false
	usedInOldIng := false
	newIngressClassName := ""
	oldIngressClassName := ""

	if ing.Spec.IngressClassName != nil {
		usedInNewIng = true
		newIngressClassName = awssdk.StringValue(ing.Spec.IngressClassName)
	}
	if oldIng != nil && oldIng.Spec.IngressClassName != nil {
		usedInOldIng = true
		oldIngressClassName = awssdk.StringValue(oldIng.Spec.IngressClassName)
	}

	if usedInNewIng {
		if !usedInOldIng || (newIngressClassName != oldIngressClassName) {
			_, err := v.classLoader.Load(ctx, ing)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// +kubebuilder:webhook:path=/validate-networking-v1-ingress,mutating=false,failurePolicy=fail,groups=networking.k8s.io,resources=ingresses,verbs=create;update,versions=v1,name=vingress.elbv2.k8s.aws,sideEffects=None,matchPolicy=Equivalent,webhookVersions=v1,admissionReviewVersions=v1beta1

func (v *ingressValidator) SetupWithManager(mgr ctrl.Manager) {
	mgr.GetWebhookServer().Register(apiPathValidateNetworkingIngress, webhook.ValidatingWebhookForValidator(v))
}
