package main

const (
	// Table Names
	userTableName       = "user"
	todoStatusTableName = "todo_status"
	todoTableName       = "todo"
)

type TableNameDataStruct struct {
	tableName  string
	dataStruct interface{}
}

//
// Basic Data Pattern
//

// User Table
type UserData struct {
	UserID   string `name:"user_id"       prefix:"userid_"`
	UserName string `name:"user_name"     prefix:"username_"`
	Email    string `name:"email"         prefix:"email@email"`
	Password string `name:"password_hash"                      fixed:"password"`
}
