package domain

//go:generate mockgen -source=${GOFILE} -destination=${ROOT_DIR}/testing/mock/mock_${GOPACKAGE}_${GOFILE} -package=mock

// Explorer is an interface of file explorer.
type Explorer interface {
	OpenFile(FileName) (*File, error)
}
