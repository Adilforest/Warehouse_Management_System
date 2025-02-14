package main

import (
    "bytes"
    "context"
    "fmt"
    "html/template"
    "net/http"
    "os"
    "regexp"
    "strconv"
    "strings"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/jung-kurt/gofpdf"
    "github.com/sirupsen/logrus"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "gopkg.in/gomail.v2"
)

// Transaction описывает транзакцию оплаты
type Transaction struct {
    ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    TransactionID string            `bson:"transaction_id" json:"transaction_id"`
    UserID       string             `bson:"user_id" json:"user_id"`
    Cart         interface{}        `bson:"cart" json:"cart"`
    Status       string             `bson:"status" json:"status"` // pending, paid, declined
    CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
    UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}

type PaymentForm struct {
    CardNumber     string `form:"cardNumber" binding:"required"`
    ExpirationDate string `form:"expirationDate" binding:"required"`
    CVV            string `form:"cvv" binding:"required"`
    Name           string `form:"name" binding:"required"`
    Address        string `form:"address" binding:"required"`
    Email          string `form:"email" binding:"required"`
}

var (
    mongoClient *mongo.Client
    db          *mongo.Database
)

// initMongo подключается к MongoDB
func initMongo() {
    mongoURI := os.Getenv("MONGO_URI")
    if mongoURI == "" {
        mongoURI = "mongodb://localhost:27017"
    }
    mongoDBName := os.Getenv("MONGO_DB")
    if mongoDBName == "" {
        mongoDBName = "warehouse"
    }
    clientOptions := options.Client().ApplyURI(mongoURI)
    client, err := mongo.Connect(context.Background(), clientOptions)
    if err != nil {
        logrus.Fatalf("Failed to connect to MongoDB: %v", err)
    }
    mongoClient = client
    db = client.Database(mongoDBName)
}

func main() {
    initMongo()
    router := gin.Default()

    // Эндпоинт для создания транзакции (вызывается основным сервером)
    router.POST("/transactions", createTransaction)

    // Эндпоинт для отображения формы оплаты (GET запрос)
    router.GET("/transactions/:id/payment", showPaymentForm)

    // Эндпоинт для обработки данных оплаты (POST запрос из формы)
    router.POST("/transactions/:id/pay", processPayment)

    // Запуск микросервиса на порту 8081
    port := os.Getenv("PORT")
    if port == "" {
        port = "8081"
    }
    router.Run(":" + port)
}

// createTransaction принимает транзакционные данные от основного сервера,
// создаёт запись со статусом "pending" и возвращает URL формы оплаты.
func createTransaction(c *gin.Context) {
    var payload map[string]interface{}
    if err := c.BindJSON(&payload); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
        return
    }
    transactionID, ok := payload["transaction_id"].(string)
    if !ok {
        c.JSON(http.StatusBadRequest, gin.H{"error": "transaction_id is required"})
        return
    }
    userID, ok := payload["user_id"].(string)
    if !ok {
        c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
        return
    }
    // По умолчанию статус "pending"
    status := "pending"
    // Если переданы данные платежа, симулируем проверку оплаты
    if pd, exists := payload["payment_details"]; exists {
        paymentDetails, ok := pd.(map[string]interface{})
        if ok {
            if cardNumber, ok := paymentDetails["cardNumber"].(string); ok && len(cardNumber) > 0 && cardNumber[0] == '4' {
                status = "paid"
            } else {
                status = "declined"
            }
        }
    }
    txn := Transaction{
        TransactionID: transactionID,
        UserID:        userID,
        Cart:          payload["cart"],
        Status:        status,
        CreatedAt:     time.Now(),
        UpdatedAt:     time.Now(),
    }
    collection := db.Collection("transactions")
    res, err := collection.InsertOne(context.Background(), txn)
    if err != nil {
        logrus.Errorf("Failed to create transaction: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction"})
        return
    }
    insertedID := res.InsertedID.(primitive.ObjectID).Hex()
    paymentURL := fmt.Sprintf("http://%s/transactions/%s/payment", c.Request.Host, insertedID)
    // Оплата считается успешной, если статус "paid"
    success := (status == "paid")
    c.JSON(http.StatusOK, gin.H{
        "message":    "Transaction processed.",
        "paymentURL": paymentURL,
        "success":    success,
    })
}

