package resolver

import (
	"testing"

	"github.com/blang/semver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/operator-framework/api/pkg/operators/v1alpha1"
	"github.com/operator-framework/operator-registry/pkg/api"
	opregistry "github.com/operator-framework/operator-registry/pkg/registry"

	"github.com/operator-framework/api/pkg/operators/v1alpha1"
	"github.com/operator-framework/operator-lifecycle-manager/pkg/api/client/clientset/versioned"
	"github.com/operator-framework/operator-lifecycle-manager/pkg/api/client/clientset/versioned/fake"
	"github.com/operator-framework/operator-lifecycle-manager/pkg/api/client/informers/externalversions"
	controllerbundle "github.com/operator-framework/operator-lifecycle-manager/pkg/controller/bundle"
	"github.com/operator-framework/operator-lifecycle-manager/pkg/controller/registry/resolver/solve"
	"github.com/operator-framework/operator-lifecycle-manager/pkg/lib/operatorlister"
)

func TestSolveOperators(t *testing.T) {
	APISet := APISet{opregistry.APIKey{"g", "v", "k", "ks"}: struct{}{}}
	Provides := APISet

func TestNamespaceResolver(t *testing.T) {
	namespace := "catsrc-namespace"
	catalog := CatalogKey{"catsrc", namespace}
	type out struct {
		steps       [][]*v1alpha1.Step
		lookups     []v1alpha1.BundleLookup
		subs        []*v1alpha1.Subscription
		err         error
		solverError solve.NotSatisfiable
	}
	nothing := out{
		steps:   [][]*v1alpha1.Step{},
		lookups: []v1alpha1.BundleLookup{},
		subs:    []*v1alpha1.Subscription{},
	}
	tests := []struct {
		name             string
		clusterState     []runtime.Object
		querier          SourceQuerier
		bundlesByCatalog map[CatalogKey][]*api.Bundle
		out              out
	}{
		{
			name: "SingleNewSubscription/NoDeps",
			clusterState: []runtime.Object{
				newSub(namespace, "a", "alpha", catalog),
			},
			bundlesByCatalog: map[CatalogKey][]*api.Bundle{
				catalog: {
					bundle("a.v1", "a", "alpha", "", nil, nil, nil, nil),
				},
			},
			out: out{
				steps: [][]*v1alpha1.Step{
					bundleSteps(bundle("a.v1", "a", "alpha", "", nil, nil, nil, nil), namespace, "", catalog),
				},
				subs: []*v1alpha1.Subscription{
					updatedSub(namespace, "a.v1", "", "a", "alpha", catalog),
				},
			},
		},
		{
			name: "SingleNewSubscription/ResolveOne",
			clusterState: []runtime.Object{
				newSub(namespace, "a", "alpha", catalog),
			},
			bundlesByCatalog: map[CatalogKey][]*api.Bundle{
				catalog: {
					bundle("b.v1", "b", "beta", "", Provides1, nil, nil, nil),
					bundle("a.v1", "a", "alpha", "", nil, Requires1, nil, nil),
				},
			},
			out: out{
				steps: [][]*v1alpha1.Step{
					bundleSteps(bundle("a.v1", "a", "alpha", "", nil, Requires1, nil, nil), namespace, "", catalog),
					bundleSteps(bundle("b.v1", "b", "beta", "", Provides1, nil, nil, nil), namespace, "", catalog),
					subSteps(namespace, "b.v1", "b", "beta", catalog),
				},
				subs: []*v1alpha1.Subscription{
					updatedSub(namespace, "a.v1", "", "a", "alpha", catalog),
				},
			},
		},
		{
			name: "SingleNewSubscription/ResolveOne/BundlePath",
			clusterState: []runtime.Object{
				newSub(namespace, "a", "alpha", catalog),
			},
			bundlesByCatalog: map[CatalogKey][]*api.Bundle{
				catalog: {
					bundle("a.v1", "a", "alpha", "", nil, Requires1, nil, nil),
					stripManifests(withBundlePath(bundle("b.v1", "b", "beta", "", Provides1, nil, nil, nil), "quay.io/test/bundle@sha256:abcd")),
				},
			},
			out: out{
				steps: [][]*v1alpha1.Step{
					bundleSteps(bundle("a.v1", "a", "alpha", "", nil, Requires1, nil, nil), namespace, "", catalog),
					subSteps(namespace, "b.v1", "b", "beta", catalog),
				},
				lookups: []v1alpha1.BundleLookup{
					{
						Path:       "quay.io/test/bundle@sha256:abcd",
						Identifier: "b.v1",
						CatalogSourceRef: &corev1.ObjectReference{
							Namespace: catalog.Namespace,
							Name:      catalog.Name,
						},
						Conditions: []v1alpha1.BundleLookupCondition{
							{
								Type:    BundleLookupConditionPacked,
								Status:  corev1.ConditionTrue,
								Reason:  controllerbundle.NotUnpackedReason,
								Message: controllerbundle.NotUnpackedMessage,
							},
							{
								Type:    v1alpha1.BundleLookupPending,
								Status:  corev1.ConditionTrue,
								Reason:  controllerbundle.JobNotStartedReason,
								Message: controllerbundle.JobNotStartedMessage,
							},
						},
					},
				},
				subs: []*v1alpha1.Subscription{
					updatedSub(namespace, "a.v1", "", "a", "alpha", catalog),
				},
			},
		},
		{
			name: "SingleNewSubscription/ResolveOne/AdditionalBundleObjects",
			clusterState: []runtime.Object{
				newSub(namespace, "a", "alpha", catalog),
			},
			bundlesByCatalog: map[CatalogKey][]*api.Bundle{
				catalog: {
					withBundleObject(bundle("b.v1", "b", "beta", "", Provides1, nil, nil, nil), u(&rbacv1.RoleBinding{TypeMeta: metav1.TypeMeta{Kind: "RoleBinding", APIVersion: "rbac.authorization.k8s.io/v1"}, ObjectMeta: metav1.ObjectMeta{Name: "test-rb"}})),
					bundle("a.v1", "a", "alpha", "", nil, Requires1, nil, nil),
				},
			},
			out: out{
				steps: [][]*v1alpha1.Step{
					bundleSteps(bundle("a.v1", "a", "alpha", "", nil, Requires1, nil, nil), namespace, "", catalog),
					bundleSteps(withBundleObject(bundle("b.v1", "b", "beta", "", Provides1, nil, nil, nil), u(&rbacv1.RoleBinding{TypeMeta: metav1.TypeMeta{Kind: "RoleBinding", APIVersion: "rbac.authorization.k8s.io/v1"}, ObjectMeta: metav1.ObjectMeta{Name: "test-rb"}})), namespace, "", catalog),
					subSteps(namespace, "b.v1", "b", "beta", catalog),
				},
				subs: []*v1alpha1.Subscription{
					updatedSub(namespace, "a.v1", "", "a", "alpha", catalog),
				},
			},
		},
		{
			name: "SingleNewSubscription/ResolveOne/AdditionalBundleObjects/Service",
			clusterState: []runtime.Object{
				newSub(namespace, "a", "alpha", catalog),
			},
			bundlesByCatalog: map[CatalogKey][]*api.Bundle{
				catalog: {
					withBundleObject(bundle("b.v1", "b", "beta", "", Provides1, nil, nil, nil), u(&corev1.Service{TypeMeta: metav1.TypeMeta{Kind: "Service", APIVersion: ""}, ObjectMeta: metav1.ObjectMeta{Name: "test-service"}})),
					bundle("a.v1", "a", "alpha", "", nil, Requires1, nil, nil),
				},
			},
			out: out{
				steps: [][]*v1alpha1.Step{
					bundleSteps(bundle("a.v1", "a", "alpha", "", nil, Requires1, nil, nil), namespace, "", catalog),
					bundleSteps(withBundleObject(bundle("b.v1", "b", "beta", "", Provides1, nil, nil, nil), u(&corev1.Service{TypeMeta: metav1.TypeMeta{Kind: "Service", APIVersion: ""}, ObjectMeta: metav1.ObjectMeta{Name: "test-service"}})), namespace, "", catalog),
					subSteps(namespace, "b.v1", "b", "beta", catalog),
				},
				subs: []*v1alpha1.Subscription{
					updatedSub(namespace, "a.v1", "", "a", "alpha", catalog),
				},
			},
		},
	}
	satResolver := SatResolver{
		cache: getFakeOperatorCache(fakeNamespacedOperatorCache),
	}

	operators, err := satResolver.SolveOperators([]string{"olm"}, csvs, subs)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(operators))
	for _, op := range operators {
		assert.Equal(t, "1.0.1", op.Version().String())
	}
}

