package payment

import (
	"crowdfunding-golang/campaign"
	"crowdfunding-golang/transaction"
	"crowdfunding-golang/user"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	midtrans "github.com/veritrans/go-midtrans"
)

type service struct {
	transactionRepository transaction.Repository
	campaignRepository    campaign.Repository
}

type Service interface {
	GetPaymentURL(transaction Transaction, user user.User) (string, error)
	ProcessPayment(notification TransactionNotificationInput) error
}

func NewService(transactionRepository transaction.Repository, campaignRepository campaign.Repository) *service {
	return &service{transactionRepository, campaignRepository}
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

func (s *service) ProcessPayment(notification TransactionNotificationInput) error {
	transaction_id, _ := strconv.Atoi(notification.OrderID)

	transaction, err := s.transactionRepository.GetByID(transaction_id)
	if err != nil {
		return err
	}

	if notification.PaymentType == "credit_card" && notification.TransactionStatus == "capture" && notification.FraudStatus == "accept" {
		transaction.Status = "paid"
	} else if notification.TransactionStatus == "settlement" {
		transaction.Status = "paid"
	} else if notification.TransactionStatus == "deny" || notification.TransactionStatus == "expire" || notification.TransactionStatus == "cancel" {
		transaction.Status = "cancelled"
	}

	updatedTransaction, err := s.transactionRepository.Update(transaction)
	if err != nil {
		return err
	}

	campaign, err := s.campaignRepository.FindByID(updatedTransaction.CampaignID)
	if err != nil {
		return err
	}

	if updatedTransaction.Status == "paid" {
		campaign.BackerCount = campaign.BackerCount + 1
		campaign.CurrentAmount = campaign.CurrentAmount + updatedTransaction.Amount

		_, err = s.campaignRepository.Update(campaign)
		if err != nil {
			return err
		}
	}

	return nil
}
