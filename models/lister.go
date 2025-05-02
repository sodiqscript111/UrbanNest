package models

import (
	"UrbanNest/db"
	"UrbanNest/utils"
	"errors"
	"fmt"
)

type Lister struct {
	Id          int    `json:"id"`
	Email       string `json:"email" binding:"required"`
	Password    string `json:"password" binding:"required"`
	FullName    string `json:"full_name"`
	Ninverified bool   `json:"ninverified"`
	CreatedAt   string `json:"created_at"`
}

func AddLister(lister Lister) error {
	exists, err := db.DB.Query("SELECT 1 FROM customers WHERE email = $1", lister.Email)
	if err != nil {
		return fmt.Errorf("DB query error: %v", err)
	}
	defer exists.Close()

	if exists.Next() {
		return errors.New("email already registered")
	}

	hashedPassword, err := utils.HashPassword(lister.Password)
	if err != nil {
		return err
	}

	query := `INSERT INTO listers (email, password_hash, full_name) VALUES ($1, $2, $3)`
	_, err = db.DB.Exec(query, lister.Email, hashedPassword, lister.FullName)
	return err

}

func (l *Lister) ValidateLister() error {
	query := `SELECT password_hash, id FROM listers WHERE email = $1`
	row := db.DB.QueryRow(query, l.Email)

	var retrievedHash string
	err := row.Scan(&retrievedHash, &l.Id)
	if err != nil {
		fmt.Println("DB Error:", err)
		return errors.New("invalid email or password")
	}

	passwordIsValid := utils.CheckPasswordHash(l.Password, retrievedHash)
	if !passwordIsValid {
		return errors.New("invalid email or password")
	}

	return nil
}
