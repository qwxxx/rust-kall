package store

import (
	"fmt"
	"reflect"

	"github.com/doug-martin/goqu/v9"
)

type CommonGateway interface {
	GetElements(tablename string, objects any, offset int, limit int, filters ...goqu.Expression) error
	GetElementsCount(tablename string, filters ...goqu.Expression) (int, error)
	GetElement(tablename string, object any) error
	EditElement(tablename string, object any) error
	AddElement(tablename string, object any) (any, error)
	DeleteElement(tablename string, object any) error
}

func _getKeyField(object any) (string, any) {

	t := reflect.TypeOf(object)
	val := reflect.ValueOf(object)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		val = val.Elem()
	}
	numfield := t.NumField()

	for i := 0; i < numfield; i++ {
		field := t.Field(i)
		tag := field.Tag.Get("qwery")
		if tag == "key" {
			tag1 := field.Tag.Get("db")

			return tag1, val.Field(i).Interface()
		}
	}
	return "", nil
}
func (db *Store) GetElementsCount(tablename string, filters ...goqu.Expression) (int, error) {
	i, err := db.Select().From(tablename).Where(filters...).Count()

	return int(i), err
}
func (db *Store) GetElements(tablename string, objects any, offset int, limit int, filters ...goqu.Expression) error {
	err := db.Select().From(tablename).Where(filters...).Offset(uint(offset)).Limit(uint(limit)).ScanStructs(objects)
	return err
}
func (db *Store) GetElement(tablename string, object any) error {
	key, value := _getKeyField(object)
	_, err := db.Select().From(tablename).Where(goqu.Ex{key: value}).ScanStruct(object)
	if err != nil {
		return fmt.Errorf("Read %s err: %v", tablename, err)
	}
	return err
}
func (db *Store) EditElement(tablename string, object any) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("Open edit %s transaction error: %v", tablename, err)
	}
	defer tx.Rollback()
	fieldname, fieldvalue := _getKeyField(object)
	update := db.Update(tablename).Set(object).Where(goqu.Ex{fieldname: fieldvalue}).Executor()

	if _, err := update.Exec(); err != nil {
		return fmt.Errorf("Update %s error: %v", tablename, err)
	}
	err = tx.Commit()
	return err
}
func (db *Store) AddElement(tablename string, object any) (any, error) {
	var id any
	tx, err := db.Begin()
	if err != nil {
		return id, fmt.Errorf("create %s error: %v", tablename, err)
	}
	defer tx.Rollback()

	fieldname, _ := _getKeyField(object)
	_, err = db.Insert(tablename).Rows(object).Returning(fieldname).Executor().ScanVal(&id)

	if err != nil {
		return id, fmt.Errorf("create %s error: %v", tablename, err)
	}

	err = tx.Commit()
	return id, err
}

func (db *Store) DeleteElement(tablename string, object any) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("Open delete %s transaction error: %v", tablename, err)
	}
	defer tx.Rollback()

	key, value := _getKeyField(object)
	del := db.Delete(tablename).Where(goqu.Ex{key: value}).Executor()

	if _, err = del.Exec(); err != nil {
		return fmt.Errorf("Delete %s error: %v", tablename, err)
	}

	err = tx.Commit()
	return err

}
