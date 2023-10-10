package endpoint

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/b-turchyn/lifetracker-backend/service"
	"github.com/gin-gonic/gin"
)

var bucketService *service.BucketService

func BucketEndpoints(r *gin.Engine, service *service.BucketService) {
	bucketService = service
	group := r.Group("/buckets")
	{
		group.GET("", bucketIndexEndpoint)
		group.GET(":bucket", bucketShowEndpoint)
	}
}

var bucketIndexEndpoint = func(c *gin.Context) {
	data, err := bucketService.IndexBucket()

	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
	} else {
		c.JSON(http.StatusOK, data)
	}
}

var bucketShowEndpoint = func(c *gin.Context) {
	bucket := c.Param("bucket")
	pivotPointStr := c.Query("pivot")
	pivotPoint, err := strconv.ParseFloat(pivotPointStr, 64)

	var data []service.BucketResultRow
	if err == nil {
		data, err = bucketService.ShowBucket(bucket, pivotPoint)
	}
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
	} else {
		c.JSON(http.StatusOK, data)
	}

}
