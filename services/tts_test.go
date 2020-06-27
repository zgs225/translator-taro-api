package services

import "testing"

func TestNewAWSPollyTTSService(t *testing.T) {
	_, err := NewAWSPollyTTSService("ap-northeast-2")

	if err != nil {
		t.Error("初始化 AWS Polly 错误: ", err)
	}
}

func TestAWSPollyTTSService_Speak(t *testing.T) {
	svc, _ := NewAWSPollyTTSService("ap-southeast-2")

	data := [][2]string{
		{"吃葡萄不吐葡萄皮，不吃葡萄倒吐葡萄皮", "zh-CHS"},
		{"My friend is afraid of spiders. This isn't very unusual; a lot of people are afraid of spiders.", "en"},
	}

	for _, ex := range data {
		file, err := svc.Speak(ex[0], ex[1])
		if err != nil {
			t.Fatal("AWSPollyTTSService_Speak error: text = ", ex[0], "; lang = ", ex[1], "; error = ", err.Error())
		}
		if file == "" {
			t.Error("AWSPollyTTSService_Speak return empty file path: text = ", ex[0], "; lang = ", ex[1])
		}
		t.Log("AWSPollyTTSService_Speak result: text = ", ex[0], "; lang = ", ex[1], "; file = ", file)
	}
}
