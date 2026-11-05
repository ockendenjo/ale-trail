package e2e

import (
	"github.com/mxschmitt/playwright-go"
	g "github.com/onsi/ginkgo/v2"
	"github.com/stretchr/testify/assert"
)

const testDisplayName = "E2E Test User"
const (
	friendName1 = "Imaginary Friend"
	friendName2 = "Jess"
	friendName3 = "Sam"
)

var _ = TestCase("001", "Register, set display name, and log a drink", func() {
	var (
		t      g.FullGinkgoTInterface
		page   playwright.Page
		lh     *LogHelper
		userID string
	)

	g.BeforeAll(func() {
		t = g.GinkgoT()
		page = LoadHomePage(t)
		lh = NewLogHelper(t)
	})

	g.It("should upload data into S3", func() {

	})

	g.It("should load the home page", func() {
		text, err := page.Locator("h1").TextContent()
		assert.NoError(t, err)
		assert.Contains(t, text, "BOGALE")
	})

	g.Context("Sign in", func() {
		g.It("should sign in via Cognito hosted UI", func() {
			SignIn(t, page, testEmail, testPassword)
		})

		g.It("should navigate to the profile page on first sign-in", func() {
			AwaitURL(t, page, "profile")
			g.GinkgoWriter.Println("Redirected to profile page after first sign-in")
		})
	})

	g.Context("Set display name", func() {
		g.It("should navigate to profile page", func() {
			NavigateTo(t, page, "profile")
		})

		g.It("should fill in the display name", func() {
			displayNameInput := page.Locator(`input[formcontrolname="displayName"]`)
			err := playwright.NewPlaywrightAssertions(10_000).Locator(displayNameInput).ToBeVisible()
			assert.NoError(t, err)
			FillFormControlInput(t, page, "displayName", testDisplayName)
		})

		g.It("should submit the display name", func() {
			SubmitForm(t, page)
		})

		g.It("should navigate back to the home page", func() {
			AwaitURL(t, page, "")
			g.GinkgoWriter.Printf("Display name set to %q\n", testDisplayName)
		})

		g.It("should persist the display name in DynamoDB", func() {
			userID = lh.AwaitDisplayName(testEmail, testDisplayName)
		})
	})

	g.Context("Log a drink", func() {
		g.It("should navigate to log drink page", func() {
			NavigateTo(t, page, "log")
			ExpectH1(t, page, "Log a drink")
		})

		g.It("should submit the when form", func() {
			// Default selection is "Now" - click Next to proceed
			nextBtn := page.Locator(`button[type="submit"]`).First()
			err := playwright.NewPlaywrightAssertions(5_000).Locator(nextBtn).ToBeEnabled()
			assert.NoError(t, err)
			err = nextBtn.Click()
			assert.NoError(t, err)
		})

		g.It("should wait for pub selections to load", func() {
			venueInput := page.Locator(`input#venue`)
			err := playwright.NewPlaywrightAssertions(10_000).Locator(venueInput).ToBeVisible()
			assert.NoError(t, err)
			g.GinkgoWriter.Println("Drink details form loaded")
		})

		g.It("should select a pub from the autocomplete", func() {
			// Click the venue input to open the autocomplete panel
			venueInput := page.Locator(`input#venue`).First()
			err := venueInput.Click()
			assert.NoError(t, err)

			// Wait for at least one enabled pub option to appear
			enabledOptions := page.Locator(`mat-option:not([aria-disabled="true"])`)
			err = playwright.NewPlaywrightAssertions(5_000).Locator(enabledOptions.First()).ToBeVisible()
			assert.NoError(t, err)

			// Select the first enabled option
			option := enabledOptions.First()
			text, err := option.TextContent()
			assert.NoError(t, err)
			err = option.Click()
			assert.NoError(t, err)
			g.GinkgoWriter.Printf("Selected pub: %s\n", text)
		})

		g.It("should fill in the drink details", func() {
			FillAutocompleteInput(t, page, "name", "Pentland IP")
			FillModelInput(t, page, "drunkWith", friendName1)
			FillTextarea(t, page, "notes", "Extra notes")
		})

		g.It("should submit the drink", func() {
			submitBtn := page.Locator(`button[type="submit"]`).First()
			err := playwright.NewPlaywrightAssertions(5_000).Locator(submitBtn).ToBeEnabled()
			assert.NoError(t, err)
			err = submitBtn.Click()
			assert.NoError(t, err)
		})

		g.It("should show drink logged successfully", func() {
			successMsg := page.GetByText("Drink logged successfully")
			err := playwright.NewPlaywrightAssertions(10_000).Locator(successMsg).ToBeVisible()
			assert.NoError(t, err)
			g.GinkgoWriter.Println("Drink logged successfully")
		})

		g.It("should persist the drink in DynamoDB", func() {
			loggedDrink := lh.AwaitDrinkLogged(testEmail, "Pentland IPA")
			exp := drinks.Drink{
				UserID:      userID,
				Name:        "Pentland IPA",
				Brewery:     "Stewart Brewing",
				UntappdID:   new(8306),
				SelectedPub: true,
				EndOfWord:   false,
				NotInWord:   false,
				With:        friendName1,
				Version:     1,
				Notes:       new("Extra notes"),
			}
			assertDrinksEqual(t, exp, *loggedDrink)
		})

		g.It("should update the user's points in DynamoDB", func() {
			lh.AwaitPoints(testEmail, lh.PointsEqual(20))
		})
	})

	g.Context("Log a second drink", func() {
		g.It("should click 'submit another' button", func() {
			err := page.Locator("button#submit-another").Click()
			assert.NoError(t, err)
		})

		g.It("should submit the when form", func() {
			nextBtn := page.Locator(`button[type="submit"]`).First()
			err := playwright.NewPlaywrightAssertions(5_000).Locator(nextBtn).ToBeEnabled()
			assert.NoError(t, err)
			err = nextBtn.Click()
			assert.NoError(t, err)
		})

		g.It("should wait for pub selections to load", func() {
			venueInput := page.Locator(`input#venue`)
			err := playwright.NewPlaywrightAssertions(10_000).Locator(venueInput).ToBeVisible()
			assert.NoError(t, err)
		})

		g.It("should select a pub from the autocomplete", func() {
			venueInput := page.Locator(`input#venue`).First()
			err := venueInput.Click()
			assert.NoError(t, err)

			enabledOptions := page.Locator(`mat-option:not([aria-disabled="true"])`)
			err = playwright.NewPlaywrightAssertions(5_000).Locator(enabledOptions.First()).ToBeVisible()
			assert.NoError(t, err)

			option := enabledOptions.First()
			text, err := option.TextContent()
			assert.NoError(t, err)
			err = option.Click()
			assert.NoError(t, err)
			g.GinkgoWriter.Printf("Selected pub: %s\n", text)
		})

		g.It("should fill in the drink details", func() {
			FillAutocompleteInput(t, page, "name", "avalan")
			FillModelInput(t, page, "drunkWith", friendName2)
		})

		g.It("should submit the drink", func() {
			submitBtn := page.Locator(`button[type="submit"]`).First()
			err := playwright.NewPlaywrightAssertions(5_000).Locator(submitBtn).ToBeEnabled()
			assert.NoError(t, err)
			err = submitBtn.Click()
			assert.NoError(t, err)
		})

		g.It("should show drink logged successfully", func() {
			successMsg := page.GetByText("Drink logged successfully")
			err := playwright.NewPlaywrightAssertions(10_000).Locator(successMsg).ToBeVisible()
			assert.NoError(t, err)
			g.GinkgoWriter.Println("Second drink logged successfully")
		})

		g.It("should persist the second drink in DynamoDB", func() {
			loggedDrink := lh.AwaitDrinkLogged(testEmail, "Avalanche")
			exp := drinks.Drink{
				UserID:      userID,
				Name:        "Avalanche",
				Brewery:     "Fyne Ales",
				UntappdID:   new(10922),
				SelectedPub: true,
				EndOfWord:   false,
				NotInWord:   false,
				With:        friendName2,
				Version:     1,
				Notes:       nil,
			}
			assertDrinksEqual(t, exp, *loggedDrink)
		})

		g.It("should update the user's points to 40 in DynamoDB", func() {
			lh.AwaitPoints(testEmail, lh.PointsEqual(40))
		})
	})

	g.Context("Log a third drink", func() {
		g.It("should navigate to log drink page", func() {
			NavigateTo(t, page, "log")
		})

		g.It("should submit the when form", func() {
			nextBtn := page.Locator(`button[type="submit"]`).First()
			err := playwright.NewPlaywrightAssertions(5_000).Locator(nextBtn).ToBeEnabled()
			assert.NoError(t, err)
			err = nextBtn.Click()
			assert.NoError(t, err)
		})

		g.It("should wait for pub selections to load", func() {
			venueInput := page.Locator(`input#venue`)
			err := playwright.NewPlaywrightAssertions(10_000).Locator(venueInput).ToBeVisible()
			assert.NoError(t, err)
		})

		g.It("should select a pub from the autocomplete", func() {
			venueInput := page.Locator(`input#venue`).First()
			err := venueInput.Click()
			assert.NoError(t, err)

			enabledOptions := page.Locator(`mat-option:not([aria-disabled="true"])`)
			err = playwright.NewPlaywrightAssertions(5_000).Locator(enabledOptions.First()).ToBeVisible()
			assert.NoError(t, err)

			option := enabledOptions.First()
			text, err := option.TextContent()
			assert.NoError(t, err)
			err = option.Click()
			assert.NoError(t, err)
			g.GinkgoWriter.Printf("Selected pub: %s\n", text)
		})

		g.It("should fill in the drink details", func() {
			FillAutocompleteInput(t, page, "name", "schiehal")
			FillModelInput(t, page, "drunkWith", friendName3)
		})

		g.It("should submit the drink", func() {
			submitBtn := page.Locator(`button[type="submit"]`).First()
			err := playwright.NewPlaywrightAssertions(5_000).Locator(submitBtn).ToBeEnabled()
			assert.NoError(t, err)
			err = submitBtn.Click()
			assert.NoError(t, err)
		})

		g.It("should show drink logged successfully", func() {
			successMsg := page.GetByText("Drink logged successfully")
			err := playwright.NewPlaywrightAssertions(10_000).Locator(successMsg).ToBeVisible()
			assert.NoError(t, err)
			g.GinkgoWriter.Println("Third drink logged successfully")
		})

		g.It("should persist the third drink in DynamoDB", func() {
			loggedDrink := lh.AwaitDrinkLogged(testEmail, "Schiehallion")
			exp := drinks.Drink{
				UserID:      userID,
				Name:        "Schiehallion",
				Brewery:     "Harviestoun",
				UntappdID:   new(8899),
				SelectedPub: true,
				EndOfWord:   false,
				NotInWord:   false,
				With:        friendName3,
				Version:     1,
				Notes:       nil,
			}
			assertDrinksEqual(t, exp, *loggedDrink)
		})
	})

	g.Context("My drinks page", func() {
		g.It("should navigate to my drinks page", func() {
			NavigateTo(t, page, "myDrinks")
		})

		g.It("should show all three drinks", func() {
			firstDrink := page.GetByText("Pentland IPA")
			err := playwright.NewPlaywrightAssertions(10_000).Locator(firstDrink.First()).ToBeVisible()
			assert.NoError(t, err)

			secondDrink := page.GetByText("Avalanche")
			err = playwright.NewPlaywrightAssertions(10_000).Locator(secondDrink.First()).ToBeVisible()
			assert.NoError(t, err)

			thirdDrink := page.GetByText("Schiehallion")
			err = playwright.NewPlaywrightAssertions(10_000).Locator(thirdDrink.First()).ToBeVisible()
			assert.NoError(t, err)

			g.GinkgoWriter.Println("All three drinks are visible on the my drinks page")
		})
	})

	g.Context("Complete word", func() {
		g.It("should open the menu on the last drink", func() {
			lastDrinkItem := page.Locator("app-drink-item").Last()
			menuBtn := lastDrinkItem.Locator("button[mat-icon-button]")
			err := playwright.NewPlaywrightAssertions(5_000).Locator(menuBtn).ToBeVisible()
			assert.NoError(t, err)
			err = menuBtn.Click()
			assert.NoError(t, err)
		})

		g.It("should click Set as end of word", func() {
			setEndBtn := page.GetByText("Set as end of word")
			err := playwright.NewPlaywrightAssertions(5_000).Locator(setEndBtn).ToBeVisible()
			assert.NoError(t, err)
			err = setEndBtn.Click()
			assert.NoError(t, err)
			g.GinkgoWriter.Println("Marked word as complete")
		})

		g.It("should show the word in the previous words section", func() {
			previousWords := page.GetByText("Previous words")
			err := playwright.NewPlaywrightAssertions(10_000).Locator(previousWords).ToBeVisible()
			assert.NoError(t, err)
		})

		g.It("should update the user's points to more than 60 in DynamoDB", func() {
			lh.AwaitPoints(testEmail, lh.PointsGreaterThan(60))
		})
	})
})
