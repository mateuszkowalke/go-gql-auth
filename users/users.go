package users

import (
	"database/sql"
	"log"

	"example.com/go-graphql-auth/database"
	"example.com/go-graphql-auth/graph/model"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	*model.NewUser
}

func (u User) Create() {
	stmt, err := database.Db.Prepare("INSERT INTO Users(Email,FirstName,LastName,Password,Role) VALUES(?,?,?,?,?)")
	if err != nil {
		log.Println(err)
	}
	hashedPassword, err := HashPassword(u.Password)
	if err != nil {
		log.Println(err)
	}
	_, err = stmt.Exec(u.Email, u.FirstName, u.LastName, hashedPassword, model.RoleUser)
	if err != nil {
		log.Println(err)
	}
	log.Println("Created new user")
}

func GetAllUsers() ([]*model.User, error) {
	stmt, err := database.Db.Prepare("SELECT ID,Email,FirstName,LastName,Role,CompanyID FROM Users")
	if err != nil {
		log.Println(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		log.Println()
	}
	defer rows.Close()
	var users []*model.User
	for rows.Next() {
		var user model.User
		err = rows.Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.Role, &user.Company)
		if err != nil {
			log.Println(err)
		}
		users = append(users, &user)
	}
	err = rows.Err()
	return users, err
}

func GetUserByEmail(email string) (*model.User, error) {
	stmt, err := database.Db.Prepare("SELECT ID,Email,FirstName,LastName,Role,CompanyID FROM Users WHERE Email = ?")
	if err != nil {
		log.Fatal(err)
	}
	row := stmt.QueryRow(email)
	var user model.User
	err = row.Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.Role, &user.Company)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Print(err)
		}
		return nil, err
	}
	return &user, nil
}

func GetUserIdByEmail(email string) (int, error) {
	stmt, err := database.Db.Prepare("select ID from Users WHERE Email = ?")
	if err != nil {
		log.Fatal(err)
	}
	row := stmt.QueryRow(email)
	var Id int
	err = row.Scan(&Id)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Print(err)
		}
		return 0, err
	}
	return Id, nil
}

func HashPassword(p string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(p), 14)
	return string(bytes), err
}

// returns true if password matches hash
func CheckPasswordHash(p, h string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(h), []byte(p))
	return err == nil
}