func TestSolveOperators_FindLatestVersionWithDependencies(t *testing.T) {
	APISet := APISet{opregistry.APIKey{"g", "v", "k", "ks"}: struct{}{}}
	Provides := APISet

	namespace := "olm"
	catalog := registry.CatalogKey{"community", namespace}

	csv := existingOperator(namespace, "packageA.v1", "packageA", "alpha", "", Provides, nil, nil, nil)
	csvs := []*v1alpha1.ClusterServiceVersion{csv}
	sub := existingSub(namespace, "packageA.v1", "packageA", "alpha", catalog)
	newSub := newSub(namespace, "packageB", "alpha", catalog)
	subs := []*v1alpha1.Subscription{sub, newSub}

	opToAddVersionDeps := []*api.Dependency{
		{
			name: "SingleNewSubscription/DependencyMissing",
			clusterState: []runtime.Object{
				newSub(namespace, "a", "alpha", catalog),
			},
			bundlesByCatalog: map[CatalogKey][]*api.Bundle{
				catalog: {
					bundle("a.v1", "a", "alpha", "", nil, Requires1, nil, nil),
				},
			},
			out: out{
				steps:   [][]*v1alpha1.Step{},
				lookups: []v1alpha1.BundleLookup{},
				subs:    []*v1alpha1.Subscription{},
				solverError: solve.NotSatisfiable([]solve.AppliedConstraint{
					{
						Installable: VirtPackageInstallable{
							identifier: "a",
							constraints: []solve.Constraint{
								solve.Mandatory(),
								solve.Dependency("catsrc/catsrc-namespace/alpha/a.v1"),
							},
						},
						Constraint: solve.Dependency("catsrc/catsrc-namespace/alpha/a.v1"),
					},
					{
						Installable: &BundleInstallable{
							identifier:  "catsrc/catsrc-namespace/alpha/a.v1",
							constraints: []solve.Constraint{solve.Prohibited()},
						},
						Constraint: solve.Prohibited(),
					},
					{
						Installable: VirtPackageInstallable{
							identifier: "a",
							constraints: []solve.Constraint{
								solve.Mandatory(),
								solve.Dependency("catsrc/catsrc-namespace/alpha/a.v1"),
							},
						},
						Constraint: solve.Mandatory(),
					},
				}),
			},
		},
		{
			name: "InstalledSub/NoUpdates",
			clusterState: []runtime.Object{
				existingSub(namespace, "a.v1", "a", "alpha", catalog),
				existingOperator(namespace, "a.v1", "a", "alpha", "", Provides1, nil, nil, nil),
			},
			bundlesByCatalog: map[CatalogKey][]*api.Bundle{
				catalog: {
					bundle("a.v1", "a", "alpha", "", Provides1, nil, nil, nil),
				},
			},
			out: nothing,
		},
		{
			name: "InstalledSub/UpdateAvailable",
			clusterState: []runtime.Object{
				existingSub(namespace, "a.v1", "a", "alpha", catalog),
				existingOperator(namespace, "a.v1", "a", "alpha", "", Provides1, nil, nil, nil),
			},
			bundlesByCatalog: map[CatalogKey][]*api.Bundle{
				catalog: {
					bundle("a.v2", "a", "alpha", "a.v1", Provides1, nil, nil, nil),
					bundle("a.v1", "a", "alpha", "", Provides1, nil, nil, nil),
				},
			},
			out: out{
				steps: [][]*v1alpha1.Step{
					bundleSteps(bundle("a.v2", "a", "alpha", "a.v1", Provides1, nil, nil, nil), namespace, "", catalog),
				},
				subs: []*v1alpha1.Subscription{
					updatedSub(namespace, "a.v2", "a.v1", "a", "alpha", catalog),
				},
			},
		},
		{
			name: "InstalledSub/UpdateAvailable/FromBundlePath",
			clusterState: []runtime.Object{
				existingSub(namespace, "a.v1", "a", "alpha", catalog),
				existingOperator(namespace, "a.v1", "a", "alpha", "", Provides1, nil, nil, nil),
			},
			bundlesByCatalog: map[CatalogKey][]*api.Bundle{catalog: {
				stripManifests(withBundlePath(bundle("a.v2", "a", "alpha", "a.v1", Provides1, nil, nil, nil), "quay.io/test/bundle@sha256:abcd"))},
			},
			out: out{
				steps: [][]*v1alpha1.Step{},
				lookups: []v1alpha1.BundleLookup{
					{
						Path:       "quay.io/test/bundle@sha256:abcd",
						Identifier: "a.v2",
						Replaces:   "a.v1",
						CatalogSourceRef: &corev1.ObjectReference{
							Namespace: catalog.Namespace,
							Name:      catalog.Name,
						},
						Conditions: []v1alpha1.BundleLookupCondition{
							{
								Type:    BundleLookupConditionPacked,
								Status:  corev1.ConditionTrue,
								Reason:  controllerbundle.NotUnpackedReason,
								Message: controllerbundle.NotUnpackedMessage,
							},
							{
								Type:    v1alpha1.BundleLookupPending,
								Status:  corev1.ConditionTrue,
								Reason:  controllerbundle.JobNotStartedReason,
								Message: controllerbundle.JobNotStartedMessage,
							},
						},
					},
				},
				subs: []*v1alpha1.Subscription{
					updatedSub(namespace, "a.v2", "a.v1", "a", "alpha", catalog),
				},
			},
		},
		{
			name: "InstalledSub/NoRunningOperator",
			clusterState: []runtime.Object{
				existingSub(namespace, "a.v1", "a", "alpha", catalog),
			},
			bundlesByCatalog: map[CatalogKey][]*api.Bundle{
				catalog: {
					bundle("a.v1", "a", "alpha", "", Provides1, nil, nil, nil),
				},
			},
			out: out{
				steps: [][]*v1alpha1.Step{
					bundleSteps(bundle("a.v1", "a", "alpha", "", Provides1, nil, nil, nil), namespace, "", catalog),
				},
				// no updated subs because existingSub already points the right CSV, it just didn't exist for some reason
				subs: []*v1alpha1.Subscription{},
			},
		},
		{
			name: "InstalledSub/UpdateFound/UpdateRequires/ResolveOne",
			clusterState: []runtime.Object{
				existingSub(namespace, "a.v1", "a", "alpha", catalog),
				existingOperator(namespace, "a.v1", "a", "alpha", "", Provides1, nil, nil, nil),
			},
			bundlesByCatalog: map[CatalogKey][]*api.Bundle{
				catalog: {
					bundle("a.v1", "a", "alpha", "", nil, nil, nil, nil),
					bundle("a.v2", "a", "alpha", "a.v1", nil, Requires1, nil, nil),
					bundle("b.v1", "b", "beta", "", Provides1, nil, nil, nil),
				},
			},
			out: out{
				steps: [][]*v1alpha1.Step{
					bundleSteps(bundle("a.v2", "a", "alpha", "a.v1", nil, Requires1, nil, nil), namespace, "", catalog),
					bundleSteps(bundle("b.v1", "b", "beta", "", Provides1, nil, nil, nil), namespace, "", catalog),
					subSteps(namespace, "b.v1", "b", "beta", catalog),
				},
				subs: []*v1alpha1.Subscription{
					updatedSub(namespace, "a.v2", "a.v1", "a", "alpha", catalog),
				},
			},
		},
	}
	satResolver := SatResolver{
		cache: getFakeOperatorCache(fakeNamespacedOperatorCache),
	}

	operators, err := satResolver.SolveOperators([]string{"olm"}, csvs, subs)
	assert.NoError(t, err)
	assert.Equal(t, 4, len(operators))
	for _, op := range operators {
		assert.Equal(t, "1.0.1", op.Version().String())
	}
}

