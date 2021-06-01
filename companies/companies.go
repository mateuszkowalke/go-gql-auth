package companies

import (
	"database/sql"
	"log"

	"example.com/go-graphql-auth/database"
	"example.com/go-graphql-auth/graph/model"
)

type Company struct {
	*model.NewCompany
}

func (c Company) Save() int64 {
	stmt, err := database.Db.Prepare("INSERT INTO Companies(Name,Email,Country) VALUES(?,?,?)")
	if err != nil {
		log.Println(err)
	}
	res, err := stmt.Exec(c.Name, c.Email, c.Country)
	if err != nil {
		log.Println(err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Println(err)
	}
	log.Printf("Added company with id: %v\n", id)
	return id
}

func GetAllCompanies() ([]*model.Company, error) {
	stmt, err := database.Db.Prepare("SELECT * FROM Companies")
	if err != nil {
		log.Println(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		log.Println()
	}
	defer rows.Close()
	var companies []*model.Company
	for rows.Next() {
		var company model.Company
		err := rows.Scan(&company.ID, &company.Name, &company.Email, &company.Country)
		if err != nil {
			log.Println(err)
		}
		companies = append(companies, &company)
	}
	err = rows.Err()
	return companies, err
}

func GetCompanyByName(name string) (*model.Company, error) {
	stmt, err := database.Db.Prepare("SELECT * FROM Companies WHERE Name=?")
	if err != nil {
		log.Println(err)
	}
	row := stmt.QueryRow(name)
	var company model.Company
	err = row.Scan(&company)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Println(err)
		}
		return nil, err
	}
	return &company, nil
}
