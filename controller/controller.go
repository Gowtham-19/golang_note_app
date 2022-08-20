package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Gowtham-19/note_golang_server/configs"
	"github.com/Gowtham-19/note_golang_server/model"
	"github.com/Gowtham-19/note_golang_server/responses"
	"github.com/bitly/go-simplejson"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//getting all notes
func GetAll_Notes(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "*")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var notes []model.Notes
	defer cancel()
	data, err := configs.Collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
		return
	}
	//if no error then returning data
	defer data.Close(ctx)
	for data.Next(ctx) {
		var singleNote model.Notes
		if err = data.Decode(&singleNote); err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
		}
		notes = append(notes, singleNote)
	}
	c.JSON(200, gin.H{
		"status": 200,
		"data":   notes,
	})
}

//creating a note
func Create_Note(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "*")
	var note model.Notes
	//validate the request body
	if err := c.BindJSON(&note); err != nil {
		c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
		return
	}
	time_value := time.Now()
	//inserting body
	newNote := model.Notes{
		Id:            primitive.NewObjectID(),
		Subject:       note.Subject,
		Content_Type:  note.Content_Type,
		Description:   note.Description,
		CreatedAt:     &time_value,
		UpdatedAt:     &time_value,
		Created_Date:  int64(time_value.Day()),
		Created_Month: int64(time_value.Month()),
		Created_Year:  int64(time_value.Year()),
		Status:        "Active",
	}
	_, err := configs.Collection.InsertOne(context.Background(), newNote)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
		return
	}
	c.JSON(200, gin.H{
		"status": 200,
		"data":   "note created of id",
	})
}

//updating a note
func Update_Note(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "*")
	var notes model.Notes
	//validating body
	if err := json.NewDecoder(c.Request.Body).Decode(&notes); err != nil {
		c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
		return
	}
	//converting struct into bytes
	data, _ := json.Marshal(notes)
	//converting bytes into json
	body, _ := simplejson.NewJson(data)
	//forming key to update notes record
	update_id, _ := primitive.ObjectIDFromHex(body.Get("_id").MustString())
	filter := bson.M{"_id": update_id}
	update_values := bson.M{}
	//forming keys to update
	for key, value := range body.MustMap() {
		if key != "_id" {
			update_values[key] = value
		}
	}
	update_values["updatedat"] = time.Now() //updating time
	update := bson.M{"$set": update_values}
	_, err := configs.Collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "Error", Data: map[string]interface{}{"data": err}})
	} else {
		c.JSON(200, gin.H{
			"status": 200,
			"data":   "data updated",
		})
	}
}

//Deleting a note
func Delete_Note(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "*")
	id, _ := c.Params.Get("id")
	//fetching id to delete note
	note_id, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": note_id}
	_, err := configs.Collection.DeleteOne(context.Background(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err}},
		)
	} else {
		c.JSON(200, gin.H{
			"status": 200,
			"data":   "deletion successfull",
		})
	}
}

//filtering notes based on date filter
func Filter_Notes(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "*")
	//defining a struct for filter body
	filter_data, _ := ioutil.ReadAll(c.Request.Body)
	body, _ := simplejson.NewJson(filter_data)
	//request model for filter notes
	var filter_type = body.Get("filter_type").MustString()
	created_date := body.Get("date").MustInt64()
	created_month := body.Get("month").MustInt64()
	created_year := body.Get("year").MustInt64()
	//response data
	var notes []model.Notes
	if filter_type == "specific_day" {
		//specific day implementation,here we need to get a exact day records
		data, err := configs.Collection.Find(context.Background(), bson.M{"created_year": created_year, "created_month": created_month, "created_date": created_date})
		if err != nil {
			c.JSON(http.StatusInternalServerError,
				responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err}},
			)
		} else {
			defer data.Close(context.Background())
			for data.Next(context.Background()) {
				var singleNote model.Notes
				if err = data.Decode(&singleNote); err != nil {
					c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				}
				notes = append(notes, singleNote)
			}
		}
	} else if filter_type == "last_week" {
		//last week implementation,here we need to get last 7 days from current date
		data, err := configs.Collection.Find(context.Background(), bson.M{
			"createdat": bson.M{"$lte": time.Now()},
			"$and": []interface{}{
				bson.M{"createdat": bson.M{"$gte": time.Now().AddDate(0, 0, -7)}},
			},
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError,
				responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err}},
			)
		} else {
			defer data.Close(context.Background())
			for data.Next(context.Background()) {
				var singleNote model.Notes
				if err = data.Decode(&singleNote); err != nil {
					c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				}
				notes = append(notes, singleNote)
			}
		}
		fmt.Println("Last week filter")
	} else if filter_type == "specific_month" {
		//specific month implementation,here we need to get exact month and exact year records
		data, err := configs.Collection.Find(context.Background(), bson.M{"created_year": created_year, "created_month": created_month})
		if err != nil {
			c.JSON(http.StatusInternalServerError,
				responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err}},
			)
		} else {
			defer data.Close(context.Background())
			for data.Next(context.Background()) {
				var singleNote model.Notes
				if err = data.Decode(&singleNote); err != nil {
					c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				}
				notes = append(notes, singleNote)
			}

		}
	} else if filter_type == "specific_year" {
		//specific year implementation
		//specific month implementation,here we need to get exact month and exact year records
		data, err := configs.Collection.Find(context.Background(), bson.M{"created_year": created_year})
		if err != nil {
			c.JSON(http.StatusInternalServerError,
				responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err}},
			)
		} else {
			defer data.Close(context.Background())
			for data.Next(context.Background()) {
				var singleNote model.Notes
				if err = data.Decode(&singleNote); err != nil {
					c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				}
				notes = append(notes, singleNote)
			}
		}
	} else if filter_type == "get_all_notes" {
		data, err := configs.Collection.Find(context.Background(), bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError,
				responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err}},
			)
		} else {
			defer data.Close(context.Background())
			for data.Next(context.Background()) {
				var singleNote model.Notes
				if err = data.Decode(&singleNote); err != nil {
					c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				}
				notes = append(notes, singleNote)
			}
		}
	}
	c.JSON(200, gin.H{
		"status": 200,
		"data":   notes,
	})
}
