package actions_test

import (
	"testing"

	"github.com/jetstack/navigator/pkg/util/ptr"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/jetstack/navigator/internal/test/unit/framework"
	"github.com/jetstack/navigator/internal/test/util/generate"
	"github.com/jetstack/navigator/pkg/controllers/cassandra/actions"
)

func TestScaleIn(t *testing.T) {
	type testT struct {
		kubeObjects         []runtime.Object
		navObjects          []runtime.Object
		cluster             generate.CassandraClusterConfig
		nodePool            generate.CassandraClusterNodePoolConfig
		expectedStatefulSet *generate.StatefulSetConfig
		expectedErr         bool
		mutator             func(*framework.StateFixture)
	}
	tests := map[string]testT{
		"Error if StatefulSet not listed": {
			cluster: generate.CassandraClusterConfig{
				Name:      "cluster1",
				Namespace: "ns1",
			},
			nodePool: generate.CassandraClusterNodePoolConfig{
				Name:     "pool1",
				Replicas: 123,
			},
			expectedErr: true,
		},
		"Error if clientset.Update fails (e.g. listed but not found)": {
			kubeObjects: []runtime.Object{
				generate.StatefulSet(
					generate.StatefulSetConfig{
						Name:      "cass-cluster1-pool1",
						Namespace: "ns1",
						Replicas:  ptr.Int32(122),
					},
				),
			},
			cluster: generate.CassandraClusterConfig{
				Name:      "cluster1",
				Namespace: "ns1",
			},
			nodePool: generate.CassandraClusterNodePoolConfig{
				Name:     "pool1",
				Replicas: 123,
			},
			expectedErr: true,
			mutator: func(f *framework.StateFixture) {
				err := f.KubeClient().
					AppsV1beta1().
					StatefulSets("ns1").
					Delete("cass-cluster1-pool1", &metav1.DeleteOptions{})
				if err != nil {
					f.T.Fatal(err)
				}
			},
		},
		"Idempotent: No error if ReplicaCount already matches the actual ReplicaCount": {
			kubeObjects: []runtime.Object{
				generate.StatefulSet(
					generate.StatefulSetConfig{
						Name:      "cass-cluster1-pool1",
						Namespace: "ns1",
						Replicas:  ptr.Int32(124),
					},
				),
			},
			cluster: generate.CassandraClusterConfig{
				Name:      "cluster1",
				Namespace: "ns1",
			},
			nodePool: generate.CassandraClusterNodePoolConfig{
				Name:     "pool1",
				Replicas: 124,
			},
			expectedStatefulSet: &generate.StatefulSetConfig{
				Name:      "cass-cluster1-pool1",
				Namespace: "ns1",
				Replicas:  ptr.Int32(124),
			},
			expectedErr: false,
		},
		"The replicas count is decremented without successfully decommissioned pilots": {
			kubeObjects: []runtime.Object{
				generate.StatefulSet(
					generate.StatefulSetConfig{
						Name:      "cass-cluster1-pool1",
						Namespace: "ns1",
						Replicas:  ptr.Int32(2),
					},
				),
			},
			navObjects: []runtime.Object{
				generate.CassPilot(
					generate.PilotConfig{
						Name:      "cass-cluster1-pool1-0",
						Namespace: "ns1",
						NodePool:  "pool1",
						Cluster:   "cluster1",
					},
				),
				generate.CassPilot(
					generate.PilotConfig{
						Name:      "cass-cluster1-pool1-1",
						Namespace: "ns1",
						NodePool:  "pool1",
						Cluster:   "cluster1",
					},
				),
			},
			cluster: generate.CassandraClusterConfig{
				Name:      "cluster1",
				Namespace: "ns1",
			},
			nodePool: generate.CassandraClusterNodePoolConfig{
				Name:     "pool1",
				Replicas: 1,
			},
			expectedStatefulSet: &generate.StatefulSetConfig{
				Name:      "cass-cluster1-pool1",
				Namespace: "ns1",
				Replicas:  ptr.Int32(2),
			},
		},
		"The replicas count is decremented with successfully decommissioned pilots": {
			kubeObjects: []runtime.Object{
				generate.StatefulSet(
					generate.StatefulSetConfig{
						Name:      "cass-cluster1-pool1",
						Namespace: "ns1",
						Replicas:  ptr.Int32(2),
					},
				),
			},
			navObjects: []runtime.Object{
				generate.CassPilot(
					generate.PilotConfig{
						Name:      "cass-cluster1-pool1-0",
						Namespace: "ns1",
						NodePool:  "pool1",
						Cluster:   "cluster1",
					},
				),
				generate.CassPilot(
					generate.PilotConfig{
						Name:                 "cass-cluster1-pool1-1",
						Namespace:            "ns1",
						NodePool:             "pool1",
						Cluster:              "cluster1",
						DecommissionedStatus: true,
						Decommissioned:       true,
					},
				),
			},
			cluster: generate.CassandraClusterConfig{
				Name:      "cluster1",
				Namespace: "ns1",
			},
			nodePool: generate.CassandraClusterNodePoolConfig{
				Name:     "pool1",
				Replicas: 1,
			},
			expectedStatefulSet: &generate.StatefulSetConfig{
				Name:      "cass-cluster1-pool1",
				Namespace: "ns1",
				Replicas:  ptr.Int32(1),
			},
		},
	}

	for name, test := range tests {
		t.Run(
			name,
			func(t *testing.T) {
				fixture := &framework.StateFixture{
					T:                t,
					KubeObjects:      test.kubeObjects,
					NavigatorObjects: test.navObjects,
				}
				fixture.Start()
				defer fixture.Stop()
				state := fixture.State()
				if test.mutator != nil {
					test.mutator(fixture)
				}

				a := &actions.ScaleIn{
					Cluster:  generate.CassandraCluster(test.cluster),
					NodePool: generate.CassandraClusterNodePool(test.nodePool),
				}
				err := a.Execute(state)
				if err != nil {
					t.Logf("The error returned by Execute was: %s", err)
				}
				if !test.expectedErr && err != nil {
					t.Errorf("Unexpected error: %s", err)
				}
				if test.expectedErr && err == nil {
					t.Errorf("Expected an error")
				}
				if test.expectedStatefulSet != nil {
					actualStatefulSet, err := fixture.KubeClient().
						AppsV1beta1().
						StatefulSets(test.expectedStatefulSet.Namespace).
						Get(test.expectedStatefulSet.Name, metav1.GetOptions{})
					if err != nil {
						t.Fatalf("Unexpected error retrieving statefulset: %v", err)
					}
					if err != nil {
						t.Fatalf("Unexpected error retrieving statefulset: %v", err)
					}
					if *test.expectedStatefulSet.Replicas != *actualStatefulSet.Spec.Replicas {
						t.Errorf(
							"Unexpected replica count. Expected: %d. Actual: %d",
							*test.expectedStatefulSet.Replicas, *actualStatefulSet.Spec.Replicas,
						)
					}
				}
			},
		)
	}
}
