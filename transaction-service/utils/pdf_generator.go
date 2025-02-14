package utils

import (
	"fmt"

	"github.com/jung-kurt/gofpdf"
	"go.mongodb.org/mongo-driver/bson/primitive/models"
)

func GenerateFiscalReceipt(transaction models.Transaction) (string, error) {
    pdf := gofpdf.New("P", "mm", "A4", "")
    pdf.AddPage()
    pdf.SetFont("Arial", "B", 16)

    // Добавление данных в PDF
    pdf.Cell(40, 10, "Фискальный чек")
    pdf.Ln(10)
    pdf.Cell(40, 10, fmt.Sprintf("Номер транзакции: %s", transaction.ID.Hex()))
    pdf.Ln(10)
    pdf.Cell(40, 10, fmt.Sprintf("Дата: %s", transaction.CreatedAt.Format("2006-01-02 15:04:05")))
    pdf.Ln(10)
    pdf.Cell(40, 10, fmt.Sprintf("Общая сумма: %.2f", transaction.TotalAmount))

    // Сохранение PDF
    filePath := fmt.Sprintf("receipts/%s.pdf", transaction.ID.Hex())
    err := pdf.OutputFileAndClose(filePath)
    if err != nil {
        return "", err
    }

    return filePath, nil
}