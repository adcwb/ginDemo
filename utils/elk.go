package utils

import (
	"bytes"
	"context"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"go.uber.org/zap"
	"io"
)

/*
	定义将数据存入elasticsearch集群中
*/

// CreateIndex 创建 index
func CreateIndex(index string, client *elasticsearch.TypedClient) bool {
	resp, err := client.Indices.Create(index).Do(context.Background())
	if err != nil {
		zap.L().Error("create index failed", zap.Error(err))
		return false
	}
	zap.L().Info("create index success", zap.String("index", resp.Index))
	return true
}

// DeleteIndex 删除 index
func DeleteIndex(index string, client *elasticsearch.TypedClient) bool {
	resp, err := client.Indices.Delete(index).Do(context.Background())
	if err != nil {
		zap.L().Error("delete index failed", zap.Error(err))
		return false
	}
	zap.L().Info("delete index success", zap.String("index", index), zap.Bool("Acknowledged", resp.Acknowledged))
	return resp.Acknowledged
}

// GetDocument 获取文档
func GetDocument(client *elasticsearch.Client, indexName, docID string, doc []byte) error {
	req := esapi.GetRequest{
		Index:      indexName,
		DocumentID: docID,
	}

	// Perform the request
	res, err := req.Do(context.Background(), client)
	if err != nil {
		return fmt.Errorf("error performing request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			zap.L().Error("close error", zap.Error(err))
		}
	}(res.Body)
	if res.IsError() {
		return fmt.Errorf("error indexing document: %s", res.String())
	}

	return nil
}

// CreateDocument 创建文档
func CreateDocument(client *elasticsearch.Client, indexName, docID string, doc []byte) error {
	req := esapi.IndexRequest{
		Index:      indexName,
		DocumentID: docID,
		Body:       bytes.NewReader(doc),
		Refresh:    "true",
	}

	// Perform the request
	res, err := req.Do(context.Background(), client)
	if err != nil {
		return fmt.Errorf("error performing request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			zap.L().Error("close error", zap.Error(err))
		}
	}(res.Body)
	if res.IsError() {
		return fmt.Errorf("error indexing document: %s", res.String())
	}

	return nil
}

// UpdateDocument 更新文档
func UpdateDocument(client *elasticsearch.Client, indexName, docID string, docNew []byte) error {
	req := esapi.UpdateRequest{
		Index:      indexName,
		DocumentID: docID,
		Body:       bytes.NewReader(docNew),
		Refresh:    "true",
	}

	// Perform the request
	res, err := req.Do(context.Background(), client)
	if err != nil {
		return fmt.Errorf("error performing request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			zap.L().Error("close error", zap.Error(err))
		}
	}(res.Body)
	if res.IsError() {
		return fmt.Errorf("error indexing document: %s", res.String())
	}

	return nil
}

// 判断文档是否存在

// DeleteDocument 删除文档
func DeleteDocument(client *elasticsearch.Client, indexName, docID string) error {
	req := esapi.DeleteRequest{
		Index:      indexName,
		DocumentID: docID,
	}

	// Perform the request
	res, err := req.Do(context.Background(), client)
	if err != nil {
		return fmt.Errorf("error performing request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			zap.L().Error("close error", zap.Error(err))
		}
	}(res.Body)
	if res.IsError() {
		return fmt.Errorf("error indexing document: %s", res.String())
	}

	return nil
}
