package e2e

import (
	"testing"
	"time"

	"github.com/mxschmitt/playwright-go"
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

}, func() {
	// Deliberately empty
})

var _ = g.SynchronizedAfterSuite(func() {}, func() {
	if !IsHeadless() {
		time.Sleep(5 * time.Second)
	}
})
