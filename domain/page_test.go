package domain_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/goropikari/simpledbgo/domain"
	"github.com/goropikari/simpledbgo/testing/mock"
)

func TestPage_NewPage(t *testing.T) {
	t.Run("test page constructor", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		bb := mock.NewMockByteBuffer(ctrl)

		domain.NewPage(bb)
	})
}
