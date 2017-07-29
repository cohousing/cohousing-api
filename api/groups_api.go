package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"github.com/cohousing/cohousing-api-utils/api"
	"github.com/cohousing/cohousing-api-utils/db"
	"github.com/cohousing/cohousing-tenant-api/domain"
	domain2 "github.com/cohousing/cohousing-api-utils/domain"
)

var (
	GroupBasePath string
)

func CreateGroupRoutes(router *gin.RouterGroup, dbFactory db.DBFactory) {
	endpoint := ConfigureBasicTenantEndpoint(router, utils.BasicEndpointConfig{
		Path:           "groups",
		Domain:         domain.Group{},
		DBFactory:      dbFactory,
		RouterHandlers: []gin.HandlerFunc{MustBeTenant(), MustAuthenticate()},
	})
	GroupBasePath = endpoint.BasePath()

	utils.AddLinkFactory(&domain.Group{}, groupLinkFactory)

	endpoint.GET("/:id/users", AuthorizeDomainObject(domain.Group{}, PERM_READ), getUsersForGroup(dbFactory))
}

func groupLinkFactory(c *gin.Context, halResource domain2.HalResource, basePath string, detailed bool) {
	u := halResource.(*domain.Group)
	u.AddLink(domain2.REL_SELF, fmt.Sprintf("%s/%d", basePath, u.ID))

	if detailed {
		permission := ResolvePermission(c)
		if permission.GlobalAdmin || permission.UpdateGroups {
			u.AddLink(domain2.REL_UPDATE, fmt.Sprintf("%s/%d", basePath, u.ID))
		}
		if permission.GlobalAdmin || permission.DeleteGroups {
			u.AddLink(domain2.REL_DELETE, fmt.Sprintf("%s/%d", basePath, u.ID))
		}
		if permission.GlobalAdmin || permission.ReadUsers {
			u.AddLink(domain.REL_USERS, fmt.Sprintf("%s/%d/users", basePath, u.ID))
		}
	}
}

func getUsersForGroup(dbFactory db.DBFactory) gin.HandlerFunc {
	return func(c *gin.Context) {
		if id, err := strconv.ParseUint(c.Param("id"), 10, 64); err == nil {
			g := domain.Group{}
			g.ID = id
			var users []domain.User
			dbFactory(c).Model(&g).Related(&users, "Users")
			utils.AddLinks(c, users, UserBasePath, false)
			c.JSON(http.StatusOK, users)
		}
	}
}
