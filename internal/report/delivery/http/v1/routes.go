package v1

// MapRoutes wallets REST API routes
func (h *reportHandlers) MapRoutes() {
	h.group.GET("/:gift_code", h.GetGiftCodeReport())
	h.group.GET("/:mobile", h.GetGiftCodeReport())
}
