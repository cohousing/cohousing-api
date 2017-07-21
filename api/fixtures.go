package api

import (
	"fmt"
	"github.com/cohousing/cohousing-api/db"
	"github.com/cohousing/cohousing-api/domain"
	"github.com/gin-gonic/gin"
	"math/rand"
)

func CreateFixtureRoutes(router *gin.RouterGroup) {

	router.GET("fixtures", MustBeTenant(), func(c *gin.Context) {
		tenantDB := db.GetTenantDB(GetTenantFromContext(c))

		tenantDB.Delete(domain.Resident{})
		tenantDB.Delete(domain.Apartment{})

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
	})

}
