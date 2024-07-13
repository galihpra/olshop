package handler

type CartRequest struct {
	ProductId uint  `json:"id_product"`
	VarianId  uint  `json:"id_varian"`
	Quantity  int16 `json:"qty"`
}
