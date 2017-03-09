package resources

import (
	"errors"
	"fmt"

	"k8s.io/client-go/pkg/api/v1"

	"github.com/Machiel/slugify"
	"github.com/generationtux/brizo/database"
	"github.com/generationtux/brizo/kube"
	"github.com/jinzhu/gorm"
	"github.com/pborman/uuid"
)

// Environment as defined by Brizo.
type Environment struct {
	database.Model
	UUID          string      `gorm:"not null;unique_index" sql:"type:varchar(36)" json:"uuid"`
	Name          string      `gorm:"not null" json:"name"`
	Slug          string      `gorm:"not null" json:"slug"`
	ApplicationID uint64      `json:"application_id,string"`
	Application   Application `json:"application,array,omitempty"`
	Versions      []Version   `json:"versions,array,omitempty"`
}

// BeforeCreate is a hook that runs before inserting a new record into the
// database
func (environment *Environment) BeforeCreate() (err error) {
	if environment.UUID == "" {
		environment.UUID = uuid.New()
	}
	if environment.Slug == "" {
		environment.Slug = slugify.Slugify(environment.Name)
	}
	return
}

// AllEnvironments will return all of the Environments
func AllEnvironments(db *gorm.DB) ([]Environment, error) {
	var environments []Environment
	result := db.Preload("Application").Find(&environments)

	return environments, result.Error
}

// CreateEnvironment will add a new Environment to Brizo
func CreateEnvironment(db *gorm.DB, client kube.APIInterface, environment *Environment, application *Application) (bool, error) {
	result := db.Create(&environment)
	service := environmentServiceDefinition(environment, application)
	err := client.CreateService(service)
	if err != nil {
		// @TODO need a better way to handle some of this upstream
		DeleteEnvironment(db, environment.Name)
		return false, err
	}

	return result.RowsAffected == 1, result.Error
}

func environmentServiceDefinition(environment *Environment, application *Application) *v1.Service {
	name := fmt.Sprintf(
		"%v-%v",
		application.Slug,
		environment.Slug,
	)
	return &v1.Service{
		ObjectMeta: v1.ObjectMeta{
			Name:      name,
			Namespace: "brizo",
			Labels: map[string]string{
				"brizoManaged": "true",
				"appUUID":      application.UUID,
				"envUUID":      environment.UUID,
			},
		},
		Spec: v1.ServiceSpec{
			Selector: map[string]string{
				"appUUID": application.UUID,
				"envUUID": environment.UUID,
			},
			Ports: []v1.ServicePort{
				v1.ServicePort{
					Protocol: v1.ProtocolTCP,
					Port:     int32(80),
				},
			},
		},
	}
}

// UpdateEnvironment will update an existing Environment
func UpdateEnvironment(db *gorm.DB, environment *Environment) (bool, error) {
	environment.Slug = slugify.Slugify(environment.Name)
	result := db.Model(environment).Where("uuid = ?", environment.UUID).
		UpdateColumns(environment)

	return result.RowsAffected == 1, result.Error
}

// UpdateEnvironmentService will update an existing Environment's service in K8s
func UpdateEnvironmentService(db *gorm.DB, client kube.APIInterface, environment *Environment, ports []ContainerPort) (bool, error) {
	var svcPorts = make([]v1.ServicePort, len(ports))
	var protocol = v1.ProtocolTCP

	// build svcPorts (service ports)
	for index, element := range ports {
		if element.Protocol == "UDP" {
			protocol = v1.ProtocolUDP
		}

		// Name is required because k8s requires each port have a name for
		// services with multiple ports
		var sp = v1.ServicePort{
			Name:     slugify.Slugify(fmt.Sprintf("%v-%v", element.Protocol, element.Port)),
			Protocol: protocol,
			Port:     int32(element.Port),
		}

		svcPorts[index] = sp
	}

	name := fmt.Sprintf(
		"%v-%v",
		environment.Application.Slug,
		environment.Slug,
	)

	service, err := client.GetService("brizo", name)

	if err != nil {
		return false, err
	}

	service.Spec.Ports = svcPorts

	err = client.UpdateService(service)

	if err != nil {
		return false, err
	}

	return true, err
}

// GetEnvironment will get an existing Environment by id
func GetEnvironment(db *gorm.DB, id string) (*Environment, error) {
	environment := new(Environment)
	if err := db.Preload("Application.Environments").Preload("Versions").Where("uuid = ?", id).First(&environment).Error; err != nil {
		return environment, err
	}

	if environment.ID == 0 {
		return new(Environment), errors.New("not-found")
	}

	return environment, nil
}

// DeleteEnvironment will delete an existing Environment by name
func DeleteEnvironment(db *gorm.DB, name string) (bool, error) {
	result := db.Delete(Environment{}, "name = ?", name)

	return result.RowsAffected == 1, result.Error
}
