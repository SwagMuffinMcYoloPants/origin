package etcd

import (
	"strconv"

	etcderr "github.com/GoogleCloudPlatform/kubernetes/pkg/api/errors/etcd"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/labels"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/runtime"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/tools"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/watch"
	"github.com/golang/glog"
	"github.com/openshift/origin/pkg/deploy/api"
)

// Etcd implements deployment.Registry and deploymentconfig.Registry interfaces.
type Etcd struct {
	tools.EtcdHelper
}

// New creates an etcd registry.
func New(helper tools.EtcdHelper) *Etcd {
	return &Etcd{
		EtcdHelper: helper,
	}
}

// ListDeployments obtains a list of Deployments.
func (r *Etcd) ListDeployments(selector labels.Selector) (*api.DeploymentList, error) {
	deployments := api.DeploymentList{}
	err := r.ExtractToList("/deployments", &deployments)
	if err != nil {
		return nil, err
	}
	filtered := []api.Deployment{}
	for _, item := range deployments.Items {
		if selector.Matches(labels.Set(item.Labels)) {
			filtered = append(filtered, item)
		}
	}

	deployments.Items = filtered
	return &deployments, err
}

func makeDeploymentKey(id string) string {
	return "/deployments/" + id
}

// GetDeployment gets a specific Deployment specified by its ID.
func (r *Etcd) GetDeployment(id string) (*api.Deployment, error) {
	var deployment api.Deployment
	key := makeDeploymentKey(id)
	err := r.ExtractObj(key, &deployment, false)
	if err != nil {
		return nil, etcderr.InterpretGetError(err, "deployment", id)
	}
	return &deployment, nil
}

// CreateDeployment creates a new Deployment.
func (r *Etcd) CreateDeployment(deployment *api.Deployment) error {
	err := r.CreateObj(makeDeploymentKey(deployment.ID), deployment, 0)
	return etcderr.InterpretCreateError(err, "deployment", deployment.ID)
}

// UpdateDeployment replaces an existing Deployment.
func (r *Etcd) UpdateDeployment(deployment *api.Deployment) error {
	err := r.SetObj(makeDeploymentKey(deployment.ID), deployment)
	return etcderr.InterpretUpdateError(err, "deployment", deployment.ID)
}

// DeleteDeployment deletes a Deployment specified by its ID.
func (r *Etcd) DeleteDeployment(id string) error {
	key := makeDeploymentKey(id)
	err := r.Delete(key, false)
	return etcderr.InterpretDeleteError(err, "deployment", id)
}

// WatchDeployments begins watching for new, changed, or deleted Deployments.
func (r *Etcd) WatchDeployments(resourceVersion string, filter func(deployment *api.Deployment) bool) (watch.Interface, error) {
	version, err := parseWatchResourceVersion(resourceVersion, "deployment")
	if err != nil {
		return nil, err
	}

	return r.WatchList("/deployments", version, func(obj runtime.Object) bool {
		deployment, ok := obj.(*api.Deployment)
		if !ok {
			glog.Errorf("Unexpected object during deployment watch: %#v", obj)
			return false
		}
		return filter(deployment)
	})
}

// ListDeploymentConfigs obtains a list of DeploymentConfigs.
func (r *Etcd) ListDeploymentConfigs(selector labels.Selector) (*api.DeploymentConfigList, error) {
	deploymentConfigs := api.DeploymentConfigList{}
	err := r.ExtractToList("/deploymentConfigs", &deploymentConfigs)
	if err != nil {
		return nil, err
	}
	filtered := []api.DeploymentConfig{}
	for _, item := range deploymentConfigs.Items {
		if selector.Matches(labels.Set(item.Labels)) {
			filtered = append(filtered, item)
		}
	}

	deploymentConfigs.Items = filtered
	return &deploymentConfigs, err
}

// TODO expose this from kubernetes.  I will do that, but I don't want this merge stuck on kubernetes refactoring
// parseWatchResourceVersion takes a resource version argument and converts it to
// the etcd version we should pass to helper.Watch(). Because resourceVersion is
// an opaque value, the default watch behavior for non-zero watch is to watch
// the next value (if you pass "1", you will see updates from "2" onwards).
func parseWatchResourceVersion(resourceVersion, kind string) (uint64, error) {
	if resourceVersion == "" || resourceVersion == "0" {
		return 0, nil
	}
	version, err := strconv.ParseUint(resourceVersion, 10, 64)
	if err != nil {
		return 0, etcderr.InterpretResourceVersionError(err, kind, resourceVersion)
	}
	return version + 1, nil
}

// WatchDeploymentConfigs begins watching for new, changed, or deleted DeploymentConfigs.
func (r *Etcd) WatchDeploymentConfigs(resourceVersion string, filter func(repo *api.DeploymentConfig) bool) (watch.Interface, error) {
	version, err := parseWatchResourceVersion(resourceVersion, "deploymentConfig")
	if err != nil {
		return nil, err
	}

	return r.WatchList("/deploymentConfigs", version, func(obj runtime.Object) bool {
		config, ok := obj.(*api.DeploymentConfig)
		if !ok {
			glog.Errorf("Unexpected object during deploymentConfig watch: %#v", obj)
			return false
		}
		return filter(config)
	})
}

func makeDeploymentConfigKey(id string) string {
	return "/deploymentConfigs/" + id
}

// GetDeploymentConfig gets a specific DeploymentConfig specified by its ID.
func (r *Etcd) GetDeploymentConfig(id string) (*api.DeploymentConfig, error) {
	var deploymentConfig api.DeploymentConfig
	key := makeDeploymentConfigKey(id)
	err := r.ExtractObj(key, &deploymentConfig, false)
	if err != nil {
		return nil, etcderr.InterpretGetError(err, "deploymentConfig", id)
	}
	return &deploymentConfig, nil
}

// CreateDeploymentConfig creates a new DeploymentConfig.
func (r *Etcd) CreateDeploymentConfig(deploymentConfig *api.DeploymentConfig) error {
	err := r.CreateObj(makeDeploymentConfigKey(deploymentConfig.ID), deploymentConfig, 0)
	return etcderr.InterpretCreateError(err, "deploymentConfig", deploymentConfig.ID)
}

// UpdateDeploymentConfig replaces an existing DeploymentConfig.
func (r *Etcd) UpdateDeploymentConfig(deploymentConfig *api.DeploymentConfig) error {
	err := r.SetObj(makeDeploymentConfigKey(deploymentConfig.ID), deploymentConfig)
	return etcderr.InterpretUpdateError(err, "deploymentConfig", deploymentConfig.ID)
}

// DeleteDeploymentConfig deletes a DeploymentConfig specified by its ID.
func (r *Etcd) DeleteDeploymentConfig(id string) error {
	key := makeDeploymentConfigKey(id)
	err := r.Delete(key, false)
	return etcderr.InterpretDeleteError(err, "deploymentConfig", id)
}
