package services

import (
	"io"
	"io/ioutil"
	"sort"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/polly"
)

var neuralSupportedRegions = sort.StringSlice([]string{"us-east-1", "us-west-2", "eu-west-1", "ap-southeast-2"})

// TTSService 语音合成服务
type TTSService interface {
	Speak(text, lang string) (string, error)
}

// TTSServiceMiddleware 语音合成服务中间件
type TTSServiceMiddleware func(TTSService) TTSService

// AWSPollyTTSService 亚马逊 Polly 语音合成
type AWSPollyTTSService struct {
	Region string

	sess   *session.Session
	client *polly.Polly
}

// NewAWSPollyTTSService 初始化 AWS Polly 语言合成服务
func NewAWSPollyTTSService(region string) (*AWSPollyTTSService, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String(region),
		},
	}))
	client := polly.New(sess)

	return &AWSPollyTTSService{
		Region: region,
		client: client,
		sess:   sess,
	}, nil
}

// Speak 将内容通过 AWS Polly 合成语言，下载对应的 mp3 文件，保存到本地后返回路径
func (s *AWSPollyTTSService) Speak(text, lang string) (string, error) {
	en, vc := s.getVoiceID(lang)

	if en == "neual" && neuralSupportedRegions.Search(en) < 0 {
		en = "standard"
	}

	input := &polly.SynthesizeSpeechInput{
		Engine:       aws.String(en),
		Text:         aws.String(text),
		VoiceId:      aws.String(vc),
		OutputFormat: aws.String("mp3"),
	}
	output, err := s.client.SynthesizeSpeech(input)
	if err != nil {
		return "", err
	}

	f, err := ioutil.TempFile("", "polly.*.mp3")
	if err != nil {
		return "", err
	}
	defer f.Close()
	if _, err = io.Copy(f, output.AudioStream); err != nil {
		return "", err
	}
	return f.Name(), nil
}

func (s *AWSPollyTTSService) getVoiceID(lang string) (string, string) {
	st := "standard"
	ne := "neural"

	switch lang {
	case "en":
		return ne, "Kendra"
	case "zh-CHS":
		return st, "Zhiyu"
	case "ja":
		return st, "Mizuki"
	case "ko":
		return st, "Seoyeon"
	case "fr":
		return st, "Celine"
	case "es":
		return ne, "Lupe"
	case "pt":
		return st, "Ines"
	case "it":
		return st, "Bianca"
	case "ru":
		return st, "Tatyana"
	case "de":
		return st, "Marlene"
	case "ar":
		return st, "Zeina"
	case "da":
		return st, "Naja"
	case "is":
		return st, "Dora"
	case "hi":
		return st, "Aditi"
	case "tr":
		return st, "Filiz"
	}
	return ne, "Kendra"
}

var (
	_ TTSService = (*AWSPollyTTSService)(nil)
)
