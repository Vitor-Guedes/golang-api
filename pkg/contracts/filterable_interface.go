package contracts

type Filterable interface {
    FilterableFields() map[string]string
}