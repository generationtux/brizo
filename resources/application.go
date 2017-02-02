package resources

import (
	"errors"

	"github.com/generationtux/brizo/database"
	"github.com/generationtux/brizo/kube"
	"github.com/jinzhu/gorm"
	"github.com/pborman/uuid"
)

// Application as defined by Brizo.
type Application struct {
	database.Model
	UUID         string        `gorm:"not null;unique_index" sql:"type:varchar(36)" json:"uuid"`
	Name         string        `gorm:"not null;unique_index" json:"name"`
	Slug         string        `gorm:"not null;unique_index" json:"slug"`
	Pods         []kube.Pod    `gorm:"-" json:"pods,array"` // gorm will ignore, but we can populate
	Environments []Environment `json:"environments,array"`
}

// BeforeCreate is a hook that runs before inserting a new record into the database
func (a *Application) BeforeCreate() (err error) {
	if a.UUID == "" {
		a.UUID = uuid.New()
	}
	return
}

// AllApplications will return all of the Applications
func AllApplications(db *gorm.DB) ([]Application, error) {
	var apps []Application
	result := db.Preload("Environments").Find(&apps)

	return apps, result.Error
}

// CreateApplication will add a new Application to Brizo
func CreateApplication(db *gorm.DB, app *Application) (bool, error) {
	result := db.Create(&app)

	return result.RowsAffected == 1, result.Error
}

// UpdateApplication will update an existing Application
func UpdateApplication(db *gorm.DB, app *Application) (bool, error) {
	result := db.Model(app).Where("uuid = ?", app.UUID).UpdateColumns(app)

	return result.RowsAffected == 1, result.Error
}

// GetApplication will get an existing Application by name
func GetApplication(db *gorm.DB, client kube.APIInterface, id string) (*Application, error) {
	app := new(Application)
	if err := db.Where("uuid = ?", id).Preload("Environments").First(&app).Error; err != nil {
		return app, err
	}

	if app.ID == 0 {
		return new(Application), errors.New("not-found")
	}

	pods, err := GetApplicationPods(client, id)
	app.Pods = pods

	return app, err
}

// DeleteApplication will delete an existing Application by name
func DeleteApplication(db *gorm.DB, name string) (bool, error) {
	result := db.Delete(Application{}, "name = ?", name)

	return result.RowsAffected == 1, result.Error
}

// GetApplicationPods returns the pods running a provided application
func GetApplicationPods(client kube.APIInterface, UUID string) ([]kube.Pod, error) {
	return client.GetPods(kube.PodOptions{})
	// if err != nil {
	// 	return []Pod{}, err
	// }
	//
	// kubePods := client.Pods("brizo")
	// list, err := kubePods.List(v1.ListOptions{
	// 	LabelSelector: fmt.Sprintf("appUUID=%v, brizoManaged=true", UUID),
	// })
	//
	// pods := []Pod{}
	// for _, kubePod := range list.Items {
	// 	pods = append(pods, Pod{Name: kubePod.GetName()})
	// }
	//
	// if len(pods) == 0 {
	// 	pods = make([]Pod, 0)
	// }
	//
	// return pods, err
}
