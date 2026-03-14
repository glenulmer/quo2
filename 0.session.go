package main

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	. "pm/lib/output"
)

func (a App_t) SessionID(r *http.Request) (id string, ok bool) {
	if r == nil { return ``, false }
	c, _ := r.Cookie(a.Session.Name)
	if c == nil || Trim(c.Value) == `` { return ``, false }
	return c.Value, true
}

func RandomSessionID() string {
	buf := make([]byte, 16)
	if _, e := rand.Read(buf); e != nil { return Str(time.Now().UnixNano()) }
	return hex.EncodeToString(buf)
}

func (a *App_t) EnsureSession(w http.ResponseWriter, r *http.Request) string {
	if w == nil { return `` }
	id, ok := a.SessionID(r)
	if ok {
		a.SessionEpochGet(id)
		return id
	}
	id = RandomSessionID()
	a.SetSessionCookie(w, id)
	a.SessionEpochGet(id)
	return id
}

func (a *App_t) SetSessionCookie(w http.ResponseWriter, id string) {
	if w == nil { return }
	id = Trim(id)
	if id == `` { return }
	maxAge := a.Session.MaxAge
	http.SetCookie(w, &http.Cookie{
		Name:     a.Session.Name,
		Value:    id,
		Path:     a.Session.Path,
		MaxAge:   maxAge,
		Expires:  time.Now().Add(time.Duration(maxAge) * time.Second),
		HttpOnly: a.Session.HttpOnly,
		Secure:   a.Session.Secure,
		SameSite: a.Session.SameSite,
	})
}

func (a *App_t) ForceNewSession(w http.ResponseWriter) string {
	id := RandomSessionID()
	a.SetSessionCookie(w, id)
	a.SessionEpochGet(id)
	return id
}

func (a *App_t) SessionEpochGet(sessionID string) int {
	sessionID = Trim(sessionID)
	if sessionID == `` { return 1 }
	epoch, ok := a.sessionEpoch[sessionID]
	if !ok || epoch <= 0 {
		epoch = 1
		a.sessionEpoch[sessionID] = epoch
	}
	return epoch
}

func (a *App_t) SessionEpochBump(sessionID string) int {
	sessionID = Trim(sessionID)
	if sessionID == `` { return 1 }
	epoch := a.SessionEpochGet(sessionID) + 1
	a.sessionEpoch[sessionID] = epoch
	return epoch
}

func (a *App_t) SimpleStateGet(r *http.Request) (SimpleState_t, bool) {
	id, ok := a.SessionID(r)
	if !ok { return SimpleState_t{}, false }
	out, ok := a.sessionState[id]
	return out, ok
}

func (a *App_t) SimpleStateSet(sessionID string, state SimpleState_t) {
	sessionID = Trim(sessionID)
	if sessionID == `` { return }
	a.sessionState[sessionID] = state
}

func (a *App_t) FilterStateGet(r *http.Request) (FilterState_t, bool) {
	id, ok := a.SessionID(r)
	if !ok { return FilterState_t{}, false }
	out, ok := a.sessionFilters[id]
	return out, ok
}

func (a *App_t) FilterStateSet(sessionID string, state FilterState_t) {
	sessionID = Trim(sessionID)
	if sessionID == `` { return }
	a.sessionFilters[sessionID] = state
}

func (a *App_t) FilterStateClear(sessionID string) {
	sessionID = Trim(sessionID)
	if sessionID == `` { return }
	delete(a.sessionFilters, sessionID)
}

func (a *App_t) CustomerStateGet(r *http.Request) (CustomerState_t, bool) {
	id, ok := a.SessionID(r)
	if !ok { return CustomerState_t{}, false }
	out, ok := a.sessionCustomers[id]
	return out, ok
}

func (a *App_t) CustomerStateSet(sessionID string, state CustomerState_t) {
	sessionID = Trim(sessionID)
	if sessionID == `` { return }
	a.sessionCustomers[sessionID] = state
}

func (a *App_t) CustomerStateClear(sessionID string) {
	sessionID = Trim(sessionID)
	if sessionID == `` { return }
	delete(a.sessionCustomers, sessionID)
}
