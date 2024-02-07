package space

import (
	"fmt"
	"time"

	toolchainv1alpha1 "github.com/codeready-toolchain/api/api/v1alpha1"
	"github.com/codeready-toolchain/toolchain-common/pkg/hash"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type Option func(space *toolchainv1alpha1.Space)

func WithoutSpecTargetCluster() Option {
	return func(space *toolchainv1alpha1.Space) {
		space.Spec.TargetCluster = ""
	}
}

func WithSpecTargetCluster(name string) Option {
	return func(space *toolchainv1alpha1.Space) {
		space.Spec.TargetCluster = name
	}
}

func WithSpecTargetClusterRoles(roles []string) Option {
	return func(space *toolchainv1alpha1.Space) {
		space.Spec.TargetClusterRoles = roles
	}
}

func WithName(name string) Option {
	return func(space *toolchainv1alpha1.Space) {
		space.ObjectMeta.Name = name
		space.ObjectMeta.GenerateName = ""
	}
}

func WithGenerateName(namePrefix string) Option {
	return func(space *toolchainv1alpha1.Space) {
		space.ObjectMeta.Name = ""
		space.ObjectMeta.GenerateName = namePrefix + "-"
	}
}

func WithSpecParentSpace(name string) Option {
	return func(space *toolchainv1alpha1.Space) {
		space.Spec.ParentSpace = name
	}
}

func WithLabel(key, value string) Option {
	return func(space *toolchainv1alpha1.Space) {
		if space.ObjectMeta.Labels == nil {
			space.ObjectMeta.Labels = map[string]string{}
		}
		space.ObjectMeta.Labels[key] = value
	}
}

func WithDefaultTier() Option {
	return func(space *toolchainv1alpha1.Space) {
		space.Spec.TierName = ""
	}
}

func WithTierName(tierName string) Option {
	return func(space *toolchainv1alpha1.Space) {
		space.Spec.TierName = tierName
	}
}

func WithDisableInheritance(disableInheritance bool) Option {
	return func(space *toolchainv1alpha1.Space) {
		space.Spec.DisableInheritance = disableInheritance
	}
}

func WithTierHashLabelFor(tier *toolchainv1alpha1.NSTemplateTier) Option {
	return func(space *toolchainv1alpha1.Space) {
		h, _ := hash.ComputeHashForNSTemplateTier(tier) // we can assume the JSON marshalling will always work
		if space.ObjectMeta.Labels == nil {
			space.ObjectMeta.Labels = map[string]string{}
		}
		space.ObjectMeta.Labels[hash.TemplateTierHashLabelKey(tier.Name)] = h
	}
}

func WithTierNameAndHashLabelFor(tier *toolchainv1alpha1.NSTemplateTier) Option {
	return func(space *toolchainv1alpha1.Space) {
		WithTierName(tier.Name)(space)
		WithTierHashLabelFor(tier)(space)
	}
}

func WithStatusTargetCluster(name string) Option {
	return func(space *toolchainv1alpha1.Space) {
		space.Status.TargetCluster = name
	}
}

func WithStatusProvisionedNamespaces(provisionedNamespaces []toolchainv1alpha1.SpaceNamespace) Option {
	return func(space *toolchainv1alpha1.Space) {
		space.Status.ProvisionedNamespaces = provisionedNamespaces
	}
}

func WithFinalizer() Option {
	return func(space *toolchainv1alpha1.Space) {
		space.Finalizers = append(space.Finalizers, toolchainv1alpha1.FinalizerName)
	}
}

func WithDeletionTimestamp() Option {
	return func(space *toolchainv1alpha1.Space) {
		now := metav1.NewTime(time.Now())
		space.DeletionTimestamp = &now
	}
}

func WithCondition(c toolchainv1alpha1.Condition) Option {
	return func(space *toolchainv1alpha1.Space) {
		space.Status.Conditions = append(space.Status.Conditions, c)
	}
}

func WithCreatorLabel(creator string) Option {
	return func(space *toolchainv1alpha1.Space) {
		if space.Labels == nil {
			space.Labels = map[string]string{}
		}
		space.Labels[toolchainv1alpha1.SpaceCreatorLabelKey] = creator
	}
}

func WithCreationTimestamp(t time.Time) Option {
	return func(space *toolchainv1alpha1.Space) {
		space.CreationTimestamp = metav1.NewTime(t)
	}
}

func WithStateLabel(stateValue string) Option {
	return func(space *toolchainv1alpha1.Space) {
		if space.Labels == nil {
			space.Labels = map[string]string{}
		}
		space.Labels[toolchainv1alpha1.SpaceStateLabelKey] = stateValue
	}
}

func CreatedBefore(before time.Duration) Option {
	return func(space *toolchainv1alpha1.Space) {
		space.ObjectMeta.CreationTimestamp = metav1.Time{Time: time.Now().Add(-before)}
	}
}

func NewSpace(namespace, name string, options ...Option) *toolchainv1alpha1.Space {
	space := &toolchainv1alpha1.Space{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
		},
		Spec: toolchainv1alpha1.SpaceSpec{
			TierName: "base1ns",
		},
	}
	for _, apply := range options {
		apply(space)
	}
	return space
}

func NewSpaceWithGeneratedName(namespace, prefix string, options ...Option) *toolchainv1alpha1.Space {
	space := &toolchainv1alpha1.Space{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:    namespace,
			GenerateName: prefix,
		},
		Spec: toolchainv1alpha1.SpaceSpec{
			TierName: "base1ns",
		},
	}
	for _, apply := range options {
		apply(space)
	}
	return space
}

func NewSpaces(size int, namespace, nameFmt string, options ...Option) []runtime.Object {
	murs := make([]runtime.Object, size)
	for i := 0; i < size; i++ {
		murs[i] = NewSpace(namespace, fmt.Sprintf(nameFmt, i), options...)
	}
	return murs
}
