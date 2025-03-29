package firebase_app

import (
	"context"
	"fmt"
	"lot/config"

	"firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

func ConnectFirebaseApp() (*firebase.App, error) {
	adminSdkFilePath, err := config.Config("firebaseAdminSdkFilePath")
	if err != nil {
		return nil, fmt.Errorf("please specify the`firebaseAdminSdkFilePath`")
	}
	opt := option.WithCredentialsFile(adminSdkFilePath)
	return firebase.NewApp(context.Background(), nil, opt)
}
