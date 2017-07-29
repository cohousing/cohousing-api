package domain

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"reflect"
	"testing"
)

func TestUser_ResolvePermissions_NoGroups(t *testing.T) {
	u := User{}

	u.GlobalAdmin = true
	assertEqualPermissions(t, u.ResolvePermissions(), Permission{
		GlobalAdmin: true,
	})

	u.ReadApartments = true
	assertEqualPermissions(t, u.ResolvePermissions(), Permission{
		GlobalAdmin:    true,
		ReadApartments: true,
	})
}

func TestUser_ResolvePermissions_WithGroups(t *testing.T) {
	u := User{}
	u.GlobalAdmin = true

	group1 := Group{
		Name: "Group 1",
	}
	group1.ReadApartments = true

	group2 := Group{
		Name: "Group 2",
	}
	group2.CreateApartments = true

	u.Groups = []Group{
		group1,
		group2,
	}

	assertEqualPermissions(t, u.ResolvePermissions(), Permission{
		GlobalAdmin:      true,
		ReadApartments:   true,
		CreateApartments: true,
	})
}

func TestUser_BeforeCreate(t *testing.T) {
	u := User{
		Password: "password",
	}

	u.BeforeCreate()

	assert.NotEqual(t, "password", u.Password)

	assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(u.Password), []byte("password")))
}

func TestPermission_BeforeSave(t *testing.T) {
	p := Permission{}
	p.CreateResidents = true
	p.ReadResidents = true
	p.UpdateResidents = true
	p.DeleteResidents = true

	assert.Equal(t, "", p.Residents)
	assert.Equal(t, "", p.Apartments)

	p.BeforeSave()

	assert.Equal(t, "crud", p.Residents)
	assert.Equal(t, "____", p.Apartments)

	p.DeleteResidents = false
	p.ReadApartments = true

	p.BeforeSave()

	assert.Equal(t, "cru_", p.Residents)
	assert.Equal(t, "_r__", p.Apartments)
}

func TestPermission_AfterFind(t *testing.T) {
	p := Permission{}
	p.Residents = "crud"
	p.Apartments = "____"
	p.Users = "____"

	assert.Equal(t, p.CreateResidents, false)
	assert.Equal(t, p.ReadResidents, false)
	assert.Equal(t, p.UpdateResidents, false)
	assert.Equal(t, p.DeleteResidents, false)
	assert.Equal(t, p.CreateApartments, false)
	assert.Equal(t, p.ReadApartments, false)
	assert.Equal(t, p.UpdateApartments, false)
	assert.Equal(t, p.DeleteApartments, false)

	p.AfterFind()

	assert.Equal(t, p.CreateResidents, true)
	assert.Equal(t, p.ReadResidents, true)
	assert.Equal(t, p.UpdateResidents, true)
	assert.Equal(t, p.DeleteResidents, true)
	assert.Equal(t, p.CreateApartments, false)
	assert.Equal(t, p.ReadApartments, false)
	assert.Equal(t, p.UpdateApartments, false)
	assert.Equal(t, p.DeleteApartments, false)

	p.Apartments = "_r__"
	p.AfterFind()

	assert.Equal(t, p.CreateResidents, true)
	assert.Equal(t, p.ReadResidents, true)
	assert.Equal(t, p.UpdateResidents, true)
	assert.Equal(t, p.DeleteResidents, true)
	assert.Equal(t, p.CreateApartments, false)
	assert.Equal(t, p.ReadApartments, true)
	assert.Equal(t, p.UpdateApartments, false)
	assert.Equal(t, p.DeleteApartments, false)
}

func TestPermission_HasPermission(t *testing.T) {
	p := Permission{
		ReadResidents: true,
	}

	val, err := p.HasPermission("CreateResidents")
	assert.NoError(t, err)
	assert.False(t, val)

	val, err = p.HasPermission("ReadResidents")
	assert.NoError(t, err)
	assert.True(t, val)

	val, err = p.HasPermission("SomethingThatDoesNotExist")
	assert.Error(t, err)
	assert.False(t, val)
}

func assertEqualPermissions(t *testing.T, actualPermission, expectedPermission Permission) {
	permissionType := reflect.TypeOf(Permission{})
	actualPermissionValue := reflect.ValueOf(actualPermission)
	expectedPermissionValue := reflect.ValueOf(expectedPermission)
	for i := 0; i < permissionType.NumField(); i++ {
		if permissionType.Field(i).Type.Kind() == reflect.Bool {
			actualPermissionBool := actualPermissionValue.Field(i).Bool()
			expectedPermissionBool := expectedPermissionValue.Field(i).Bool()
			if actualPermissionBool != expectedPermissionBool {
				t.Error(fmt.Sprintf("Permission %s wasn't the same. Expected %t but got %t", permissionType.Field(i).Name, expectedPermissionBool, actualPermissionBool))
			}
		}
	}
}
