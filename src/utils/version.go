package utils

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

func GetGoVersion() (int64, int64, int64, error) {
	version := runtime.Version()
	v := version

	if len(v) < 2 {
		return 0, 0, 0, fmt.Errorf("invalid version: %q", version)
	}

	if strings.HasPrefix(v, "go") {
		v = v[2:]
	}

	if len(v) == 0 {
		return 0, 0, 0, fmt.Errorf("invalid version: %q", version)
	}

	vLstStr := strings.Split(v, ".")
	vLst := make([]int64, len(vLstStr))

	for i, j := range vLstStr {
		var err error
		vLst[i], err = strconv.ParseInt(j, 10, 64)
		if err != nil {
			return 0, 0, 0, fmt.Errorf("invalid version: %q", version)
		}
	}

	if len(vLst) == 0 {
		return 0, 0, 0, fmt.Errorf("invalid version: %q", version)
	} else if len(vLst) == 1 {
		return vLst[0], 0, 0, nil
	} else if len(vLst) == 2 {
		return vLst[0], vLst[1], 0, nil
	} else if len(vLst) == 3 {
		return vLst[0], vLst[1], vLst[2], nil
	} else {
		return 0, 0, 0, fmt.Errorf("invalid version: %q", version)
	}
}

func GetGoVersionMajor() (int64, error) {
	major, _, _, err := GetGoVersion()
	if err != nil {
		return 0, err
	}

	return major, nil
}

func GetGoVersionMinor() (int64, error) {
	_, minor, _, err := GetGoVersion()
	if err != nil {
		return 0, err
	}

	return minor, nil
}

func GetGoVersionPatch() (int64, error) {
	_, _, patch, err := GetGoVersion()
	if err != nil {
		return 0, err
	}

	return patch, nil
}

func GetGoVersionMajorMust() int64 {
	major, err := GetGoVersionMajor()
	if err != nil {
		panic(err)
	}
	return major
}

func GetGoVersionMinorMust() int64 {
	minor, err := GetGoVersionMinor()
	if err != nil {
		panic(err)
	}
	return minor
}

func GetGoVersionPatchMust() int64 {
	patch, err := GetGoVersionPatch()
	if err != nil {
		panic(err)
	}
	return patch
}
