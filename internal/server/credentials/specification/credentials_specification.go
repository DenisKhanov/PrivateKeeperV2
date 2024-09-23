package specification

import (
	"strings"

	pb "github.com/DenisKhanov/PrivateKeeperV2/internal/proto/credentials"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/model"
)

type CredentialsSpecification struct {
	Login    string
	Password string
	Metadata string
}

func NewCredentialsSpecification(in *pb.GetCredentialsRequest) (CredentialsSpecification, error) {
	spec := CredentialsSpecification{
		Login:    in.Login,
		Password: in.Password,
		Metadata: in.Metadata,
	}

	return spec, nil
}

type Predicate func(spec CredentialsSpecification, resp model.Credentials) bool

func (t *CredentialsSpecification) MakeFilterPredicates() []Predicate {
	var predicates []Predicate
	if t.Login != "" {
		predicates = append(predicates, makeLoginPredicate())
	}
	if t.Password != "" {
		predicates = append(predicates, makePasswordPredicate())
	}
	if t.Metadata != "" {
		predicates = append(predicates, makeMetaDataPredicate())
	}
	return predicates
}

func makeMetaDataPredicate() Predicate {
	return func(spec CredentialsSpecification, resp model.Credentials) bool {
		return strings.Contains(strings.ToLower(resp.MetaData), strings.ToLower(spec.Metadata))
	}
}

func makePasswordPredicate() Predicate {
	return func(spec CredentialsSpecification, resp model.Credentials) bool {
		return strings.Contains(strings.ToLower(resp.Password), strings.ToLower(spec.Password))
	}
}

func makeLoginPredicate() Predicate {
	return func(spec CredentialsSpecification, resp model.Credentials) bool {
		return strings.Contains(strings.ToLower(resp.Login), strings.ToLower(spec.Login))
	}
}