func TestSolveOperators_FindLatestVersionWithDependencies_LargeCatalogSet(t *testing.T) {
	APISet := APISet{opregistry.APIKey{"g", "v", "k", "ks"}: struct{}{}}
	Provides := APISet

	namespace := "olm"
	catalog := registry.CatalogKey{"community", namespace}

	csv := existingOperator(namespace, "packageA.v1", "packageA", "alpha", "", Provides, nil, nil, nil)
	csvs := []*v1alpha1.ClusterServiceVersion{csv}
	sub := existingSub(namespace, "packageA.v1", "packageA", "alpha", catalog)
	newSub := newSub(namespace, "packageB", "alpha", catalog)
	subs := []*v1alpha1.Subscription{sub, newSub}

	opToAddVersionDeps := []*api.Dependency{
		{
			name: "InstalledSub/UpdateFound/UpdateRequires/ResolveOne/APIServer",
			clusterState: []runtime.Object{
				existingSub(namespace, "a.v1", "a", "alpha", catalog),
				existingOperator(namespace, "a.v1", "a", "alpha", "", nil, nil, Provides1, nil),
			},
			bundlesByCatalog: map[CatalogKey][]*api.Bundle{
				catalog: {
					bundle("a.v1", "a", "alpha", "", nil, nil, nil, nil),
					bundle("a.v2", "a", "alpha", "a.v1", nil, nil, nil, Requires1),
					bundle("b.v1", "b", "beta", "", nil, nil, Provides1, nil),
				},
			},
			out: out{
				steps: [][]*v1alpha1.Step{
					bundleSteps(bundle("a.v2", "a", "alpha", "a.v1", nil, nil, nil, Requires1), namespace, "", catalog),
					bundleSteps(bundle("b.v1", "b", "beta", "", nil, nil, Provides1, nil), namespace, "", catalog),
					subSteps(namespace, "b.v1", "b", "beta", catalog),
				},
				subs: []*v1alpha1.Subscription{
					updatedSub(namespace, "a.v2", "a.v1", "a", "alpha", catalog),
				},
			},
		},
		{
			name: "InstalledSub/SingleNewSubscription/UpdateAvailable/ResolveOne",
			clusterState: []runtime.Object{
				existingSub(namespace, "a.v1", "a", "alpha", catalog),
				existingOperator(namespace, "a.v1", "a", "alpha", "", Provides1, nil, nil, nil),
				newSub(namespace, "b", "beta", catalog),
			},
			bundlesByCatalog: map[CatalogKey][]*api.Bundle{
				catalog: {
					bundle("a.v1", "a", "alpha", "", nil, nil, nil, nil),
					bundle("a.v2", "a", "alpha", "a.v1", nil, nil, nil, nil),
					bundle("b.v1", "b", "beta", "", nil, nil, nil, nil),
				},
			},
			out: out{
				steps: [][]*v1alpha1.Step{
					bundleSteps(bundle("a.v2", "a", "alpha", "a.v1", nil, nil, nil, nil), namespace, "", catalog),
					bundleSteps(bundle("b.v1", "b", "beta", "", nil, nil, nil, nil), namespace, "", catalog),
				},
				subs: []*v1alpha1.Subscription{
					updatedSub(namespace, "a.v2", "a.v1", "a", "alpha", catalog),
					updatedSub(namespace, "b.v1", "", "b", "beta", catalog),
				},
			},
		},
		{
			name: "InstalledSub/SingleNewSubscription/NoRunningOperator/ResolveOne",
			clusterState: []runtime.Object{
				existingSub(namespace, "a.v1", "a", "alpha", catalog),
				newSub(namespace, "b", "beta", catalog),
			},
			bundlesByCatalog: map[CatalogKey][]*api.Bundle{
				catalog: {
					bundle("a.v1", "a", "alpha", "", Provides1, nil, nil, nil),
					bundle("b.v1", "b", "beta", "", nil, nil, nil, nil),
				},
			},
			out: out{
				steps: [][]*v1alpha1.Step{
					bundleSteps(bundle("a.v1", "a", "alpha", "", Provides1, nil, nil, nil), namespace, "", catalog),
					bundleSteps(bundle("b.v1", "b", "beta", "", nil, nil, nil, nil), namespace, "", catalog),
				},
				subs: []*v1alpha1.Subscription{
					updatedSub(namespace, "b.v1", "", "b", "beta", catalog),
				},
			},
		},
	}
	satResolver := SatResolver{
		cache: getFakeOperatorCache(fakeNamespacedOperatorCache),
	}

	operators, err := satResolver.SolveOperators([]string{"olm", "ns2"}, csvs, subs)
	assert.NoError(t, err)
	assert.Equal(t, 4, len(operators))
	for _, op := range operators {
		assert.Equal(t, "community", op.SourceInfo().Catalog.Name)
		assert.Equal(t, "1.0.1", op.Version().String())
	}
}

