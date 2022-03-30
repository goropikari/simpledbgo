package fake

import "github.com/goropikari/simpledb_go/backend/domain"

func Page(bb domain.ByteBuffer) *domain.Page {
	return domain.NewPage(bb)
}
