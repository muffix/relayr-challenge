package database

import (
	"database/sql"

	"github.com/pkg/errors"
)

const (
	createOffersTableStmt = "CREATE TABLE IF NOT EXISTS offers (product TEXT NOT NULL, category TEXT NOT NULL, supplier TEXT NOT NULL, price REAL NOT NULL, PRIMARY KEY (product, category, supplier))"
	insertOfferStmt       = "INSERT INTO offers (product, category, supplier, price) VALUES (?, ?, ?, ?) ON CONFLICT(product, category, supplier) DO UPDATE SET price=EXCLUDED.price"
	getOfferQuery         = "SELECT product, category, supplier, price FROM offers WHERE product=? AND category=? ORDER BY price ASC"
)

// Offers is an interface for a database client
type Offers interface {
	Insert(productName, categoryName, supplierName string, price float32) error
	InsertMultiple(offers []Offer) error
	Get(productName, categoryName string) ([]Offer, error)
	Close() error
}

// OffersSQLiteDatabase is a database client using SQLite
type OffersSQLiteDatabase sql.DB

// Offer is a struct representing an offer for a product by a supplier
type Offer struct {
	Product, Category, Supplier string
	Price                       float32
}

// InitSQLiteDatabase opens the database and sets it up if needed.
// Returns a database handle.
func InitSQLiteDatabase(dbPath string) (*OffersSQLiteDatabase, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, errors.Wrap(err, "error opening database")
	}

	err = execStatement(db, createOffersTableStmt)
	if err != nil {
		return nil, errors.Wrap(err, "error creating table")
	}

	return (*OffersSQLiteDatabase)(db), nil
}

// Insert inserts an offer into the database
//
// If an offer for an existing product, category and supplier exists, the offer is updated.
func (d *OffersSQLiteDatabase) Insert(productName, categoryName, supplierName string, price float32) error {
	return d.InsertMultiple([]Offer{
		{
			Product:  productName,
			Category: categoryName,
			Supplier: supplierName,
			Price:    price,
		},
	})
}

// InsertMultiple inserts multiple offers into the database in a transaction
//
// If an offer for an existing product, category and supplier exists, the offer is updated.
func (d *OffersSQLiteDatabase) InsertMultiple(offers []Offer) (err error) {
	db := (*sql.DB)(d)
	tx, err := db.Begin()
	if err != nil {
		return errors.Wrap(err, "error beginning transaction")
	}

	// Make sure that we commit the transaction or rollback in case of an error
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		commitErr := tx.Commit()
		if commitErr != nil {
			err = errors.Wrap(commitErr, "error committing transaction")
		}
	}()

	stmt, err := tx.Prepare(insertOfferStmt)
	if err != nil {
		return errors.Wrap(err, "error preparing insert statement")
	}

	for _, offer := range offers {
		_, err = tx.Stmt(stmt).Exec(offer.Product, offer.Category, offer.Supplier, offer.Price)
		if err != nil {
			return errors.Wrap(err, "error inserting offer")
		}
	}

	return
}

// Get returns all offers for a given product in a category
func (d *OffersSQLiteDatabase) Get(productName, categoryName string) ([]Offer, error) {
	db := (*sql.DB)(d)
	rows, err := db.Query(getOfferQuery, productName, categoryName)
	if err != nil {
		return []Offer{}, err
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	var offers []Offer

	for rows.Next() {
		offer := Offer{}
		err = rows.Scan(&offer.Product, &offer.Category, &offer.Supplier, &offer.Price)
		if err != nil {
			return []Offer{}, errors.Wrap(err, "error retrieving row")
		}
		offers = append(offers, offer)
	}
	return offers, nil
}

// Close closes the database connection
func (d *OffersSQLiteDatabase) Close() error {
	return (*sql.DB)(d).Close()
}

func execStatement(db *sql.DB, stmt string) error {
	statement, err := db.Prepare(stmt)
	if err != nil {
		return errors.Wrap(err, "error preparing statement")
	}

	_, err = statement.Exec()

	if err != nil {
		return errors.Wrap(err, "error executing statement")
	}

	return nil
}
