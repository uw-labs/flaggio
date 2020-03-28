package operator

import (
	"net"
)

// InNetwork operator will check if the value from the user context is an ip
// that is included in any of the networks configured on the flag.
func InNetwork(usrValue interface{}, validValues []interface{}) (bool, error) {
	for _, v := range validValues {
		ok, err := inNetwork(v, usrValue)
		if err != nil {
			return false, err
		}
		if ok {
			return true, nil
		}
	}
	return false, nil
}

func inNetwork(cnstrnValue, userValue interface{}) (bool, error) {
	u, err := toString(userValue)
	if err != nil {
		return false, err
	}
	v, ok := cnstrnValue.(string)
	if !ok {
		return false, nil
	}
	userIP := net.ParseIP(u)
	_, ipnet, err := net.ParseCIDR(v)
	if err != nil {
		return false, err
	}

	return ipnet.Contains(userIP), nil
}
