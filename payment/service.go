package payment

import (
	"crowdfunding-golang/user"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	midtrans "github.com/veritrans/go-midtrans"
)

type service struct {
}

type Service interface {
	GetPaymentURL(transaction Transaction, user user.User) (string, error)
}

func NewService() *service {
	return &service{}
}

func (s *service) GetPaymentURL(transaction Transaction, user user.User) (string, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return "", err
	}

	midtrans_client_key := os.Getenv("MIDTRANS_CLIENT_KEY")
	midtrans_server_key := os.Getenv("MIDTRANS_SERVER_KEY")

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
