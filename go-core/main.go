package main

import (
	"encoding/json"
	"notification-center-go-api/middlewares"
	"notification-center-go-api/models"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	api := router.Group("/api/v3")
	api.Use(middlewares.JWTAuth())
	api.POST("notifications/:id", CreateActivity)

	router.Run(":8080")
}

// CreateActivity Handle activity request
func CreateActivity(c *gin.Context) {

	var notification models.Notification
	id := c.Params.ByName("id")
	models.GetDB().Where("slug = ?", id).First(&notification)

	deliver := models.FindDeliverByNotificationID(notification.ID)

	if deliver.ID != 0 {

		// Parse request contet
		var content models.RawContent
		c.BindJSON(&content)

		data, _ := json.Marshal(content)

		if content.Delivery.Date == "" {
			var activity = models.Activity{NotificationDeliverID: deliver.ID, Status: "pending", SendAt: time.Now(), RequestContent: data}
			models.GetDB().Create(&activity)
		} else {
			// TODO: Add zone support content.Delivery.Zone

			layout := "2006-01-02 15:04"
			sendTime := content.Delivery.Date + " " + content.Delivery.Time
			parsedTime, err := time.Parse(layout, sendTime)

			if err != nil {
				c.JSON(422, gin.H{"error": err})
				return
			} else {
				var activity = models.Activity{NotificationDeliverID: deliver.ID, Status: "pending", SendAt: parsedTime, RequestContent: data}
				models.GetDB().Create(&activity)
			}
		}

		c.JSON(200, gin.H{"info": "Success!"})
	} else {
		c.JSON(422, gin.H{"error": "Email type deliver not found!"})
	}
}