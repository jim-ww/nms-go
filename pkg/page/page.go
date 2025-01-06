package page

type Order int

const (
	ASC Order = iota
	DESC
)

var orderString = map[Order]string{
	ASC:  "ASC",
	DESC: "DESC",
}

func (o Order) String() string {
	return orderString[o]
}

type PaginationParams struct {
	PageNumber int
	PageSize   int
	Order      Order
	OrderBy    string
}

func NewPaginationParams(pageNumber, pageSize int, order Order, orderBy string) *PaginationParams {
	return &PaginationParams{
		PageNumber: pageNumber,
		PageSize:   pageSize,
		Order:      order,
		OrderBy:    orderBy,
	}
}

type Page struct {
	TotalPages    int
	TotalElements int
	PageNumber    int
}

func NewPage(totalPages, totalElements, pageNumber int) *Page {
	return &Page{
		TotalPages:    totalPages,
		TotalElements: totalElements,
		PageNumber:    pageNumber,
	}
}
