package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Represents the response structure for the /order endpoint
type orderResponse struct {
	Orders []order `json:"Orders,omitempty"`
	Order  *order  `json:"Order,omitempty"`
}

// Represents an order with its items and correct packs
type order struct {
	Items        int               `json:"Items ordered"`
	CorrectPacks []correctPackItem `json:"Correct number of Packs"`
	Error        string            `json:"error,omitempty"`
}

// Represents an item in the correct packs
type correctPackItem struct {
	NumberOfPacks int `json:"Number of Packs"`
	PackSize      int `json:"Pack Size"`
}

// Calculates the correct number of packs needed to fulfill an order
func correctCalculatePacks(items int) []correctPackItem {
	packSizes := []int{5000, 2000, 1000, 500, 250}
	correctPacks := []correctPackItem{}

	for _, packSize := range packSizes {
		packs := items / packSize // Calculate the number of packs for the current pack size

		if packs > 0 {
			if packSize == 250 && items > 250 && items < 500 {
				// Special case for pack size 250 when items are between 250 and 500
				correctPacks = append(correctPacks, correctPackItem{
					NumberOfPacks: 1,
					PackSize:      500,
				})
				items -= 500 // Reduce the remaining items by 500
			} else {
				correctPacks = append(correctPacks, correctPackItem{
					NumberOfPacks: packs,
					PackSize:      packSize,
				})
				items -= packs * packSize // Reduce the remaining items by the total value of the calculated packs
			}
		}
	}

	if items > 0 {
		// If there are remaining items, add a correctPackItem with 1 pack of size 250
		correctPacks = append(correctPacks, correctPackItem{
			NumberOfPacks: 1,
			PackSize:      250,
		})
	}

	return correctPacks
}

// Validates user input for the number of items and checks for invalid inputs
func getInputAndValidate(c *gin.Context) ([]int, error) {
	// Retrieve the items from the request form
	itemsStr := c.PostFormArray("items")

	if len(itemsStr) == 0 {
		return nil, fmt.Errorf("No items provided. Please provide a valid list of items") // Return an error message
	}

	// Validate each item in the input
	items := make([]int, len(itemsStr)) // Create a slice to store the validated items
	for i, itemStr := range itemsStr {
		item, err := strconv.Atoi(itemStr) // Convert the item string to an integer
		if err != nil {
			return nil, fmt.Errorf("Invalid input. Please provide a valid number of items") // Return an error message for invalid input
		}

		// Check for specific conditions related to the input item
		if item == 0 {
			return nil, fmt.Errorf("Invalid input. You cannot order 0 items. Please provide a valid number of items")
		}

		if item < 0 {
			return nil, fmt.Errorf("Invalid input. You cannot order negative items. Please provide a valid number of items")
		}

		items[i] = item // Store the validated item in the items slice
	}

	return items, nil // Return the validated items slice and nil to indicate no errors
}

func main() {
	router := gin.Default()

	// Define the "/order" endpoint
	router.POST("/order", func(c *gin.Context) {
		// Get and validate the input items from the request
		items, err := getInputAndValidate(c)
		if err != nil {
			// If there's an error in the input, return a JSON response with the error message and a 400 status code
			c.IndentedJSON(http.StatusBadRequest, gin.H{"Error message": err.Error()})
			return
		}

		if len(items) == 1 {
			// If there's only one item in the input, calculate the correct packs for that item
			item := items[0]
			correctPacks := correctCalculatePacks(item)
			order := order{
				Items:        item,
				CorrectPacks: correctPacks,
			}

			response := orderResponse{
				Order: &order,
			}

			// Return a JSON response with the single order and a 200 status code
			c.IndentedJSON(http.StatusOK, response)
		} else {
			// If there are multiple items in the input, calculate the correct packs for each item
			orders := make([]order, len(items))
			for i, item := range items {
				correctPacks := correctCalculatePacks(item)
				order := order{
					Items:        item,
					CorrectPacks: correctPacks,
				}
				orders[i] = order
			}

			response := orderResponse{
				Orders: orders,
			}

			// Return a JSON response with the multiple orders and a 200 status code
			c.IndentedJSON(http.StatusOK, response)
		}
	})

	// Start the server and listen on port 8080
	router.Run(":8080")
}
