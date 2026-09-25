package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/artvault/artvault/internal/constants"
	"github.com/artvault/artvault/internal/repository"
	"github.com/artvault/artvault/internal/util"
	"github.com/gin-gonic/gin"
)

// ContentReview 内容审核中间件：未审核通过的作品/展览不公开返回。
func ContentReview(artworkRepo *repository.ArtworkRepository, exhibitionRepo *repository.ExhibitionRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 仅对公开只读接口生效
		if c.Request.Method != http.MethodGet {
			c.Next()
			return
		}
		ctx := context.Background()
		if strings.Contains(c.Request.URL.Path, "/artworks") {
			id := c.Param("id")
			if id != "" {
				a, err := artworkRepo.FindByID(ctx, id)
				if err == nil && a.ReviewStatus != "Approved" {
					role, _ := c.Get("role")
					if role != "Admin" && role != "Curator" {
						util.Fail(c, http.StatusNotFound, constants.CodeNotFound, constants.MsgNotFound)
						return
					}
				}
			}
		}
		if strings.Contains(c.Request.URL.Path, "/exhibitions") {
			id := c.Param("id")
			if id != "" {
				e, err := exhibitionRepo.FindByID(ctx, id)
				if err == nil && e.ReviewStatus != "Approved" {
					role, _ := c.Get("role")
					if role != "Admin" && role != "Curator" {
						util.Fail(c, http.StatusNotFound, constants.CodeNotFound, constants.MsgNotFound)
						return
					}
				}
			}
		}
		c.Next()
	}
}
