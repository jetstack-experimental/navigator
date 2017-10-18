package cassandra

import (
	"github.com/golang/glog"
	v1alpha1 "github.com/jetstack-experimental/navigator/pkg/apis/navigator/v1alpha1"
	navigatorclientset "github.com/jetstack-experimental/navigator/pkg/client/clientset_generated/clientset"
	"k8s.io/client-go/kubernetes"
)

const (
	errorSync = "ErrSync"

	successSync = "SuccessSync"

	messageErrorSyncServiceAccount = "Error syncing service account: %s"
	messageErrorSyncConfigMap      = "Error syncing config map: %s"
	messageErrorSyncService        = "Error syncing service: %s"
	messageErrorSyncNodePools      = "Error syncing node pools: %s"
	messageSuccessSync             = "Successfully synced CassandraCluster"
)

type ControlInterface interface {
	Sync(*v1alpha1.CassandraCluster) error
}

var _ ControlInterface = &defaultCassandraClusterControl{}

type defaultCassandraClusterControl struct {
	navigatorClient navigatorclientset.Interface
	kubeClient      kubernetes.Interface
}

func NewController(
	navigatorClient navigatorclientset.Interface,
	kubeClient kubernetes.Interface,
) ControlInterface {
	return &defaultCassandraClusterControl{
		navigatorClient: navigatorClient,
		kubeClient:      kubeClient,
	}
}

func (e *defaultCassandraClusterControl) Sync(c *v1alpha1.CassandraCluster) error {
	c = c.DeepCopy()
	glog.V(4).Infof("defaultCassandraClusterControl.Sync")
	_, err := e.navigatorClient.
		NavigatorV1alpha1().
		CassandraClusters(c.Namespace).
		UpdateStatus(c)
	return err
}