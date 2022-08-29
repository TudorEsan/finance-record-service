package controller

import (
	"App/database"
	"App/helpers"
	"App/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var NetWorthCollection *mongo.Collection = database.OpenCollection(database.Client, "NetWorth")
var InfoCollection *mongo.Collection = database.OpenCollection(database.Client, "Info")

func InitNetWort() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func GetRecord() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		user, err := helpers.GetUserFromContext(c)
		if err != nil {
			helpers.ReturnError(c, http.StatusBadRequest, err)
			return
		}
		record, err := helpers.GetRecord(user.ID, id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"record": record})

	}
}

func AddRecord() gin.HandlerFunc {
	return func(c *gin.Context) {
		var recordBody models.Record
		if err := c.BindJSON(&recordBody); err != nil {
			helpers.ReturnError(c, http.StatusBadRequest, err)
			return
		}
		if err := validate.Struct(recordBody); err != nil {
			helpers.ReturnError(c, http.StatusBadRequest, err)
			return
		}
		user, err := helpers.GetUserFromContext(c)
		if err != nil {
			helpers.ReturnError(c, http.StatusInternalServerError, err)
			return
		}
		err = helpers.AddRecord(user.ID, recordBody)
		if err != nil {
			helpers.ReturnError(c, http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	}
}

func GetRecords() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := helpers.GetUserFromContext(c)
		if err != nil {
			helpers.ReturnError(c, http.StatusInternalServerError, err)
			return
		}
		page, err := strconv.Atoi(c.DefaultQuery("page", "0"))
		if err != nil {
			helpers.ReturnError(c, http.StatusBadRequest, err)
			return
		}
		perPage, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
		if err != nil {
			helpers.ReturnError(c, http.StatusBadRequest, err)
			return
		}
		records, err := helpers.GetRecords(user.ID, int(page), perPage)
		if err != nil {
			helpers.ReturnError(c, http.StatusBadRequest, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"records": records})

	}
}

func GetRecordCount() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := helpers.GetUserFromContext(c)
		if err != nil {
			helpers.ReturnError(c, http.StatusInternalServerError, err)
			return
		}
		count, err := helpers.GetRecordCount(user.ID)
		if err != nil {
			helpers.ReturnError(c, http.StatusBadRequest, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"recordCount": count})
	}
}

func DeleteRecord() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := helpers.GetUserFromContext(c)
		if err != nil {
			helpers.ReturnError(c, http.StatusBadRequest, err)
			return
		}

		id := c.Param("id")
		if id == "" {
			helpers.ReturnError(c, http.StatusBadRequest, fmt.Errorf("ID is required"))
			return
		}

		netWorthId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			helpers.ReturnError(c, http.StatusBadRequest, err)
		}
		err = helpers.DeleteRecord(user.ID, netWorthId)
		if err != nil {
			helpers.ReturnError(c, http.StatusBadRequest, err)
		}
		c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
	}
}

func UpdateRecord() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := helpers.GetUserFromContext(c)
		if err != nil {
			helpers.ReturnError(c, http.StatusBadRequest, err)
			return
		}
		id := c.Param("id")
		if id == "" {
			helpers.ReturnError(c, http.StatusBadRequest, fmt.Errorf("ID is required"))
			return
		}
		var recordBody models.Record
		if err := c.BindJSON(&recordBody); err != nil {
			helpers.ReturnError(c, http.StatusBadRequest, err)
			return
		}
		if err := validate.Struct(recordBody); err != nil {
			helpers.ReturnError(c, http.StatusBadRequest, err)
			return
		}
		err = helpers.UpdateRecord(user.ID, recordBody)
		if err != nil {
			helpers.ReturnError(c, http.StatusBadRequest, err)
		}
		c.JSON(http.StatusOK, gin.H{"message": "Updated"})
	}
}
