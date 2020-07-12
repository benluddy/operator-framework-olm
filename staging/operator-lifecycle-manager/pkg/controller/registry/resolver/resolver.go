package resolver

import (
	"context"
	"fmt"
	"github.com/operator-framework/operator-lifecycle-manager/pkg/controller/registry"

	"github.com/blang/semver"
	"github.com/sirupsen/logrus"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"

	"github.com/operator-framework/api/pkg/operators/v1alpha1"
	"github.com/operator-framework/operator-lifecycle-manager/pkg/controller/registry/resolver/solver"
)

type OperatorResolver interface {
	SolveOperators(csvs []*v1alpha1.ClusterServiceVersion, subs []*v1alpha1.Subscription, add map[OperatorSourceInfo]struct{}) (OperatorSet, error)
}

type SatResolver struct {
	cache OperatorCacheProvider
	log logrus.FieldLogger
}

type OperatorsV1alpha1Resolver struct {
	subLister              v1alpha1listers.SubscriptionLister
	csvLister              v1alpha1listers.ClusterServiceVersionLister
	ipLister               v1alpha1listers.InstallPlanLister
	client                 versioned.Interface
	kubeclient             kubernetes.Interface
	globalCatalogNamespace string
	satResolver            *SatResolver
	updatedResolution      bool
}

type debugWriter struct {
	logrus.FieldLogger
}

func NewOperatorsV1alpha1Resolver(lister operatorlister.OperatorLister, client versioned.Interface, kubeclient kubernetes.Interface, globalCatalogNamespace string, updatedResolution bool) *OperatorsV1alpha1Resolver {
	return &OperatorsV1alpha1Resolver{
		subLister:              lister.OperatorsV1alpha1().SubscriptionLister(),
		csvLister:              lister.OperatorsV1alpha1().ClusterServiceVersionLister(),
		ipLister:               lister.OperatorsV1alpha1().InstallPlanLister(),
		client:                 client,
		kubeclient:             kubeclient,
		globalCatalogNamespace: globalCatalogNamespace,
		satResolver:            NewDefaultSatResolver(NewDefaultRegistryClientProvider(client)),
		updatedResolution:      updatedResolution,
	}
	w.Debug(b)
	return n, nil
}

