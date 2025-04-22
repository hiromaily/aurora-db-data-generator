package main

const (
	// table name
	userTableName       = "user"
	todoStatusTableName = "todo_status"
	todoTableName       = "todo"
)

// data pattern struct. set only 1 property
type DataTag struct {
	Prefix   string
	Fixed    string
	Iterator Iterator
}

func (d *DataTag) isValid() bool {
	count := 0
	if d.Iterator != nil {
		count++
	}
	if d.Prefix != "" {
		count++
	}
	if d.Fixed != "" {
		count++
	}

	if count != 1 {
		return false
	}
	return true
}

type DataTagMap map[string]DataTag

type TableDataTagMap struct {
	tableName  string
	dataTagMap DataTagMap
}

//
// Basic Pattern
//

func getUesrDataTagMap() DataTagMap {
	return DataTagMap{
		"user_id":       {Prefix: "userid_"},
		"user_name":     {Prefix: "username_"},
		"email":         {Prefix: "email@email"},
		"password_hash": {Prefix: "", Fixed: "password"},
	}
}
