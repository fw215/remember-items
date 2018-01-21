package main

import (
	"testing"
)

func Test_SetGoogleConfig(t *testing.T) {
	err := SetGoogleConfig()
	if err != nil {
		t.Error(err)
	} else {
		t.Log(appConf)
	}
}

func Test_GetGoogleAuthURL(t *testing.T) {
	url, err := GetGoogleAuthURL()
	if err != nil {
		t.Error(url, err)
	} else {
		t.Log(url)
	}
}
