package paginate

import (
	"fmt"
	"strings"
	"reflect"

	"gorm.io/gorm"

	"golang-crud/internal/database"
)

type SearchRequest struct {
	AdvancedSearch [][]Search `json:"advanced_search,omitempty"`
	Search   *Search `json:"search"`
	Limit    int     `json:"limit"`
	Page     int     `json:"page"`
}

type Search struct {
	Field    string `json:"field"`
	Operator string `json:"operator"`
	Value    any    `json:"value"`
}

type Wrapper struct {
	Data           any   	         `json:"data"`
	HasMore        bool  	         `json:"has_more"`
	Current        int   	         `json:"current"`
	Next           int   	         `json:"next"`
	Totals         int64 	         `json:"totals"`
	AllowedFilters map[string]string `json:"allowed_filters"`
}

var allowedOperators = map[string]bool{
	"=":       true,
	"!=":      true,
	">":       true,
	"<":       true,
	"LIKE":    true,
	"IN":      true,
	"BETWEEN": true,
}

var filterable map[string]string 

func (filter *Search) GetCondition () (string, []any, error) {
	operator := strings.ToUpper(filter.Operator)
	field    := filter.Field

	switch operator {
		case "IN", "NOT IN":
			values, ok := filter.Value.([]any)
			if ! ok || len(values) == 0 {
				return "", nil, fmt.Errorf("%s operator requires a non-empty array", operator)
			}

			condition := fmt.Sprintf("%s %s ?", field, operator)
			return condition, []any{values}, nil

		case "BETWEEN":
			values, ok := filter.Value.([]any)
			if ! ok || len(values) != 2 {
				return "", nil, fmt.Errorf("BETWEEN operator requires exactly two values")
			}

			condition := fmt.Sprintf("%s BETWEEN ? AND ?", field)
			return condition, values, nil

		default:
			condition := fmt.Sprintf("%s %s ?", field, operator)
			return condition, []any{filter.Value}, nil
		}
}

func ApplyAdvancedFilters(
	query *gorm.DB,
	groups [][]Search,
	allowedFilters map[string]string,
) (*gorm.DB, error) {

	if len(groups) == 0 {
		return query, nil
	}

	var orQueries []*gorm.DB

	for _, group := range groups {
		groupQuery := database.DB

		for _, filter := range group {
			var err error

			groupQuery, err = ApplyFilter(groupQuery, filter, allowedFilters)
			
			if err != nil {
				return nil, err
			}
		}

		orQueries = append(orQueries, groupQuery)
	}

	if len(orQueries) > 0 {
		query = query.Where(orQueries[0])

		for i := 1; i < len(orQueries); i++ {
			query = query.Or(orQueries[i])
		}
	}

	return query, nil
}

func ApplyFilter(query *gorm.DB, filter Search, allowedFilters map[string]string) (*gorm.DB, error) {
	if filter.Field == "" {
		return query, nil
	}

	column, ok := allowedFilters[filter.Field]
	if ! ok {
		return nil, fmt.Errorf("field '%s' is not allowed for filtering", filter.Field)
	}

	filter.Field = column

	operator := strings.ToUpper(filter.Operator)
	if !allowedOperators[operator] {
		return nil, fmt.Errorf("invalid operator: %s", filter.Operator)
	}

	filter.Operator = operator

	condition, args, err := filter.GetCondition()
	if err != nil {
		return nil, err
	}

	return query.Where(condition, args...), nil
}

func GetFilterableFields(model any) (map[string]string, error) {
    valueOf := reflect.ValueOf(model)

	// Remove ponteiros
	for valueOf.Kind() == reflect.Ptr {
		valueOf = valueOf.Elem()
	}

	// Se for slice, pega o tipo do elemento
	if valueOf.Kind() == reflect.Slice {
		valueOf = reflect.New(valueOf.Type().Elem()).Elem()
	}

	filterable, ok := valueOf.Interface().(interface {
		FilterableFields() map[string]string
	})

	if ! ok {
		return nil, fmt.Errorf("model does not implement FilterableFields")
	}

	return filterable.FilterableFields(), nil
}

var filteableFields map[string]string

func (search *SearchRequest) Apply (model any) (*Wrapper, error) {
	// Pega a query para ser atualizada
	query := database.DB.Model(model)

	// Garante que o Page, tenha valor que pode ser usado
	if (search.Page <= 0) {
		search.Page = 1 
	}

	// Garante que o Limit, tenha valor que pode ser usado
	if (search.Limit <= 0) {
		search.Limit = 10
	}

	// Aplica o offset e limit a query
	offset := (search.Page - 1) * search.Limit
	query  = query.Offset(offset).Limit(search.Limit)

	// Campos que podem ser filtrados
	allowedFilters, _ := GetFilterableFields(model)

	if search.Search != nil {
		var err error

		query, err = ApplyFilter(query, *search.Search, allowedFilters)
		
		if err != nil {
			return nil, err
		}
	}

	if len(search.AdvancedSearch) > 0 {
		var err error

		query, err = ApplyAdvancedFilters(query, search.AdvancedSearch, allowedFilters)
		
		if err != nil {
			return nil, err
		}
	}

	// Obtem a quantidade de registro totais com condicoes. Retorna o erro
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	stmt := query.Session(&gorm.Session{DryRun: true}).Find(model).Statement

	fmt.Println("SQL:", stmt.SQL.String())

	// Faz o fetch. ou retorna erro
	if err := query.Find(model).Error; err != nil {
		return nil, err
	}

	// Verifica se teria outra pagina
	hasMore := int64(search.Page * search.Limit) < total

	// Define a proxima pagina
	next := search.Page
	if hasMore {
		next = search.Page + 1
	}
 
	return &Wrapper{
		Data:           model,
		HasMore:        hasMore,
		Current:        search.Page,
		Next:           next,
		Totals:         total,
		AllowedFilters: allowedFilters,
	}, nil
} 