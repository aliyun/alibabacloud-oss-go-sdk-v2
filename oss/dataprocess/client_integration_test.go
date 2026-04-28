//go:build integration

package dataprocess

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/stretchr/testify/assert"
)

var (
	region_    = os.Getenv("OSS_TEST_REGION")
	endpoint_  = os.Getenv("OSS_TEST_ENDPOINT")
	accessID_  = os.Getenv("OSS_TEST_ACCESS_KEY_ID")
	accessKey_ = os.Getenv("OSS_TEST_ACCESS_KEY_SECRET")
	bucket_    = os.Getenv("OSS_TEST_DATAPROCESS_BUCKET")

	instance_ *Client
	testOnce_ sync.Once
)

var (
	datasetNamePrefix = "go-sdk-test-ds-"
	letters           = []rune("abcdefghijklmnopqrstuvwxyz")
)

func getDefaultClient() *Client {
	testOnce_.Do(func() {
		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessID_, accessKey_)).
			WithRegion(region_).
			WithEndpoint(endpoint_)

		instance_ = NewClient(cfg)
	})
	return instance_
}

func getInvalidAkClient() *Client {
	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider("invalid-ak", "invalid-sk")).
		WithRegion(region_).
		WithEndpoint(endpoint_)

	return NewClient(cfg)
}

func randStr(n int) string {
	b := make([]rune, n)
	randMarker := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range b {
		b[i] = letters[randMarker.Intn(len(letters))]
	}
	return string(b)
}

func genDatasetName() string {
	return datasetNamePrefix + fmt.Sprintf("%d", time.Now().UnixMilli()) + "-" + randStr(5)
}

func cleanDatasets(prefix string, t *testing.T) {
	c := getDefaultClient()
	request := &ListDatasetsRequest{
		Bucket: oss.Ptr(bucket_),
		Prefix: oss.Ptr(prefix),
	}
	result, err := c.ListDatasets(context.TODO(), request)
	if err != nil {
		return
	}
	for _, ds := range result.Datasets {
		if ds.DatasetName != nil && strings.HasPrefix(*ds.DatasetName, prefix) {
			_, _ = c.DeleteDataset(context.TODO(), &DeleteDatasetRequest{
				Bucket:      oss.Ptr(bucket_),
				DatasetName: ds.DatasetName,
			})
		}
	}
}

func dumpErrIfNotNil(err error) {
	if err != nil {
		fmt.Printf("error:%s\n", err.Error())
	}
}

