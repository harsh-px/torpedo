package postgres

import (
	"fmt"

	"github.com/portworx/torpedo/drivers/scheduler/k8s/spec/factory"
	"k8s.io/apimachinery/pkg/api/resource"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/apis/apps/v1beta1"
	storage_v1beta1 "k8s.io/client-go/pkg/apis/storage/v1beta1"
)

type postgres struct{}

func (p *postgres) ID() string {
	return "postgres"
}

func (p *postgres) Core(name string) []interface{} {
	var coreComponents []interface{}
	deployment := &v1beta1.Deployment{
		ObjectMeta: meta_v1.ObjectMeta{
			Name: name,
			Namespace: v1.NamespaceDefault,
		},
		Spec: v1beta1.DeploymentSpec{
			Strategy: v1beta1.DeploymentStrategy{
				Type: v1beta1.RollingUpdateDeploymentStrategyType,
				RollingUpdate: &v1beta1.RollingUpdateDeployment{
					MaxSurge: &intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: 1,
					},
					MaxUnavailable: &intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: 1,
					},
				},
			},
			Replicas: int32Ptr(1),
			Template: v1.PodTemplateSpec{
				ObjectMeta: meta_v1.ObjectMeta{
					Labels: map[string]string{
						"app": name,
					},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:            "postgres",
							Image:           "postgres:9.5",
							ImagePullPolicy: v1.PullIfNotPresent,
							Ports: []v1.ContainerPort{
								{
									Name:          "http",
									Protocol:      v1.ProtocolTCP,
									ContainerPort: 5432,
								},
							},
							Env: []v1.EnvVar{
								{
									Name:  "POSTGRES_USER",
									Value: "pgbench",
								},
								{
									Name:  "POSTGRES_PASSWORD",
									Value: "superpostgres",
								},
								{
									Name:  "PGBENCH_PASSWORD",
									Value: "superpostgres",
								},
								{
									Name:  "PGDATA",
									Value: "/var/lib/postgresql/data/pgdata",
								},
							},
							VolumeMounts: []v1.VolumeMount{
								{
									MountPath: "/var/lib/postgresql/data",
									Name:      "postgredb",
								},
							},
						},
					},
					Volumes: []v1.Volume{
						{
							Name: "postgredb",
							VolumeSource: v1.VolumeSource{
								PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
									ClaimName: "postgredb",
								},
							},
						},
					},
				},
			},
		},
	}

	coreComponents = append(coreComponents, deployment)
	return coreComponents
}

func (p *postgres) Storage() []interface{} {
	var storageComponents []interface{}

	sc := &storage_v1beta1.StorageClass{
		ObjectMeta: meta_v1.ObjectMeta{
			Name: "postgres-sc-pxd",
		},
		Provisioner: "kubernetes.io/portworx-volume",
		Parameters: map[string]string{
			"repl":       "3",
			"io_profile": "db",
		},
	}

	storageComponents = append(storageComponents, sc)

	pvc := &v1.PersistentVolumeClaim{
		ObjectMeta: meta_v1.ObjectMeta{
			Name: "postgredb",
			Annotations: map[string]string{
				"volume.beta.kubernetes.io/storage-class": "postgres-sc-pxd",
			},
			Namespace: v1.NamespaceDefault,
		},
		Spec: v1.PersistentVolumeClaimSpec{
			AccessModes: []v1.PersistentVolumeAccessMode{
				v1.ReadWriteOnce,
			},
			Resources: v1.ResourceRequirements{
				Requests: v1.ResourceList{
					v1.ResourceName(v1.ResourceStorage): resource.MustParse(fmt.Sprintf("%dGi", 2)),
				},
			},
		},
	}
	storageComponents = append(storageComponents, pvc)

	return storageComponents
}

func int32Ptr(i int32) *int32 { return &i }

func init() {
	p := &postgres{}
	factory.Register(p.ID(), p)
}
