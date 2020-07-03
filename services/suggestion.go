package services

import (
	"math"
	"strings"

	"github.com/zgs225/go-ecdict/dict"
)

// Suggestion 提示
type Suggestion struct {
	Label string   `json:"label"`
	Value string   `json:"value"`
	Tags  []string `json:"tags"`
}

// SuggestionService 单词输入提示
type SuggestionService interface {
	Suggest(string) ([]*Suggestion, error)
}

// EcDictSuggestionService 使用 ECDict 进行提示
type EcDictSuggestionService struct {
	Dict dict.Interface
	Max  int
}

// Suggest 获取输入提示
func (s *EcDictSuggestionService) Suggest(k string) ([]*Suggestion, error) {
	v, err := s.Dict.Like(k)
	if err != nil {
		return nil, err
	}

	max := int(math.Min(math.Max(10, float64(s.Max)), float64(len(v))))

	ss := make([]*Suggestion, max)

	for i := 0; i < max; i++ {
		r := v[i]
		ss[i] = &Suggestion{
			Label: r.Translation,
			Value: r.Word,
			Tags:  strings.Split(r.Tag, " "),
		}
	}

	return ss, nil
}
