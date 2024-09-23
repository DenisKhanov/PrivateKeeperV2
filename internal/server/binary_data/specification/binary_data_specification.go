package specification

import (
	"strings"

	pb "github.com/DenisKhanov/PrivateKeeperV2/internal/proto/binary_data"
	"github.com/DenisKhanov/PrivateKeeperV2/internal/server/model"
)

type BinaryDataSpecification struct {
	Name     string
	MetaData string
}

func NewTextDataSpecification(in *pb.GetBinaryDataRequest) (BinaryDataSpecification, error) {
	spec := BinaryDataSpecification{
		Name:     in.Name,
		MetaData: in.Metadata,
	}

	return spec, nil
}

type Predicate func(spec BinaryDataSpecification, resp model.BinaryData) bool

func (t *BinaryDataSpecification) MakeFilterPredicates() []Predicate {
	var predicates []Predicate
	if t.MetaData != "" {
		predicates = append(predicates, makeMetaDataPredicate())
	}
	if t.Name != "" {
		predicates = append(predicates, makeNamePredicate())
	}
	return predicates
}

func makeMetaDataPredicate() Predicate {
	return func(spec BinaryDataSpecification, resp model.BinaryData) bool {
		return strings.Contains(strings.ToLower(resp.MetaData), strings.ToLower(spec.MetaData))
	}
}

func makeNamePredicate() Predicate {
	return func(spec BinaryDataSpecification, resp model.BinaryData) bool {
		return strings.Contains(strings.ToLower(resp.Name), strings.ToLower(spec.Name))
	}
}
