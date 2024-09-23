package specification

import (
	"errors"
	"strings"
	"time"

	pb "github.com/DenisKhanov/PrivateKeeperV2/internal/proto/credit_card"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/model"
)

const expiresAtLayout = "02-01-2006"

type CreditCardSpecification struct {
	OwnerID       string
	Number        string
	OwnerName     string
	CVV           string
	PinCode       string
	MetaData      string
	ExpiresAfter  time.Time
	ExpiresBefore time.Time
}

func NewCreditCardSpecification(in *pb.GetCreditCardRequest) (CreditCardSpecification, error) {
	spec := CreditCardSpecification{
		Number:    in.Number,
		OwnerName: in.Owner,
		CVV:       in.CvvCode,
		PinCode:   in.PinCode,
		MetaData:  in.Metadata,
	}

	if in.ExpiresAfter != "" {
		expiresAfter, err := time.Parse(expiresAtLayout, in.ExpiresAfter)
		if err != nil {
			return CreditCardSpecification{}, errors.New("expires after must be in format 'DD-MM-YYYY'")
		}
		spec.ExpiresAfter = expiresAfter
	}

	if in.ExpiresBefore != "" {
		expiresBefore, err := time.Parse(expiresAtLayout, in.ExpiresBefore)
		if err != nil {
			return CreditCardSpecification{}, errors.New("expires before must be in format 'DD-MM-YYYY'")
		}
		spec.ExpiresBefore = expiresBefore
	}

	return spec, nil
}

type Predicate func(spec CreditCardSpecification, resp model.CreditCard) bool

func (t *CreditCardSpecification) MakeFilterPredicates() []Predicate {
	var predicates []Predicate
	if t.Number != "" {
		predicates = append(predicates, makeNumberPredicate())
	}
	if t.OwnerName != "" {
		predicates = append(predicates, makeOwnerNamePredicate())
	}
	if t.CVV != "" {
		predicates = append(predicates, makeCVVPredicate())
	}
	if t.PinCode != "" {
		predicates = append(predicates, makePINPredicate())
	}
	if t.MetaData != "" {
		predicates = append(predicates, makeMetaDataPredicate())
	}
	if !t.ExpiresAfter.IsZero() {
		predicates = append(predicates, makeExpiresAfterPredicate())
	}
	if !t.ExpiresBefore.IsZero() {
		predicates = append(predicates, makeExpiresBeforePredicate())
	}
	return predicates
}

func makeExpiresAfterPredicate() Predicate {
	return func(spec CreditCardSpecification, resp model.CreditCard) bool {
		expiresAt, err := time.Parse(expiresAtLayout, resp.ExpiresAt)
		if err != nil {
			return false
		}

		return expiresAt.After(spec.ExpiresAfter)
	}
}

func makeExpiresBeforePredicate() Predicate {
	return func(spec CreditCardSpecification, resp model.CreditCard) bool {
		expiresAt, err := time.Parse(expiresAtLayout, resp.ExpiresAt)
		if err != nil {
			return false
		}

		return expiresAt.Before(spec.ExpiresBefore)
	}
}

func makeNumberPredicate() Predicate {
	return func(spec CreditCardSpecification, resp model.CreditCard) bool {
		return strings.Contains(strings.ToLower(resp.Number), strings.ToLower(spec.Number))
	}
}

func makeOwnerNamePredicate() Predicate {
	return func(spec CreditCardSpecification, resp model.CreditCard) bool {
		return strings.Contains(strings.ToLower(resp.OwnerName), strings.ToLower(spec.OwnerName))
	}
}

func makeCVVPredicate() Predicate {
	return func(spec CreditCardSpecification, resp model.CreditCard) bool {
		return strings.Contains(strings.ToLower(resp.CVV), strings.ToLower(spec.CVV))
	}
}

func makePINPredicate() Predicate {
	return func(spec CreditCardSpecification, resp model.CreditCard) bool {
		return strings.Contains(strings.ToLower(resp.PinCode), strings.ToLower(spec.PinCode))
	}
}

func makeMetaDataPredicate() Predicate {
	return func(spec CreditCardSpecification, resp model.CreditCard) bool {
		return strings.Contains(strings.ToLower(resp.MetaData), strings.ToLower(spec.MetaData))
	}
}
