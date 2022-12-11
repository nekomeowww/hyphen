package urls

import (
	"bytes"
	"context"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/nekomeowww/hyphen/internal/dao"
	"github.com/nekomeowww/hyphen/internal/lib"
	"github.com/nekomeowww/hyphen/pkg/types/dao/bbolt/keys"
	urlsType "github.com/nekomeowww/hyphen/pkg/types/dao/bbolt/urls"
	"github.com/nekomeowww/hyphen/pkg/utils"
	"github.com/samber/lo"
	"github.com/samber/mo"
	"go.etcd.io/bbolt"
	"go.uber.org/fx"
)

type NewURLModelParam struct {
	fx.In

	Logger *lib.Logger
	BBolt  *dao.BBolt
}

type URLModel struct {
	Logger    *lib.Logger
	BBolt     *dao.BBolt
	URLBucket *bbolt.Bucket
}

func NewURLModel() func(NewURLModelParam) *URLModel {
	return func(param NewURLModelParam) *URLModel {
		return &URLModel{
			Logger: param.Logger,
			BBolt:  param.BBolt,
		}
	}
}

func (m *URLModel) NormalizeURL(urlString string) (string, error) {
	escapedURL, err := url.QueryUnescape(urlString)
	if err != nil {
		return "", err
	}

	parsedURL, err := url.Parse(escapedURL)
	if err != nil {
		return "", err
	}

	return parsedURL.String(), nil
}

func (m *URLModel) saveNewURL(tx *bbolt.Tx, bucket *bbolt.Bucket, fullURL string) mo.Result[string] {
	hash := fmt.Sprintf("%x", sha512.Sum512([]byte(fullURL)))[0:10]
	base64URL := base64.StdEncoding.EncodeToString([]byte(fullURL))

	base64EncodedShortURL := base64.StdEncoding.EncodeToString(lo.Must(json.Marshal(urlsType.ShortURL{ShortURL: hash, FullURL: fullURL})))
	err := bucket.Put(keys.FullURL1.Format(base64URL), []byte(base64EncodedShortURL))
	if err != nil {
		return mo.Err[string](err)
	}

	base64EncodedFullURL := base64.StdEncoding.EncodeToString(lo.Must(json.Marshal(urlsType.FullURL{ShortURL: hash, FullURL: fullURL})))
	err = bucket.Put(keys.ShortenedURL1.Format(hash), []byte(base64EncodedFullURL))
	if err != nil {
		return mo.Err[string](err)
	}

	return mo.Ok(hash)
}

func (m *URLModel) New(fullURL string) mo.Result[string] {
	tx, err := m.BBolt.Begin(true)
	if err != nil {
		return mo.Err[string](err)
	}

	bucket, err := tx.CreateBucketIfNotExists(keys.URLBucket.Format())
	if err != nil {
		_ = tx.Rollback()
		return mo.Err[string](err)
	}

	result := m.saveNewURL(tx, bucket, fullURL)
	if result.IsError() {
		_ = tx.Rollback()
		return result
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	err = utils.Invoke0(func() error {
		err = tx.Commit()
		if err != nil {
			return err
		}

		return nil
	}, utils.WithContext(ctx))
	if err != nil {
		return mo.Err[string](err)
	}

	return result
}

func (m *URLModel) FindOneShortURLByURL(url string) mo.Result[urlsType.ShortURL] {
	tx, err := m.BBolt.Begin(false)
	if err != nil {
		return mo.Err[urlsType.ShortURL](err)
	}
	defer func() {
		err = tx.Rollback()
		if err != nil {
			m.Logger.Error(err.Error())
			// PASS
		}
	}()

	bucket := tx.Bucket(keys.URLBucket.Format())
	if bucket == nil {
		return mo.Ok(urlsType.ShortURL{})
	}

	var foundShortURL []byte
	prefix := keys.FullURL1.Format(base64.StdEncoding.EncodeToString([]byte(url)))
	cursor := bucket.Cursor()
	for k, v := cursor.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = cursor.Next() {
		foundShortURL = v
		break
	}
	if string(foundShortURL) == "" {
		return mo.Ok(urlsType.ShortURL{})
	}

	decodedShortURL, err := base64.StdEncoding.DecodeString(string(foundShortURL))
	if err != nil {
		return mo.Err[urlsType.ShortURL](err)
	}

	var shortURL urlsType.ShortURL
	lo.Must0(json.Unmarshal(decodedShortURL, &shortURL))
	return mo.Ok(shortURL)
}

func (m *URLModel) FindOneURLByShortURL(shortURL string) mo.Result[urlsType.FullURL] {
	tx, err := m.BBolt.Begin(false)
	if err != nil {
		return mo.Err[urlsType.FullURL](err)
	}
	defer func() {
		err = tx.Rollback()
		if err != nil {
			m.Logger.Error(err.Error())
			// PASS
		}
	}()

	bucket := tx.Bucket(keys.URLBucket.Format())
	if bucket == nil {
		return mo.Ok(urlsType.FullURL{})
	}

	foundUrl := bucket.Get(keys.ShortenedURL1.Format(shortURL))
	if string(foundUrl) == "" {
		return mo.Ok(urlsType.FullURL{})
	}

	decodedURL, err := base64.StdEncoding.DecodeString(string(foundUrl))
	if err != nil {
		return mo.Err[urlsType.FullURL](err)
	}

	var fullURL urlsType.FullURL
	lo.Must0(json.Unmarshal(decodedURL, &fullURL))
	return mo.Ok(fullURL)
}

func (m *URLModel) RevokeOneShortURL(shortURL string) mo.Result[bool] {
	tx, err := m.BBolt.Begin(true)
	if err != nil {
		return mo.Err[bool](err)
	}

	bucket, err := tx.CreateBucketIfNotExists(keys.URLBucket.Format())
	if err != nil {
		_ = tx.Rollback()
		return mo.Err[bool](err)
	}

	result := m.FindOneURLByShortURL(shortURL)
	if result.IsError() {
		_ = tx.Rollback()
		return mo.Err[bool](result.Error())
	}
	if result.MustGet().FullURL == "" {
		_ = tx.Rollback()
		return mo.Ok(true)
	}

	base64URL := base64.StdEncoding.EncodeToString([]byte(result.MustGet().FullURL))
	err = bucket.Delete(keys.FullURL1.Format(base64URL))
	if err != nil {
		_ = tx.Rollback()
		return mo.Err[bool](err)
	}

	err = bucket.Delete(keys.ShortenedURL1.Format(shortURL))
	if err != nil {
		_ = tx.Rollback()
		return mo.Err[bool](err)
	}

	err = tx.Commit()
	if err != nil {
		return mo.Err[bool](err)
	}

	return mo.Ok(true)
}
