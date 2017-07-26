package tenant

import (
	"fmt"
	"github.com/cohousing/cohousing-api/api/utils"
	"github.com/cohousing/cohousing-api/db"
	"github.com/cohousing/cohousing-api/domain/tenant"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"math/rand"
)

func CreateFixtureRoutes(router *gin.RouterGroup) {

	router.GET("fixtures", utils.MustBeTenant(), func(c *gin.Context) {
		tenantDB := db.GetTenantDB(utils.GetTenantFromContext(c))

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
	admin := tenant.User{
		Username: "admin",
		Password: "admin",
	}
	admin.GlobalAdmin = true

	tenantDB.Create(&admin)
	// Create Admin Group
	adminGroup := tenant.Group{
		Name: "Admins",
	}
	adminGroup.GlobalAdmin = true
	tenantDB.Create(&adminGroup)

	// Create "Moderator User" Group
	moderatorGroup := tenant.Group{
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
	userGroup := tenant.Group{
		Name: "Users",
	}
	userGroup.ReadResidents = true
	userGroup.ReadApartments = true
	tenantDB.Create(&userGroup)

	// Create a Guest Group
	guestGroup := tenant.Group{
		Name: "Guests",
	}
	tenantDB.Create(&guestGroup)

	// Create 5 users with admin rights
	for i := 1; i <= 5; i++ {
		u := tenant.User{
			Username: fmt.Sprintf("admin_user%d", i),
			Password: fmt.Sprintf("admin_user%d", i),
			Groups:   []tenant.Group{adminGroup, moderatorGroup, userGroup},
		}
		tenantDB.Create(&u)
	}

	// Create 20 users with moderator rights
	for i := 1; i <= 50; i++ {
		u := tenant.User{
			Username: fmt.Sprintf("moderator%d", i),
			Password: fmt.Sprintf("moderator%d", i),
			Groups:   []tenant.Group{moderatorGroup, userGroup},
		}
		tenantDB.Create(&u)
	}

	// Create 50 users with only user rights
	for i := 1; i <= 50; i++ {
		u := tenant.User{
			Username: fmt.Sprintf("user%d", i),
			Password: fmt.Sprintf("user%d", i),
			Groups:   []tenant.Group{userGroup},
		}
		tenantDB.Create(&u)
	}

	// Create 10 guest users with only guest rights
	for i := 1; i <= 10; i++ {
		u := tenant.User{
			Username: fmt.Sprintf("guest%d", i),
			Password: fmt.Sprintf("guest%d", i),
			Groups:   []tenant.Group{guestGroup},
		}
		tenantDB.Create(&u)
	}
}
func createApartmentAndResidents(tenantDB *gorm.DB) {
	for i := 1; i <= 100; i++ {
		a := tenant.Apartment{
			Address: fmt.Sprintf("Apartment %d", i),
		}
		tenantDB.Create(&a)

		for j := 1; j <= rand.Intn(4)+1; j++ {
			r := tenant.Resident{
				Name:         fmt.Sprintf("Resident %d.%d", i, j),
				PhoneNumber:  fmt.Sprintf("12345%d", i+j),
				EmailAddress: fmt.Sprintf("resident%d.%d@example.com", i, j),
				ApartmentID:  &a.ID,
			}

			tenantDB.Create(&r)
		}
	}
}
