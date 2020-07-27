package agones_shutdown

import (
	"testing"
	"time"

	agonesv1 "agones.dev/agones/pkg/apis/agones/v1"
	corev1 "k8s.io/api/core/v1"
)

func newGameServerSpec(name, image string) agonesv1.GameServerSpec {
	return agonesv1.GameServerSpec{
		Ports: []agonesv1.GameServerPort{
			{
				ContainerPort: 7000,
				Name:          "tcp",
				PortPolicy:    agonesv1.Dynamic,
				Protocol:      corev1.ProtocolTCP,
			}},
		Template: corev1.PodTemplateSpec{
			Spec: createPodSpec(name, image),
		},
	}
}

func TestDeployNewVersion(t *testing.T) {
	namespace := newRandomNamespace(t)
	v1spec := newGameServerSpec("gameserver", "gameserver:v1")
	replicas := 10
	flt := createFleet("simple-fleet", namespace, replicas, &v1spec)
	framework.AssertFleetCondition(t, flt, func(flt *agonesv1.Fleet) bool {
		return flt.Status.ReadyReplicas == int32(replicas)
	})

	// deploy new version image
	if err := updateContainerImage(flt, "gameserver:v2"); err != nil {
		t.Fatal(err)
	}
	time.Sleep(10 * time.Second)
}
