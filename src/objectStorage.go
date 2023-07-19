package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"log"
    "io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func initMinio() *minio.Client {
	// MinIO connection
	var endpoint = os.Getenv("MINIO_ENDPOINT")
	var accessKeyID = os.Getenv("MINIO_ACCESS_KEY_ID")
	var secretAccessKey = os.Getenv("MINIO_SECRET_ACCESS_KEY")
	var bucketName = os.Getenv("MINIO_BUCKET")
	// var bucketRegion = os.Getenv("MINIO_REGION")
	useSSL, _ := strconv.ParseBool(os.Getenv("MINIO_SSL_ENABLED"))

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})

	if err != nil {
		log.Fatalln(err)
	}

	found, err := minioClient.BucketExists(context.Background(), bucketName)
	if err != nil {
		log.Printf(err.Error())
	    log.Printf("Bucket %+v not accessible, try to create it\n", bucketName) // Region: bucketRegion,
	    err = minioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{ObjectLocking: true})
		if err != nil {
			log.Printf("error during bucket creation: %+v\n", err.Error())
		    log.Fatalln(err.Error())
		}
		log.Printf("Successfully created mybucket %+v.\n", bucketName)
	}
	if found {
	    log.Printf("Bucket %+v found\n", bucketName)
	}


	log.Printf("Connected to object storage %+v\n", endpoint)

	return minioClient
}

func putObject(minioClient *minio.Client, buffer io.Reader, object string) {
	var bucketName = os.Getenv("MINIO_BUCKET")

	// TODO: optimize it
	// bufferSize = int64(buffer.Len())
	// BufferedSize := bufferedWriter.Buffered()
	// fmt.Println(buffer.String())

	// putToMinIO(buffer, minioClient, bucketName, object)

	// bufferStat, err := buffer.Size()
	// if err != nil {
	//     fmt.Println(err)
	//     return
	// }

	// buf := new(bytes.Buffer)
    // buf.ReadFrom(buffer)
    // bufferSize := int64(buf.Len())

	// TODO: add bufferSize for memory optimization
	_, err := minioClient.PutObject(context.Background(), bucketName, object, buffer, -1, minio.PutObjectOptions{ContentType:"application/octet-stream"})
	if err != nil {
	    fmt.Println(err)
	    return
	}
}