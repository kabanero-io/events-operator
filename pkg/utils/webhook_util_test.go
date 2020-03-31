package utils_test

import (
	"github.com/kabanero-io/events-operator/pkg/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestWebhookUtil(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Webhook Utils Suite")
}

var _ = Describe("TestWebhookUtil", func() {
	Context("Payload validation", func() {
		payload := `{"msg: "Hello, world!"}`
		secret := "my-super-secret-secret"
		expectedSig := "886baa20847c41b910b5c4f85b3303ac49538fc1"
		It("should have the expected hash", func() {
			Expect(utils.ValidatePayload("sha1", secret, expectedSig, []byte(payload))).ToNot(HaveOccurred())
		})
	})
})
