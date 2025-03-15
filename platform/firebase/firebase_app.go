package firebase_app

import (
	"context"
	"lot/config"

	"firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

func ConnectFirebaseApp() (*firebase.App, error) {
	opt := option.WithCredentialsFile(config.Config("firebaseAdminSdkFilePath"))
	return firebase.NewApp(context.Background(), nil, opt)
}
