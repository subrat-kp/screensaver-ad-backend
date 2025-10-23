package models

func Models() []interface{} {
	return []interface{}{
		&Asset{},
		&Template{},
		&Task{},
	}
}
