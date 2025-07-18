package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ahsmha/discount-service/internal/models"
	"github.com/ahsmha/discount-service/internal/repository"
	"github.com/ahsmha/discount-service/internal/service"
	"github.com/ahsmha/discount-service/testdata"
)

func main() {
	repo := repository.NewMemoryDiscountRepository()
	memoryRepo, ok := repo.(interface{ SeedDiscounts([]models.Discount) error })
	if !ok {
		log.Fatal("Repository does not support seeding")
	}

	err := memoryRepo.SeedDiscounts(testdata.GetSampleDiscounts())
	if err != nil {
		log.Fatalf("Failed to seed discounts: %v", err)
	}

	discountService := service.NewDiscountService(repo)

	runDemonstration(discountService)
}

func runDemonstration(discountService service.DiscountService) {
	ctx := context.Background()

	fmt.Println("\n📋 Running Multiple Discount Scenario Demonstration")
	fmt.Println("Scenario: PUMA T-shirt with brand, category, and bank discounts")
	fmt.Println("---------------------------------------------------------------")

	cartItems, customer, paymentInfo := testdata.GetMultipleDiscountScenario()

	// Reset the PUMA T-shirt price to base price for demonstration
	cartItems[0].Product.CurrentPrice = cartItems[0].Product.BasePrice

	fmt.Printf("Cart Items:\n")
	for _, item := range cartItems {
		fmt.Printf("- %s %s (%s) x%d @ ₹%s each\n",
			item.Product.Brand.ID,
			item.Product.Category.ID,
			item.Size,
			item.Quantity,
			item.Product.CurrentPrice.String())
	}

	fmt.Printf("\nCustomer: %s (Tier: %s)\n", customer.ID, customer.Tier)
	if paymentInfo != nil {
		fmt.Printf("Payment: %s", paymentInfo.Method)
		if paymentInfo.BankName != nil {
			fmt.Printf(" (%s)", *paymentInfo.BankName)
		}
		fmt.Println()
	}

	// Calculate discounts
	result, err := discountService.CalculateCartDiscounts(ctx, cartItems, customer, paymentInfo)
	if err != nil {
		log.Fatalf("Failed to calculate discounts: %v", err)
	}

	// Display results
	fmt.Println("\n💰 Discount Calculation Results")
	fmt.Println("------------------------------")
	fmt.Printf("Original Price: ₹%s\n", result.OriginalPrice.String())
	fmt.Printf("Final Price: ₹%s\n", result.FinalPrice.String())
	fmt.Printf("Total Savings: ₹%s (%.2s%%)\n",
		result.GetTotalDiscount().String(),
		result.GetDiscountPercentage().String())

	fmt.Println("\n🎯 Applied Discounts:")
	for name, amount := range result.AppliedDiscounts {
		fmt.Printf("- %s: ₹%s\n", name, amount.String())
	}

	fmt.Printf("\nMessage: %s\n", result.Message)

	// Test discount code validation
	fmt.Println("\n🔍 Testing Discount Code Validation")
	fmt.Println("-----------------------------------")

	testCodes := []string{"SUPER69", "PREMIUM15", "INVALID123", ""}

	for _, code := range testCodes {
		if code == "" {
			fmt.Printf("Testing empty code: ")
		} else {
			fmt.Printf("Testing code '%s': ", code)
		}

		isValid, err := discountService.ValidateDiscountCode(ctx, code, cartItems, customer)
		if err != nil {
			fmt.Printf("Error - %v\n", err)
		} else if isValid {
			fmt.Println("✅ Valid")
		} else {
			fmt.Println("❌ Invalid")
		}
	}

	fmt.Println("\n✨ Demonstration completed successfully!")
}