func (r *SatResolver) SolveOperators(namespaces []string, csvs []*v1alpha1.ClusterServiceVersion, subs []*v1alpha1.Subscription) (OperatorSet, error) {
	var errs []error

	installables := make([]solver.Installable, 0)
	visited := make(map[OperatorSurface]*BundleInstallable, 0)

	// TODO: better abstraction
	startingCSVs := make(map[string]struct{})

	namespacedCache := r.cache.Namespaced(namespaces...)

	// build constraints for each Subscription
	for _, sub := range subs {
		pkg := sub.Spec.Package
		catalog := registry.CatalogKey{
			Name:      sub.Spec.CatalogSource,
			Namespace: sub.Spec.CatalogSourceNamespace,
		}
		predicates := []OperatorPredicate{InChannel(pkg, sub.Spec.Channel)}

		// find the currently installed operator (if it exists)
		var current *Operator
		for _, csv := range csvs {
			if csv.Name == sub.Status.InstalledCSV {
				op, err := NewOperatorFromV1Alpha1CSV(csv)
				if err != nil {
					return nil, err
				}
				current = op
				break
			}
		}

		channelFilter := []OperatorPredicate{}

		// if we found an existing installed operator, we should filter the channel by operators that can replace it
		if current != nil {
			channelFilter = append(channelFilter, Or(SkipRangeIncludes(*current.Version()), Replaces(current.Identifier())))
		}

		// if no operator is installed and we have a startingCSV, filter for it
		if current == nil && len(sub.Spec.StartingCSV) > 0 {
			channelFilter = append(channelFilter, WithCSVName(sub.Spec.StartingCSV))
			startingCSVs[sub.Spec.StartingCSV] = struct{}{}
		}

		// find operators, in channel order, that can skip from the current version or list the current in "replaces"
		replacementInstallables, err := r.getSubscriptionInstallables(pkg, current, catalog, predicates, channelFilter, namespacedCache, visited)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		for _, repInstallable := range replacementInstallables {
			installables = append(installables, repInstallable)
		}
	}

	// TODO: Consider csvs not attached to subscriptions

	if len(errs) > 0 {
		return nil, utilerrors.NewAggregate(errs)
	}
	s, err := solver.New(solver.WithInput(installables), solver.WithTracer(solver.LoggingTracer{&debugWriter{r.log}}))
	if err != nil {
		return nil, err
	}
	solvedInstallables, err := s.Solve(context.TODO())
	if err != nil {
		return nil, err
	}

	// get the set of bundle installables from the result solved installables
	operatorInstallables := make([]BundleInstallable, 0)
	for _, installable := range solvedInstallables {
		if bundleInstallable, ok := installable.(BundleInstallable); ok {
			operatorInstallables = append(operatorInstallables, bundleInstallable)
		}
		if bundleInstallable, ok := installable.(*BundleInstallable); ok {
			operatorInstallables = append(operatorInstallables, *bundleInstallable)
		}
	}

	operators := make(map[string]OperatorSurface, 0)
	for _, installableOperator := range operatorInstallables {
		csvName, channel, catalog, err := installableOperator.BundleSourceInfo()
		if err != nil {
			errs = append(errs, err)
			continue
		}
		op, err := ExactlyOne(namespacedCache.Catalog(catalog).Find(WithCSVName(csvName), WithChannel(channel)))
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if len(installableOperator.Replaces) > 0 {
			op.replaces = installableOperator.Replaces
		}

		// lookup if this installable came from a starting CSV
		if _, ok := startingCSVs[csvName]; ok {
			op.sourceInfo.StartingCSV = csvName
		}

		operators[csvName] = op
	}

	// create a map of operatorsourceinfo (subscription+catalogsource data) to the original subscriptions
	subMap := r.sourceInfoToSubscriptions(subs)
	// get a list of new operators to add to the generation
	add := r.sourceInfoForNewSubscriptions(namespace, subMap)

	var operators OperatorSet
	if !r.updatedResolution {
		operators, err = r.generateOperators(csvs, subs, sourceQuerier, add)
		if err != nil {
			return nil, nil, nil, err
		}
	} else {
		// new dependency resolution
		namespaces := []string{namespace, r.globalCatalogNamespace}
		operators, err = r.satResolver.SolveOperators(namespaces, csvs, subs, add)
		if err != nil {
			return nil, nil, nil, err
		}
	}

	// if there's no error, we were able to satisfy all constraints in the subscription set, so we calculate what
	// changes to persist to the cluster and write them out as `steps`
	steps := []*v1alpha1.Step{}
	updatedSubs := []*v1alpha1.Subscription{}
	bundleLookups := []v1alpha1.BundleLookup{}
	for name, op := range operators {
		_, isAdded := add[*op.SourceInfo()]
		existingSubscription, subExists := subMap[*op.SourceInfo()]

	// all candidates added as options for this constraint
	subInstallable.AddDependency(depIds)

	return installables, nil
}

func (r *SatResolver) getBundleInstallables(catalog registry.CatalogKey, predicates []OperatorPredicate, preferredCatalog registry.CatalogKey, namespacedCache MultiCatalogOperatorFinder, visited map[OperatorSurface]*BundleInstallable) (map[solver.Identifier]struct{}, map[solver.Identifier]*BundleInstallable, error) {
	var errs []error
	installables := make(map[solver.Identifier]*BundleInstallable, 0) // aggregate all of the installables at every depth
	identifiers := make(map[solver.Identifier]struct{}, 0)            // keep track of depth + 1 dependencies

	var finder OperatorFinder = namespacedCache
	if !catalog.IsEmpty() {
		finder = namespacedCache.Catalog(catalog)
	}

	bundleStack := finder.Find(predicates...)
	for _, bundle := range bundleStack {
		// pop from the stack
		bundleStack = bundleStack[:len(bundleStack)-1]

		bundleSource := bundle.SourceInfo()
		if bundleSource == nil {
			err := fmt.Errorf("unable to resolve the source of bundle %s, invalid cache", bundle.Identifier())
			errs = append(errs, err)
			continue
		}

		if b, ok := visited[bundle]; ok {
			installables[b.identifier] = b
			identifiers[b.Identifier()] = struct{}{}
			continue
		}

		bundleInstallable := NewBundleInstallable(bundle.Identifier(), bundle.bundle.ChannelName, bundleSource.Catalog)
		visited[bundle] = &bundleInstallable

		dependencyPredicates, err := bundle.DependencyPredicates()
		if err != nil {
			errs = append(errs, err)
			continue
		}
		for _, d := range dependencyPredicates {
			candidateBundles, err := AtLeast(1, namespacedCache.Find(d))
			if err != nil {
				// If there are no candidates for a dependency, it means this bundle can't be resolved
				bundleInstallable.MakeProhibited()
				continue
			}

			bundleDependencies := make(map[solver.Identifier]struct{}, 0)
			for _, dep := range candidateBundles {
				// TODO: search in preferred catalog
				candidateBundles := finder.Find(WithCSVName(dep.Identifier()))

				sortedCandidates := r.sortByVersion(candidateBundles)

				for _, b := range sortedCandidates {
					src := b.SourceInfo()
					if src == nil {
						err := fmt.Errorf("unable to resolve the source of bundle %s, invalid cache", bundle.Identifier())
						errs = append(errs, err)
						continue
					}

					i := NewBundleInstallable(b.Identifier(), b.bundle.ChannelName, bundleSource.Catalog)
					installables[i.Identifier()] = &i
					bundleDependencies[i.Identifier()] = struct{}{}
					bundleStack = append(bundleStack, b)
				}
			}

			// TODO: IMPORTANT: current a solver bug will skip later dependency clauses
			bundleInstallable.AddDependencyFromSet(bundleDependencies)
		}

		installables[bundleInstallable.Identifier()] = &bundleInstallable
		identifiers[bundleInstallable.Identifier()] = struct{}{}
	}

	if len(errs) > 0 {
		return nil, nil, utilerrors.NewAggregate(errs)
	}

	return identifiers, installables, nil
}

func (r *SatResolver) sortByVersion(bundles []*Operator) []*Operator {
	versionMap := make(map[string]*Operator, 0)
	versionSlice := make([]semver.Version, 0)
	unsortableList := make([]*Operator, 0)

	zeroVersion, _ := semver.Make("")

	for _, bundle := range bundles {
		version := bundle.Version() // initialized to zero value if not set in CSV
		if version.Equals(zeroVersion) {
			unsortableList = append(unsortableList, bundle)
			continue
		}

		versionMap[version.String()] = bundle
		versionSlice = append(versionSlice, *version)
	}

	semver.Sort(versionSlice)

	// todo: if len(versionSlice == 0) then try to build the graph and sort that way

	sortedBundles := make([]*Operator, 0)
	for _, sortedVersion := range versionSlice {
		sortedBundles = append(sortedBundles, versionMap[sortedVersion.String()])
	}
	for _, unsortable := range unsortableList {
		sortedBundles = append(sortedBundles, unsortable)
	}

	return sortedBundles
}

func (r *OperatorsV1alpha1Resolver) generateOperators(csvs []*v1alpha1.ClusterServiceVersion, subs []*v1alpha1.Subscription, sourceQuerier SourceQuerier, add map[OperatorSourceInfo]struct{}) (OperatorSet, error) {
	gen, err := NewGenerationFromCluster(csvs, subs)
	if err != nil {
		return nil, err
	}

	// evolve a generation by resolving the set of subscriptions (in `add`) by querying with `source`
	// and taking the current generation (in `gen`) into account
	if err := NewNamespaceGenerationEvolver(sourceQuerier, gen).Evolve(add); err != nil {
		return nil, err
	}

	return gen.Operators(), nil
}

func (r *OperatorsV1alpha1Resolver) sourceInfoForNewSubscriptions(namespace string, subs map[OperatorSourceInfo]*v1alpha1.Subscription) (add map[OperatorSourceInfo]struct{}) {
	add = make(map[OperatorSourceInfo]struct{})
	for key, sub := range subs {
		if sub.Status.CurrentCSV == "" {
			add[key] = struct{}{}
			continue
		}
		if r, ok := bundleLookup[b.replaces]; ok {
			replacedBy[r] = b
			replaces[b] = r
		}
	}

	// a bundle without a replacement is a channel head, but if we find more than one of those something is weird
	headCandidates := []*Operator{}
	for _, b := range bundles {
		if _, ok := replacedBy[b]; !ok {
			headCandidates = append(headCandidates, b)
		}
	}

	if len(headCandidates) != 1 {
		// TODO: more context in error
		return nil, fmt.Errorf("found more than one head for channel")
	}

	head := headCandidates[0]
	current := head
	for {
		channel = append(channel, current)
		next, ok := replaces[current]
		if !ok {
			break
		}
		current = next
	}

	// TODO: do we care if the channel doesn't include every bundle in the input?

	return channel, nil
}
