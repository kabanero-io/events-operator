package semverimage

import (
    "fmt"
    "strings"
    "strconv"
)

type Version struct {
    Major int  /* major version */
    Minor int /* minor version, or -1 */
    Patch int /* patch version, or -1 */
}

/* Return true if current version is compatible with other version.
   "1" is compatible with "1.x.y" for any x, y
   "1.2" is compatible with "1.2.y", for an y
   "1.2.3" is compatible with "1.2.3" only
*/
func (ver *Version) IsCompatible(otherVer *Version) bool {
    if ver.Major != otherVer.Major  {
        return false
    }

    /* same major version. Check minor */
    if ver.Minor < 0 && otherVer.Minor >= 0 {
        return true
    }

    if ver.Minor >= 0 && otherVer.Minor < 0 {
        return true
    }

    /* both < 0, or both > 0 */
    if ver.Minor != otherVer.Minor {
        return false
    }

    /* same minor version. Check patch */
    if ver.Patch < 0 && otherVer.Patch >= 0 {
        return true
    }
    if ver.Patch >= 0 && otherVer.Patch < 0 {
        return true
    }
    return ver.Patch == otherVer.Patch
}


/* Return true if this version is greater than other version*/
func (ver *Version) GreaterThan(otherVer *Version) bool {
    if ver.Major > otherVer.Major {
        return true
    }
    if ver.Major < otherVer.Major {
        return false
    }

    /* same major */
    if ver.Minor > otherVer.Minor {
        return true
    }
    if ver.Minor < otherVer.Minor {
        return false
    }

    /* same minor */
    return ver.Patch > otherVer.Patch
}

/* Parse version string for form "major.minor.patch", where minor and patch are optional */
func NewVersion(str string) (*Version, error ) {
    ret := &Version {
        Minor: -1, 
        Patch: -1,
    }

    var err error
    var remainder string
    ret.Major, remainder, err = parseInt(str, ".")
    if err != nil {
        return nil, fmt.Errorf("%v not a semantic version of form major.minor.patch", str)
    }
    if remainder == ""  {
        return ret, nil
    }

    ret.Minor, remainder, err = parseInt(remainder, ".")
    if err != nil {
        return nil, fmt.Errorf("%v not a semantic version of form major.minor.patch", str)
    }
    if remainder == ""  {
        return ret, nil
    }

    ret.Patch, err = strconv.Atoi(remainder)
    if err != nil {
        return nil, fmt.Errorf("%v not a semantic version of form major.minor.patch", str)
    }
    return ret, nil
}

func (ver *Version) String() string {
    var ret string
    if ver.Minor < 0 {
        ret = fmt.Sprintf("%d", ver.Major)
    } else {
        if ver.Patch < 0 {
            ret = fmt.Sprintf("%d.%d", ver.Major, ver.Minor)
        } else {
            ret = fmt.Sprintf("%d.%d.%d", ver.Major, ver.Minor, ver.Patch)
        }
    }
    return ret
}

/*  Fetch first integer component from the input string, up to the optional separator.
 Return:
   ret: the integer
   remainder: what remains after separtor
   error: any error
*/
func parseInt(str string, separator string) (ret int, remainder string, err error ) {
    remainder = str
    index := strings.Index(str, separator)
    if index == 0 {
        return ret, remainder, fmt.Errorf("separtor %s first component of string %s", separator, str)
    }
    if index < 0 {
        /* no separator found */
        remainder = ""
        ret, err = strconv.Atoi(str)
        return ret, remainder, err
    }

    intStr := str[0:index]
    if index == len(str) -1 {
       remainder = ""
    } else {
       remainder = str[index+1:]
    }
    ret, err = strconv.Atoi(intStr)
    return ret, remainder, err
}
