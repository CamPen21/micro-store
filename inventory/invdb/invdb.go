package invdb

import (
	"database/sql"
	"errors"
	"log"

	"github.com/mattn/go-sqlite3"
)

type InventoryItem struct {
	Id        int64
	Name      string
	Stock     int64
	Allocated int64
}

type QueryRower interface {
	QueryRow(query string, args ...interface{}) *sql.Row
}

func Initialize(db *sql.DB) {
	_, err := db.Exec(`
		CREATE TABLE inventory_item(
			id			INTEGER PRIMARY KEY,
			name 		TEXT UNIQUE,
			stock		INTEGER,
			allocated 	INTEGER
		);
	`)
	if err != nil {
		if sqlError, ok := err.(sqlite3.Error); ok {
			if sqlError.Code != 1 {
				log.Fatal(sqlError)
			}
		} else {
			log.Fatal(err)
		}
	}
}

func inventoryItemFromRow(row *sql.Row) (*InventoryItem, error) {
	var id int64
	var name string
	var stock int64
	var allocated int64

	err := row.Scan(&id, &name, &stock, &allocated)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &InventoryItem{id, name, stock, allocated}, nil
}

func fetchInventoryItem(cx QueryRower, itemId int64) (*InventoryItem, error) {
	// Query the database
	row := cx.QueryRow(`
		SELECT
			id,
			name,
			stock,
			allocated
		FROM inventory_item	
		WHERE id = ?
		LIMIT 1
	`, itemId)
	// Check if the query returned an error
	if item, err := inventoryItemFromRow(row); err != nil {
		// Error returned from row parsing
		log.Println(err)
		return nil, err
	} else {
		// Successful parsing
		return item, nil
	}

}

func GetInventoryItem(db *sql.DB, itemId int64) (*InventoryItem, error) {
	// Opens a new transaction
	item, err := fetchInventoryItem(db, itemId)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return item, nil
}

func AllocateInventoryItemAmount(db *sql.DB, itemId, amount int64) error {
	// Opens a new transaction
	tx, err := db.Begin()
	if err != nil {
		log.Println(err)
		return err
	}
	// Defers a rollback in case something goes wrong
	defer tx.Rollback()
	item, err := fetchInventoryItem(tx, itemId)
	if err != nil {
		log.Println(err)
		return err
	}
	enoughStock := (item.Stock - item.Allocated) >= amount
	enoughAllocated := item.Allocated >= amount
	if enoughStock && enoughAllocated {
		// Execute the update
		_, err = tx.Exec(`
			UPDATE inventory_item
			SET allocated = allocated + ?
			WHERE id = ?;
		`, amount, itemId)
		if err != nil {
			log.Println(err)
			return err
		}
		// Commit the transaction
		if err = tx.Commit(); err != nil {
			log.Println(err)
			return err
		}
		return nil
	} else {
		return errors.New("Insufficient stock")
	}
}
