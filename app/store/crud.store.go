package store

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"de.amplifonx/app/model"
	"github.com/davecgh/go-spew/spew"
	"gorm.io/gorm"
)

type CrudStore[T model.Model] struct {
	TableName string
	DB        *gorm.DB
}

func (store *CrudStore[T]) AddOne(data *T) (*T, error) {
	spew.Dump(data)
	result := store.DB.Debug().Table(store.TableName).Create(data)
	return data, result.Error
}

func (store *CrudStore[T]) FindOne(data *T, preload *[]string) (*T, error) {
	spew.Dump(data)
	var response T
	result := store.DB.Debug().Table(store.TableName).Scopes(Preload(preload)).First(&response, data)
	return &response, result.Error
}

func (store *CrudStore[T]) FindMany(r *http.Request, data *T, preload *[]string) (*[]T, error) {
	spew.Dump(data)
	var elements []T
	result := store.DB.Debug().Table(store.TableName).Scopes(Paginate(r), Preload(preload)).Find(&elements, data)
	return &elements, result.Error
}

func (store *CrudStore[T]) Filter(r *http.Request, preload *[]string, filters *[]string) (*[]T, error) {
	fmt.Println("filtering for " + strings.Join(*filters, " and "))
	var elements []T
	result := store.DB.Debug().Table(store.TableName).Scopes(Paginate(r), FilterStrings(r, filters), Preload(preload)).Find(&elements)
	return &elements, result.Error
}

func (store *CrudStore[T]) Delete(data *T) (*T, error) {
	spew.Dump(data)
	if data != nil {
		result := store.DB.Debug().Table(store.TableName).Delete(data)
		return data, result.Error
	}
	return data, fmt.Errorf("nothing to delete")
}

func (store *CrudStore[T]) Update(data *T) (*T, error) {
	spew.Dump(data)
	if data != nil {
		result := store.DB.Debug().Table(store.TableName).Save(data)
		return data, result.Error
	}
	return data, fmt.Errorf("nothing to update")
}

func Paginate(r *http.Request) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		q := r.URL.Query()
		page, err := strconv.Atoi(q.Get("page"))
		if err != nil {
			return db
		}
		if page <= 0 {
			page = 1
		}

		pageSize, err := strconv.Atoi(q.Get("page_size"))
		if err != nil {
			return db
		}
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func Preload(preload *[]string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if preload != nil {
			for _, join := range *preload {
				db = db.Preload(join)
			}
		}
		return db
	}
}

func FilterStrings(r *http.Request, filters *[]string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		q := r.URL.Query()
		fmt.Println(filters)
		for _, filter := range *filters {
			value := q.Get(filter)

			if value != "" {
				_, err := strconv.Atoi(value)
				// add other cases for different filter types
				if err == nil {
					fmt.Println("only strings are accepted")
				} else {
					db.Where("LOWER(" + filter + ") like LOWER('%" + value + "%')")
				}
			}
		}
		return db
	}
}