func TestSolveOperators_FindLatestVersionWithNestedDependencies(t *testing.T) {
	APISet := APISet{opregistry.APIKey{"g", "v", "k", "ks"}: struct{}{}}
	Provides := APISet

	namespace := "olm"
	catalog := registry.CatalogKey{"community", namespace}

	csv := existingOperator(namespace, "packageA.v1", "packageA", "alpha", "", Provides, nil, nil, nil)
	csvs := []*v1alpha1.ClusterServiceVersion{csv}
	sub := existingSub(namespace, "packageA.v1", "packageA", "alpha", catalog)
	newSub := newSub(namespace, "packageB", "alpha", catalog)
	subs := []*v1alpha1.Subscription{sub, newSub}

	opToAddVersionDeps := []*api.Dependency{
		{
			name: "InstalledSub/SingleNewSubscription/NoRunningOperator/ResolveOne/APIServer",
			clusterState: []runtime.Object{
				existingSub(namespace, "a.v1", "a", "alpha", catalog),
				newSub(namespace, "b", "beta", catalog),
			},
			bundlesByCatalog: map[CatalogKey][]*api.Bundle{
				catalog: {
					bundle("a.v1", "a", "alpha", "", nil, nil, Provides1, nil),
					bundle("b.v1", "b", "beta", "", nil, nil, nil, nil),
				},
			},
			out: out{
				steps: [][]*v1alpha1.Step{
					bundleSteps(bundle("a.v1", "a", "alpha", "", nil, nil, Provides1, nil), namespace, "", catalog),
					bundleSteps(bundle("b.v1", "b", "beta", "", nil, nil, nil, nil), namespace, "", catalog),
				},
				subs: []*v1alpha1.Subscription{
					updatedSub(namespace, "b.v1", "", "b", "beta", catalog),
				},
			},
		},
		{
			// This test verifies that version deadlock that could happen with the previous algorithm can't happen here
			name: "NoMoreVersionDeadlock",
			clusterState: []runtime.Object{
				existingSub(namespace, "a.v1", "a", "alpha", catalog),
				existingOperator(namespace, "a.v1", "a", "alpha", "", Provides1, Requires2, nil, nil),
				existingSub(namespace, "b.v1", "b", "alpha", catalog),
				existingOperator(namespace, "b.v1", "b", "alpha", "", Provides2, Requires1, nil, nil),
			},
			bundlesByCatalog: map[CatalogKey][]*api.Bundle{
				catalog: {
					bundle("a.v2", "a", "alpha", "a.v1", Provides3, Requires4, nil, nil),
					bundle("b.v2", "b", "alpha", "b.v1", Provides4, Requires3, nil, nil),
				},
			},
			out: out{
				steps: [][]*v1alpha1.Step{
					bundleSteps(bundle("a.v2", "a", "alpha", "a.v1", Provides3, Requires4, nil, nil), namespace, "", catalog),
					bundleSteps(bundle("b.v2", "b", "alpha", "b.v1", Provides4, Requires3, nil, nil), namespace, "", catalog),
				},
				subs: []*v1alpha1.Subscription{
					updatedSub(namespace, "a.v2", "a.v1", "a", "alpha", catalog),
					updatedSub(namespace, "b.v2", "b.v1", "b", "alpha", catalog),
				},
			},
		},
		{
			// This test verifies that ownership of an api can be migrated between two operators
			name: "OwnedAPITransfer",
			clusterState: []runtime.Object{
				existingSub(namespace, "a.v1", "a", "alpha", catalog),
				existingOperator(namespace, "a.v1", "a", "alpha", "", Provides1, nil, nil, nil),
				existingSub(namespace, "b.v1", "b", "alpha", catalog),
				existingOperator(namespace, "b.v1", "b", "alpha", "", nil, Requires1, nil, nil),
			},
			bundlesByCatalog: map[CatalogKey][]*api.Bundle{
				catalog: {
					bundle("a.v2", "a", "alpha", "a.v1", nil, nil, nil, nil),
					bundle("b.v2", "b", "alpha", "b.v1", Provides1, nil, nil, nil),
				},
			},
			out: out{
				steps: [][]*v1alpha1.Step{
					bundleSteps(bundle("a.v2", "a", "alpha", "a.v1", nil, nil, nil, nil), namespace, "", catalog),
					bundleSteps(bundle("b.v2", "b", "alpha", "b.v1", Provides1, nil, nil, nil), namespace, "", catalog),
				},
				subs: []*v1alpha1.Subscription{
					updatedSub(namespace, "a.v2", "a.v1", "a", "alpha", catalog),
					updatedSub(namespace, "b.v2", "b.v1", "b", "alpha", catalog),
				},
			},
		},
		{
			name: "PicksOlderProvider",
			clusterState: []runtime.Object{
				newSub(namespace, "b", "alpha", catalog),
			},
			bundlesByCatalog: map[CatalogKey][]*api.Bundle{
				catalog: {
					bundle("a.v1", "a", "alpha", "", Provides1, nil, nil, nil),
					bundle("a.v2", "a", "alpha", "a.v1", nil, nil, nil, nil),
					bundle("b.v1", "b", "alpha", "", nil, Requires1, nil, nil),
				},
			},
			out: out{
				steps: [][]*v1alpha1.Step{
					bundleSteps(bundle("a.v1", "a", "alpha", "", Provides1, nil, nil, nil), namespace, "", catalog),
					bundleSteps(bundle("b.v1", "b", "alpha", "", nil, Requires1, nil, nil), namespace, "", catalog),
					subSteps(namespace, "a.v1", "a", "alpha", catalog),
				},
				subs: []*v1alpha1.Subscription{
					updatedSub(namespace, "b.v1", "", "b", "alpha", catalog),
				},
			},
		},
	}
	satResolver := SatResolver{
		cache: getFakeOperatorCache(fakeNamespacedOperatorCache),
	}

	operators, err := satResolver.SolveOperators([]string{"olm"}, csvs, subs)
	assert.NoError(t, err)
	assert.Equal(t, 5, len(operators))
	for _, op := range operators {
		assert.Equal(t, "1.0.1", op.Version().String())
	}
}

