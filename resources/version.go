package resources

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

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

// Container individual container for a version
type Container struct {
	Name         string                 `json:"name"`
	Image        string                 `json:"image"`
	AlwaysPull   bool                   `json:"alwaysPull"`
	Args         []string               `json:"args"`
	VolumeMounts []ContainerVolumeMount `json:"volumeMounts"`
	Ports        []ContainerPort        `json:"ports"`
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
	UUID          string      `gorm:"not null;unique_index" sql:"type:varchar(36)" json:"uuid"`
	Name          string      `gorm:"not null" json:"name"`
	Slug          string      `gorm:"not null" json:"slug"`
	Replicas      int         `gorm:"not null" sql:"DEFAULT:'0'" json:"replicas"`
	EnvironmentID uint        `gorm:"not null" json:"environment_id"`
	Environment   Environment `gorm:"not null" json:"environment"`
	Volumes       []Volume    `gorm:"-" json:"volumes"`
	Containers    []Container `gorm:"-" json:"containers"`
	Spec          string      `gorm:"type:json" json:"-"`
}

type Spec struct {
	Name              string
	Labels            []string
	Namespace         string
	creationTimestamp string
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

	err := client.CreateOrUpdateDeployment(deployment)
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

	// update environment service
	ports := gatherContainerPorts(version.Containers)
	UpdateEnvironmentService(db, client, &version.Environment, ports)

	return persist.RowsAffected == 1, persist.Error
}

func gatherContainerPorts(containers []Container) []ContainerPort {
	ports := make([]ContainerPort, 0)
	for _, container := range containers {
		ports = append(ports, container.Ports...)
	}
	return ports
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

func createContainerSpec(container Container, environmentUUID string) v1.Container {
	policy := v1.PullAlways
	if !container.AlwaysPull {
		policy = v1.PullIfNotPresent
	}

	k8sPorts := make([]v1.ContainerPort, len(container.Ports))
	for i, port := range container.Ports {
		protocol := v1.ProtocolTCP
		if port.Protocol == "UDP" {
			protocol = v1.ProtocolUDP
		}
		k8sPorts[i] = v1.ContainerPort{
			Protocol:      protocol,
			ContainerPort: int32(port.Port),
		}
	}

	k8sVolumeMounts := make([]v1.VolumeMount, len(container.VolumeMounts))
	for i, mount := range container.VolumeMounts {
		k8sVolumeMounts[i] = v1.VolumeMount{
			Name:      mount.Name,
			MountPath: mount.Path,
		}
	}

	db, err := database.Connect()
	defer db.Close()
	if err != nil {
		log.Printf("Database error: '%s'\n", err)
	}
	environmentVars, err := GetEnvironmentConfig(db, environmentUUID)
	if err != nil {
		log.Printf("Error retrieving environment configs: '%s'\n", err)
	}

	var k8sEnvVars []v1.EnvVar
	for _, environmentVar := range *environmentVars {
		k8sEnvVars = append(k8sEnvVars, v1.EnvVar{
			Name:  environmentVar.Name,
			Value: environmentVar.Value,
		})
	}

	return v1.Container{
		Name:            container.Name,
		Image:           container.Image,
		ImagePullPolicy: policy,
		Args:            container.Args,
		Ports:           k8sPorts,
		VolumeMounts:    k8sVolumeMounts,
		Env:             k8sEnvVars,
	}
}

// versionDeploymentDefinition builds a deployment spec for the provided version
func versionDeploymentDefinition(version *Version) *v1beta1.Deployment {
	replicas := int32(version.Replicas)
	name := fmt.Sprintf(
		"%v-%v",
		version.Environment.Application.Slug,
		version.Environment.Slug,
	)

	var k8sVolumes []v1.Volume
	for index := 0; index < len(version.Volumes); index++ {
		k8sVolumes = append(k8sVolumes, convertVolume(version.Volumes[index]))
	}

	var k8sContainers []v1.Container
	for index := 0; index < len(version.Containers); index++ {
		k8sContainers = append(k8sContainers, createContainerSpec(version.Containers[index], version.Environment.UUID))
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
					Volumes:    k8sVolumes,
					Containers: k8sContainers,
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
func GetVersion(db *gorm.DB, id string, client *kube.Client) (*Version, error) {
	version := new(Version)
	if err := db.Preload("Environment.Application").Where("uuid = ?", id).First(version).Error; err != nil {
		return version, err
	}

	if version.ID == 0 {
		return new(Version), errors.New("not-found")
	}

	var spec map[string]map[string]interface{}
	err := json.Unmarshal([]byte(version.Spec), &spec)
	if err != nil {
		return new(Version), err
	}

	specName := spec["metadata"]["name"].(string)
	specNS := spec["metadata"]["namespace"].(string)

	deployment, err := client.FindDeploymentByName(specName, specNS)
	if err != nil {
		return new(Version), err
	}

	// gather container information
	for i := 0; i < len(deployment.Spec.Template.Spec.Containers); i++ {
		// determine pull policy
		pullPolicy := true
		if deployment.Spec.Template.Spec.Containers[i].ImagePullPolicy != "Always" {
			pullPolicy = false
		}

		container := Container{
			Name:       deployment.Spec.Template.Spec.Containers[i].Name,
			Image:      deployment.Spec.Template.Spec.Containers[i].Image,
			AlwaysPull: pullPolicy,
			Args:       deployment.Spec.Template.Spec.Containers[i].Args,
		}
		version.Containers = append(version.Containers, container)
	}

	// gather volumes information
	for i := 0; i < len(deployment.Spec.Template.Spec.Volumes); i++ {
		volume := Volume{
			Name: deployment.Spec.Template.Spec.Volumes[i].Name,
		}
		version.Volumes = append(version.Volumes, volume)
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
