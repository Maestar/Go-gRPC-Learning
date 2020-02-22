package main

import (
	"BasicAPI/proto"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func main() {
	//the client needs to 'call up' the server, so he dials the port we told the server to listen on, which is 4040
	//its not a secure call, not using HTTPS, so the call is made with WithInsecure()
	connection, err := grpc.Dial("localhost:4040", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	//pass in out connection and now we have a client
	client := proto.NewAddServiceClient(connection)

	//the library called gin lets us make api endpoints pretty easily
	//first we make a gin server like so:
	ginServer := gin.Default()

	//then we tell the ginserver's get function what to do with our two functions

	//this one is for add, so the url for it is add/int a / int b
	ginServer.GET("add/:a/:b", func(contxt *gin.Context) {
		//so now we need to get A and B from the api call the user types. We can get it from the gin contxt which is the url the user puts in.
		// we tell contxt to get the a from the url, tell it that its an int, base10 and 64-bit
		a, err := strconv.ParseUint(contxt.Param("a"), 10, 64)
		if err != nil {
			//if we get a bad parameter with log it with a json and return
			contxt.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid Parameter A"})
			return
		}

		b, err := strconv.ParseUint(contxt.Param("b"), 10, 64)
		if err != nil {
			contxt.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid Parameter B"})
			return
		}

		//now that we have a and b, we cast them to int64 and we can assign them to proto's request
		request := &proto.Request{A: int64(a), B: int64(b)}

		//then we call client.add on the client we created in this if statement for error checking.
		//If error is nill we'll have gotten back a response from the server since the client 'called' the server asking for it to add.
		//we save that response in the json and move on.
		if response, err := client.Add(contxt, request); err == nil {
			contxt.JSON(http.StatusOK, gin.H{
				"result": fmt.Sprint(response.Result),
			})
		} else { //err was not nil here so we have a server error, so we log the error in json
			contxt.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	})

	//this one is for Multiply
	ginServer.GET("mult/:a/:b", func(contxt *gin.Context) {

		a, err := strconv.ParseUint(contxt.Param("a"), 10, 64)
		if err != nil {
			//if we get a bad parameter with log it with a json and return
			contxt.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid Parameter A"})
			return
		}

		b, err := strconv.ParseUint(contxt.Param("b"), 10, 64)
		if err != nil {
			contxt.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid Parameter B"})
			return
		}

		request := &proto.Request{A: int64(a), B: int64(b)}

		if response, err := client.Multiply(contxt, request); err == nil {
			contxt.JSON(http.StatusOK, gin.H{
				"result": fmt.Sprint(response.Result),
			})
		} else {
			contxt.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

	})

	//now we run the ginServer so we can access the api, one again using an if to catch errors
	if err := ginServer.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}

}
