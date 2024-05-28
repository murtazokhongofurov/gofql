package gofql

import (
	"context"
	"fmt"
	"reflect"
)

func (o *ORM) Insert(model interface{}) error {
	v := reflect.ValueOf(model).Elem()
	t := v.Type()

	columns := ""
	values := ""
	args := []interface{}{}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		columns += field.Name + ","
		values += "?,"
		args = append(args, v.Field(i).Interface())
	}
	columns = columns[:len(columns)-1] // Remove the trailing comma
	values = values[:len(values)-1]    // Remove the trailing comma

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", t.Name(), columns, values)

	_, err := o.db.Exec(query, args...)
	return err
}

func (o *ORM) FindByID(model interface{}, id int) error {
	v := reflect.ValueOf(model).Elem()
	t := v.Type()

	query := fmt.Sprintf("SELECT * FROM %s WHERE id = ?", t.Name())

	row, err := o.db.QueryContext(context.TODO(), query, id)
	columns, err := row.Columns()
	values := make([]interface{}, len(columns))
	scanArgs := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	err = row.Scan(scanArgs...)
	if err != nil {
		return err
	}

	for i, col := range columns {
		field := v.FieldByName(col)
		if field.IsValid() && field.CanSet() {
			field.Set(reflect.ValueOf(values[i]).Convert(field.Type()))
		}
	}

	return nil
}

func (o *ORM) Update(model interface{}, id int) error {
	v := reflect.ValueOf(model).Elem()
	t := v.Type()

	setClauses := ""
	args := []interface{}{}
	var idValue interface{}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i).Interface()

		if field.Name == "ID" {
			idValue = value
			continue
		}

		setClauses += fmt.Sprintf("%s = ?,", field.Name)
		args = append(args, value)
	}

	setClauses = setClauses[:len(setClauses)-1] // Remove the trailing comma
	args = append(args, idValue)

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = ?", t.Name(), setClauses)

	_, err := o.db.Exec(query, args...)
	return err
}

func (o *ORM) Delete(model interface{}, id int) error {
	t := reflect.TypeOf(model).Elem()

	query := fmt.Sprintf("DELETE FROM %s WHERE id = ?", t.Name())

	_, err := o.db.Exec(query, id)
	return err
}
