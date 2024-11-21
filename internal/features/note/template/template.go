package note

import (
	"github.com/jim-ww/nms-go/internal/features/note/dtos"
	userDTO "github.com/jim-ww/nms-go/internal/features/user/dtos"
)

type SearchIn string
type Order string
type SortByField string

const (
	Title       SearchIn    = "title"
	Content     SearchIn    = "content"
	Asc         Order       = "ascending"
	Desc        Order       = "descending"
	ByTitle     SortByField = "title"
	ByCreatedAt SortByField = "created at"
	ByUpdatedAt SortByField = "updated at"
)

type PaginationData struct {
	SearchIn
	Order
	SortByField
	NotesPerPage int
	Page         int
}

type DashboardData struct {
	Notes        *[]dtos.NoteSummaryDTO
	SelectedNote *dtos.NoteDetailDTO
	userDTO.UserProfileDTO
	PanelClosed         bool
	SearchOptionsClosed bool
	UserProfileClosed   bool
}

func New(notes *[]dtos.NoteSummaryDTO, selectedNote *dtos.NoteDetailDTO, userDTO userDTO.UserProfileDTO) *DashboardData {
	return &DashboardData{
		PanelClosed:         false,
		SearchOptionsClosed: true,
		UserProfileClosed:   true,
		Notes:               notes,
		SelectedNote:        selectedNote,
		UserProfileDTO:      userDTO,
	}
}
