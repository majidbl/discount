package v1

// MapRoutes wallets REST API routes
func (h *discountHandlers) MapRoutes() {
	h.group.POST("", h.DiscountRequest())
}