func TestSolveOperators_WithDependencies(t *testing.T) {
	APISet := APISet{opregistry.APIKey{"g", "v", "k", "ks"}: struct{}{}}
	Provides := APISet

	namespace := "olm"
	catalog := registry.CatalogKey{"community", namespace}

	csv := existingOperator(namespace, "packageA.v1", "packageA", "alpha", "", Provides, nil, nil, nil)
	csvs := []*v1alpha1.ClusterServiceVersion{csv}
	sub := existingSub(namespace, "packageA.v1", "packageA", "alpha", catalog)
	newSub := newSub(namespace, "packageB", "alpha", catalog)
	subs := []*v1alpha1.Subscription{sub, newSub}

	opToAddVersionDeps := []*api.Dependency{
		{
			name: "InstalledSub/UpdateInHead/SkipRange",
			clusterState: []runtime.Object{
				existingSub(namespace, "a.v1", "a", "alpha", catalog),
				existingOperator(namespace, "a.v1", "a", "alpha", "", Provides1, nil, nil, nil),
			},
			bundlesByCatalog: map[CatalogKey][]*api.Bundle{catalog: {
				bundle("a.v3", "a", "alpha", "a.v2", nil, nil, nil, nil, withVersion("1.0.0"), withSkipRange("< 1.0.0")),
			}},
			out: out{
				steps: [][]*v1alpha1.Step{
					bundleSteps(bundle("a.v3", "a", "alpha", "a.v2", nil, nil, nil, nil), namespace, "a.v1", catalog),
				},
				subs: []*v1alpha1.Subscription{
					updatedSub(namespace, "a.v3", "a.v1", "a", "alpha", catalog),
				},
			},
		},
		{
			// This test uses logic that implements the FakeSourceQuerier to ensure
			// that the required API is provided by the new Operator.
			//
			// Background:
			// OLM used to add the new operator to the generation before removing
			// the old operator from the generation. The logic that removes an operator
			// from the current generation removes the APIs it provides from the list of
			// "available" APIs. This caused OLM to search for an operator that provides the API.
			// If the operator that provides the API uses a skipRange rather than the Spec.Replaces
			// field, the Replaces field is set to an empty string, causing OLM to fail to upgrade.
			name: "InstalledSubs/ExistingOperators/OldCSVsReplaced",
			clusterState: []runtime.Object{
				existingSub(namespace, "a.v1", "a", "alpha", catalog),
				existingSub(namespace, "b.v1", "b", "beta", catalog),
				existingOperator(namespace, "a.v1", "a", "alpha", "", nil, Requires1, nil, nil),
				existingOperator(namespace, "b.v1", "b", "beta", "", Provides1, nil, nil, nil),
			},
			bundlesByCatalog: map[CatalogKey][]*api.Bundle{
				catalog: {
					bundle("a.v1", "a", "alpha", "", nil, nil, nil, nil),
					bundle("a.v2", "a", "alpha", "a.v1", nil, Requires1, nil, nil),
					bundle("b.v1", "b", "beta", "", Provides1, nil, nil, nil),
					bundle("b.v2", "b", "beta", "b.v1", Provides1, nil, nil, nil),
				},
			},
			out: out{
				steps: [][]*v1alpha1.Step{
					bundleSteps(bundle("a.v2", "a", "alpha", "a.v1", nil, Requires1, nil, nil), namespace, "", catalog),
					bundleSteps(bundle("b.v2", "b", "beta", "b.v1", Provides1, nil, nil, nil), namespace, "", catalog),
				},
				subs: []*v1alpha1.Subscription{
					updatedSub(namespace, "a.v2", "a.v1", "a", "alpha", catalog),
					updatedSub(namespace, "b.v2", "b.v1", "b", "beta", catalog),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stopc := make(chan struct{})
			defer func() {
				stopc <- struct{}{}
			}()
			expectedSteps := []*v1alpha1.Step{}
			for _, steps := range tt.out.steps {
				expectedSteps = append(expectedSteps, steps...)
			}
			clientFake, informerFactory, _ := StartResolverInformers(namespace, stopc, tt.clusterState...)
			lister := operatorlister.NewLister()
			lister.OperatorsV1alpha1().RegisterSubscriptionLister(namespace, informerFactory.Operators().V1alpha1().Subscriptions().Lister())
			lister.OperatorsV1alpha1().RegisterClusterServiceVersionLister(namespace, informerFactory.Operators().V1alpha1().ClusterServiceVersions().Lister())
			kClientFake := k8sfake.NewSimpleClientset()

			resolver := NewOperatorsV1alpha1Resolver(lister, clientFake, kClientFake, "", false)

			tt.querier = NewFakeSourceQuerier(tt.bundlesByCatalog)
			steps, lookups, subs, err := resolver.ResolveSteps(namespace, tt.querier)
			require.Equal(t, tt.out.err, err)
			RequireStepsEqual(t, expectedSteps, steps)
			require.ElementsMatch(t, tt.out.lookups, lookups)
			require.ElementsMatch(t, tt.out.subs, subs)

			// todo -- factor this out
			stubSnapshot := &CatalogSnapshot{}
			for _, bundles := range tt.bundlesByCatalog {
				for _, bundle := range bundles {
					op, err := NewOperatorFromBundle(bundle, "", catalog)
					if err != nil {
						t.Fatalf("unexpected error: %v", err)
					}
					op.replaces = bundle.Replaces
					stubSnapshot.operators = append(stubSnapshot.operators, op)
				}
			}
			stubCache := &stubOperatorCacheProvider{
				noc: &NamespacedOperatorCache{
					snapshots: map[CatalogKey]*CatalogSnapshot{
						catalog: stubSnapshot,
					},
				},
			}
			satresolver := &SatResolver{
				cache: stubCache,
			}
			resolver.satResolver = satresolver
			resolver.updatedResolution = true

			steps, lookups, subs, err = resolver.ResolveSteps(namespace, tt.querier)
			if tt.out.solverError == nil {
				require.Equal(t, tt.out.err, err, "%s", err)
			} else {
				// the solver outputs useful information on a failed resolution, which is different from the old resolver
				require.ElementsMatch(t, tt.out.solverError, err.(solve.NotSatisfiable))
			}
			RequireStepsEqual(t, expectedSteps, steps)
			require.ElementsMatch(t, tt.out.lookups, lookups)
			require.ElementsMatch(t, tt.out.subs, subs)
		})
	}

	operators, err := satResolver.SolveOperators([]string{"olm"}, csvs, subs)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(operators))
}

