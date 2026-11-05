package e2e

import (
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/failsafe-go/failsafe-go"
	"github.com/failsafe-go/failsafe-go/retrypolicy"
	"github.com/ockendenjo/bogale/pkg/drinks"
	"github.com/ockendenjo/bogale/pkg/user"
	"github.com/ockendenjo/handler"
	g "github.com/onsi/ginkgo/v2"
	"github.com/stretchr/testify/assert"
)

type LogHelper struct {
	t            g.FullGinkgoTInterface
	userClient   user.Client
	drinksClient drinks.Client
}

func NewLogHelper(t g.FullGinkgoTInterface) *LogHelper {
	awsConfig, err := config.LoadDefaultConfig(t.Context())
	assert.NoError(t, err)

	ddbClient := dynamodb.NewFromConfig(awsConfig)
	userClient := user.NewClient(ddbClient, handler.MustGetEnv("USER_TABLE_NAME"))
	drinksClient := drinks.NewClient(ddbClient, handler.MustGetEnv("DRINKS_TABLE_NAME"))

	return &LogHelper{
		t:            t,
		userClient:   userClient,
		drinksClient: drinksClient,
	}
}

// AwaitDisplayName polls DynamoDB until the user's display name matches the expected value.
func (h *LogHelper) AwaitDisplayName(email, expectedDisplayName string) string {
	check := func() (*user.User, error) {
		u, err := h.userClient.GetUserByEmail(h.t.Context(), email)
		if err != nil {
			return nil, err
		}
		if u == nil {
			return nil, errors.New("user not found in DynamoDB")
		}
		if u.DisplayName != expectedDisplayName {
			return nil, fmt.Errorf("display name is %q, expected %q", u.DisplayName, expectedDisplayName)
		}
		return u, nil
	}

	u, err := failsafe.With(retryPolicy[*user.User]()).Get(check)
	assert.NoError(h.t, err)
	g.GinkgoWriter.Printf("DynamoDB confirms display name is %q for %s\n", expectedDisplayName, email)
	return u.ID
}

// AwaitDrinkLogged polls DynamoDB until a drink with the given name appears for the user.
func (h *LogHelper) AwaitDrinkLogged(email, drinkName string) *drinks.Drink {
	u := h.mustGetUser(email)

	getDrink := func() (*drinks.Drink, error) {
		all, err := h.drinksClient.ListDrinksForUser(h.t.Context(), u.ID).Run()
		if err != nil {
			return nil, err
		}
		for _, d := range all {
			if d.Name == drinkName {
				return d, nil
			}
		}
		return nil, fmt.Errorf("no drink with name %q found for user %s", drinkName, u.ID)
	}

	drink, err := failsafe.With(retryPolicy[*drinks.Drink]()).Get(getDrink)
	assert.NoError(h.t, err)
	if drink != nil {
		g.GinkgoWriter.Printf("DynamoDB confirms drink %q logged for %s (pub %d, %d points)\n", drinkName, email, drink.PubID, drink.Points)
	}
	return drink
}

// AwaitPoints polls DynamoDB until the user's total_points satisfies the predicate.
// Points are updated asynchronously via the DynamoDB stream after a drink is logged.
func (h *LogHelper) AwaitPoints(email string, predicate func(points int) bool) {
	check := func() error {
		u, err := h.userClient.GetUserByEmail(h.t.Context(), email)
		if err != nil {
			return err
		}
		if u == nil {
			return errors.New("user not found in DynamoDB")
		}
		if !predicate(u.TotalPoints) {
			return fmt.Errorf("total_points %d does not satisfy predicate", u.TotalPoints)
		}
		return nil
	}

	err := failsafe.With(retryPolicy[any]()).Run(check)
	assert.NoError(h.t, err)

	u, err := h.userClient.GetUserByEmail(h.t.Context(), email)
	assert.NoError(h.t, err)
	if u != nil {
		g.GinkgoWriter.Printf("DynamoDB confirms total_points=%d for %s\n", u.TotalPoints, email)
	}
}

// AwaitVerificationCode polls DynamoDB until a verification code is present for the given email.
func (h *LogHelper) AwaitVerificationCode(email string) string {
	getCode := func() (string, error) {
		u, err := h.userClient.GetUserByEmail(h.t.Context(), email)
		if err != nil {
			return "", err
		}
		if u == nil {
			return "", errors.New("user not found in DynamoDB")
		}
		if u.Code == "" {
			return "", errors.New("verification code not yet set")
		}
		return u.Code, nil
	}

	code, err := failsafe.With(retryPolicy[string]()).Get(getCode)
	assert.NoError(h.t, err)
	g.GinkgoWriter.Printf("Retrieved verification code from DynamoDB for %s\n", email)
	return code
}

func (h *LogHelper) mustGetUser(email string) *user.User {
	u, err := h.userClient.GetUserByEmail(h.t.Context(), email)
	assert.NoError(h.t, err)
	assert.NotNil(h.t, u, "user %s not found in DynamoDB", email)
	return u
}

func (h *LogHelper) PointsEqual(expected int) func(int) bool {
	return func(points int) bool { return points == expected }
}

func (h *LogHelper) PointsGreaterThan(threshold int) func(int) bool {
	return func(points int) bool { return points > threshold }
}

func retryPolicy[T any]() retrypolicy.RetryPolicy[T] {
	return retrypolicy.NewBuilder[T]().
		HandleIf(func(_ T, err error) bool {
			return err != nil
		}).
		WithMaxRetries(10).
		WithBackoff(100*time.Millisecond, time.Second).
		Build()
}
