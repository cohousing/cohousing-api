package tenant

import (
	"errors"
	"fmt"
	"github.com/cohousing/cohousing-api/domain"
	"golang.org/x/crypto/bcrypt"
	"reflect"
	"strings"
)

const (
	REL_USERS  domain.RelType = "users"
	REL_GROUPS domain.RelType = "groups"
)

type User struct {
	domain.BaseModel
	Username string  `json:"username"`
	Password string  `gorm:"size:60" json:"-"`
	Groups   []Group `gorm:"many2many:users_groups;" json:"-"`
	Permission
	domain.DefaultHalResource
}

func (u *User) ResolvePermissions() Permission {
	resolvedPermission := Permission{}

	overridePermissions(&resolvedPermission, u.Permission)
	for i := 0; i < len(u.Groups); i++ {
		overridePermissions(&resolvedPermission, u.Groups[i].Permission)
	}

	return resolvedPermission
}

func overridePermissions(permission *Permission, overridePermission Permission) {
	permissionType := reflect.TypeOf(Permission{})
	permissionValue := reflect.ValueOf(permission).Elem()
	overridePermissionValue := reflect.ValueOf(overridePermission)
	for i := 0; i < permissionType.NumField(); i++ {
		if permissionType.Field(i).Type.Kind() == reflect.Bool && overridePermissionValue.Field(i).Bool() {
			permissionValue.Field(i).SetBool(overridePermissionValue.Field(i).Bool())
		}
	}
}

func (u *User) BeforeCreate() (err error) {
	if u.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), 4)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
	}

	return nil
}

type Group struct {
	domain.BaseModel
	Name  string `json:"name"`
	Users []User `gorm:"many2many:users_groups;" json:"-"`
	Permission
	domain.DefaultHalResource
}

type Permission struct {
	GlobalAdmin      bool   `json:"global_admin"`
	Residents        string `gorm:"size:4;column:perm_residents" json:"-"`
	CreateResidents  bool   `gorm:"-" json:"create_residents"`
	ReadResidents    bool   `gorm:"-" json:"read_residents"`
	UpdateResidents  bool   `gorm:"-" json:"update_residents"`
	DeleteResidents  bool   `gorm:"-" json:"delete_residents"`
	Apartments       string `gorm:"size:4;column:perm_apartments" json:"-"`
	CreateApartments bool   `gorm:"-" json:"create_apartments"`
	ReadApartments   bool   `gorm:"-" json:"read_apartments"`
	UpdateApartments bool   `gorm:"-" json:"update_apartments"`
	DeleteApartments bool   `gorm:"-" json:"delete_apartments"`
	Users            string `gorm:"size:4;column:perm_users" json:"-"`
	CreateUsers      bool   `gorm:"-" json:"create_users"`
	ReadUsers        bool   `gorm:"-" json:"read_users"`
	UpdateUsers      bool   `gorm:"-" json:"update_users"`
	DeleteUsers      bool   `gorm:"-" json:"delete_users"`
	Groups           string `gorm:"size:4;column:perm_groups" json:"-"`
	CreateGroups     bool   `gorm:"-" json:"create_groups"`
	ReadGroups       bool   `gorm:"-" json:"read_groups"`
	UpdateGroups     bool   `gorm:"-" json:"update_groups"`
	DeleteGroups     bool   `gorm:"-" json:"delete_groups"`
}

func (p *Permission) BeforeSave() (err error) {
	permissionType := reflect.TypeOf(p).Elem()
	permissionValue := reflect.ValueOf(p).Elem()

	permissionRender := func(fieldPrefix, permissionFieldName, renderValue string) string {
		fieldValue := permissionValue.FieldByName(fmt.Sprintf("%s%s", fieldPrefix, permissionFieldName)).Bool()
		if fieldValue {
			return renderValue
		} else {
			return "_"
		}
	}

	for i := 0; i < permissionType.NumField(); i++ {
		permissionFieldName := permissionType.Field(i).Name
		if permissionFieldName != "GlobalAdmin" &&
			!strings.HasPrefix(permissionFieldName, "Create") &&
			!strings.HasPrefix(permissionFieldName, "Read") &&
			!strings.HasPrefix(permissionFieldName, "Update") &&
			!strings.HasPrefix(permissionFieldName, "Delete") {

			permissionValue.Field(i).SetString(fmt.Sprintf("%s%s%s%s",
				permissionRender("Create", permissionFieldName, "c"),
				permissionRender("Read", permissionFieldName, "r"),
				permissionRender("Update", permissionFieldName, "u"),
				permissionRender("Delete", permissionFieldName, "d"),
			))
		}
	}
	return nil
}

func (p *Permission) AfterFind() (err error) {
	permissionType := reflect.TypeOf(p).Elem()
	permissionValue := reflect.ValueOf(p).Elem()

	permissionSetter := func(fieldPrefix, permissionFieldName, value, truthValue string) {
		hasPermission := false
		if value == truthValue {
			hasPermission = true
		}

		permissionValue.FieldByName(fmt.Sprintf("%s%s", fieldPrefix, permissionFieldName)).SetBool(hasPermission)
	}

	for i := 0; i < permissionType.NumField(); i++ {
		permissionFieldName := permissionType.Field(i).Name
		if permissionFieldName != "GlobalAdmin" &&
			!strings.HasPrefix(permissionFieldName, "Create") &&
			!strings.HasPrefix(permissionFieldName, "Read") &&
			!strings.HasPrefix(permissionFieldName, "Update") &&
			!strings.HasPrefix(permissionFieldName, "Delete") {

			permissionString := permissionValue.Field(i).String()

			permissionSetter("Create", permissionFieldName, string(permissionString[0]), "c")
			permissionSetter("Read", permissionFieldName, string(permissionString[1]), "r")
			permissionSetter("Update", permissionFieldName, string(permissionString[2]), "u")
			permissionSetter("Delete", permissionFieldName, string(permissionString[3]), "d")
		}
	}
	return nil
}

func (p *Permission) HasPermission(permissionString string) (bool, error) {
	permissionField := reflect.ValueOf(p).Elem().FieldByName(permissionString)
	if permissionField == (reflect.Value{}) {
		return false, errors.New(fmt.Sprintf("No such permission field %s on permission %v\n", permissionString, p))
	}
	return permissionField.Bool(), nil
}
