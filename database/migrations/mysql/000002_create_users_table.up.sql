CREATE TABLE IF NOT EXISTS Users(
    ID INT NOT NULL UNIQUE AUTO_INCREMENT,
    Email VARCHAR (127) NOT NULL UNIQUE,
    FirstName VARCHAR (127) NOT NULL,
    LastName VARCHAR (127) NOT NULL,
    Role ENUM('user', 'superadmin'),
    Password VARCHAR (127) NOT NULL,
    CompanyID INT,
    FOREIGN KEY (CompanyID) REFERENCES Companies(ID),
    PRIMARY KEY (ID)
)