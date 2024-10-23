package minio

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	transport "github.com/aws/smithy-go/endpoints"
	"log"
	"mime/multipart"
	"net/url"
	"os"
	"time"
)

type Resolver struct {
	URL *url.URL
}

type minioCreds struct {
	key, secret, endpoint, bucket string
}

type Database struct {
	client *s3.Client
	creds  minioCreds
}

func GetConnection() (*Database, error) {
	db := &Database{}
	var err error

	db.creds.key = os.Getenv("AWS_ACCESS_KEY")
	db.creds.secret = os.Getenv("AWS_SECRET_KEY")
	db.creds.endpoint = os.Getenv("AWS_ENDPOINT")
	db.creds.bucket = os.Getenv("AWS_BUCKET")

	if db.creds.endpoint == "" || db.creds.key == "" || db.creds.secret == "" || db.creds.bucket == "" {
		return nil, errors.New("Missing environment variables for Minio connection")
	}

	// need to set up creds first
	db.client, err = db.createSession()

	if err != nil {
		return nil, err
	}
	return db, nil
}

func (db *Database) GetPresignedURL(key string, duration time.Duration) (string, error) {
	presignedClient := s3.NewPresignClient(db.client)

	req, err := presignedClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(db.creds.bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(duration))

	if err != nil {
		log.Println(err)
		return "", errors.New("Could not generate URL")
	}
	return req.URL, nil
}

func (db *Database) PutObject(file multipart.File, fileHeader *multipart.FileHeader) error {
	fileSize := fileHeader.Size
	fileName := fileHeader.Filename

	userProfilePrefix := "profile_pics/"

	_, err := db.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:        aws.String(db.creds.bucket),
		Key:           aws.String(userProfilePrefix + fileName),
		Body:          file,
		ContentLength: aws.Int64(fileSize),
	})

	if err != nil {
		return err
	}
	return nil
}

func (r *Resolver) ResolveEndpoint(_ context.Context, params s3.EndpointParameters) (transport.Endpoint, error) {
	u := *r.URL

	u.Path += "/" + *params.Bucket
	return transport.Endpoint{URI: u}, nil
}

func (db *Database) createSession() (*s3.Client, error) {
	if db.creds.key == "" || db.creds.secret == "" {
		log.Fatal("AWS_ACCESS_KEY and AWS_SECRET_KEY are required")
	}
	creds := credentials.NewStaticCredentialsProvider(db.creds.key, db.creds.secret, "")

	endpointURL, err := url.Parse(db.creds.endpoint)

	if err != nil {
		return nil, err
	}

	client := s3.New(s3.Options{
		EndpointResolverV2: &Resolver{URL: endpointURL},
		Credentials:        creds,
	})

	return client, nil
}