type stubOperatorCacheProvider struct {
	noc *NamespacedOperatorCache
}

func (stub *stubOperatorCacheProvider) Namespaced(namespaces ...string) MultiCatalogOperatorFinder {
	return stub.noc
}

func TestNamespaceResolverRBAC(t *testing.T) {
	namespace := "catsrc-namespace"
	catalog := CatalogKey{"catsrc", namespace}

	namespace := "olm"
	catalog := registry.CatalogKey{"community", namespace}

	csv := existingOperator(namespace, "packageA.v1", "packageA", "alpha", "", Provides, nil, nil, nil)
	csvs := []*v1alpha1.ClusterServiceVersion{csv}
	sub := existingSub(namespace, "packageA.v1", "packageA", "alpha", catalog)
	newSub := newSub(namespace, "packageB", "alpha", catalog)
	subs := []*v1alpha1.Subscription{sub, newSub}

	fakeNamespacedOperatorCache := NamespacedOperatorCache{
		snapshots: map[registry.CatalogKey]*CatalogSnapshot{
			registry.CatalogKey{
				Namespace: "olm",
				Name:      "community",
			}: &CatalogSnapshot{
				operators: []*Operator{
					genOperator("packageA.v1", "0.0.1", "packageA", "alpha", "community", "olm", nil, nil, nil),
					genOperator("packageB.v1", "1.0.0", "packageB", "alpha", "community", "olm", Provides, nil, nil),
					genOperator("packageC.v1", "0.1.0", "packageC", "alpha", "community", "olm", nil, Provides, nil),
				},
			},
		},
	}
	satResolver := SatResolver{
		cache: getFakeOperatorCache(fakeNamespacedOperatorCache),
	}
	tests := []struct {
		name             string
		clusterState     []runtime.Object
		bundlesInCatalog []*api.Bundle
		out              out
	}{
		{
			name: "NewSubscription/Permissions/ClusterPermissions",
			clusterState: []runtime.Object{
				newSub(namespace, "a", "alpha", catalog),
			},
			bundlesInCatalog: []*api.Bundle{bundle},
			out: out{
				steps: [][]*v1alpha1.Step{
					bundleSteps(bundle, namespace, "", catalog),
				},
				subs: []*v1alpha1.Subscription{
					updatedSub(namespace, "a.v1", "", "a", "alpha", catalog),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stopc := make(chan struct{})
			defer func() {
				stopc <- struct{}{}
			}()
			expectedSteps := []*v1alpha1.Step{}
			for _, steps := range tt.out.steps {
				expectedSteps = append(expectedSteps, steps...)
			}
			kClientFake := k8sfake.NewSimpleClientset()
			clientFake, informerFactory, _ := StartResolverInformers(namespace, stopc, tt.clusterState...)
			lister := operatorlister.NewLister()
			lister.OperatorsV1alpha1().RegisterSubscriptionLister(namespace, informerFactory.Operators().V1alpha1().Subscriptions().Lister())
			lister.OperatorsV1alpha1().RegisterClusterServiceVersionLister(namespace, informerFactory.Operators().V1alpha1().ClusterServiceVersions().Lister())

			resolver := NewOperatorsV1alpha1Resolver(lister, clientFake, kClientFake, "", false)
			querier := NewFakeSourceQuerier(map[CatalogKey][]*api.Bundle{catalog: tt.bundlesInCatalog})
			steps, _, subs, err := resolver.ResolveSteps(namespace, querier)
			require.Equal(t, tt.out.err, err)
			RequireStepsEqual(t, expectedSteps, steps)
			require.ElementsMatch(t, tt.out.subs, subs)
		})
	}
}

func TestSolveOperators_DependenciesMultiCatalog(t *testing.T) {
	APISet := APISet{opregistry.APIKey{"g", "v", "k", "ks"}: struct{}{}}
	Provides := APISet

	namespace := "olm"
	catalog := registry.CatalogKey{"community", namespace}

	csv := existingOperator(namespace, "packageA.v1", "packageA", "alpha", "", Provides, nil, nil, nil)
	csvs := []*v1alpha1.ClusterServiceVersion{csv}
	sub := existingSub(namespace, "packageA.v1", "packageA", "alpha", catalog)
	newSub := newSub(namespace, "packageB", "alpha", catalog)
	subs := []*v1alpha1.Subscription{sub, newSub}

	opToAddVersionDeps := []*api.Dependency{
		{
			Type:  "olm.gvk",
			Value: `{"packageName": "packageC", "version"": "0.1.0"}`,
		},
	}
	fakeNamespacedOperatorCache := NamespacedOperatorCache{
		snapshots: map[registry.CatalogKey]*CatalogSnapshot{
			registry.CatalogKey{
				Namespace: "olm",
				Name:      "community",
			}: &CatalogSnapshot{
				operators: []*Operator{
					genOperator("packageA.v1", "0.0.1", "packageA", "alpha", "community", "olm", nil, nil, nil),
					genOperator("packageB.v1", "1.0.0", "packageB", "alpha", "community", "olm", nil, nil, opToAddVersionDeps),
					genOperator("packageC.v1", "0.1.0", "packageC", "alpha", "community", "olm", nil, nil, nil),
				},
			},
			registry.CatalogKey{
				Namespace: "olm",
				Name:      "certified",
			}: &CatalogSnapshot{
				operators: []*Operator{
					genOperator("packageA.v1", "0.0.1", "packageA", "alpha", "certified", "olm", nil, nil, nil),
					genOperator("packageB.v1", "1.0.0", "packageB", "alpha", "certified", "olm", nil, nil, opToAddVersionDeps),
					genOperator("packageC.v1", "0.1.0", "packageC", "alpha", "certified", "olm", nil, nil, nil),
				},
			},
		},
	}
	satResolver := SatResolver{
		cache: getFakeOperatorCache(fakeNamespacedOperatorCache),
	}

	operators, err := satResolver.SolveOperators([]string{"olm"}, csvs, subs)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(operators))
	for _, op := range operators {
		assert.Equal(t, "community", op.SourceInfo().Catalog.Name)
	}
}

