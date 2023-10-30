package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/julienschmidt/httprouter"
)

func newRouter() *httprouter.Router {
	mux := httprouter.New()

	mux.GET("/youtube/channel/stats", getChannelStats())
	mux.GET("/savelogs", saveLogs())

	return mux
}

func getChannelStats() httprouter.Handle {
	customer := os.Getenv("CUSTOMER")
	environment := os.Getenv("ENVIRONMENT")

	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.Write([]byte("response!" + customer + "/" + environment))
	}
}

func saveLogs() httprouter.Handle {

	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		log.Println("Request in save logs")

		result, err := getStorageClient("cyucrastorage")
		if err != nil {
			w.Write([]byte("error " + err.Error()))
		}
		w.Write([]byte("result " + result))

	}
}

func getStorageClient(accountName string) (string, error) {
	// login with azure identity
	// ref: https://learn.microsoft.com/en-us/azure/developer/go/azure-sdk-authentication?tabs=bash#workload-identity
	log.Println("getStorageClient: Precredentials")
	credentials, err := azidentity.NewDefaultAzureCredential(nil)
	log.Println("getStorageClient: Poscredentials")

	if err != nil {
		return "", errors.New(err.Error())
	}
	url := fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)

	log.Println("getStorageClient: PreNewClient")
	client, err := azblob.NewClient(url, credentials, nil)
	log.Println("getStorageClient: PosNewClient")
	if err != nil {
		return "", errors.New(err.Error())
	}

	blobName := "sample-blob"
	container := "images"
	data := []byte("\nHello, world! This is a blob.\n")

	log.Println("getStorageClient: PreUpload")
	_, err = client.UploadBuffer(context.TODO(), container, blobName, data, &azblob.UploadBufferOptions{})
	log.Println("getStorageClient: PosUpload")

	if err != nil {
		return "UploadBuffer: ", errors.New(err.Error())
	}

	log.Println("getStorageClient: PreNewContainer")
	containerClient := client.ServiceClient().NewContainerClient(container)
	blobClient := containerClient.NewBlobClient(blobName)

	log.Println("getStorageClient: PosNewContainer")
	return blobClient.URL(), nil

}

func printConfiguration() {
	log.Println("AZURE_CLIENT_ID: " + os.Getenv("AZURE_CLIENT_ID"))
	log.Println("AZURE_TENANT_ID: " + os.Getenv("AZURE_TENANT_ID"))
}

func main() {
	printConfiguration()

	srv := &http.Server{
		Addr:    ":10101",
		Handler: newRouter(),
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		signal.Notify(sigint, syscall.SIGTERM)
		<-sigint

		log.Println("service interrupt received")

		log.Println("http server shutting down")
		time.Sleep(5 * time.Second)

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("http server shutdown error: %v", err)
		}

		log.Println("shutdown complete")

		close(idleConnsClosed)

	}()

	log.Printf("Starting server on port 10101")
	if err := srv.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("fatal http server failed to start: %v", err)
		}
	}

	<-idleConnsClosed
	log.Println("Service Stop")

}
