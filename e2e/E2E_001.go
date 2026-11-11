package e2e

import (
	"time"

	"github.com/mxschmitt/playwright-go"
	"github.com/ockendenjo/ale-trail/pkg/data"
	"github.com/ockendenjo/ale-trail/pkg/generator"
	g "github.com/onsi/ginkgo/v2"
	"github.com/stretchr/testify/assert"
)

var _ = TestCase("001", "Perform user actions on website", func() {
	var (
		t        g.FullGinkgoTInterface
		page     playwright.Page
		lh       *LogHelper
		file     *data.File
		now      time.Time
		adminKey string
		s4Key    string
	)

	g.BeforeAll(func() {
		t = g.GinkgoT()
		lh = NewLogHelper(t)
	})

	g.Context("reset data file", func() {
		g.It("should upload data into S3", func() {
			now = time.Now()
			tm := now.AddDate(0, 0, -3)
			df := generator.NewFile(tm)
			dfp := &df
			setPreviousDays(dfp)
			lh.PutDataFile(dfp)
			file = dfp
			adminKey = df.AdminKey
		})

		g.It("should wait for the same file to be returned", func() {
			lh.AwaitDataFile(func(file *data.File) bool {
				return file.AdminKey == adminKey
			})
		})
	})

	g.It("should load the home page", func() {
		page = LoadHomePage(t)
		assert.NotNil(t, file)

		text, err := page.GetByText("Every day between now and the wassail").TextContent()
		assert.NoError(t, err)
		assert.Contains(t, text, "goodies")
	})

	g.Context("Set stasher details", func() {
		var dateStr string
		var card playwright.Locator

		g.BeforeAll(func() {
			dateStr = now.AddDate(0, 0, 1).Format(time.DateOnly)
			card = page.GetByTestId(dateStr)
		})

		g.It("should check card contains expected text", func() {
			text, err := card.InnerText()
			assert.NoError(t, err)
			assert.Contains(t, text, "Stasher needed")
		})

		g.It("should click a card", func() {
			err := card.Click()
			assert.NoError(t, err)
		})

		g.It("should check URL", func() {
			AwaitURL(t, page, dateStr+"/stasher")
		})

		g.It("should enter details in the form", func() {
			FillFormControlInput(t, page, "name", "Harry Trotter")
			FillFormControlInput(t, page, "email", "harry@example.com")
			SubmitForm(t, page)
		})

		g.It("should show saved message", func() {
			AwaitVisible(t, page, "div[data-testid='state-loaded']")
		})

		g.It("should validate the data file has been updated", func() {
			file := lh.AwaitDataFile(func(file *data.File) bool {
				return file.Stashes[4].Stasher != nil
			})
			s4 := file.Stashes[4]
			assert.Equal(t, "Harry Trotter", s4.Stasher.Name)
			assert.Equal(t, "harry@example.com", s4.Stasher.Email)
			s4Key = s4.StashKey
		})

		g.It("should validate card on home page is updated", func() {
			NavigateTo(t, page, "")
			card = page.GetByTestId(dateStr)
			text, err := card.InnerText()
			assert.NoError(t, err)
			assert.NotContains(t, text, "Stasher needed")
		})
	})

	g.Context("Set stash location", func() {
		const (
			description = "Base of tree at west end of the meadows"
			w3w         = "///assets.grab.clouds"
			lat         = "55.941921"
			lon         = "-3.199647"
			contents    = "2 tins of Jarl"
		)

		var dateStr string

		g.BeforeAll(func() {
			dateStr = now.AddDate(0, 0, 1).Format(time.DateOnly)
		})

		g.It("Should go to stash location page", func() {
			urlPath := dateStr + "/location?key=" + s4Key
			NavigateTo(t, page, urlPath)

			AwaitVisible(t, page, "input#lat")
			AwaitVisible(t, page, "input#lon")
		})

		g.It("Should fill form", func() {
			FillFormControlInput(t, page, "w3w", w3w)
			FillFormControlInput(t, page, "lat", lat)
			FillFormControlInput(t, page, "lon", lon)
			FillTextarea(t, page, "description", description)
			FillTextarea(t, page, "contents", contents)
			SubmitForm(t, page)
		})

		g.It("should show success message", func() {
			AwaitVisible(t, page, "div[data-testid='state-saved']")
		})

		g.It("should validate the data file has been updated", func() {
			file := lh.AwaitDataFile(func(file *data.File) bool {
				return file.Stashes[4].Details != nil
			})
			deets := file.Stashes[4].Details

			exp := data.Details{
				Contents:    contents,
				Description: description,
				Lat:         55.941921,
				Lon:         -3.199647,
				W3W:         "assets.grab.clouds",
			}
			assert.Equal(t, exp, *deets)
		})

		g.It("should click edit again", func() {
			err := page.GetByTestId("edit-again").Click()
			assert.NoError(t, err)

			AwaitVisible(t, page, "input#lat")
		})

		g.It("should go back to home page", func() {
			err := page.GetByTestId("home-link").Click()
			assert.NoError(t, err)
		})
	})

	g.Context("Set stasher (conflict)", func() {
		var dateStr string
		var card playwright.Locator

		g.BeforeAll(func() {
			dateStr = now.AddDate(0, 0, 2).Format(time.DateOnly)
			card = page.GetByTestId(dateStr)
		})

		g.It("should check card contains expected text", func() {
			text, err := card.InnerText()
			assert.NoError(t, err)
			assert.Contains(t, text, "Stasher needed")
		})

		g.It("should click a card", func() {
			err := card.Click()
			assert.NoError(t, err)
		})

		g.It("should check URL", func() {
			AwaitURL(t, page, dateStr+"/stasher")
		})

		g.It("should update data file directly", func() {
			file = lh.GetDataFile()
			s5 := file.Stashes[5]
			s5.Stasher = &data.Stasher{
				Name:  "Stasher Baron Cohen",
				Email: "sbc@example.com",
			}
			lh.PutDataFile(file)
		})

		g.It("should wait for the same file to be returned", func() {
			lh.AwaitDataFile(func(file *data.File) bool {
				return file.Stashes[5].Stasher != nil
			})
		})

		g.It("should enter details in the form", func() {
			FillFormControlInput(t, page, "name", "Harry Trotter")
			FillFormControlInput(t, page, "email", "harry@example.com")
			SubmitForm(t, page)
		})

		g.It("should show saved message", func() {
			AwaitVisible(t, page, "div[data-testid='state-conflict']")
		})
	})
})
