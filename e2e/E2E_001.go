package e2e

import (
	"time"

	"github.com/mxschmitt/playwright-go"
	"github.com/ockendenjo/ale-trail/pkg/data"
	"github.com/ockendenjo/ale-trail/pkg/generator"
	g "github.com/onsi/ginkgo/v2"
	"github.com/stretchr/testify/assert"
)

var _ = TestCase("001", "Register, set display name, and log a drink", func() {
	var (
		t    g.FullGinkgoTInterface
		page playwright.Page
		lh   *LogHelper
		file *data.File
	)

	g.BeforeAll(func() {
		t = g.GinkgoT()
		page = LoadHomePage(t)
		lh = NewLogHelper(t)
	})

	g.It("should upload data into S3", func() {
		tm := time.Now().AddDate(0, 0, -3)
		df := new(generator.NewFile(tm))
		setPreviousDays(df)
		lh.PutDataFile(df)
		file = df
	})

	g.It("should load the home page", func() {
		assert.NotNil(t, file)

		text, err := page.GetByText("Every day between now and the wassail").TextContent()
		assert.NoError(t, err)
		assert.Contains(t, text, "goodies")
	})

	g.Context("Set stasher details", func() {
		g.It("should click a card", func() {
			err := page.GetByText("Stasher needed").First().Click()
			assert.NoError(t, err)
		})

		g.It("should enter details in the form", func() {
			FillFormControlInput(t, page, "name", "Harry Trotter")
			FillFormControlInput(t, page, "email", "harry@example.com")
			SubmitForm(t, page)
		})

		g.It("should update the data file", func() {
			file := lh.AwaitDataFile(func(file *data.File) bool {
				return file.Stashes[4].Details != nil
			})
			s4 := file.Stashes[4]
			assert.Equal(t, "Harry Trotter", s4.Stasher.Name)
			assert.Equal(t, "harry@example.com", s4.Stasher.Email)
		})
	})
})
