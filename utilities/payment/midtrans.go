package payment

import (
	"errors"
	"fmt"
	"olshop/config"
	"olshop/features/transactions"
	"time"

	mdt "github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
)

type Midtrans interface {
	NewTransactionPayment(data transactions.Transaction) (*transactions.Payment, error)
}

func NewMidtrans(config config.Midtrans) Midtrans {
	var client coreapi.Client
	client.New(config.ApiKey, config.Env)

	return &midtrans{
		config: config,
		client: client,
	}
}

type midtrans struct {
	config config.Midtrans
	client coreapi.Client
}

func (pay *midtrans) NewTransactionPayment(data transactions.Transaction) (*transactions.Payment, error) {
	req := new(coreapi.ChargeReq)
	req.TransactionDetails = mdt.TransactionDetails{
		OrderID:  fmt.Sprintf("%d", data.Invoice),
		GrossAmt: int64(data.Total),
	}

	req.CustomerDetails = &mdt.CustomerDetails{
		FName: data.User.Name,
		Email: data.User.Email,
	}

	switch data.Payment.Bank {
	case "bca":
		req.PaymentType = coreapi.PaymentTypeBankTransfer
		req.BankTransfer = &coreapi.BankTransferDetails{
			Bank: mdt.BankBca,
		}
	case "bni":
		req.PaymentType = coreapi.PaymentTypeBankTransfer
		req.BankTransfer = &coreapi.BankTransferDetails{
			Bank: mdt.BankBni,
		}
	case "bri":
		req.PaymentType = coreapi.PaymentTypeBankTransfer
		req.BankTransfer = &coreapi.BankTransferDetails{
			Bank: mdt.BankBri,
		}
	default:
		return nil, errors.New("unsupported payment")
	}

	res, _ := pay.client.ChargeTransaction(req)
	if res.StatusCode != "201" {
		return nil, errors.New(res.StatusMessage)
	}

	if res.BillKey != "" {
		data.Payment.BillKey = res.BillKey
	}

	if res.BillerCode != "" {
		data.Payment.BillCode = res.BillerCode
	}

	if len(res.VaNumbers) == 1 {
		data.Payment.VirtualNumber = res.VaNumbers[0].VANumber
	}

	if res.PermataVaNumber != "" {
		data.Payment.VirtualNumber = res.PermataVaNumber
	}

	if res.PaymentType != "" {
		data.Payment.Method = res.PaymentType
	}

	if res.TransactionStatus != "" {
		data.Payment.Status = res.TransactionStatus
	}

	if expiredAt, err := time.Parse("2006-01-02 15:04:05", res.ExpiryTime); err != nil {
		return nil, err
	} else {
		data.Payment.ExpiredAt = expiredAt
	}

	data.Payment.TransactionTotal = data.Total

	return &data.Payment, nil
}
