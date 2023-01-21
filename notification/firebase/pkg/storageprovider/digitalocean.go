package storageprovider

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	bylogs "github.com/cuemby/by-go-utils/pkg/bylogger"
	"github.com/igorariza/golang-api-gozero-grpc/notification/firebase/pkg/file"
)

type DigitalOceanSpaceProvider struct {
	endpoint   string
	key        string
	secret     string
	bucketName string
	region     string
	urlBase    string

	client *s3.S3
}

const (
	digitalOceanPermissionPublicMetadataKey = "x-amz-acl"
)

var (
	digitalOceanPermissionPublicMetadataValue = "public-read"
)

func (d *DigitalOceanSpaceProvider) ProviderName() string {
	return string(DigitalOceanSpace)
}
func (d *DigitalOceanSpaceProvider) Connect() error {
	if d == nil {
		d = &DigitalOceanSpaceProvider{}
	}
	d.endpoint = os.Getenv("DO_SPACES_ENDPOINT")
	d.key = os.Getenv("DO_SPACES_KEY")
	d.secret = os.Getenv("DO_SPACES_SECRET")
	d.bucketName = os.Getenv("DO_SPACES_NAME")
	d.region = os.Getenv("DO_SPACES_REGION")
	d.urlBase = fmt.Sprintf("https://%s.%s/", d.bucketName, d.endpoint)

	s3Config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(d.key, d.secret, ""),
		Endpoint:    aws.String(d.endpoint),
		Region:      aws.String(d.region),
	}

	newSession, err := session.NewSession(s3Config)

	if err != nil {
		log.Println("Error: ", err)
		return errors.New("DigitalOceanSpaceProvider - error creating session")
	}

	d.client = s3.New(newSession)
	if d.client == nil {
		return errors.New("DigitalOceanSpaceProvider - error creating client instance")
	}
	return nil
}

func (d *DigitalOceanSpaceProvider) UploadFile(doFile *file.GenericFile) (err error) {
	fileBytes := bytes.NewReader(doFile.Data)

	//save file to s3
	object := s3.PutObjectInput{
		Bucket: aws.String(d.bucketName),
		Key:    aws.String(doFile.FileName),
		Body:   fileBytes,
	}

	// add public permission when private is false
	if !doFile.Private {
		doFile.Metadata[digitalOceanPermissionPublicMetadataKey] = digitalOceanPermissionPublicMetadataValue
		object.ACL = aws.String(digitalOceanPermissionPublicMetadataValue)
	}

	metadata := map[string]*string{}
	for key, value := range doFile.Metadata {
		pointerValue := value
		metadata[key] = &pointerValue
	}
	// attach metadata
	object.Metadata = metadata

	_, err = d.client.PutObject(&object)
	if err != nil {
		bylogs.LogErr(err)
		return err
	}

	//Save file to mongo
	//GET file from s3
	doFile.URL = d.urlBase + doFile.FileName
	return nil
}

func (d *DigitalOceanSpaceProvider) GetFile(filename string) (*file.GenericFile, error) {
	req, _ := d.client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(d.bucketName),
		Key:    aws.String(filename),
	})

	url, err := req.Presign(5 * time.Minute)
	if err != nil {
		bylogs.LogErr(err)
		return nil, err
	}

	//get metadata
	metadata, _ := d.client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(d.bucketName),
		Key:    aws.String(filename),
	})

	doFile := &file.GenericFile{
		FileName: filename,
		URL:      url,
		Metadata: map[string]string{},
		Private:  false,
		Data:     nil,
	}

	for key := range metadata.Metadata {
		if key == digitalOceanPermissionPublicMetadataKey {
			doFile.Private = true
		} else {
			doFile.Metadata[key] = *metadata.Metadata[key]
		}

	}
	doFile.URL = fmt.Sprint(d.urlBase, doFile.FileName)
	return doFile, nil
}

func (d *DigitalOceanSpaceProvider) ListFiles() ([]*file.GenericFile, error) {
	var fileList []*file.GenericFile

	input := &s3.ListObjectsInput{
		Bucket: aws.String(d.bucketName),
		//MaxKeys: aws.Int64(10),
	}

	objects, err := d.client.ListObjects(input)
	if err != nil {
		bylogs.LogErr(err)
		return fileList, err
	}

	for _, obj := range objects.Contents {
		//fmt.Println(aws.StringValue(obj.Key))
		fileList = append(fileList, d.ParseS3toDoFile(obj))
	}

	return fileList, nil
}

func (d *DigitalOceanSpaceProvider) ListFilesPaginated(page int64, size int64) ([]*file.GenericFile, error) {
	var fileList []*file.GenericFile

	input := &s3.ListObjectsInput{
		Bucket:  aws.String(d.bucketName),
		MaxKeys: aws.Int64(size),
		Marker:  aws.String(fmt.Sprintf("%d", page)),
	}

	i := 0
	err := d.client.ListObjectsPages(input, func(p *s3.ListObjectsOutput, last bool) (shouldContinue bool) {

		i++

		if p.IsTruncated != nil && *p.IsTruncated {
			for _, obj := range p.Contents {
				fileList = append(fileList, d.ParseS3toDoFile(obj))
			}
			return true

		} else {
			return false
		}
	})

	if err != nil {
		bylogs.LogErr(err)
		return nil, err
	}

	return fileList, nil
}

func (d *DigitalOceanSpaceProvider) DeleteFile(filename string) (bool, error) {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(d.bucketName),
		Key:    aws.String(filename),
	}

	resultDelete, err := d.client.DeleteObject(input)
	if err != nil {
		bylogs.LogErr(err)
		return false, err
	}
	return aws.BoolValue(resultDelete.DeleteMarker), nil

}

func (d *DigitalOceanSpaceProvider) ParseS3toDoFile(obj *s3.Object) *file.GenericFile {
	doFile := &file.GenericFile{
		FileName: aws.StringValue(obj.Key),
		URL:      d.urlBase + aws.StringValue(obj.Key),
		Metadata: map[string]string{},
	}
	return doFile
}
