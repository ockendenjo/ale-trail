package e2e

import (
	"fmt"
	"net/http"
	"os"

	"github.com/mxschmitt/playwright-go"
	"github.com/ockendenjo/bogale/pkg/drinks"
	g "github.com/onsi/ginkgo/v2"
	"github.com/stretchr/testify/assert"
)

const testEmail = "e2e-test@ockenden.io"

var testPassword string

func TestCase(id, title string, args ...any) bool {
	fullTitle := fmt.Sprintf("%s - %s", id, title)
	return g.Describe(fullTitle, g.Label(id), g.Ordered, args)
}

func IsHeadless() bool {
	return os.Getenv("SHOW_BROWSER") != "1"
}

func baseURL() string {
	if url := os.Getenv("BASE_URL"); url != "" {
		return url
	}
	return "https://dev.bogale.ockenden.io"
}

func LoadHomePage(t g.FullGinkgoTInterface) playwright.Page {
	pw, err := playwright.Run(&playwright.RunOptions{Stdout: g.GinkgoWriter, Stderr: g.GinkgoWriter})
	assert.NoError(t, err)

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: new(IsHeadless())})
	assert.NoError(t, err)

	page, err := browser.NewPage(playwright.BrowserNewPageOptions{BaseURL: new(baseURL())})
	assert.NoError(t, err)

	res, err := page.Goto("/")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.Status(), "loading home page returned HTTP %d", res.Status())
	return page
}

func NavigateTo(t g.FullGinkgoTInterface, page playwright.Page, path string) {
	url := fmt.Sprintf("%s/%s", baseURL(), path)
	res, err := page.Goto(url)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.Status(), "navigating to %s returned HTTP %d", path, res.Status())
	g.GinkgoWriter.Printf("Navigated to %s\n", page.URL())
}

func ExpectH1(t g.FullGinkgoTInterface, page playwright.Page, expected string) {
	h1, err := page.Locator("h1").TextContent()
	assert.NoError(t, err)
	assert.Contains(t, h1, expected)
}

func AwaitURL(t g.FullGinkgoTInterface, page playwright.Page, path string) {
	url := fmt.Sprintf("%s/%s", baseURL(), path)
	err := page.WaitForURL(url)
	assert.NoError(t, err)
	g.GinkgoWriter.Printf("Page URL is now: %s\n", url)
}

func SignIn(t g.FullGinkgoTInterface, page playwright.Page, email, password string) {
	err := page.Locator("button#sign-in").Click()
	assert.NoError(t, err)
	err = page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{State: playwright.LoadStateLoad, Timeout: new(15_000.0)})
	assert.NoError(t, err)

	loc := page.Locator(`input[name="username"]`)
	err = playwright.NewPlaywrightAssertions(10_000).Locator(loc).ToBeEnabled()
	assert.NoError(t, err)
	err = loc.Fill(email)
	assert.NoError(t, err)

	err = page.Locator(`input[name="password"]`).Fill(password)
	assert.NoError(t, err)

	err = page.Locator(`button[type="submit"]`).Click()
	assert.NoError(t, err)

	g.GinkgoWriter.Println("Submitted Cognito sign-in form")
}

func FillFormControlInput(t g.FullGinkgoTInterface, page playwright.Page, formControlName, value string) {
	locator := page.Locator(fmt.Sprintf(`input[formcontrolname="%s"]`, formControlName)).First()
	err := playwright.NewPlaywrightAssertions(5_000).Locator(locator).ToBeEnabled()
	assert.NoError(t, err, "form field %q is not enabled", formControlName)
	err = locator.Fill(value)
	assert.NoError(t, err)
	g.GinkgoWriter.Printf("Filled formControlName=%q\n", formControlName)
}

func FillModelInput(t g.FullGinkgoTInterface, page playwright.Page, name, value string) {
	locator := page.Locator(fmt.Sprintf(`input#%s`, name)).First()
	fillLocatorInput(t, locator, name, value)
}

func fillLocatorInput(t g.FullGinkgoTInterface, locator playwright.Locator, name string, value string) {
	err := playwright.NewPlaywrightAssertions(5_000).Locator(locator).ToBeEnabled()
	assert.NoError(t, err, "form field name=%q is not enabled", name)
	err = locator.Fill(value)
	assert.NoError(t, err)
	g.GinkgoWriter.Printf("Filled name=%q\n", name)
}

func FillTextarea(t g.FullGinkgoTInterface, page playwright.Page, name, value string) {
	locator := page.Locator(fmt.Sprintf(`textarea#%s`, name)).First()
	fillLocatorInput(t, locator, name, value)
}

func FillAutocompleteInput(t g.FullGinkgoTInterface, page playwright.Page, name, value string) {
	locator := page.Locator(fmt.Sprintf(`input#%s`, name)).First()
	err := playwright.NewPlaywrightAssertions(5_000).Locator(locator).ToBeEnabled()
	assert.NoError(t, err, "autocomplete field name=%q is not enabled", name)
	err = locator.Fill(value)
	assert.NoError(t, err)

	option := page.Locator("mat-option").First()
	err = playwright.NewPlaywrightAssertions(5_000).Locator(option).ToBeVisible()
	assert.NoError(t, err, "no mat-option visible after typing %q into name=%q", value, name)
	err = option.Click()
	assert.NoError(t, err)
	g.GinkgoWriter.Printf("Filled autocomplete name=%q with %q\n", name, value)
}

func SubmitForm(t g.FullGinkgoTInterface, page playwright.Page) {
	locator := page.Locator(`button[type="submit"]`).First()
	err := playwright.NewPlaywrightAssertions(5_000).Locator(locator).ToBeEnabled()
	assert.NoError(t, err, "submit button not enabled")
	err = locator.ScrollIntoViewIfNeeded()
	assert.NoError(t, err)
	err = locator.Click()
	assert.NoError(t, err)
	g.GinkgoWriter.Println("Submitted form")
}

func assertDrinksEqual(t g.FullGinkgoTInterface, expected, actual drinks.Drink) {
	expected.Timestamp = nil
	actual.Timestamp = nil

	expected.PubID = 0
	actual.PubID = 0

	expected.Points = 0
	actual.Points = 0

	expected.Letter = ""
	actual.Letter = ""

	assert.Equal(t, expected, actual)
}
