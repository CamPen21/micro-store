package orderdb

import (
	"database/sql"
	"log"
	"strings"
	"time"

	"github.com/mattn/go-sqlite3"
)

type OrderPlacement struct {
	Id        int64
	Status    string
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

type OrderedItem struct {
	ItemId   int64
	Quantity int64
}

type OrderItem struct {
	OrderedItem
	Id      int64
	OrderId int64
}

type QueryRower interface {
	QueryRow(query string, args ...interface{}) *sql.Row
}

func Initialize(db *sql.DB) {
	_, err := db.Exec(`
		CREATE TABLE order_placement(
			id			INTEGER PRIMARY KEY,
			status 		TEXT,
			created_at	INTEGER DEFAULT (STRFTIME('%s', 'now')),
			updated_at 	INTEGER DEFAULT (STRFTIME('%s', 'now'))
		);
		CREATE TABLE order_item(
			id			INTEGER PRIMARY KEY,
			order_id 	INTEGER ,
			item_id		INTEGER,
			quantity 	INTEGER,
			FOREIGN KEY(order_id) REFERENCES order_placement(id)
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

func orderPlacementFromRow(row *sql.Row) (*OrderPlacement, error) {
	var id int64
	var name string
	var createdAt int64
	var updatedAt int64

	err := row.Scan(&id, &name, &createdAt, &updatedAt)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	timeCreatedAt := time.Unix(createdAt, 0)
	timeUpdatedAt := time.Unix(updatedAt, 0)
	return &OrderPlacement{id, name, &timeCreatedAt, &timeUpdatedAt}, nil
}

func orderItemFromRow(row *sql.Row) (*OrderItem, error) {
	var id int64
	var orderId int64
	var itemId int64
	var quantity int64
	err := row.Scan(&id, &orderId, &itemId, &quantity)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &OrderItem{OrderedItem{itemId, quantity}, id, orderId}, nil
}

func CreateOrder(tx *sql.Tx, orderedItems []*OrderedItem) (int64, error) {
	defer tx.Rollback()
	res, err := tx.Exec(`
		INSERT INTO order_placement(status, updated_at)
		VALUES(?, STRFTIME('%s', 'now'));
	`, "RECEIVED")
	if err != nil {
		log.Println(err)
		return -1, err
	}
	orderId, err := res.LastInsertId()
	if err != nil {
		log.Println(err)
		return -1, err
	}
	if err = AddOrderItems(tx, orderId, orderedItems); err != nil {
		log.Println(err)
		return -1, err
	} else {
		tx.Commit()
		return orderId, nil
	}
}

func AddOrderItems(tx *sql.Tx, orderId int64, orderItems []*OrderedItem) error {
	orderItemsLen := len(orderItems)
	inserts := make([]string, orderItemsLen)
	values := make([]interface{}, orderItemsLen*3)
	for i := 0; i < orderItemsLen; i++ {
		if i == orderItemsLen-1 {
			inserts[i] = "(?, ?, ?)\n"
		} else {
			inserts[i] = "(?, ?, ?),\n"
		}
		valuesIndex := i * 3
		values[valuesIndex] = orderId
		values[valuesIndex+1] = orderItems[i].ItemId
		values[valuesIndex+2] = orderItems[i].Quantity

	}
	valuesSection := strings.Join(inserts, "")
	log.Println(values)
	_, err := tx.Exec("INSERT INTO order_item(order_id, item_id, quantity)\nVALUES"+valuesSection+";", values...)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func CancelOrder(tx *sql.Tx, orderId int64) error {
	_, err := tx.Exec(`
		UPDATE order_placement
		SET status = 'CANCELED'
		WHERE id = ?;
	`, orderId)
	if err != nil {
		log.Println(err)
		return err
	}
	tx.Commit()
	return nil
}
