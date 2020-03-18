package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/twinj/uuid"
	"log"
	"time"
)

type Student struct {
	Id      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Surname string    `json:"surname"`
}

var studentlist []Student = []Student{}
var client *redis.Client

func main() {
	doconfig()
	engine := gin.Default()
	studentgroup := engine.Group("/students")

	studentgroup.POST("/", func(c *gin.Context) {
		var student Student
		if err := c.ShouldBind(&student); err != nil {
			c.Status(500)
			return
		}
		student.Id = uuid.NewV1()

		studentlist = append(studentlist, student)
		go func() {
			client.Set("students"+student.Id.String(), marshall(student), time.Hour*24)
			client.Del("students")
		}()
		c.String(201, "object created succesfully ")
	})

	studentgroup.GET("/", func(c *gin.Context) {
		result, err := client.Get("students").Result()
		if err == nil {
			c.Header("Content-Type", "application-json")
			c.String(200, result)
			return
		}
		client.Set("students", marshall(studentlist), time.Hour*24)
		time.Sleep(time.Second * 1)
		c.JSON(200, studentlist)
	})

	studentgroup.GET("/:uuid", func(context *gin.Context) {

		uuid, _ := context.Params.Get("uuid")
		result, err := client.Get("students" + uuid).Result()

		if err == nil {
			context.Header("Content-Type", "application-json")
			context.String(200, result)
			return
		}
		context.String(404, "resource not found123 ")
	})

	engine.Run(":8080")
}
func marshall(data interface{}) []byte {
	marshal, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	return marshal
}
func doconfig() {

	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := client.Ping().Result()
	if err != nil {
		log.Fatal(err)
	}

}
