package domain_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/goropikari/simpledb_go/backend/domain"
	"github.com/goropikari/simpledb_go/testing/mock"
)

func TestPage_NewPage(t *testing.T) {
	t.Run("test page constructor", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		bb := mock.NewMockByteBuffer(ctrl)

		domain.NewPage(bb)
	})
}
