package api

import (
	"fmt"
	"github.com/cohousing/cohousing-tenant-api/db"
	"github.com/cohousing/cohousing-tenant-api/domain"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"math/rand"
)

func CreateFixtureRoutes(router *gin.RouterGroup) {

	router.GET("fixtures", MustBeTenant(), MustAuthenticate(), MustBeGlobalAdmin(), func(c *gin.Context) {
		tenantDB := db.GetTenantDB(GetTenantFromContext(c))

		tenantDB.Exec("DELETE FROM residents")
		tenantDB.Exec("DELETE FROM apartments")
		tenantDB.Exec("DELETE FROM users_groups")
		tenantDB.Exec("DELETE FROM users")
		tenantDB.Exec("DELETE FROM groups")

		createApartmentAndResidents(tenantDB)

		createUsersAndGroups(tenantDB)
	})

}
func createUsersAndGroups(tenantDB *gorm.DB) {
	// Create Admin User
	admin := domain.User{
		Username: "admin",
		Password: "admin",
	}
	admin.GlobalAdmin = true

	tenantDB.Create(&admin)
	// Create Admin Group
	adminGroup := domain.Group{
		Name: "Admins",
	}
	adminGroup.CreateResidents = true
	adminGroup.ReadResidents = true
	adminGroup.UpdateResidents = true
	adminGroup.DeleteResidents = true
	adminGroup.CreateApartments = true
	adminGroup.ReadApartments = true
	adminGroup.UpdateApartments = true
	adminGroup.DeleteApartments = true
	adminGroup.CreateUsers = true
	adminGroup.ReadUsers = true
	adminGroup.UpdateUsers = true
	adminGroup.DeleteUsers = true
	tenantDB.Create(&adminGroup)

	// Create "Moderator User" Group
	moderatorGroup := domain.Group{
		Name: "Moderators",
	}
	moderatorGroup.CreateResidents = true
	moderatorGroup.ReadResidents = true
	moderatorGroup.UpdateResidents = true
	moderatorGroup.CreateApartments = true
	moderatorGroup.ReadApartments = true
	moderatorGroup.UpdateApartments = true
	tenantDB.Create(&moderatorGroup)

	// Create "Normal User" Group
	userGroup := domain.Group{
		Name: "Users",
	}
	userGroup.ReadResidents = true
	userGroup.ReadApartments = true
	tenantDB.Create(&userGroup)

	// Create a Guest Group
	guestGroup := domain.Group{
		Name: "Guests",
	}
	tenantDB.Create(&guestGroup)

	// Create 5 users with admin rights
	for i := 1; i <= 5; i++ {
		u := domain.User{
			Username: fmt.Sprintf("admin_user%d", i),
			Password: fmt.Sprintf("admin_user%d", i),
			Groups:   []domain.Group{adminGroup, moderatorGroup, userGroup},
		}
		tenantDB.Create(&u)
	}

	// Create 20 users with moderator rights
	for i := 1; i <= 50; i++ {
		u := domain.User{
			Username: fmt.Sprintf("moderator%d", i),
			Password: fmt.Sprintf("moderator%d", i),
			Groups:   []domain.Group{moderatorGroup, userGroup},
		}
		tenantDB.Create(&u)
	}

	// Create 50 users with only user rights
	for i := 1; i <= 50; i++ {
		u := domain.User{
			Username: fmt.Sprintf("user%d", i),
			Password: fmt.Sprintf("user%d", i),
			Groups:   []domain.Group{userGroup},
		}
		tenantDB.Create(&u)
	}

	// Create 10 guest users with only guest rights
	for i := 1; i <= 10; i++ {
		u := domain.User{
			Username: fmt.Sprintf("guest%d", i),
			Password: fmt.Sprintf("guest%d", i),
			Groups:   []domain.Group{guestGroup},
		}
		tenantDB.Create(&u)
	}
}
func createApartmentAndResidents(tenantDB *gorm.DB) {
	for i := 1; i <= 100; i++ {
		a := domain.Apartment{
			Address: fmt.Sprintf("Apartment %d", i),
		}
		tenantDB.Create(&a)

		for j := 1; j <= rand.Intn(4)+1; j++ {
			r := domain.Resident{
				Name:         fmt.Sprintf("Resident %d.%d", i, j),
				PhoneNumber:  fmt.Sprintf("12345%d", i+j),
				EmailAddress: fmt.Sprintf("resident%d.%d@example.com", i, j),
				ApartmentID:  &a.ID,
			}

			tenantDB.Create(&r)
		}
	}
}
