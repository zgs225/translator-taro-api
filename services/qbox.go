package services

import (
	"bytes"
	"context"
	"sync"
	"time"

	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/storage"
	"github.com/sirupsen/logrus"
)

// QBox 七牛服务
type QBox interface {
	UploadFile(local string, key ...string) (*storage.PutRet, error)
	MakePrivateURL(key, domain string, duration time.Duration) string
}

// SimpleQBox 实现基本的七牛服务
type SimpleQBox struct {
	mu sync.RWMutex

	accessKey string
	secretKey string
	bucket    string
	zone      *storage.Zone
	logger    logrus.FieldLogger

	uploader *storage.ResumeUploader
	token    string
}

// NewSimpleQBox 初始化 QBox
func NewSimpleQBox(ak, sk, bucket string, logger logrus.FieldLogger) (*SimpleQBox, error) {
	qb := SimpleQBox{
		accessKey: ak,
		bucket:    bucket,
		secretKey: sk,
	}
	qb.logger = logger.WithFields(logrus.Fields{
		"access_key": ak[:5] + string(bytes.Repeat([]byte{'*'}, len(ak)-10)) + ak[len(ak)-5:],
		"secret_key": sk[:5] + string(bytes.Repeat([]byte{'*'}, len(sk)-10)) + ak[len(sk)-5:],
		"bucket":     bucket,
	})

	zone, err := storage.GetZone(ak, bucket)
	if err != nil {
		return nil, err
	}
	qb.zone = zone

	cfg := storage.Config{
		Zone:          zone,
		UseHTTPS:      true,
		UseCdnDomains: false,
	}

	qb.uploader = storage.NewResumeUploader(&cfg)

	return &qb, nil
}

// UploadFile 上传本地文件到指定的桶
func (q *SimpleQBox) UploadFile(local string, key ...string) (ret *storage.PutRet, err error) {
	var k string
	if len(key) > 0 {
		k = key[0]
	}

	defer func(b time.Time) {
		l := q.logger.WithField("file", local).WithField("duration", time.Since(b))
		if len(k) > 0 {
			l = l.WithField("key", k)
		}
		if ret != nil {
			l = l.WithFields(logrus.Fields{
				"hash": ret.Hash,
				"key":  ret.Key,
				"id":   ret.PersistentID,
			})
		}
		if err == nil {
			l.Info("文件上传成功")
		} else {
			l.WithError(err).Error("文件上传失败")
		}
	}(time.Now())

	token := q.getToken()
	ret = &storage.PutRet{}
	err = q.uploader.PutFile(context.Background(), ret, token, k, local, &storage.RputExtra{})

	if err != nil && err == storage.ErrBadToken {
		q.clearToken()
		return q.UploadFile(local, key...)
	}

	return ret, err
}

// MakePrivateURL 生成私有空间访问链接
func (q *SimpleQBox) MakePrivateURL(key, domain string, duration time.Duration) string {
	mac := qbox.NewMac(q.accessKey, q.secretKey)
	return storage.MakePrivateURL(mac, domain, key, time.Now().Add(duration).Unix())
}

func (q *SimpleQBox) getToken() string {
	q.mu.RLock()
	if len(q.token) > 0 {
		defer q.mu.RUnlock()
		return q.token
	}
	q.mu.RUnlock()

	var expires uint64 = 86400
	mac := qbox.NewMac(q.accessKey, q.secretKey)
	policy := storage.PutPolicy{
		Scope:   q.bucket,
		Expires: expires,
	}
	q.mu.Lock()
	q.token = policy.UploadToken(mac)
	q.mu.Unlock()
	timer := time.NewTimer(time.Duration(expires-60) * time.Second)

	go func() {
		<-timer.C
		q.clearToken()
	}()

	return q.token
}

func (q *SimpleQBox) clearToken() {
	q.mu.Lock()
	q.token = ""
	q.mu.Unlock()
}
