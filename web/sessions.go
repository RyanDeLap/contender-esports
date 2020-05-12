package web

import (
	"net/http"

	"github.com/gorilla/sessions"
)

var store *sessions.FilesystemStore

type session struct {
	request        *http.Request
	responseWriter http.ResponseWriter
}

func NewSession(request *http.Request, responseWriter http.ResponseWriter) session {
	if store == nil {
		store = sessions.NewFilesystemStore("", []byte("my spooky secret session key, this should be hidden btw, like not in the source code hidden btw"))
		store.Options = &sessions.Options{
			Path:     "/",
			Domain:   "localhost",
			MaxAge:   86400 * 7,
			HttpOnly: true,
		}
	}

	return session{
		request,
		responseWriter,
	}
}

func (s session) Set(key string, value interface{}) error {
	sessionStore, err := store.Get(s.request, "s")
	if err != nil {
		return err
	}

	sessionStore.Values[key] = value
	if err = store.Save(s.request, s.responseWriter, sessionStore); err != nil {
		return err
	}
	return nil
}

func (s session) Unset(key string) error {
	sessionStore, err := store.Get(s.request, "s")
	if err != nil {
		return err
	}

	delete(sessionStore.Values, key)
	if err = store.Save(s.request, s.responseWriter, sessionStore); err != nil {
		return err
	}
	return nil
}

func (s session) Clear() error {
	sessionStore, err := store.Get(s.request, "s")
	if err != nil {
		return err
	}

	sessionStore.Values = make(map[interface{}]interface{})
	if err = store.Save(s.request, s.responseWriter, sessionStore); err != nil {
		return err
	}
	return nil
}

func (s session) Get(key string) (interface{}, error) {
	sessionStore, err := store.Get(s.request, "s")
	if err != nil {
		return nil, err
	}
	return sessionStore.Values[key], nil
}

func (s session) SetFlash(key string, value interface{}) error {
	sessionStore, err := store.Get(s.request, "s")
	if err != nil {
		return err
	}
	sessionStore.AddFlash(value, key)
	return store.Save(s.request, s.responseWriter, sessionStore)
}

func (s session) GetFlash(key string) interface{} {
	sessionStore, err := store.Get(s.request, "s")
	if err != nil {
		return nil
	}
	flashes := sessionStore.Flashes(key)
	var flash interface{}

	if len(flashes) == 0 {
		flash = nil
	} else {
		flash = flashes[0]
	}
	err = store.Save(s.request, s.responseWriter, sessionStore)
	return flash
}
