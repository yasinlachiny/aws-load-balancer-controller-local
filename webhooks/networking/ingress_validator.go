package networking

import (
	"context"

	elbv2sdk "github.com/aws/aws-sdk-go/service/elbv2"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/aws-load-balancer-controller/pkg/aws/services"

	//elbv2deploy "sigs.k8s.io/aws-load-balancer-controller/pkg/deploy/elbv2"
	"sigs.k8s.io/aws-load-balancer-controller/pkg/service"
	//	"sigs.k8s.io/aws-load-balancer-controller/pkg/service"

	//"sigs.k8s.io/aws-load-balancer-controller/controllers/service"

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
)

// NewIngressValidator returns a validator for Ingress API.
func NewIngressValidator(elbv2Client services.ELBV2, client client.Client, ingConfig config.IngressConfig, logger logr.Logger) *ingressValidator {
	return &ingressValidator{
		annotationParser:              annotations.NewSuffixAnnotationParser(annotations.AnnotationPrefixIngress),
		classAnnotationMatcher:        ingress.NewDefaultClassAnnotationMatcher(ingConfig.IngressClass),
		classLoader:                   ingress.NewDefaultClassLoader(client),
		disableIngressClassAnnotation: ingConfig.DisableIngressClassAnnotation,
		disableIngressGroupAnnotation: ingConfig.DisableIngressGroupNameAnnotation,
		logger:                        logger,
		elbv2Client:                   elbv2Client,
	}
}

var _ webhook.Validator = &ingressValidator{}

type ingressValidator struct {
	elbv2Client                   services.ELBV2
	annotationParser              annotations.Parser
	classAnnotationMatcher        ingress.ClassAnnotationMatcher
	classLoader                   ingress.ClassLoader
	disableIngressClassAnnotation bool
	disableIngressGroupAnnotation bool
	logger                        logr.Logger
	modelBuilder                  service.ModelBuilder
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

type defaultLoadBalancerManager struct {
	elbv2Client services.ELBV2
	logger      logr.Logger
}

func (m *defaultLoadBalancerManager) FindLoadBalancerByDNSName(ctx context.Context, dnsName string) (string, error) {
	req := &elbv2sdk.DescribeLoadBalancersInput{}
	lbs, err := m.elbv2Client.DescribeLoadBalancersAsList(ctx, req)
	if err != nil {
		return "", err
	}
	for _, lb := range lbs {
		if awssdk.StringValue(lb.DNSName) == dnsName {
			return awssdk.StringValue(lb.LoadBalancerArn), nil
		}
	}
	return "", errors.Errorf("couldn't find LoadBalancer with dnsName: %v", dnsName)
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
	fmt.Print("testwetwetwsdfsdf123123123")
	req := &elbv2sdk.DescribeLoadBalancersInput{}
	lbs, err := v.elbv2Client.DescribeLoadBalancersAsList(ctx, req)
	fmt.Println(lbs)
	fmt.Println(err)
	svc := &corev1.Service{}
	v.modelBuilder.Build(ctx, svc)
	//	lbs, err := v.DescribeLoadBalancersPagesWithContext(ctx,req)
	//	if err != nil {
	//		fmt.Println( err)
	//	}
	fmt.Print("testwetwetwsdfsdf123123123")
	fmt.Print("testwetwetwsdfsdf123123123")
	fmt.Print("testwetwetwsdfsdf123123123")
	fmt.Print("testwetwetwsdfsdf123123123")
	fmt.Print("testwetwetwsdfsdf123123123")
	fmt.Print("testwetwetwsdfsdf123123123111111111111111111111111111111111111111111111111")

	//sdkLBs, err := findSDKLoadBalancers(ctx)
	//        fmt.Println(ctx)
	//	var result []*elbv2sdk.Listener

	//        if err := v.elbv2Client.DescribeLoadBalancersPagesWithContext(ctx , req,func(output *elbv2sdk.DescribeLoadBalancersOutput, _ bool) bool {
	//		fmt.Print(output)
	//		return true

	//	}); err != nil {
	//		fmt.Print("133333333333332222222222222222222222222222222222222222222222222222")
	//	}
	//	return result, nil

	//	if err := c.DescribeListenersPagesWithContext(ctx, input, func(output *elbv2.DescribeListenersOutput, _ bool) bool {
	//		result = append(result, output.Listeners...)
	//		return true
	//	}); err != nil {
	//		return nil, err
	//	}
	//	return result, nil

	fmt.Print("testwetwetwsdfsdf123123123")
	fmt.Print("testwetwetwsdfsdf123123123")
	fmt.Print("testwetwetwsdfsdf123123123")
	fmt.Print("testwetwetwsdfsdf123123123")
	fmt.Print("testwetwetwsdfsdf123123123")
	fmt.Print("testwetwetwsdfsdf123123123")
	fmt.Print("testwetwetwsdfsdf123123123")
	fmt.Print("testwetwetwsdfsdf123123123")
	fmt.Print("testwetwetwsdfsdf123123123")
	fmt.Print("testwetwetwsdfsdf123123123")

	return nil
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
