package src

import (
	"github.com/afistapratama12/micli/src/repo"
	"github.com/afistapratama12/micli/src/service"
	"github.com/afistapratama12/micli/src/view"
)

func NewCrypto() view.CryptoView {
	cryptoRepo := repo.NewCryptoRepo()
	cryptoService := service.NewCryptoService(cryptoRepo)
	cryptoView := view.NewCryptoView(cryptoService)

	return *cryptoView
}
