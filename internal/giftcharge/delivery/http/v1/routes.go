package v1

// MapRoutes giftCharges REST API routes
func (h *giftHandlers) MapRoutes() {
	h.group.POST("", h.Create())
	h.group.GET("/:gift_id", h.GetByID())
	h.group.GET("", h.GetList())
}
