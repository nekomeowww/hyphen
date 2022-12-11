package dao

import (
	"os"
	"path/filepath"

	"github.com/nekomeowww/hyphen/pkg/utils"
	"go.etcd.io/bbolt"
	"go.uber.org/fx"
)

type NewBBoltParam struct {
	fx.In
}

type BBolt struct {
	*bbolt.DB
}

func NewBBolt(path string) func(NewBBoltParam) (*BBolt, error) {
	return func(param NewBBoltParam) (*BBolt, error) {
		bboltDB, err := bbolt.Open(path, 0600, nil)
		if err != nil {
			return nil, err
		}

		return &BBolt{
			DB: bboltDB,
		}, nil
	}
}

func NewTestBBolt() (*BBolt, func()) {
	testdataDir := utils.RelativePathOf("../../testdata")
	if _, err := os.Stat(testdataDir); os.IsNotExist(err) {
		err := os.Mkdir(testdataDir, 0755)
		if err != nil {
			panic(err)
		}
	}

	bboltDBPath := filepath.Join(testdataDir, "test.db")
	bbolt, err := NewBBolt(bboltDBPath)(NewBBoltParam{})
	if err != nil {
		panic(err)
	}

	return bbolt, func() {
		err := os.RemoveAll(bboltDBPath)
		if err != nil {
			panic(err)
		}
	}
}
