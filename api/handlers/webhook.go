package http

import (
	"fmt"
	"io"
	"log"
	H "net/http"
)

func StartWebhookListener(port int) {
	H.HandleFunc("/webhook", func(w H.ResponseWriter, r *H.Request) {
		body, _ := io.ReadAll(r.Body)
		fmt.Println("Webhook received:", string(body))
		w.WriteHeader(H.StatusOK)
	})
	addr := fmt.Sprintf(":%d", port)
	log.Printf("Webhook listener running on %s\n", addr)
	log.Fatal(H.ListenAndServe(addr, nil))
}
