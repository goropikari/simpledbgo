package log_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/goropikari/simpledb_go/backend/log"
	"github.com/goropikari/simpledb_go/infra"
	"github.com/goropikari/simpledb_go/testing/mock"
	"github.com/stretchr/testify/require"
)

func TestManager_Error(t *testing.T) {
	t.Run("error PreparePage", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileMgr := mock.NewMockFileManager(ctrl)
		fileMgr.EXPECT().PreparePage().Return(nil, errors.New("errors"))

		_, err := log.NewManager(fileMgr, infra.Config{})
		require.Error(t, err)
	})

}
