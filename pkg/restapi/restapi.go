package restapi

import (
	"encoding/json"
	"github.com/function61/pi-security-module/pkg/command"
	"github.com/function61/pi-security-module/pkg/commandhandlers"
	"github.com/function61/pi-security-module/pkg/domain"
	"github.com/function61/pi-security-module/pkg/httputil"
	"github.com/function61/pi-security-module/pkg/physicalauth"
	"github.com/function61/pi-security-module/pkg/state"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strings"
	"time"
)

func Define(router *mux.Router, st *state.State) {
	router.HandleFunc("/auditlog", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httputil.RespondHttpJson(st.State.AuditLog, http.StatusOK, w)
	}))

	router.HandleFunc("/command/{commandName}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		commandName := mux.Vars(r)["commandName"]

		// only command able to be invoked anonymously is the Unseal command
		commandNeedsAuthorization := commandName != "database.Unseal"

		if commandNeedsAuthorization && httputil.ErrorIfSealed(w, r, st.IsUnsealed()) {
			return
		}

		cmdStructBuilder, commandExists := commandhandlers.StructBuilders[commandName]
		if !commandExists {
			httputil.CommandCustomError(w, r, "unsupported_command", nil, http.StatusBadRequest)
			return
		}

		ctx := &command.Ctx{
			State: st,
			Meta:  domain.Meta(time.Now(), "2"),
		}

		cmdStruct := cmdStructBuilder()

		// FIXME: assert application/json
		if errJson := json.NewDecoder(r.Body).Decode(cmdStruct); errJson != nil {
			httputil.CommandCustomError(w, r, "json_parsing_failed", errJson, http.StatusBadRequest)
			return
		}

		if errValidate := cmdStruct.Validate(); errValidate != nil {
			httputil.CommandCustomError(w, r, "command_validation_failed", errValidate, http.StatusBadRequest)
			return
		}

		if errInvoke := cmdStruct.Invoke(ctx); errInvoke != nil {
			httputil.CommandCustomError(w, r, "command_failed", errInvoke, http.StatusBadRequest)
			return
		}

		raisedEvents := ctx.GetRaisedEvents()

		log.Printf("Command %s raised %d event(s)", commandName, len(raisedEvents))

		st.EventLog.AppendBatch(raisedEvents)

		httputil.RespondHttpJson(httputil.GenericSuccess(), http.StatusOK, w)
	}))

	router.HandleFunc("/folder/{folderId}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if httputil.ErrorIfSealed(w, r, st.IsUnsealed()) {
			return
		}

		folder := st.FolderById(mux.Vars(r)["folderId"])

		accounts := st.AccountsByFolder(folder.Id)
		subFolders := st.SubfoldersById(folder.Id)
		parentFolders := []state.Folder{}

		parentId := folder.ParentId
		for parentId != "" {
			parent := st.FolderById(parentId)

			parentFolders = append(parentFolders, *parent)

			parentId = parent.ParentId
		}

		type FolderResponse struct {
			Folder        *state.Folder
			SubFolders    []state.Folder
			ParentFolders []state.Folder
			Accounts      []state.SecureAccount
		}

		resp := FolderResponse{folder, subFolders, parentFolders, accounts}

		httputil.RespondHttpJson(resp, http.StatusOK, w)
	}))

	router.HandleFunc("/accounts", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if httputil.ErrorIfSealed(w, r, st.IsUnsealed()) {
			return
		}

		search := strings.ToLower(r.URL.Query().Get("search"))
		sshkey := strings.ToLower(r.URL.Query().Get("sshkey"))

		w.Header().Set("Content-Type", "application/json")

		matches := []state.SecureAccount{}

		if sshkey == "y" {
			for _, account := range st.State.Accounts {
				for _, secret := range account.Secrets {
					if secret.SshPublicKeyAuthorized == "" {
						continue
					}

					matches = append(matches, account.ToSecureAccount())
				}
			}
		} else if search == "" { // no filter => return all
			for _, s := range st.State.Accounts {
				matches = append(matches, s.ToSecureAccount())
			}
		} else { // search filter
			for _, s := range st.State.Accounts {
				if !strings.Contains(strings.ToLower(s.Title), search) {
					continue
				}

				matches = append(matches, s.ToSecureAccount())
			}
		}

		httputil.RespondHttpJson(matches, http.StatusOK, w)
	}))

	router.HandleFunc("/accounts/{accountId}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if httputil.ErrorIfSealed(w, r, st.IsUnsealed()) {
			return
		}

		account := st.AccountById(mux.Vars(r)["accountId"])

		if account == nil {
			httputil.CommandCustomError(w, r, "account_not_found", nil, http.StatusNotFound)
			return
		}

		httputil.RespondHttpJson(account, http.StatusOK, w)
	}))

	router.HandleFunc("/accounts/{accountId}/secrets", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if httputil.ErrorIfSealed(w, r, st.IsUnsealed()) {
			return
		}

		account := st.AccountById(mux.Vars(r)["accountId"])

		if account == nil {
			httputil.CommandCustomError(w, r, "account_not_found", nil, http.StatusNotFound)
			return
		}

		authorized, err := physicalauth.Dummy()
		if err != nil {
			httputil.CommandCustomError(w, r, "technical_error_in_physical_authorization", err, http.StatusInternalServerError)
			return
		}

		// FIXME: PasswordExposed via enum

		if !authorized {
			httputil.CommandCustomError(w, r, "did_not_receive_physical_authorization", nil, http.StatusForbidden)
			return
		}

		st.EventLog.Append(domain.NewAccountSecretUsed(
			account.Id,
			"PasswordExposed",
			domain.Meta(time.Now(), "2")))

		httputil.RespondHttpJson(account.GetSecrets(), http.StatusOK, w)
	}))
}