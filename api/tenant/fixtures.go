package tenant

import (
	"fmt"
	"github.com/cohousing/cohousing-api/api/utils"
	"github.com/cohousing/cohousing-api/db"
	"github.com/cohousing/cohousing-api/domain/tenant"
	"github.com/gin-gonic/gin"
	"math/rand"
)

func CreateFixtureRoutes(router *gin.RouterGroup) {

	router.GET("fixtures", utils.MustBeTenant(), func(c *gin.Context) {
		tenantDB := db.GetTenantDB(utils.GetTenantFromContext(c))

		tenantDB.Delete(tenant.Resident{})
		tenantDB.Delete(tenant.Apartment{})

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
	})

}
