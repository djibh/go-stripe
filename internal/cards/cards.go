package cards

import (
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/paymentintent"
)

type Card struct {
	Secret   string
	Key      string
	Currency string
}

type Transaction struct {
	TransactionStatusId int
	Amount              int
	Currency            string
	LastFour            int
	BankReturnCode      string
}

// create an alias to the function for more clarity and that shouldn't change even if payment provider changes
func (c *Card) Charge(currency string, amount int) (*stripe.PaymentIntent, string, error) {
	return c.CreatePaymentIntent(currency, amount)
}

// this func name is a bit tied to the payment provider (Stripe) which could change over time
func (c *Card) CreatePaymentIntent(currency string, amount int) (*stripe.PaymentIntent, string, error) {
	stripe.Key = c.Secret

	// create a payment intent
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(amount)),
		Currency: stripe.String(currency),
	}

	// here we also have the possibility to add metadata: params.AddMetadata("key", "value")

	pi, err := paymentintent.New(params)
	if err != nil {
		msg := ""
		if stripeErr, ok := err.(*stripe.Error); ok {
			msg = cardErrorMessage(stripeErr.Code)
		}
		return nil, msg, err
	}
	return pi, "", nil
}

func cardErrorMessage(code stripe.ErrorCode) string {
	var msg = ""
	switch code {
	case stripe.ErrorCodeCardDeclined:
		msg = "Your card was declined."
	case stripe.ErrorCodeExpiredCard:
		msg = "Your card is expired."
	case stripe.ErrorCodeIncorrectCVC:
		msg = "Incorrect CVC code."
	case stripe.ErrorCodeIncorrectZip:
		msg = "Incorrect ZIP code."
	case stripe.ErrorCodeAmountTooLarge:
		msg = "The amount is too large to charge your card."
	case stripe.ErrorCodeAmountTooSmall:
		msg = "The amount is too small to charge your card."
	case stripe.ErrorCodeBalanceInsufficient:
		msg = "Insufficient balance."
	case stripe.ErrorCodePostalCodeInvalid:
		msg = "Incorrect Postal code."
	default:
		msg = "Your card was declined."
	}

	return msg
}
