package models

import (
	"UrbanNest/db"
	"UrbanNest/utils"
	"errors"
	"fmt"
)

type Customer struct {
	Id        int
	Email     string `binding:"required"`
	Password  string `binding:"required"` // only for input, not stored
	FullName  string `binding:"required" json:"full_name"`
	CreatedAt string
}

func AddCustomer(customer Customer) error {
	exists, err := db.DB.Query("SELECT 1 FROM customers WHERE email = $1", customer.Email)
	if err != nil {
		return fmt.Errorf("DB query error: %v", err)
	}
	defer exists.Close()

	if exists.Next() {
		return errors.New("email already registered")
	}

	hashedPassword, err := utils.HashPassword(customer.Password)
	if err != nil {
		return err
	}

	query := `INSERT INTO customers (email, password_hash, full_name) VALUES ($1, $2, $3)`
	_, err = db.DB.Exec(query, customer.Email, hashedPassword, customer.FullName)
	return err
}

func (c *Customer) ValidateCustomer() error {
	query := `SELECT password_hash, id FROM customers WHERE email = $1`
	row := db.DB.QueryRow(query, c.Email)

	var retrievedHash string
	err := row.Scan(&retrievedHash, &c.Id)
	if err != nil {
		fmt.Println("DB Error:", err)
		return errors.New("invalid email or password")
	}

	passwordIsValid := utils.CheckPasswordHash(c.Password, retrievedHash)
	if !passwordIsValid {
		return errors.New("invalid email or password")
	}

	return nil
}
