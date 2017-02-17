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

// Version as defined by Brizo.
type Version struct {
	database.Model
	UUID          string      `gorm:"not null;unique_index" sql:"type:varchar(36)" json:"uuid"`
	Name          string      `gorm:"not null" json:"name"`
	Slug          string      `gorm:"not null" json:"slug"`
	Image         string      `gorm:"not null" json:"image"`
	Replicas      int         `gorm:"not null" sql:"DEFAULT:'0'" json:"replicas"`
	EnvironmentID uint        `gorm:"not null" json:"environment_id"`
	Environment   Environment `gorm:"not null" json:"environment"`
	Spec          string      `gorm:"type:json" json:"-"`
}

// BeforeCreate is a hook that runs before inserting a new record into the
// database
func (version *Version) BeforeCreate() (err error) {
	if version.UUID == "" {
		version.UUID = uuid.New()
	}
	if version.Slug == "" {
		version.Slug = slugify.Slugify(version.Name)
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
func CreateVersion(db *gorm.DB, client kube.APIInterface, version *Version, config *[]EnvironmentConfig) (bool, error) {
	deployment := versionDeploymentDefinition(version, *config)
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

func buildEnvironmentConfig(config EnvironmentConfig) v1.EnvVar {
	return v1.EnvVar{
		Name:  config.Name,
		Value: config.Value,
	}
}

// versionDeploymentDefinition builds a deployment spec for the provided version
func versionDeploymentDefinition(version *Version, configs []EnvironmentConfig) *v1beta1.Deployment {
	replicas := int32(version.Replicas)
	name := fmt.Sprintf(
		"%v-%v-%v",
		version.Environment.Application.Slug,
		version.Environment.Slug,
		version.Slug,
	)

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
					Containers: []v1.Container{
						v1.Container{
							Name:  "app",
							Image: version.Image,
							Env:   envVars,
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
