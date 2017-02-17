package resources

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Machiel/slugify"
	"github.com/generationtux/brizo/database"
	"github.com/generationtux/brizo/kube"
	"github.com/jinzhu/gorm"
	"github.com/pborman/uuid"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/apis/extensions/v1beta1"
	metav1 "k8s.io/client-go/pkg/apis/meta/v1"
)

// Volume as defined by Brizo.
type Volume struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// ContainerPort exposed container port
type ContainerPort struct {
	Protocol string `json:"protocol"`
	Port     int    `json:"port"`
}

// ContainerVolumeMount mount configuration for an available volume
type ContainerVolumeMount struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

// Version as defined by Brizo.
type Version struct {
	database.Model
	UUID          string                 `gorm:"not null;unique_index" sql:"type:varchar(36)" json:"uuid"`
	Name          string                 `gorm:"not null" json:"name"`
	Slug          string                 `gorm:"not null" json:"slug"`
	Image         string                 `gorm:"not null" json:"image"`
	Replicas      int                    `gorm:"not null" sql:"DEFAULT:'0'" json:"replicas"`
	EnvironmentID uint                   `gorm:"not null" json:"environment_id"`
	Environment   Environment            `gorm:"not null" json:"environment"`
	Volumes       []Volume               `gorm:"-" json:"volumes"`
	VolumeMounts  []ContainerVolumeMount `gorm:"-" json:"volumeMounts"`
	PullPolicy    string                 `gorm:"-" json:"pullPolicy"`
	Args          []string               `gorm:"-" json:"args"`
	Ports         []ContainerPort        `gorm:"-" json:"port"`
	Spec          string                 `gorm:"type:json" json:"-"`
}

// BeforeCreate is a hook that runs before inserting a new record into the
// database
func (version *Version) BeforeCreate() (err error) {
	if version.UUID == "" {
		version.UUID = uuid.New()
	}

	return
}

// AllVersions will return all of the Versions
func AllVersions(db *gorm.DB) ([]Version, error) {
	var versions []Version
	result := db.Find(&versions)

	return versions, result.Error
}

// CreateVersion will deploy a new version Brizo
func CreateVersion(db *gorm.DB, client kube.APIInterface, version *Version) (bool, error) {
	deployment := versionDeploymentDefinition(version)

	err := client.CreateDeployment(deployment)
	if err != nil {
		return false, err
	}

	spec, err := json.Marshal(deployment)
	if err != nil {
		client.DeleteDeployment(deployment)
		return false, err
	}

	version.Spec = string(spec)
	persist := db.Create(&version)

	return persist.RowsAffected == 1, persist.Error
}

func convertVolume(volume Volume) v1.Volume {
	var k8sVol v1.Volume

	switch volume.Type {
	case "temp":
		k8sVol = v1.Volume{
			Name: volume.Name,
			VolumeSource: v1.VolumeSource{
				EmptyDir: &v1.EmptyDirVolumeSource{
					Medium: v1.StorageMediumDefault,
				},
			},
		}
	default:
		// @TODO handle unsupported volume error
		panic(volume.Type + " is not a supported volume type")
	}

	return k8sVol
}

// versionDeploymentDefinition builds a deployment spec for the provided version
func versionDeploymentDefinition(version *Version, configs []EnvironmentConfig) *v1beta1.Deployment {
	replicas := int32(version.Replicas)
	name := fmt.Sprintf(
		"%v-%v",
		version.Environment.Application.Slug,
		version.Environment.Slug,
	)

	var policy v1.PullPolicy
	switch version.PullPolicy {
	case "Always":
		policy = v1.PullAlways
	case "IfNotPresent":
		policy = v1.PullIfNotPresent
	case "Never":
		policy = v1.PullNever
	default:
		policy = v1.PullAlways
	}

	var k8sVolumes []v1.Volume
	for index := 0; index < len(version.Volumes); index++ {
		k8sVolumes = append(k8sVolumes, convertVolume(version.Volumes[index]))
	}

	k8sPorts := make([]v1.ContainerPort, len(version.Ports))
	for i, port := range version.Ports {
		protocol := v1.ProtocolTCP
		if port.Protocol == "UDP" {
			protocol = v1.ProtocolUDP
		}
		k8sPorts[i] = v1.ContainerPort{
			Protocol:      protocol,
			ContainerPort: int32(port.Port),
		}
	}

	k8sVolumeMounts := make([]v1.VolumeMount, len(version.VolumeMounts))
	for i, mount := range version.VolumeMounts {
		k8sVolumeMounts[i] = v1.VolumeMount{
			Name:      mount.Name,
			MountPath: mount.Path,
		}
	}

	deployment := &v1beta1.Deployment{
		ObjectMeta: v1.ObjectMeta{
			Name:      name,
			Namespace: "brizo",
			Labels: map[string]string{
				"brizoManaged": "true",
				"appUUID":      version.Environment.Application.UUID,
				"envUUID":      version.Environment.UUID,
			},
		},
		Spec: v1beta1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"envUUID":     version.Environment.UUID,
					"versionUUID": version.UUID,
				},
			},
			Replicas: &replicas,
			Template: v1.PodTemplateSpec{
				ObjectMeta: v1.ObjectMeta{
					Labels: map[string]string{
						"brizoManaged": "true",
						"appUUID":      version.Environment.Application.UUID,
						"envUUID":      version.Environment.UUID,
						"versionUUID":  version.UUID,
					},
				},
				Spec: v1.PodSpec{
					Volumes: k8sVolumes,
					Containers: []v1.Container{
						v1.Container{
							Name:            "app",
							Image:           version.Image,
							ImagePullPolicy: policy,
							Args:            version.Args,
							Ports:           k8sPorts,
							VolumeMounts:    k8sVolumeMounts,
						},
					},
				},
			},
		},
	}

	v1beta1.SetObjectDefaults_Deployment(deployment)
	return deployment
}

// UpdateVersion will update an existing Version
func UpdateVersion(db *gorm.DB, version *Version) (bool, error) {
	version.Slug = slugify.Slugify(version.Name)
	result := db.Model(version).Where("uuid = ?", version.UUID).
		UpdateColumns(version)

	return result.RowsAffected == 1, result.Error
}

// GetVersion will get an existing Version by id
func GetVersion(db *gorm.DB, id string) (*Version, error) {
	version := new(Version)
	if err := db.Preload("Environment.Application").Where("uuid = ?", id).First(version).Error; err != nil {
		return version, err
	}

	if version.ID == 0 {
		return new(Version), errors.New("not-found")
	}

	return version, nil
}

// GetVersionsByEnvironmentUUID will get an existing version using an
// environment's UUID
func GetVersionsByEnvironmentUUID(db *gorm.DB, uuid string) (*[]Version, error) {
	var versions []Version
	environment, err := GetEnvironment(db, uuid)
	if err != nil {
		return &versions, err
	}

	if err := db.Preload("Environment.Application").Where("environment_id = ?", environment.ID).Find(&versions).Error; err != nil {
		return &versions, err
	}

	return &versions, nil
}

// DeleteVersion will delete an existing Version by name
func DeleteVersion(db *gorm.DB, name string) (bool, error) {
	result := db.Delete(Version{}, "name = ?", name)

	return result.RowsAffected == 1, result.Error
}
