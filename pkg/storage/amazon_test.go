package storage

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type AmazonTestSuite struct {
	suite.Suite
	BrokenAmazonS3Backend   *AmazonS3Backend
	NoPrefixAmazonS3Backend *AmazonS3Backend
}

func (suite *AmazonTestSuite) SetupSuite() {
	backend := NewAmazonS3Backend("fake-bucket-cant-exist-fbce123", "", "us-east-1")
	suite.BrokenAmazonS3Backend = backend

	s3Bucket := os.Getenv("TEST_STORAGE_AMAZON_BUCKET")
	s3Region := os.Getenv("TEST_STORAGE_AMAZON_REGION")
	backend = NewAmazonS3Backend(s3Bucket, "", s3Region)
	suite.NoPrefixAmazonS3Backend = backend

	data := []byte("some object")
	path := "deleteme.txt"
	err := suite.NoPrefixAmazonS3Backend.PutObject(path, data)
	suite.Nil(err, "no error putting deleteme.txt using AmazonS3 backend")
}

func (suite *AmazonTestSuite) TearDownSuite() {
	err := suite.NoPrefixAmazonS3Backend.DeleteObject("deleteme.txt")
	suite.Nil(err, "no error deleting deleteme.txt using AmazonS3 backend")
}

func (suite *AmazonTestSuite) TestListObjects() {
	_, err := suite.BrokenAmazonS3Backend.ListObjects()
	suite.NotNil(err, "cannot list objects with bad bucket")

	_, err = suite.NoPrefixAmazonS3Backend.ListObjects()
	suite.Nil(err, "can list objects with good bucket, no prefix")
}

func (suite *AmazonTestSuite) TestGetObject() {
	_, err := suite.BrokenAmazonS3Backend.GetObject("this-file-cannot-possibly-exist.tgz")
	suite.NotNil(err, "cannot get objects with bad bucket")
}

func (suite *AmazonTestSuite) TestPutObject() {
	err := suite.BrokenAmazonS3Backend.PutObject("this-file-will-not-upload.txt", []byte{})
	suite.NotNil(err, "cannot put objects with bad bucket")
}

func TestAmazonStorageTestSuite(t *testing.T) {
	if os.Getenv("TEST_CLOUD_STORAGE") == "1" {
		suite.Run(t, new(AmazonTestSuite))
	}
}
