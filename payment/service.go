package payment

import (
	"crowdfunding-golang/transaction"
	"crowdfunding-golang/user"
	"errors"
	"os"
	"strconv"

	"github.com/veritrans/go-midtrans"
)

type service struct {
}

type Service interface {
	GetToken(transaction transaction.Transaction, user user.User) (string, error)
}

func NewService() *service {
	return &service{}
}

func (s *service) GetToken(transaction transaction.Transaction, user user.User) (string, error) {
	midtrans_client_key := os.Getenv("MIDTRANS_CLIENT_KEY")
	midtrans_server_key := os.Getenv("MIDTRANS_SERVER_KEY")

	if midtrans_client_key == "" || midtrans_server_key == "" {
		return "", errors.New("midtrans client key or server key is empty")
	}

	midclient := midtrans.NewClient()
    midclient.ServerKey = midtrans_server_key
    midclient.ClientKey = midtrans_client_key
    midclient.APIEnvType = midtrans.Sandbox

    snapGateway := midtrans.SnapGateway{
        Client: midclient,
    }

	snapReq := &midtrans.SnapReq{
		CustomerDetail: &midtrans.CustDetail{
			Email: user.Email,
			FName: user.Name,
		},
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  strconv.Itoa(transaction.ID),
			GrossAmt: int64(transaction.Amount),
		},
	}

	snapTokenResp, err := snapGateway.GetToken(snapReq)
	if err != nil {
		return "", err
	}

	return snapTokenResp.RedirectURL, nil
}