package handler

import (
	"context"
	"fmt"
	"testing"

	"github.com/rsocket/rsocket-go"
	"github.com/rsocket/rsocket-go/payload"
	"github.com/rsocket/rsocket-go/rx/flux"
)

func init() {
	go StartServer(nil)
}

func connect(t *testing.T) rsocket.Client {
	// Connect to server
	client, err := rsocket.Connect().
		Transport(rsocket.TCPClient().
			SetHostAndPort("localhost", 8000).
			Build()).
		Start(context.Background())
	if err != nil {
		panic("Ошибка при подключении к серверу")
	}
	return client
}

func TestRequestResponse(t *testing.T) {
	cli := connect(t)
	defer cli.Close()
	requestPayload := payload.NewString("https://example.com", "")

	response := cli.RequestResponse(requestPayload)

	response.
		DoOnSuccess(func(input payload.Payload) error {
			return nil
		}).
		DoOnError(func(e error) {
			fmt.Printf("Возникла ошибка (%s)", e)
		}).
		Subscribe(context.Background())
}

func TestRequestStream(t *testing.T) {
	cli := connect(t)
	defer cli.Close()
	testPayload := payload.NewString("sample_request", "")

	response := cli.RequestStream(testPayload)
	var receivedPayloads []payload.Payload

	response.
		DoOnNext(func(input payload.Payload) error {
			// Process each payload received in the stream
			receivedPayloads = append(receivedPayloads, input)
			return nil
		}).
		Subscribe(context.Background())
}

func createMockFlux(mockData []string) flux.Flux {
	return flux.Create(func(ctx context.Context, sink flux.Sink) {
		for _, data := range mockData {
			sink.Next(payload.NewString(data, ""))
		}
		sink.Complete()
	})
}

func TestRequestChannel(t *testing.T) {
	cli := connect(t)
	defer cli.Close()
	mockPayloads := []string{"url1", "url2", "url3"}
	mockFlux := createMockFlux(mockPayloads)

	response := cli.RequestChannel(mockFlux)

	var receivedPayloads []payload.Payload

	response.
		DoOnNext(func(input payload.Payload) error {
			// Process each payload received in the flux
			receivedPayloads = append(receivedPayloads, input)
			return nil
		}).
		Subscribe(context.Background())
}

func TestFireAndForget(t *testing.T) {
	cli := connect(t)
	defer cli.Close()

	mockOriginalUrl := "https://example.com"

	payloadToDelete := payload.NewString(mockOriginalUrl, "")

	cli.FireAndForget(payloadToDelete)

}