func TestDatasetLifecycle(t *testing.T) {
	client := getDefaultClient()
	dsName := genDatasetName()

	// 1. Create dataset
	createResult, err := client.CreateDataset(context.TODO(), &CreateDatasetRequest{
		Bucket:      oss.Ptr(bucket_),
		DatasetName: oss.Ptr(dsName),
		Description: oss.Ptr("integration test dataset"),
	})
	dumpErrIfNotNil(err)
	assert.Nil(t, err)
	assert.Equal(t, 200, createResult.StatusCode)
	assert.NotNil(t, createResult.Dataset)
	assert.Equal(t, dsName, *createResult.Dataset.DatasetName)

	defer func() {
		_, _ = client.DeleteDataset(context.TODO(), &DeleteDatasetRequest{
			Bucket:      oss.Ptr(bucket_),
			DatasetName: oss.Ptr(dsName),
		})
	}()

	// 2. Get dataset
	getResult, err := client.GetDataset(context.TODO(), &GetDatasetRequest{
		Bucket:      oss.Ptr(bucket_),
		DatasetName: oss.Ptr(dsName),
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, getResult.StatusCode)
	assert.NotNil(t, getResult.Dataset)
	assert.Equal(t, dsName, *getResult.Dataset.DatasetName)
	assert.Equal(t, "integration test dataset", *getResult.Dataset.Description)

	// 3. Get dataset with statistics
	getWithStatsResult, err := client.GetDataset(context.TODO(), &GetDatasetRequest{
		Bucket:         oss.Ptr(bucket_),
		DatasetName:    oss.Ptr(dsName),
		WithStatistics: oss.Ptr(true),
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, getWithStatsResult.StatusCode)
	assert.NotNil(t, getWithStatsResult.Dataset)

	// 4. Update dataset
	updateResult, err := client.UpdateDataset(context.TODO(), &UpdateDatasetRequest{
		Bucket:      oss.Ptr(bucket_),
		DatasetName: oss.Ptr(dsName),
		Description: oss.Ptr("updated description 1"),
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, updateResult.StatusCode)
	assert.NotNil(t, updateResult.Dataset)

	// 5. Verify update
	getAfterUpdate, err := client.GetDataset(context.TODO(), &GetDatasetRequest{
		Bucket:      oss.Ptr(bucket_),
		DatasetName: oss.Ptr(dsName),
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, getAfterUpdate.StatusCode)
	assert.Equal(t, "updated description 1", *getAfterUpdate.Dataset.Description)

	// 6. Delete dataset
	deleteResult, err := client.DeleteDataset(context.TODO(), &DeleteDatasetRequest{
		Bucket:      oss.Ptr(bucket_),
		DatasetName: oss.Ptr(dsName),
	})
	assert.Nil(t, err)
	assert.True(t, deleteResult.StatusCode == 200 || deleteResult.StatusCode == 204)
}

func TestCreateAndDeleteDataset(t *testing.T) {
	client := getDefaultClient()
	dsName := genDatasetName()

	// Create
	createResult, err := client.CreateDataset(context.TODO(), &CreateDatasetRequest{
		Bucket:      oss.Ptr(bucket_),
		DatasetName: oss.Ptr(dsName),
	})
	dumpErrIfNotNil(err)
	assert.Nil(t, err)
	assert.Equal(t, 200, createResult.StatusCode)

	// Delete
	deleteResult, err := client.DeleteDataset(context.TODO(), &DeleteDatasetRequest{
		Bucket:      oss.Ptr(bucket_),
		DatasetName: oss.Ptr(dsName),
	})
	assert.Nil(t, err)
	assert.True(t, deleteResult.StatusCode == 200 || deleteResult.StatusCode == 204)
}

func TestGetNonExistentDataset(t *testing.T) {
	client := getDefaultClient()

	_, err := client.GetDataset(context.TODO(), &GetDatasetRequest{
		Bucket:      oss.Ptr(bucket_),
		DatasetName: oss.Ptr("non-existent-dataset-" + fmt.Sprintf("%d", time.Now().UnixMilli())),
	})
	assert.NotNil(t, err)

	var serr *oss.ServiceError
	errors.As(err, &serr)
	assert.NotNil(t, serr)
	assert.True(t, serr.StatusCode == 404 || serr.StatusCode == 400,
		"Expected 404 or 400 status, got %d", serr.StatusCode)
}

func TestUpdateDatasetWithConfig(t *testing.T) {
	client := getDefaultClient()
	dsName := genDatasetName()

	// Create dataset
	createResult, err := client.CreateDataset(context.TODO(), &CreateDatasetRequest{
		Bucket:      oss.Ptr(bucket_),
		DatasetName: oss.Ptr(dsName),
	})
	dumpErrIfNotNil(err)
	assert.Nil(t, err)
	assert.Equal(t, 200, createResult.StatusCode)

	defer func() {
		_, _ = client.DeleteDataset(context.TODO(), &DeleteDatasetRequest{
			Bucket:      oss.Ptr(bucket_),
			DatasetName: oss.Ptr(dsName),
		})
	}()

	// Update with DatasetConfig (Language = "en")
	configJSON := `{"Insights":{"Language":"en"}}`
	updateResult, err := client.UpdateDataset(context.TODO(), &UpdateDatasetRequest{
		Bucket:        oss.Ptr(bucket_),
		DatasetName:   oss.Ptr(dsName),
		DatasetConfig: oss.Ptr(configJSON),
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, updateResult.StatusCode)

	// Verify DatasetConfig is returned correctly
	getResult, err := client.GetDataset(context.TODO(), &GetDatasetRequest{
		Bucket:      oss.Ptr(bucket_),
		DatasetName: oss.Ptr(dsName),
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, getResult.StatusCode)
	assert.NotNil(t, getResult.Dataset.DatasetConfig)
	assert.NotNil(t, getResult.Dataset.DatasetConfig.Insights)
	assert.Equal(t, "en", *getResult.Dataset.DatasetConfig.Insights.Language)
}

func TestCreateDatasetWithDatasetConfig(t *testing.T) {
	client := getDefaultClient()
	dsName := genDatasetName()

	configJSON := `{"Insights":{"Language":"ch"}}`
	createResult, err := client.CreateDataset(context.TODO(), &CreateDatasetRequest{
		Bucket:        oss.Ptr(bucket_),
		DatasetName:   oss.Ptr(dsName),
		DatasetConfig: oss.Ptr(configJSON),
	})
	dumpErrIfNotNil(err)
	assert.Nil(t, err)
	assert.Equal(t, 200, createResult.StatusCode)

	defer func() {
		_, _ = client.DeleteDataset(context.TODO(), &DeleteDatasetRequest{
			Bucket:      oss.Ptr(bucket_),
			DatasetName: oss.Ptr(dsName),
		})
	}()

	// Verify DatasetConfig is returned in get
	getResult, err := client.GetDataset(context.TODO(), &GetDatasetRequest{
		Bucket:      oss.Ptr(bucket_),
		DatasetName: oss.Ptr(dsName),
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, getResult.StatusCode)
	assert.NotNil(t, getResult.Dataset.DatasetConfig)
	assert.NotNil(t, getResult.Dataset.DatasetConfig.Insights)
	assert.Equal(t, "ch", *getResult.Dataset.DatasetConfig.Insights.Language)
}

func TestCreateDatasetWithWorkflowParameters(t *testing.T) {
	client := getDefaultClient()
	dsName := genDatasetName()

	workflowParamsJSON := `[{"Name":"VideoInsightEnable","Value":"true"}]`
	createResult, err := client.CreateDataset(context.TODO(), &CreateDatasetRequest{
		Bucket:             oss.Ptr(bucket_),
		DatasetName:        oss.Ptr(dsName),
		Description:        oss.Ptr("test with workflow parameters"),
		WorkflowParameters: oss.Ptr(workflowParamsJSON),
	})
	dumpErrIfNotNil(err)
	assert.Nil(t, err)
	assert.Equal(t, 200, createResult.StatusCode)
	assert.NotNil(t, createResult.Dataset)
	assert.Equal(t, dsName, *createResult.Dataset.DatasetName)

	defer func() {
		_, _ = client.DeleteDataset(context.TODO(), &DeleteDatasetRequest{
			Bucket:      oss.Ptr(bucket_),
			DatasetName: oss.Ptr(dsName),
		})
	}()

	// Verify workflow parameters are returned in get
	getResult, err := client.GetDataset(context.TODO(), &GetDatasetRequest{
		Bucket:      oss.Ptr(bucket_),
		DatasetName: oss.Ptr(dsName),
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, getResult.StatusCode)
	assert.NotNil(t, getResult.Dataset)
	assert.NotNil(t, getResult.Dataset.WorkflowParameters, "WorkflowParameters should be returned")

	returnedParams := getResult.Dataset.WorkflowParameters.WorkflowParameter
	assert.NotNil(t, returnedParams, "WorkflowParameter list should be returned")
	assert.Equal(t, 1, len(returnedParams))
	assert.Equal(t, "VideoInsightEnable", *returnedParams[0].Name)
	assert.Equal(t, "true", *returnedParams[0].Value)
}

func TestUpdateDatasetWithWorkflowParameters(t *testing.T) {
	client := getDefaultClient()
	dsName := genDatasetName()

	// Create dataset without workflow parameters
	createResult, err := client.CreateDataset(context.TODO(), &CreateDatasetRequest{
		Bucket:      oss.Ptr(bucket_),
		DatasetName: oss.Ptr(dsName),
	})
	dumpErrIfNotNil(err)
	assert.Nil(t, err)
	assert.Equal(t, 200, createResult.StatusCode)

	defer func() {
		_, _ = client.DeleteDataset(context.TODO(), &DeleteDatasetRequest{
			Bucket:      oss.Ptr(bucket_),
			DatasetName: oss.Ptr(dsName),
		})
	}()

	// Update with workflow parameters
	workflowParamsJSON := `[{"Name":"VideoInsightEnable","Value":"true"}]`
	updateResult, err := client.UpdateDataset(context.TODO(), &UpdateDatasetRequest{
		Bucket:             oss.Ptr(bucket_),
		DatasetName:        oss.Ptr(dsName),
		WorkflowParameters: oss.Ptr(workflowParamsJSON),
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, updateResult.StatusCode)

	// Verify update by getting the dataset
	getResult, err := client.GetDataset(context.TODO(), &GetDatasetRequest{
		Bucket:      oss.Ptr(bucket_),
		DatasetName: oss.Ptr(dsName),
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, getResult.StatusCode)
	assert.NotNil(t, getResult.Dataset.WorkflowParameters, "WorkflowParameters should be returned after update")

	returnedParams := getResult.Dataset.WorkflowParameters.WorkflowParameter
	assert.NotNil(t, returnedParams, "WorkflowParameter list should be returned after update")
	assert.Equal(t, 1, len(returnedParams))
	assert.Equal(t, "VideoInsightEnable", *returnedParams[0].Name)
	assert.Equal(t, "true", *returnedParams[0].Value)
}

func TestListDatasets(t *testing.T) {
	client := getDefaultClient()

	// Use a unique prefix for this test to isolate from other datasets
	testPrefix := datasetNamePrefix + "list-" + fmt.Sprintf("%d", time.Now().UnixMilli()) + "-"
	dsName1 := testPrefix + "a"
	dsName2 := testPrefix + "b"
	dsName3 := testPrefix + "c"

	// Create 3 datasets
	for _, name := range []string{dsName1, dsName2, dsName3} {
		cr, err := client.CreateDataset(context.TODO(), &CreateDatasetRequest{
			Bucket:      oss.Ptr(bucket_),
			DatasetName: oss.Ptr(name),
		})
		dumpErrIfNotNil(err)
		assert.Nil(t, err)
		assert.Equal(t, 200, cr.StatusCode, "create %s should return 200", name)
	}

	defer func() {
		for _, name := range []string{dsName1, dsName2, dsName3} {
			_, _ = client.DeleteDataset(context.TODO(), &DeleteDatasetRequest{
				Bucket:      oss.Ptr(bucket_),
				DatasetName: oss.Ptr(name),
			})
		}
	}()

	// 1. List with prefix, verify all 3 datasets are returned
	listAll, err := client.ListDatasets(context.TODO(), &ListDatasetsRequest{
		Bucket: oss.Ptr(bucket_),
		Prefix: oss.Ptr(testPrefix),
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, listAll.StatusCode)
	assert.NotNil(t, listAll.Datasets, "datasets should not be nil")
	assert.Equal(t, 3, len(listAll.Datasets), "should list exactly 3 datasets with prefix")

	listedNames := make(map[string]bool)
	for _, ds := range listAll.Datasets {
		assert.NotNil(t, ds.DatasetName, "dataset name should not be nil")
		listedNames[*ds.DatasetName] = true
	}
	assert.True(t, listedNames[dsName1], "dsName1 should be in list")
	assert.True(t, listedNames[dsName2], "dsName2 should be in list")
	assert.True(t, listedNames[dsName3], "dsName3 should be in list")

	// 2. Paginate with maxResults=1
	paginatedNames := make(map[string]bool)
	var nextToken *string
	pageCount := 0

	for {
		req := &ListDatasetsRequest{
			Bucket:     oss.Ptr(bucket_),
			Prefix:     oss.Ptr(testPrefix),
			MaxResults: oss.Ptr(int64(1)),
			NextToken:  nextToken,
		}

		pageResult, err := client.ListDatasets(context.TODO(), req)
		assert.Nil(t, err)
		assert.Equal(t, 200, pageResult.StatusCode)
		assert.NotNil(t, pageResult.Datasets)
		assert.Equal(t, 1, len(pageResult.Datasets), "each page should have exactly 1 dataset")

		paginatedNames[*pageResult.Datasets[0].DatasetName] = true
		nextToken = pageResult.NextToken
		pageCount++

		assert.True(t, pageCount <= 10, "pagination should not exceed 10 pages")

		if nextToken == nil || *nextToken == "" {
			break
		}
	}

	assert.Equal(t, 3, pageCount, "should have paginated through 3 pages")
	assert.True(t, paginatedNames[dsName1], "paginated dsName1 should be found")
	assert.True(t, paginatedNames[dsName2], "paginated dsName2 should be found")
	assert.True(t, paginatedNames[dsName3], "paginated dsName3 should be found")

	// 3. Verify nextToken is nil when all results returned
	fullPage, err := client.ListDatasets(context.TODO(), &ListDatasetsRequest{
		Bucket:     oss.Ptr(bucket_),
		Prefix:     oss.Ptr(testPrefix),
		MaxResults: oss.Ptr(int64(100)),
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, fullPage.StatusCode)
	assert.Equal(t, 3, len(fullPage.Datasets))
	assert.True(t, fullPage.NextToken == nil || *fullPage.NextToken == "",
		"nextToken should be nil or empty when all results returned")
}

func TestListDatasetsPaginator(t *testing.T) {
	client := getDefaultClient()

	testPrefix := datasetNamePrefix + "pag-" + fmt.Sprintf("%d", time.Now().UnixMilli()) + "-"
	dsNames := make([]string, 3)
	for i := 0; i < 3; i++ {
		dsNames[i] = testPrefix + string(rune('a'+i))
		cr, err := client.CreateDataset(context.TODO(), &CreateDatasetRequest{
			Bucket:      oss.Ptr(bucket_),
			DatasetName: oss.Ptr(dsNames[i]),
		})
		dumpErrIfNotNil(err)
		assert.Nil(t, err)
		assert.Equal(t, 200, cr.StatusCode)
	}

	defer func() {
		for _, name := range dsNames {
			_, _ = client.DeleteDataset(context.TODO(), &DeleteDatasetRequest{
				Bucket:      oss.Ptr(bucket_),
				DatasetName: oss.Ptr(name),
			})
		}
	}()

	// Use paginator with limit=1
	request := &ListDatasetsRequest{
		Bucket:     oss.Ptr(bucket_),
		Prefix:     oss.Ptr(testPrefix),
		MaxResults: oss.Ptr(int64(1)),
	}
	paginator := client.NewListDatasetsPaginator(request)

	pageCount := 0
	totalDatasets := 0
	for paginator.HasNext() {
		pageCount++
		page, err := paginator.NextPage(context.TODO())
		assert.Nil(t, err)
		assert.True(t, len(page.Datasets) > 0)
		totalDatasets += len(page.Datasets)
	}
	assert.Equal(t, 3, pageCount, "should have 3 pages")
	assert.Equal(t, 3, totalDatasets, "should have 3 total datasets")
}

func TestDatasetServerErrors(t *testing.T) {
	client := getInvalidAkClient()

	var serr *oss.ServiceError

	// Create with invalid AK
	_, err := client.CreateDataset(context.TODO(), &CreateDatasetRequest{
		Bucket:      oss.Ptr(bucket_),
		DatasetName: oss.Ptr("test-invalid-ak"),
	})
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, 403, serr.StatusCode)
	assert.NotEmpty(t, serr.RequestID)

	// Get with invalid AK
	serr = nil
	_, err = client.GetDataset(context.TODO(), &GetDatasetRequest{
		Bucket:      oss.Ptr(bucket_),
		DatasetName: oss.Ptr("test-invalid-ak"),
	})
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, 403, serr.StatusCode)
	assert.NotEmpty(t, serr.RequestID)

	// List with invalid AK
	serr = nil
	_, err = client.ListDatasets(context.TODO(), &ListDatasetsRequest{
		Bucket: oss.Ptr(bucket_),
	})
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, 403, serr.StatusCode)
	assert.NotEmpty(t, serr.RequestID)

	// Delete with invalid AK
	serr = nil
	_, err = client.DeleteDataset(context.TODO(), &DeleteDatasetRequest{
		Bucket:      oss.Ptr(bucket_),
		DatasetName: oss.Ptr("test-invalid-ak"),
	})
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, 403, serr.StatusCode)
	assert.NotEmpty(t, serr.RequestID)
}

func TestSimpleQueryBasic(t *testing.T) {
	client := getDefaultClient()

	queryJSON := `{"Field":"Filename","Value":"test-media","Operation":"prefix"}`
	result, err := client.SimpleQuery(context.TODO(), &SimpleQueryRequest{
		Bucket:      oss.Ptr(bucket_),
		DatasetName: oss.Ptr("test-dataset"),
		Query:       oss.Ptr(queryJSON),
		MaxResults:  oss.Ptr(int32(10)),
	})
	dumpErrIfNotNil(err)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)

	assert.NotNil(t, result.Files, "files should not be nil")
	assert.True(t, len(result.Files) > 0, "files should not be empty")

	f := result.Files[0]
	assert.NotNil(t, f.URI, "URI should not be nil")
	assert.NotNil(t, f.Filename, "Filename should not be nil")
	assert.NotNil(t, f.MediaType, "MediaType should not be nil")
	assert.NotNil(t, f.ContentType, "ContentType should not be nil")
	assert.NotNil(t, f.Size, "Size should not be nil")
	assert.True(t, *f.Size > 0, "Size should be > 0")

	// Verify labels parsing
	anyFileHasLabels := false
	for _, file := range result.Files {
		if oss.ToInt64(file.OSSTaggingCount) > 0 {
			assert.Equal(t, *file.OSSTaggingCount, (int64)(len(file.OSSTagging)))
		}

		if len(file.Labels) > 0 {
			anyFileHasLabels = true
			assert.NotNil(t, file.Labels[0].LabelName, "Label.LabelName should not be nil")
			break
		}
	}
	assert.True(t, anyFileHasLabels, "At least one file should have labels")
}

func TestSimpleQueryWithAggregations(t *testing.T) {
	client := getDefaultClient()

	queryJSON := `{"Field":"Filename","Value":"test-media","Operation":"prefix"}`
	aggregationsJSON := `[{"Field":"Size","Operation":"sum"}]`

	result, err := client.SimpleQuery(context.TODO(), &SimpleQueryRequest{
		Bucket:       oss.Ptr(bucket_),
		DatasetName:  oss.Ptr("test-dataset"),
		Query:        oss.Ptr(queryJSON),
		Aggregations: oss.Ptr(aggregationsJSON),
		MaxResults:   oss.Ptr(int32(10)),
	})
	dumpErrIfNotNil(err)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)

	// Verify aggregations
	assert.NotNil(t, result.Aggregations, "aggregations should not be nil")
	assert.True(t, len(result.Aggregations) > 0, "aggregations should not be empty")

	agg := result.Aggregations[0]
	assert.Equal(t, "Size", *agg.Field)
	assert.Equal(t, "sum", *agg.Operation)
	assert.NotNil(t, agg.Value, "aggregation value should not be nil")
}

func TestSimpleQueryWithSortAndOrder(t *testing.T) {
	client := getDefaultClient()

	queryJSON := `{"Field":"Filename","Value":"test-media","Operation":"prefix"}`
	result, err := client.SimpleQuery(context.TODO(), &SimpleQueryRequest{
		Bucket:           oss.Ptr(bucket_),
		DatasetName:      oss.Ptr("test-dataset"),
		Query:            oss.Ptr(queryJSON),
		Sort:             oss.Ptr("Filename"),
		Order:            oss.Ptr("asc"),
		MaxResults:       oss.Ptr(int32(10)),
		WithoutTotalHits: oss.Ptr(false),
	})
	dumpErrIfNotNil(err)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)

	assert.NotNil(t, result.Files, "files should not be nil")
	assert.True(t, len(result.Files) > 0, "files should not be empty")

	// Verify sorted by Filename ascending
	files := result.Files
	for i := 1; i < len(files); i++ {
		prev := *files[i-1].Filename
		curr := *files[i].Filename
		assert.True(t, prev <= curr,
			"files should be sorted by Filename asc: %s <= %s", prev, curr)
	}
}

func TestSimpleQueryWithFields(t *testing.T) {
	client := getDefaultClient()

	queryJSON := `{"Field":"Filename","Value":"test-media","Operation":"prefix"}`
	withFieldsJSON := `["Filename","Size","ContentType"]`

	result, err := client.SimpleQuery(context.TODO(), &SimpleQueryRequest{
		Bucket:      oss.Ptr(bucket_),
		DatasetName: oss.Ptr("test-dataset"),
		Query:       oss.Ptr(queryJSON),
		WithFields:  oss.Ptr(withFieldsJSON),
		MaxResults:  oss.Ptr(int32(10)),
	})
	dumpErrIfNotNil(err)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)

	assert.NotNil(t, result.Files, "files should not be nil")
	assert.True(t, len(result.Files) > 0, "files should not be empty")

	// Verify requested fields are populated
	f := result.Files[0]
	assert.NotNil(t, f.Filename, "Filename should not be nil")
	assert.NotNil(t, f.Size, "Size should not be nil")
	assert.NotNil(t, f.ContentType, "ContentType should not be nil")
}

func TestSemanticQueryBasic(t *testing.T) {
	client := getDefaultClient()

	result, err := client.SemanticQuery(context.TODO(), &SemanticQueryRequest{
		Bucket:      oss.Ptr(bucket_),
		DatasetName: oss.Ptr("test-dataset"),
		Query:       oss.Ptr("雪景"),
		MaxResults:  oss.Ptr(int32(10)),
	})
	dumpErrIfNotNil(err)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)

	assert.NotNil(t, result.Files, "files should not be nil")
	assert.True(t, len(result.Files) > 0, "files should not be empty")

	f := result.Files[0]
	assert.NotNil(t, f.URI, "URI should not be nil")
	assert.NotNil(t, f.Filename, "Filename should not be nil")
	assert.NotNil(t, f.MediaType, "MediaType should not be nil")
	assert.NotNil(t, f.Size, "Size should not be nil")
}

func TestSemanticQueryWithMediaTypes(t *testing.T) {
	client := getDefaultClient()

	mediaTypesJSON := `["image"]`
	withFieldsJSON := `["Filename","Size","MediaType"]`

	result, err := client.SemanticQuery(context.TODO(), &SemanticQueryRequest{
		Bucket:      oss.Ptr(bucket_),
		DatasetName: oss.Ptr("test-dataset"),
		Query:       oss.Ptr("雪景"),
		MediaTypes:  oss.Ptr(mediaTypesJSON),
		WithFields:  oss.Ptr(withFieldsJSON),
		MaxResults:  oss.Ptr(int32(10)),
	})
	dumpErrIfNotNil(err)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)

	assert.NotNil(t, result.Files, "files should not be nil")
	assert.True(t, len(result.Files) > 0, "files should not be empty")

	// Verify all returned files are images
	for _, f := range result.Files {
		assert.Equal(t, "image", *f.MediaType, "MediaType should be image")
	}
}
