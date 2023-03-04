package helpers

import (
	"be-no-stunting-v2/configs"
	"context"
	"fmt"
	"log"

	"firebase.google.com/go/messaging"
)

func SendToToken(registrationToken []string, title string, body string) {
	app, _, _ := configs.SetupFirebase()
	ctx := context.Background()
	client, err := app.Messaging(ctx)
	if err != nil {
		log.Fatalf("error getting Messaging client: %v\n", err)
	}

	response, err := client.SendMulticast(context.Background(), &messaging.MulticastMessage{
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Tokens: registrationToken, // it's an array of device tokens
	})

	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Successfully sent message:", response)
}
