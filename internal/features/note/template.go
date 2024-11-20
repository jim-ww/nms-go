package note

import "github.com/jim-ww/nms-go/internal/features/note/dtos"

type DashboardData struct {
	Notes        *[]dtos.NoteSummaryDTO
	SelectedNote *dtos.NoteDetailDTO
}
