package api

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/ini.v1"
)

type apiCredential struct {
	UUID   string
	Secret string
}

func newApiCredentialFromIni(r io.Reader) (*apiCredential, error) {
	inif, err := ini.Load(r)
	if err != nil {
		return nil, err
	}
	section := inif.Section("general")
	if section == nil {
		return nil, nil
	}
	uuid := section.Key("apiuuid").String()
	secret := section.Key("apisecret").String()
	if len(uuid) > 0 && len(secret) > 0 {
		return &apiCredential{
			UUID:   uuid,
			Secret: secret,
		}, nil
	}
	return nil, nil
}

func newApiCredentialFromEnv() *apiCredential {
	uuid := os.Getenv("RACKCORP_API_UUID")
	secret := os.Getenv("RACKCORP_API_SECRET")

	if len(uuid) > 0 && len(secret) > 0 {
		return &apiCredential{
			UUID:   uuid,
			Secret: secret,
		}
	}

	uuid = os.Getenv("RACKCORP_APIUUID")
	secret = os.Getenv("RACKCORP_APISECRET")

	if len(uuid) > 0 && len(secret) > 0 {
		return &apiCredential{
			UUID:   uuid,
			Secret: secret,
		}
	}

	home := os.Getenv("HOME")
	if len(home) == 0 {
		return nil
	}

	paths := []string{
		filepath.Join(home, ".rackcorp"),
		filepath.Join(home, ".config", ".rackcorp", "config"),
	}
	for _, path := range paths {
		f, err := os.Open(path)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			// TODO log?
			continue
		}
		defer f.Close()
		cred, err := newApiCredentialFromIni(f)
		if err != nil {
			// TODO log?
			continue
		}
		if cred != nil {
			return cred
		}
	}

	return nil
}
