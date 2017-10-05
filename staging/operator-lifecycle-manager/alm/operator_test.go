package alm

import (
	"testing"

	"github.com/golang/mock/gomock"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"

	"github.com/coreos-inc/alm/client"

	"github.com/coreos-inc/alm/apis/clusterserviceversion/v1alpha1"
	"github.com/coreos-inc/alm/queueinformer"
)

type MockListWatcher struct {
}

func (l *MockListWatcher) List(options v1.ListOptions) (runtime.Object, error) {
	return nil, nil
}

func (l *MockListWatcher) Watch(options v1.ListOptions) (watch.Interface, error) {
	return nil, nil
}

type MockALMOperator struct {
	ALMOperator
	MockCSVClient     *client.MockClusterServiceVersionInterface
	TestQueueInformer queueinformer.TestQueueInformer
}

func NewMockALMOperator(gomockCtrl *gomock.Controller) *MockALMOperator {
	mockCSVClient := client.NewMockClusterServiceVersionInterface(gomockCtrl)

	almOperator := ALMOperator{
		csvClient: mockCSVClient,
	}

	csvQueueInformer := queueinformer.NewTestQueueInformer(
		"test-clusterserviceversions",
		cache.NewSharedIndexInformer(&MockListWatcher{}, &v1alpha1.ClusterServiceVersion{}, 0, nil),
		almOperator.syncClusterServiceVersion,
		nil,
	)

	qOp := queueinformer.NewMockOperator(gomockCtrl, csvQueueInformer)
	almOperator.Operator = &qOp.Operator

	return &MockALMOperator{
		ALMOperator:       almOperator,
		MockCSVClient:     mockCSVClient,
		TestQueueInformer: *csvQueueInformer,
	}
}

func TestTransitionNoneToPending(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockOp := NewMockALMOperator(ctrl)

	csv := v1alpha1.ClusterServiceVersion{
		ObjectMeta: v1.ObjectMeta{
			Name:     "test-csv",
			SelfLink: "/link/test-csv",
		},
		Spec: v1alpha1.ClusterServiceVersionSpec{
			DisplayName: "Test",
		},
	}

	mockOp.MockCSVClient.EXPECT().
		TransitionPhase(&csv, v1alpha1.CSVPhasePending, v1alpha1.CSVReasonRequirementsUnkown, "requirements not yet checked").
		Return(&csv, nil)
	mockOp.syncClusterServiceVersion(&csv)
}