// processPayment обрабатывает данные, введённые в форме оплаты,
// симулирует проверку и обновляет статус транзакции.
func processPayment(c *gin.Context) {
    transactionID := c.Param("id")
    var txn Transaction
    collection := db.Collection("transactions")
    objID, err := primitive.ObjectIDFromHex(transactionID)
    if err != nil {
        c.String(http.StatusBadRequest, "Invalid transaction ID")
        return
    }
    err = collection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&txn)
    if err != nil {
        c.String(http.StatusNotFound, "Transaction not found")
        return
    }
    var form PaymentForm
    if err := c.ShouldBind(&form); err != nil {
        c.String(http.StatusBadRequest, "Invalid form data")
        return
    }
    // Проверка CVV: должно состоять ровно из 3 цифр
    cvvRegex := regexp.MustCompile(`^\d{3}$`)
    if !cvvRegex.MatchString(form.CVV) {
        c.String(http.StatusBadRequest, "Invalid CVV: must consist of exactly 3 digits")
        return
    }
    // Проверка срока действия: формат MM/YY, где MM от 01 до 12
    expirationRegex := regexp.MustCompile(`^(0[1-9]|1[0-2])\/\d{2}$`)
    if !expirationRegex.MatchString(form.ExpirationDate) {
        c.String(http.StatusBadRequest, "Invalid expiration date format. Use MM/YY.")
        return
    }
    // Проверка номера карты: удаляем пробелы и проверяем, что состоит только из цифр и имеет длину от 13 до 19
    normalizedCard := strings.ReplaceAll(form.CardNumber, " ", "")
    if len(normalizedCard) < 13 || len(normalizedCard) > 19 {
        c.String(http.StatusBadRequest, "Invalid card number length.")
        return
    }
    if _, err := strconv.ParseInt(normalizedCard, 10, 64); err != nil {
        c.String(http.StatusBadRequest, "Invalid card number: must be numeric.")
        return
    }
    // Симуляция проверки оплаты: если номер карты (без пробелов) начинается с "4", считаем оплату успешной.
    success := false
    if len(normalizedCard) > 0 && normalizedCard[0] == '4' {
        success = true
    }
    newStatus := "declined"
    if success {
        newStatus = "paid"
    }
    update := bson.M{
        "$set": bson.M{
            "status":     newStatus,
            "updated_at": time.Now(),
        },
    }
    _, err = collection.UpdateOne(context.Background(), bson.M{"_id": objID}, update)
    if err != nil {
        c.String(http.StatusInternalServerError, "Failed to update transaction status")
        return
    }
    // При успешной оплате – генерируем PDF чека и отправляем на email, используя form.Email
    if success {
        logrus.Infof("Payment successful for transaction %s. Generating receipt and sending email.", txn.TransactionID)
        receiptData := map[string]interface{}{
            "total": 100.00, // Пример суммы, заменить на актуальную
        }
        filename, err := generateReceiptPDF(txn, receiptData)
        if err != nil {
            c.String(http.StatusInternalServerError, "Failed to generate receipt")
            return
        }
        smtpHost := os.Getenv("SMTP_HOST")
        smtpPort, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
        smtpEmail := os.Getenv("SMTP_EMAIL")
        smtpPassword := os.Getenv("SMTP_PASSWORD")

        err = sendReceiptEmail(smtpHost, smtpPort, smtpEmail, smtpPassword, form.Email, "Your Payment Receipt", "Please find attached your receipt.", filename)
        if err != nil {
            c.String(http.StatusInternalServerError, "Failed to send receipt email")
            return
        }
        logrus.Infof("Receipt generated and sent to %s for transaction %s", form.Email, txn.TransactionID)
    }
    statusText := "failed"
    if success {
        statusText = "successful"
    }
    resultHTML := `Payment %sTransaction %s has been %s.`
    c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(fmt.Sprintf(resultHTML, statusText, txn.TransactionID, statusText)))
}

// generateReceiptPDF генерирует PDF-чек и сохраняет его в файл, возвращая имя файла.
func generateReceiptPDF(txn Transaction, receiptData map[string]interface{}) (string, error) {
    pdf := gofpdf.New("P", "mm", "A4", "")
    pdf.AddPage()
    pdf.SetFont("Arial", "B", 16)
    pdf.Cell(40, 10, "Marketplace Receipt")
    pdf.Ln(12)
    pdf.SetFont("Arial", "", 12)
    pdf.Cell(40, 10, fmt.Sprintf("Transaction ID: %s", txn.TransactionID))
    pdf.Ln(8)
    pdf.Cell(40, 10, fmt.Sprintf("Date: %s", time.Now().Format("02 Jan 2006 15:04")))
    pdf.Ln(8)
    // Здесь можно добавить детали транзакции, товаров, суммы, ФИО клиента, способ оплаты и т.д.
    pdf.Cell(40, 10, fmt.Sprintf("Total Amount: $%.2f", receiptData["total"].(float64)))
    // Дополнительно: список товаров, данные карты (зашифрованные), ФИО и т.д.
    filename := fmt.Sprintf("receipt_%s.pdf", txn.TransactionID)
    err := pdf.OutputFileAndClose(filename)
    if err != nil {
        return "", err
    }
    return filename, nil
}

func sendReceiptEmail(host string, port int, fromEmail, password, toEmail, subject, body, attachmentPath string) error {
    m := gomail.NewMessage()
    m.SetHeader("From", fromEmail)
    m.SetHeader("To", toEmail)
    m.SetHeader("Subject", subject)
    m.SetBody("text/plain", body)
    m.Attach(attachmentPath)

    d := gomail.NewDialer(host, port, fromEmail, password)
    // Для Gmail возможно потребуется использовать App Password и разрешить "less secure apps"
    return d.DialAndSend(m)
}

func showPaymentForm(c *gin.Context) {
    transactionID := c.Param("id")
    var txn Transaction
    collection := db.Collection("transactions")
    objID, err := primitive.ObjectIDFromHex(transactionID)
    if err != nil {
        c.String(http.StatusBadRequest, "Invalid transaction ID")
        return
    }
    err = collection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&txn)
    if err != nil {
        c.String(http.StatusNotFound, "Transaction not found")
        return
    }
    html := `
        Payment Form
Payment Form for Transaction {{.TransactionID}}
Card Number:
Expiration Date (MM/YY):
CVV:
Name:
Address:
Email:
    `
    tmpl, err := template.New("paymentForm").Parse(html)
    if err != nil {
        c.String(http.StatusInternalServerError, "Error rendering form")
        return
    }
    var buf bytes.Buffer
    tmpl.Execute(&buf, txn)
    c.Data(http.StatusOK, "text/html; charset=utf-8", buf.Bytes())
}