func TestSolveOperators_IgnoreUnsatisfiableDependencies(t *testing.T) {
	APISet := APISet{opregistry.APIKey{"g", "v", "k", "ks"}: struct{}{}}
	Provides := APISet

	namespace := "olm"
	catalog := registry.CatalogKey{"community", namespace}

	csv := existingOperator(namespace, "packageA.v1", "packageA", "alpha", "", Provides, nil, nil, nil)
	csvs := []*v1alpha1.ClusterServiceVersion{csv}
	sub := existingSub(namespace, "packageA.v1", "packageA", "alpha", catalog)
	newSub := newSub(namespace, "packageB", "alpha", catalog)
	subs := []*v1alpha1.Subscription{sub, newSub}

	opToAddVersionDeps := []*api.Dependency{
		{
			Type:  "olm.gvk",
			Value: `{"packageName": "packageC", "version"": "0.1.0"}`,
		},
	}
	unsatisfiableVersionDeps := []*api.Dependency{
		{
			Type:  "olm.gvk",
			Value: `{"packageName": "packageD", "version"": "0.1.0"}`,
		},
	}

func updatedSub(namespace, currentOperatorName, installedOperatorName, pkg, channel string, catalog CatalogKey) *v1alpha1.Subscription {
	return &v1alpha1.Subscription{
		ObjectMeta: metav1.ObjectMeta{
			Name:      pkg + "-" + channel,
			Namespace: namespace,
		},
		Spec: &v1alpha1.SubscriptionSpec{
			Package:                pkg,
			Channel:                channel,
			CatalogSource:          catalog.Name,
			CatalogSourceNamespace: catalog.Namespace,
		},
		Status: v1alpha1.SubscriptionStatus{
			CurrentCSV:   currentOperatorName,
			InstalledCSV: installedOperatorName,
		},
	}
	satResolver := SatResolver{
		cache: getFakeOperatorCache(fakeNamespacedOperatorCache),
	}

func existingSub(namespace, operatorName, pkg, channel string, catalog CatalogKey) *v1alpha1.Subscription {
	return &v1alpha1.Subscription{
		ObjectMeta: metav1.ObjectMeta{
			Name:      pkg + "-" + channel,
			Namespace: namespace,
		},
		Spec: &v1alpha1.SubscriptionSpec{
			Package:                pkg,
			Channel:                channel,
			CatalogSource:          catalog.Name,
			CatalogSourceNamespace: catalog.Namespace,
		},
		Status: v1alpha1.SubscriptionStatus{
			CurrentCSV:   operatorName,
			InstalledCSV: operatorName,
		},
	}
}

type FakeOperatorCache struct {
	fakedNamespacedOperatorCache NamespacedOperatorCache
}

func (f *FakeOperatorCache) Namespaced(namespaces ...string) MultiCatalogOperatorFinder {
	return &f.fakedNamespacedOperatorCache
}

func getFakeOperatorCache(fakedNamespacedOperatorCache NamespacedOperatorCache) OperatorCacheProvider {
	return &FakeOperatorCache{
		fakedNamespacedOperatorCache: fakedNamespacedOperatorCache,
	}
}

func genOperator(name, version, pkg, channel, catalogName, catalogNamespace string, requiredAPIs, providedAPIs APISet, dependencies []*api.Dependency) *Operator {
	semversion, _ := semver.Make(version)
	return &Operator{
		name:    name,
		version: &semversion,
		bundle: &api.Bundle{
			PackageName: pkg,
			ChannelName: channel,
		},
		dependencies: dependencies,
		sourceInfo: &OperatorSourceInfo{
			Catalog: registry.CatalogKey{
				Name:      catalogName,
				Namespace: catalogNamespace,
			},
		},
		providedAPIs: providedAPIs,
		requiredAPIs: requiredAPIs,
	}
}
