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
		It("should return an error when passing an invalid hash type of sha2", func() {
			expectedSig := "886baa20847c41b910b5c4f85b3303ac49538fc1"
			Expect(utils.ValidatePayload("sha2", expectedSig, secret, []byte(payload))).To(HaveOccurred())
		})

		It("should have the expected sha1 hash", func() {
			expectedSig := "886baa20847c41b910b5c4f85b3303ac49538fc1"
			Expect(utils.ValidatePayload("sha1", expectedSig, secret, []byte(payload))).ToNot(HaveOccurred())
		})

		It("should have the expected sha256 hash", func() {
			expectedSig := "efa0d498cbfa1396d97d149ab64a3a9ced7922e828bc7cb0e72a564bec3fffb2"
			Expect(utils.ValidatePayload("sha256", expectedSig, secret, []byte(payload))).ToNot(HaveOccurred())
		})
	})
})
