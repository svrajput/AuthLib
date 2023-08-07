package domain

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/svrajput/golang-banking/errs"
	"github.com/svrajput/golang-banking/logger"
)

type CustomerRepositoryDb struct {
	client *sqlx.DB
}

// helper method.
func NewCustomerRepositoryDb(dbclient *sqlx.DB) CustomerRepositoryDb {
	return CustomerRepositoryDb{dbclient}
}

// Filter active or inactive - Status Filter added.
func (d CustomerRepositoryDb) FindAll(status string) ([]Customer, *errs.AppError) {
	var err error
	customers := make([]Customer, 0)

	if status == "" {
		findAllSql := "select customer_id, name, city, zipcode, date_of_birth, status from customers"
		err = d.client.Select(&customers, findAllSql)
	} else {
		findAllSql := "select customer_id, name, city, zipcode, date_of_birth, status from customers where status = ?"
		err = d.client.Select(&customers, findAllSql, status)
	}

	if err != nil {
		logger.Error("Error while querying customers table " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}

	return customers, nil
}

func (d CustomerRepositoryDb) FindById(id string) (*Customer, *errs.AppError) {
	var c Customer
	findCustomerSql := "select customer_id, name, city, zipcode, date_of_birth, status from customers where customer_id = ?"

	err := d.client.Get(&c, findCustomerSql, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.NewNotFoundError("Customer not found")
		} else {
			logger.Error("Error while scanning customer " + err.Error())
			return nil, errs.NewUnexpectedError("Unexpected database error")
		}
	}
	return &c, nil
}
