package agones_shutdown

import (
	"testing"

	agonesv1 "agones.dev/agones/pkg/apis/agones/v1"
	corev1 "k8s.io/api/core/v1"
)

func TestDeployNewVersion(t *testing.T) {
	namespace := newRandomNamespace(t)

	gsSpec := &agonesv1.GameServerSpec{
		Ports: []agonesv1.GameServerPort{
			{
				ContainerPort: 7000,
				Name:          "tcp",
				PortPolicy:    agonesv1.Dynamic,
				Protocol:      corev1.ProtocolTCP,
			}},
		Template: corev1.PodTemplateSpec{
			Spec: createPodSpec("gameserver", "gameserver:v1"),
		},
	}
	flt := createFleet("simple-fleet", namespace, gsSpec)
	framework.AssertFleetCondition(t, flt, func(flt *agonesv1.Fleet) bool {
		return flt.Status.ReadyReplicas == 1
	})
}
