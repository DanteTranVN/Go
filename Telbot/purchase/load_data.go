package purchase

import (
	"encoding/csv"
	"io"
	"os"
	"strconv"
	"time"
)

func loadPurchases() ([]Purchase, error) {
	file, err := os.Open("purchase_records.csv")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var purchases []Purchase
	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		idTele, _ := strconv.Atoi(record[0])
		amount, _ := strconv.Atoi(record[2])
		createdTime, _ := time.Parse("2006-01-02", record[4])

		purchases = append(purchases, Purchase{
			IDTele:      idTele,
			AccountName: record[1],
			Amount:      amount,
			Target:      record[3],
			CreatedTime: createdTime,
		})
	}
	return purchases, nil
}
