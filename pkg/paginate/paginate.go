package paginate

import (
	"golang-crud/internal/database"
)

type SearchRequest struct {
	Search Search `json:"search"`
	Limit  int    `json:"limit"`
	Page   int    `json:"page"`
}

type Search struct {
	Field    string `json:"field"`
	Operator string `json:"operator"`
	Value    any    `json:"value"`
}

type Wrapper struct {
	Data    any   `json:"data"`
	HasMore bool  `json:"has_more"`
	Current int   `json:"current"`
	Next    int   `json:"next"`
	Totals  int64 `json:"totals"`
}

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

	// Obtem a quantidade de registro totais. Retorna o erro
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// Aplica o offset e limit a query
	offset := (search.Page - 1) * search.Limit
	query  = query.Offset(offset).Limit(search.Limit)
	
	// Aplica uma condicao, caso tenha
	if search.Search.Field != "" {
		condition := search.Search.Field + " " + search.Search.Operator + " ?"
		query     = query.Where(condition, search.Search.Value)
	}

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
		Data:    model,
		HasMore: hasMore,
		Current: search.Page,
		Next:    next,
		Totals:  total,
	}, nil
} 