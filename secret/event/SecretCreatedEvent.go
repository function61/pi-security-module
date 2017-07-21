package event

import (
	"github.com/function61/pi-security-module/state"
)

type SecretCreated struct {
	Id       string
	FolderId string
	Title    string
}

func (e *SecretCreated) Apply() {
	secret := state.InsecureSecret{
		Id:       e.Id,
		FolderId: e.FolderId,
		Title:    e.Title,
	}

	state.Inst.State.Secrets = append(state.Inst.State.Secrets, secret)
}