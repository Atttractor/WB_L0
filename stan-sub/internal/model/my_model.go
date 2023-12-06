package model

type Delivery struct {
	Delivery_uuid string `json:"delivery_uuid" validate:"required"`
	Name          string `json:"name" validate:"required"`
	Phone         string `json:"phone" validate:"required,e164"`
	Zip           string `json:"zip" validate:"required"`
	City          string `json:"city" validate:"required"`
	Address       string `json:"address" validate:"required"`
	Region        string `json:"region" validate:"required"`
	Email         string `json:"email" validate:"required,email"`
}

type Payment struct {
	Transaction   string `json:"transaction" validate:"required"`
	Request_id    string `json:"request_id" validate:""`
	Currency      string `json:"currency" validate:"required"`
	Provider      string `json:"provider" validate:"required"`
	Amount        int    `json:"amount" validate:"required,min=0"`
	Payment_dt    int    `json:"payment_dt" validate:"required,min=0"`
	Bank          string `json:"bank" validate:"required"`
	Delivery_cost int    `json:"delivery_cost" validate:"required,min=0"`
	Goods_total   int    `json:"goods_total" validate:"required,min=0"`
	Custom_fee    int    `json:"custom_fee" validate:"min=0"`
}

type Item struct {
	Chrt_id      int    `json:"chrt_id" validate:"required,min=0"`
	Track_number string `json:"track_number" validate:"required"`
	Price        int    `json:"price" validate:"required,min=0"`
	Rid          string `json:"rid" validate:"required"`
	Name         string `json:"name" validate:"required"`
	Sale         int    `json:"sale" validate:"required,min=0"`
	Size         string `json:"size" validate:"required,number"`
	Total_price  int    `json:"total_price" validate:"required,min=0"`
	Nm_id        int    `json:"nm_id" validate:"required,min=0"`
	Brand        string `json:"brand" validate:"required"`
	Status       int    `json:"status" validate:"required,min=0"`
}

type MyModel struct {
	Order_uid          string   `json:"order_uid" validate:"required"`
	Track_number       string   `json:"track_number" validate:"required"`
	Entry              string   `json:"entry" validate:"required"`
	Delivery           Delivery `json:"delivery"`
	Payment            Payment  `json:"payment"`
	Items              []Item   `json:"items"`
	Locale             string   `json:"locale" validate:"required"`
	Internal_signature string   `json:"internal_signature"`
	Customer_id        string   `json:"customer_id" validate:"required"`
	Delivery_service   string   `json:"delivery_service" validate:"required"`
	Shardkey           string   `json:"shardkey" validate:"required,number"`
	Sm_id              int      `json:"sm_id" validate:"required,min=0"`
	Date_created       string   `json:"date_created" validate:"required"`
	Oof_shard          string   `json:"oof_shard" validate:"required,number"`
}
