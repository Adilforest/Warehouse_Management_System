package models

import (
    "go.mongodb.org/mongo-driver/bson/primitive"
    "time"
)

type Transaction struct {
    ID           primitive.ObjectID `bson:"_id,omitempty"`
    CartID       string             `bson:"cart_id"` // ID корзины из основного сервера
    UserID       string             `bson:"user_id"` // ID пользователя
    Status       string             `bson:"status"`  // "в ожидании оплаты", "оплачено", "отклонено"
    TotalAmount  float64            `bson:"total_amount"`
    CreatedAt    time.Time          `bson:"created_at"`
    UpdatedAt    time.Time          `bson:"updated_at"`
    PaymentData  PaymentDetails     `bson:"payment_data"`
    FiscalReceipt string            `bson:"fiscal_receipt,omitempty"` // URL или путь к PDF-чеку
}

type PaymentDetails struct {
    CardNumber    string `bson:"card_number"`
    ExpirationDate string `bson:"expiration_date"`
    CVV           string `bson:"cvv"`
    Name          string `bson:"name"`
    Address       string `bson:"address"`
}