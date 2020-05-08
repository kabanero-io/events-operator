package semverimage_test

import (
    "testing"

    "github.com/kabanero-io/events-operator/pkg/semverimage"
    "github.com/kabanero-io/events-operator/pkg/semverimage/model"
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

func TestEvent(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "Event Suite")
}

var _ = Describe("Semverimage", func() {

    var (
        strToVersions   []model.StringToVersion
        compatData      []model.CompatStruct
        greaterThanData []model.GreaterThanStruct
    )

    strToVersions = []model.StringToVersion{
        {
            Str: "0.0.0",
            Version: &semverimage.Version{
                Major: 0,
                Minor: 0,
                Patch: 0,
            },
        },
        {
            Str: "1.2.3",
            Version: &semverimage.Version{
                Major: 1,
                Minor: 2,
                Patch: 3,
            },
        },
        {
            Str: "1",
            Version: &semverimage.Version{
                Major: 1,
                Minor: -1,
                Patch: -1,
            },
        },
        {
            Str: "0.3",
            Version: &semverimage.Version{
                Major: 0,
                Minor: 3,
                Patch: -1,
            },
        },
    }

    compatData = []model.CompatStruct{
        {
            Version1:   "1.2.3",
            Version2:   "1.2.3",
            Compatible: true,
        },
        {
            Version1:   "1.2.3",
            Version2:   "1.2.4",
            Compatible: false,
        },
        {
            Version1:   "1.2.3",
            Version2:   "1.3.3",
            Compatible: false,
        },
        {
            Version1:   "1.2.3",
            Version2:   "0.2.3",
            Compatible: false,
        },
        {
            Version1:   "1.3",
            Version2:   "1.3.0",
            Compatible: true,
        },
        {
            Version1:   "1.3",
            Version2:   "1.3.99",
            Compatible: true,
        },
        {
            Version1:   "1.3",
            Version2:   "1.2.3",
            Compatible: false,
        },
        {
            Version1:   "1.3",
            Version2:   "1.4.0",
            Compatible: false,
        },
        {
            Version1:   "1.3",
            Version2:   "2.3.0",
            Compatible: false,
        },
        {
            Version1:   "1",
            Version2:   "1.2.3",
            Compatible: true,
        },
        {
            Version1:   "1",
            Version2:   "1.3.3",
            Compatible: true,
        },
        {
            Version1:   "1",
            Version2:   "2.0.0",
            Compatible: false,
        },
        {
            Version1:   "0",
            Version2:   "1.2.3",
            Compatible: false,
        },
        {
            Version1:   "0",
            Version2:   "0.5.6",
            Compatible: true,
        },
    }

    greaterThanData = []model.GreaterThanStruct{
        {
            Version1:    "1.2.3",
            Version2:    "1.0.0",
            GreaterThan: true,
        },
        {
            Version1:    "1.2.3",
            Version2:    "1.2.3",
            GreaterThan: false,
        },
        {
            Version1:    "1.2.3",
            Version2:    "1.1.0",
            GreaterThan: true,
        },
        {
            Version1:    "1.2.3",
            Version2:    "1.2.4",
            GreaterThan: false,
        },
        {
            Version1:    "1.2.3",
            Version2:    "1.3.0",
            GreaterThan: false,
        },
        {
            Version1:    "1.2.3",
            Version2:    "2.0.0",
            GreaterThan: false,
        },
    }

    It("should return a version without a patch level", func() {
        for _, strToVer := range strToVersions {
            ver, err := semverimage.NewVersion(strToVer.Str)
            Expect(err).Should(BeNil())
            Expect(*ver).Should(Equal(*strToVer.Version))
        }
    })

    It("should return true is the two versions are compatible", func() {
        for _, compat := range compatData {
            version1, err := semverimage.NewVersion(compat.Version1)
            Expect(err).Should(BeNil())
            version2, err := semverimage.NewVersion(compat.Version2)
            Expect(err).Should(BeNil())
            Expect(version1.IsCompatible(version2)).Should(Equal(compat.Compatible))
        }

    })

    It("should test if a version is greater that the other version", func() {
        for _, greater := range greaterThanData {
            version1, err := semverimage.NewVersion(greater.Version1)
            Expect(err).Should(BeNil())
            version2, err := semverimage.NewVersion(greater.Version2)
            Expect(err).Should(BeNil())
            Expect(version1.GreaterThan(version2)).Should(Equal(greater.GreaterThan))
        }
    })
})
