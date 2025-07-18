package repositories

import (
	"context"
	"sync"

	"github.com/ahsmha/discounts/internal/interfaces"
	"github.com/ahsmha/discounts/internal/models"
	"github.com/ahsmha/discounts/pkg/errors"
)

// InMemoryDiscountRepository implements DiscountRepository using in-memory storage
type InMemoryDiscountRepository struct {
	discounts map[string]*models.Discount
	codeIndex map[string]string // code -> id mapping
	mu        sync.RWMutex
}

// NewInMemoryDiscountRepository creates a new in-memory discount repository
func NewInMemoryDiscountRepository() interfaces.IDiscountRepository {
	return &InMemoryDiscountRepository{
		discounts: make(map[string]*models.Discount),
		codeIndex: make(map[string]string),
	}
}

// GetActiveDiscounts retrieves all active discounts
func (r *InMemoryDiscountRepository) GetActiveDiscounts(ctx context.Context) ([]models.Discount, error) {
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
func (r *InMemoryDiscountRepository) GetDiscountByCode(ctx context.Context, code string) (*models.Discount, error) {
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
func (r *InMemoryDiscountRepository) GetDiscountByID(ctx context.Context, id string) (*models.Discount, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	discount, exists := r.discounts[id]
	if !exists {
		return nil, errors.NewNotFoundError("discount not found: " + id)
	}

	return discount, nil
}

// CreateDiscount creates a new discount
func (r *InMemoryDiscountRepository) CreateDiscount(ctx context.Context, discount *models.Discount) error {
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
func (r *InMemoryDiscountRepository) UpdateDiscount(ctx context.Context, discount *models.Discount) error {
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
func (r *InMemoryDiscountRepository) DeleteDiscount(ctx context.Context, id string) error {
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
func (r *InMemoryDiscountRepository) IncrementUsageCount(ctx context.Context, id string) error {
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
func (r *InMemoryDiscountRepository) SeedDiscounts(discounts []models.Discount) error {
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

func (r *InMemoryDiscountRepository) ClearDiscounts() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.discounts = make(map[string]*models.Discount)
	r.codeIndex = make(map[string]string)
	return nil
}
