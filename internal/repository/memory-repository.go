package repository

import (
	"context"
	"sync"

	"github.com/ahsmha/discount-service/internal/models"
	"github.com/ahsmha/discount-service/pkg/errors"
)

// MemoryDiscountRepository implements DiscountRepository using in-memory storage
type MemoryDiscountRepository struct {
	discounts map[string]*models.Discount
	codeIndex map[string]string // code -> id mapping
	mu        sync.RWMutex
}

// NewMemoryDiscountRepository creates a new in-memory discount repository
func NewMemoryDiscountRepository() DiscountRepository {
	return &MemoryDiscountRepository{
		discounts: make(map[string]*models.Discount),
		codeIndex: make(map[string]string),
	}
}

// GetActiveDiscounts retrieves all active discounts
func (r *MemoryDiscountRepository) GetActiveDiscounts(ctx context.Context) ([]models.Discount, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var activeDiscounts []models.Discount
	for _, discount := range r.discounts {
		if discount.IsValid() {
			activeDiscounts = append(activeDiscounts, *discount)
		}
	}

	return activeDiscounts, nil
}

// GetDiscountByCode retrieves a discount by its code
func (r *MemoryDiscountRepository) GetDiscountByCode(ctx context.Context, code string) (*models.Discount, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	id, exists := r.codeIndex[code]
	if !exists {
		return nil, errors.NewNotFoundError("discount code not found: " + code)
	}

	discount, exists := r.discounts[id]
	if !exists {
		return nil, errors.NewNotFoundError("discount not found for code: " + code)
	}

	return discount, nil
}

// GetDiscountByID retrieves a discount by its ID
func (r *MemoryDiscountRepository) GetDiscountByID(ctx context.Context, id string) (*models.Discount, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	discount, exists := r.discounts[id]
	if !exists {
		return nil, errors.NewNotFoundError("discount not found: " + id)
	}

	return discount, nil
}

// CreateDiscount creates a new discount
func (r *MemoryDiscountRepository) CreateDiscount(ctx context.Context, discount *models.Discount) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if ID already exists
	if _, exists := r.discounts[discount.ID]; exists {
		return errors.NewValidationError("discount already exists: " + discount.ID)
	}

	// Check if code already exists (for voucher discounts)
	if discount.Code != "" {
		if _, exists := r.codeIndex[discount.Code]; exists {
			return errors.NewValidationError("discount code already exists: " + discount.Code)
		}
	}

	// Create a copy to avoid external modifications
	discountCopy := *discount
	r.discounts[discount.ID] = &discountCopy

	// Update code index if applicable
	if discount.Code != "" {
		r.codeIndex[discount.Code] = discount.ID
	}

	return nil
}

// UpdateDiscount updates an existing discount
func (r *MemoryDiscountRepository) UpdateDiscount(ctx context.Context, discount *models.Discount) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if discount exists
	existingDiscount, exists := r.discounts[discount.ID]
	if !exists {
		return errors.NewNotFoundError("discount not found: " + discount.ID)
	}

	// Handle code changes
	if existingDiscount.Code != discount.Code {
		// Remove old code index
		if existingDiscount.Code != "" {
			delete(r.codeIndex, existingDiscount.Code)
		}

		// Add new code index
		if discount.Code != "" {
			if _, exists := r.codeIndex[discount.Code]; exists {
				return errors.NewValidationError("discount code already exists: " + discount.Code)
			}
			r.codeIndex[discount.Code] = discount.ID
		}
	}

	// Update the discount
	discountCopy := *discount
	r.discounts[discount.ID] = &discountCopy

	return nil
}

// DeleteDiscount deletes a discount by ID
func (r *MemoryDiscountRepository) DeleteDiscount(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	discount, exists := r.discounts[id]
	if !exists {
		return errors.NewNotFoundError("discount not found: " + id)
	}

	// Remove from code index if applicable
	if discount.Code != "" {
		delete(r.codeIndex, discount.Code)
	}

	// Remove from main storage
	delete(r.discounts, id)

	return nil
}

// IncrementUsageCount increments the usage count for a discount
func (r *MemoryDiscountRepository) IncrementUsageCount(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	discount, exists := r.discounts[id]
	if !exists {
		return errors.NewNotFoundError("discount not found: " + id)
	}

	// Create a copy with incremented usage count
	updatedDiscount := *discount
	updatedDiscount.UsedCount++
	r.discounts[id] = &updatedDiscount

	return nil
}

// SeedDiscounts seeds the repository with initial discount data
func (r *MemoryDiscountRepository) SeedDiscounts(discounts []models.Discount) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, discount := range discounts {
		discountCopy := discount
		r.discounts[discount.ID] = &discountCopy

		if discount.Code != "" {
			r.codeIndex[discount.Code] = discount.ID
		}
	}

	return nil
}

func (r *MemoryDiscountRepository) ClearDiscounts() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.discounts = make(map[string]*models.Discount)
	r.codeIndex = make(map[string]string)
	return nil
}
