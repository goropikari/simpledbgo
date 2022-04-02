package domain

type TransactionNumber int32

func (txnum TransactionNumber) LSN() LSN {
	return LSN(txnum)
}
