package e2e

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/mxschmitt/playwright-go"
	"github.com/ockendenjo/bogale/pkg/user"
	"github.com/ockendenjo/bogale/pkg/userpool"
	"github.com/ockendenjo/handler"
	g "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
)

func Test_E2E(t *testing.T) {
	gomega.RegisterFailHandler(g.Fail)
	g.RunSpecs(t, "E2E")
}

var _ = g.SynchronizedBeforeSuite(func() {
	t := g.GinkgoT()
	headless := IsHeadless()
	err := playwright.Install(&playwright.RunOptions{
		Browsers:         []string{"chromium"},
		OnlyInstallShell: headless,
		Stdout:           g.GinkgoWriter,
		Stderr:           g.GinkgoWriter,
	})
	assert.NoError(t, err)

	awsConfig, err := config.LoadDefaultConfig(t.Context())
	assert.NoError(t, err)

	userPoolID := handler.MustGetEnv("USER_POOL_ID")
	cognitoClient := cognitoidentityprovider.NewFromConfig(awsConfig)
	poolClient := userpool.NewClient(cognitoClient, userPoolID)

	testPassword = userpool.RandomPassword(16)

	err = poolClient.AddUser(t.Context(), testEmail, testPassword)
	assert.NoError(t, err)
	g.GinkgoWriter.Printf("Created Cognito user %s\n", testEmail)

	// Delete DynamoDB user record to ensure clean state
	tableName := handler.MustGetEnv("USER_TABLE_NAME")
	ddbClient := dynamodb.NewFromConfig(awsConfig)
	userClient := user.NewClient(ddbClient, tableName)
	err = userClient.DeleteByEmail(t.Context(), testEmail)
	assert.NoError(t, err)
	g.GinkgoWriter.Printf("Cleaned up DynamoDB user %s\n", testEmail)
}, func() {
	// Deliberately empty
})

var _ = g.SynchronizedAfterSuite(func() {}, func() {
	if !IsHeadless() {
		time.Sleep(5 * time.Second)
	}
})
