package semverimage
import (
    "testing"
)

type strToVersion struct {
    str string
    version *Version
}

var strToVersions [] strToVersion = []strToVersion{
    { 
    str: "0.0.0",
    version: &Version {
        Major: 0,
        Minor: 0,
        Patch: 0,
        },
    },
    { 
    str: "1.2.3",
    version: &Version {
        Major: 1,
        Minor: 2,
        Patch: 3,
        },
    },
    { 
    str: "1",
    version: &Version {
        Major: 1,
        Minor: -1,
        Patch: -1,
        },
    },
    { 
    str: "0.3",
    version: &Version {
        Major: 0,
        Minor: 3,
        Patch: -1,
        },
    },
}

func TestNewVersion(t *testing.T) {
    for _, strToVer := range strToVersions {
        ver, err := NewVersion(strToVer.str)
        if err != nil {
            t.Fatalf("Unable to convert %s to semverimage. Error: %v", strToVer.str, err)
        }
        if *ver != *strToVer.version {
            t.Fatalf("%s parsed incorrectly. Expecting: %v, but got: %v", strToVer.str, *strToVer.version, *ver)
        }
    }
}

type compatStruct struct {
    version1 string
    version2 string
    compatible bool
}

var compatData []compatStruct = []compatStruct {
    {
         version1: "1.2.3",
         version2: "1.2.3",
         compatible: true,
    },
    {
         version1: "1.2.3",
         version2: "1.2.4",
         compatible: false,
    },
    {
         version1: "1.2.3",
         version2: "1.3.3",
         compatible: false,
    },
    {
         version1: "1.2.3",
         version2: "0.2.3",
         compatible: false,
    },
    {
         version1: "1.3",
         version2: "1.3.0",
         compatible: true,
    },
    {
         version1: "1.3",
         version2: "1.3.99",
         compatible: true,
    },
    {
         version1: "1.3",
         version2: "1.2.3",
         compatible: false,
    },
    {
         version1: "1.3",
         version2: "1.4.0",
         compatible: false,
    },
    {
         version1: "1.3",
         version2: "2.3.0",
         compatible: false,
    },
    {
         version1: "1",
         version2: "1.2.3",
         compatible: true,
    },
    {
         version1: "1",
         version2: "1.3.3",
         compatible: true,
    },
    {
         version1: "1",
         version2: "2.0.0",
         compatible: false,
    },
    {
         version1: "0",
         version2: "1.2.3",
         compatible: false,
    },
    {
         version1: "0",
         version2: "0.5.6",
         compatible: true,
    },
}


func TestCompatible(t *testing.T) {
    for _, compat := range compatData {
        ver1, err := NewVersion(compat.version1)
        if err != nil {
            t.Fatal(err)
        }
        ver2, err := NewVersion(compat.version2)
        if err != nil {
            t.Fatal(err)
        }
        compatible := ver1.IsCompatible(ver2)
        if compatible != compat.compatible {
            t.Fatalf("Compatibility test failed for versions %s and %s. Expecting compatibility: %v", compat.version1, compat.version2, compat.compatible)
        }
    }
}

type greaterThanStruct struct {
    version1 string
    version2 string
    greaterThan bool
}

var greaterThanData []greaterThanStruct = []greaterThanStruct {
    {
         version1: "1.2.3",
         version2: "1.0.0",
         greaterThan: true,
    },
    {
         version1: "1.2.3",
         version2: "1.2.3",
         greaterThan: false,
    },
    {
         version1: "1.2.3",
         version2: "1.1.0",
         greaterThan: true,
    },
    {
         version1: "1.2.3",
         version2: "1.2.4",
         greaterThan: false,
    },
    {
         version1: "1.2.3",
         version2: "1.3.0",
         greaterThan: false,
    },
    {
         version1: "1.2.3",
         version2: "2.0.0",
         greaterThan: false,
    },
}

func TestGreaterThan(t *testing.T) {
    for _, greater := range greaterThanData {
        ver1, err := NewVersion(greater.version1)
        if err != nil {
            t.Fatal(err)
        }
        ver2, err := NewVersion(greater.version2)
        if err != nil {
            t.Fatal(err)
        }
        isGreater := ver1.GreaterThan(ver2)
        if isGreater != greater.greaterThan {
            t.Fatalf("GreaterhThan test failed for versions %s and %s. Expecting: %v", greater.version1, greater.version2, greater.greaterThan)
        }
    }
}
