// +build !deploy

package vpanel

import "errors"

func Asset(_ string) ([]byte, error) {
	return []byte{}, errors.New("Cannot serve static content from a development instance")
}
