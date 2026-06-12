package main

import (
	"fmt"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/simonvetter/modbus"
	"net/http"
	"strconv"
)

// IO represents the expected JSON payload matching your existing struct
type IO struct {
	Value string `json:"value"` // Ensure json tag matches your incoming payload casing
}

// MakeModbusWriteHandler creates a Gin handler tailored to a specific Modbus register address.
func MakeModbusWriteHandler(client *modbus.ModbusClient, registerAddress uint16) gin.HandlerFunc {
	return func(c *gin.Context) {
		var DIO IO

		// 1. Bind JSON payload to the struct
		if err := c.ShouldBindJSON(&DIO); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format: " + err.Error()})
			return
		}

		// 2. Convert incoming string value to integer
		modbusPayload, converr := strconv.Atoi(DIO.Value)
		if converr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Value must be a valid integer string"})
			return
		}

		// 3. Write value to the specified Modbus address
		err := client.WriteRegister(registerAddress, uint16(modbusPayload))
		if err != nil {
			fmt.Errorf("Modbus write failed at address %d: %v\n", registerAddress, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write data to Modbus device"})
			return
		}

		// Log success on the server terminal
		fmt.Printf("Successfully wrote value %d to Modbus register %d\n", modbusPayload, registerAddress)

		// 4. Return HTTP 201 Created on success
		c.JSON(http.StatusCreated, gin.H{
			"status":           "Modbus Value Changed",
			"register_address": registerAddress,
			"value":            DIO.Value,
		})
	}
}


func main() {
    // ... setup your gin router 'r' and your modbus 'client' ...

    // Deploy endpoints pointing to different Modbus registers effortlessly
    r.POST("/ModuleOneWord", MakeModbusWriteHandler(client, 1))
    r.POST("/ModuleTwoWord", MakeModbusWriteHandler(client, 2))
    r.POST("/ModuleThreeWord", MakeModbusWriteHandler(client, 3))

    
    r.Run(":8080")
}
