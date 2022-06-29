package storage

import (
	"reflect"

	"github.com/Masterminds/squirrel"
)

type Filters struct {
	Condition squirrel.And
}

func NewFilters() Filters {
	return Filters{
		Condition: squirrel.And{},
	}
}

func (f *Filters) AddEqual(field string, value any) {
	if value == nil || reflect.ValueOf(value).IsNil() {
		return
	}

	f.Condition = append(f.Condition, squirrel.Eq{field: value})
}

func (f *Filters) AddNil(field string) {
	f.Condition = append(f.Condition, squirrel.Eq{field: nil})
}

func (f *Filters) AddGreaterThan(field string, value any) {
	if value == nil || reflect.ValueOf(value).IsNil() {
		return
	}

	f.Condition = append(f.Condition, squirrel.Gt{field: value})
}

func (f *Filters) AddLessThan(field string, value any) {
	if value == nil || reflect.ValueOf(value).IsNil() {
		return
	}

	f.Condition = append(f.Condition, squirrel.Lt{field: value})
}

func (f *Filters) AddInRange(start_field string, end_field string, value any) {
	f.AddLessThan(start_field, value)
	f.AddGreaterThan(end_field, value)
}
