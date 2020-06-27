package services

import (
	"fmt"
	"translator-api/hash"
)

// UploadToQBoxTTSService 将语言文件上传到七牛中
type UploadToQBoxTTSService struct {
	next TTSService
	qbox QBox
	set  Set
}

// NewUploadToQBoxTTSServiceMiddleware 生成使用七牛上传语言文件的中间件
func NewUploadToQBoxTTSServiceMiddleware(q QBox, set Set) TTSServiceMiddleware {
	return func(next TTSService) TTSService {
		return &UploadToQBoxTTSService{
			next: next,
			set:  set,
			qbox: q,
		}
	}
}

// Speak 将生成的本地语音文件上传到七牛云上
func (s *UploadToQBoxTTSService) Speak(text, lang string) (string, error) {
	f := s.toFileName(text, lang)
	if s.set.Exists(f) {
		return f, nil
	}

	local, err := s.next.Speak(text, lang)
	if err != nil {
		return "", err
	}

	ret, err := s.qbox.UploadFile(local, f)
	if err != nil {
		return "", err
	}
	s.set.Add(ret.Key)
	return ret.Key, nil
}

func (s *UploadToQBoxTTSService) toFileName(text, lang string) string {
	str := fmt.Sprintf("%s@%s", text, lang)
	hashStr := hash.SHA256(str)
	return fmt.Sprintf("polly/%s.mp3", hashStr)
}
