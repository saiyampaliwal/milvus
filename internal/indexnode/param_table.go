// Copyright (C) 2019-2020 Zilliz. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance
// with the License. You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed under the License
// is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
// or implied. See the License for the specific language governing permissions and limitations under the License.

package indexnode

import (
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/milvus-io/milvus/internal/log"
	"github.com/milvus-io/milvus/internal/util/paramtable"
)

const (
	StartParamsKey = "START_PARAMS"
)

// ParamTable is used to record configuration items.
type ParamTable struct {
	paramtable.BaseTable

	IP      string
	Address string
	Port    int

	NodeID int64
	Alias  string

	EtcdEndpoints []string
	MetaRootPath  string
	IndexRootPath string

	MinIOAddress         string
	MinIOAccessKeyID     string
	MinIOSecretAccessKey string
	MinIOUseSSL          bool
	MinioBucketName      string

	SimdType string

	CreatedTime time.Time
	UpdatedTime time.Time
}

// Params is an alias for ParamTable.
var Params ParamTable
var once sync.Once

func (pt *ParamTable) InitAlias(alias string) {
	pt.Alias = alias
}

// Init is used to initialize configuration items.
func (pt *ParamTable) Init() {
	pt.BaseTable.Init()
	if err := pt.LoadYaml("advanced/knowhere.yaml"); err != nil {
		panic(err)
	}

	// TODO, load index_node.yaml
	/*err := pt.LoadYaml("advanced/index_node.yaml")
	if err != nil {
		panic(err)
	}*/

	pt.initParams()
	pt.initKnowhereSimdType()
}

// InitOnce is used to initialize configuration items, and it will only be called once.
func (pt *ParamTable) InitOnce() {
	once.Do(func() {
		pt.Init()
	})
}

func (pt *ParamTable) initParams() {
	pt.initMinIOAddress()
	pt.initMinIOAccessKeyID()
	pt.initMinIOSecretAccessKey()
	pt.initMinIOUseSSL()
	pt.initMinioBucketName()
	pt.initEtcdEndpoints()
	pt.initMetaRootPath()
	pt.initIndexRootPath()
	pt.initRoleName()
}

func (pt *ParamTable) initMinIOAddress() {
	ret, err := pt.Load("_MinioAddress")
	if err != nil {
		panic(err)
	}
	pt.MinIOAddress = ret
}

func (pt *ParamTable) initMinIOAccessKeyID() {
	ret, err := pt.Load("minio.accessKeyID")
	if err != nil {
		panic(err)
	}
	pt.MinIOAccessKeyID = ret
}

func (pt *ParamTable) initMinIOSecretAccessKey() {
	ret, err := pt.Load("minio.secretAccessKey")
	if err != nil {
		panic(err)
	}
	pt.MinIOSecretAccessKey = ret
}

func (pt *ParamTable) initMinIOUseSSL() {
	ret, err := pt.Load("minio.useSSL")
	if err != nil {
		panic(err)
	}
	pt.MinIOUseSSL, err = strconv.ParseBool(ret)
	if err != nil {
		panic(err)
	}
}

func (pt *ParamTable) initEtcdEndpoints() {
	endpoints, err := pt.Load("_EtcdEndpoints")
	if err != nil {
		panic(err)
	}
	pt.EtcdEndpoints = strings.Split(endpoints, ",")
}

func (pt *ParamTable) initMetaRootPath() {
	rootPath, err := pt.Load("etcd.rootPath")
	if err != nil {
		panic(err)
	}
	subPath, err := pt.Load("etcd.metaSubPath")
	if err != nil {
		panic(err)
	}
	pt.MetaRootPath = path.Join(rootPath, subPath)
}

func (pt *ParamTable) initIndexRootPath() {
	rootPath, err := pt.Load("minio.rootPath")
	if err != nil {
		panic(err)
	}
	pt.IndexRootPath = path.Join(rootPath, "index_files")
}

func (pt *ParamTable) initMinioBucketName() {
	bucketName, err := pt.Load("minio.bucketName")
	if err != nil {
		panic(err)
	}
	pt.MinioBucketName = bucketName
}

func (pt *ParamTable) initRoleName() {
	pt.RoleName = "indexnode"
}

func (pt *ParamTable) initKnowhereSimdType() {
	simdType, err := pt.LoadWithDefault("knowhere.simdType", "auto")
	if err != nil {
		log.Error("failed to initialize the simd type",
			zap.Error(err))

		panic(err)
	}

	pt.SimdType = simdType

	log.Debug("initialize the knowhere simd type",
		zap.String("simd_type", pt.SimdType))
}
