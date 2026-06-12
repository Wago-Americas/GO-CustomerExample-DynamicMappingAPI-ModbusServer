package main

import (
	"fmt"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/simonvetter/modbus"
	"net/http"
	"strconv"
)

type IO struct {

	Value string `json:"value" binding:"required"`
}


func main() {

	var client *modbus.ModbusClient
	var err error
	

	// 1. Initialize client with the Modbus TCP URL and connection timeout
	client, err = modbus.NewClient(&modbus.ClientConfiguration{
		URL:     "tcp://192.168.1.35:502", // Replace with your target device IP/Port
		Timeout: 1 * time.Second,
	})
	if err != nil {
		fmt.Printf("failed to create client: %v\n", err)
		return
	}

	// 2. Open the connection
	err = client.Open()
	if err != nil {
		fmt.Printf("failed to open connection: %v\n", err)
		return
	}
	// Always close the connection when done
	defer client.Close()

	

	// Initialize the Gin router
	r := gin.Default()

	// Define the GET endpoint
	r.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello from Gin!",
		})
	})

	r.GET("/OutputTotal", func(c *gin.Context) {
		// 1. Run your Modbus query first
		outputnum, err := client.ReadRegister(0, modbus.INPUT_REGISTER)
		if err != nil {
			fmt.Printf("read failed: %v\n", err)
			// Return an error JSON response to the client so it doesn't hang
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read Modbus register"})
			return
		}

		// 2. Safely print the value (note: you fetched address 0, not 100-109)
		fmt.Printf("Register 0 value: %v\n", outputnum)

		// 3. Convert uint16 to standard int, then to a string
		iOutNum := strconv.Itoa(int(outputnum))

		// 4. Send the successful JSON response using clean key-value syntax
		c.JSON(http.StatusOK, gin.H{
			"OutputNumber": iOutNum,
		})
	})



	r.POST("/ModuleOneWord", func(c *gin.Context) {
		var DIO IO

		// Bind JSON payload to the struct; return 400 if validation fails
		if err := c.ShouldBindJSON(&DIO); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ModbusPayload,converr := strconv.Atoi(DIO.Value)
		if converr !=nil {
			fmt.Println("Conversion Failed:", converr)
			return
		}


		testerr := client.WriteRegister(0, uint16(ModbusPayload))
		if testerr != nil {
			fmt.Printf("Write failed: %v\n", testerr)
		} else {
			fmt.Printf("Write DOModuleOne: %v\n", ModbusPayload)
		}

		// Return 201 Created on success
		c.JSON(http.StatusCreated, gin.H{
			"status": "Modbus Value Changed",
			"value":   DIO.Value,
		})
	})

	// Start the server on port 8080
	r.Run(":8080")
}
