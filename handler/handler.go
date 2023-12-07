package handler

import (
	"Service/database"
	"Service/shortUrl"
	"context"
	"fmt"
	"github.com/rsocket/rsocket-go"
	"github.com/rsocket/rsocket-go/payload"
	"github.com/rsocket/rsocket-go/rx/flux"
	"github.com/rsocket/rsocket-go/rx/mono"
	"log"
	"sync"
)

func StartServer(storage database.URLStorage) {
	//var err error
	wg := &sync.WaitGroup{}
	wg.Add(1)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		defer wg.Done()
		err := rsocket.Receive().
			OnStart(func() {
				log.Println("Server Started")
			}).
			Acceptor(func(_ context.Context, _ payload.SetupPayload, _ rsocket.CloseableRSocket) (rsocket.RSocket, error) {
				return rsocket.NewAbstractSocket(
					rsocket.RequestResponse(func(c payload.Payload) mono.Mono {
						originalUrl := c.DataUTF8()
						if originalUrl == "" {
							return mono.Error(fmt.Errorf("invalid url format"))
						}

						shortUrl, err := shortUrl.GenerateShortUrl(originalUrl, storage)
						if err != nil {
							return mono.Error(err)
						}

						return mono.Just(payload.NewString(shortUrl, ""))
					}),

					rsocket.RequestStream(func(c payload.Payload) flux.Flux {
						allUrl, err := storage.GetAllUrls()
						if err != nil {
							return flux.Error(fmt.Errorf("with getting links"))
						}

						return flux.Create(func(_ context.Context, s flux.Sink) {
							for _, Url := range allUrl {
								s.Next(payload.NewString(Url, ""))
							}
							s.Complete()
						})
					}),

					rsocket.RequestChannel(func(c flux.Flux) flux.Flux {
						shorts := make(chan string)
						originals := make(chan string)

						c.DoOnComplete(func() {
							close(shorts)
						}).DoOnNext(func(msg payload.Payload) error {
							short := msg.DataUTF8()
							shorts <- short
							return nil
						}).Subscribe(context.Background())

						go func() {
							for short := range shorts {
								original, _ := storage.GetOriginalUrl(short)
								originals <- original
							}
							close(originals)
						}()

						return flux.Create(func(_ context.Context, s flux.Sink) {
							for original := range originals {
								s.Next(payload.NewString(original, ""))
							}
							s.Complete()
						})
					}),

					rsocket.FireAndForget(func(c payload.Payload) {
						originalUrl := c.DataUTF8()
						if originalUrl != "" {
							storage.DeleteOriginalUrl(originalUrl)
						}
					}),
				), nil
			}).
			Transport(rsocket.TCPServer().SetAddr(":8000").Build()).
			Serve(ctx)

		if err != nil {
			log.Fatalln(err)
		}
	}()

	wg.Wait()
	cancel()
}
