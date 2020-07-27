package agones_shutdown

import (
	"fmt"
	"log"
	"os"
	"testing"

	"k8s.io/apimachinery/pkg/types"

	"github.com/google/uuid"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"

	agonesv1 "agones.dev/agones/pkg/apis/agones/v1"
	e2eframework "agones.dev/agones/test/e2e/framework"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	namespacePrefix = "gameserver-test"
)

var (
	framework *e2eframework.Framework
)

func TestMain(m *testing.M) {
	fw, err := e2eframework.NewFromFlags()
	if err != nil {
		log.Fatalf("failed to init e2e e2eframework: %+v", err)
	}
	framework = fw

	var exitCode int
	defer func() {
		os.Exit(exitCode)
	}()
	exitCode = m.Run()
}

func newFleet(name, namespace string, replicas int, gsSpec *agonesv1.GameServerSpec) *agonesv1.Fleet {
	return &agonesv1.Fleet{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace},
		Spec: agonesv1.FleetSpec{
			Replicas: int32(replicas),
			Template: agonesv1.GameServerTemplateSpec{
				Spec: *gsSpec,
			},
		},
	}
}

func createFleet(name, namespace string, replicas int, gsSpec *agonesv1.GameServerSpec) *agonesv1.Fleet {
	flt, err := framework.AgonesClient.AgonesV1().
		Fleets(namespace).
		Create(newFleet(name, namespace, replicas, gsSpec))
	if err != nil {
		panic(err)
	}
	return flt
}

func createPodSpec(name, image string) corev1.PodSpec {
	return corev1.PodSpec{
		Containers: []corev1.Container{{
			Name:            name,
			Image:           image,
			ImagePullPolicy: corev1.PullNever, // using local image
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceCPU:    resource.MustParse("30m"),
					corev1.ResourceMemory: resource.MustParse("32Mi"),
				},
				Limits: corev1.ResourceList{
					corev1.ResourceCPU:    resource.MustParse("30m"),
					corev1.ResourceMemory: resource.MustParse("32Mi"),
				},
			},
		}},
	}
}

func newRandomNamespace(t *testing.T) string {
	namespace := fmt.Sprintf("%s-%s", namespacePrefix, uuid.Must(uuid.NewRandom()))
	if err := framework.CreateNamespace(namespace); err != nil {
		panic(err)
	}
	t.Cleanup(func() {
		_ = framework.DeleteNamespace(namespace)
	})
	return namespace
}

func updateContainerImage(flt *agonesv1.Fleet, image string) error {
	patch := fmt.Sprintf(`[{"op": "replace", "path": "/spec/template/spec/template/spec/containers/0/image", "value": "%s"}]`, "gameserver:v2")
	if _, err := framework.AgonesClient.AgonesV1().Fleets(flt.Namespace).Patch(flt.Name, types.JSONPatchType, []byte(patch)); err != nil {
		return err
	}
	return nil
}
