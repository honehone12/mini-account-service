package session

import (
	"errors"
	"mini-account-service/server/quick"
	"time"

	"github.com/google/uuid"
	gorillasession "github.com/gorilla/sessions"
	echosession "github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

const (
	sessionName  = "mini-session"
	versionKey   = "session-version"
	uuidKey      = "user-uuid"
	createdAtKey = "created-at"
)

const (
	CurrentSessionFuncVersion uint32 = 1
)

var (
	ErrorSessionNotStored        = errors.New("session not stored")
	ErrorSessionVersionMissMatch = errors.New("session version is not updated")
	ErrorSessionExpired          = errors.New("session is expired")
	ErrorSessionParseError       = errors.New("could not parsed stored value")
)

type SessionStatus uint8

const (
	Ok SessionStatus = 1 + iota
	NotStored
	Rejected
	Error
)

type SessionData struct {
	Version uint32
	Status  SessionStatus
	Uuid    string
}

func NewSessionStroe(secret string) gorillasession.Store {
	return gorillasession.NewCookieStore([]byte(secret))
}

func Set(c echo.Context, uuid string) error {
	sess, err := echosession.Get(sessionName, c)
	if err != nil {
		c.Logger().Error(err)
		return quick.ServiceError()
	}

	now := time.Now().Unix()
	sess.Options = NewSessionOptions()
	sess.Values[versionKey] = CurrentSessionFuncVersion
	sess.Values[uuidKey] = uuid
	sess.Values[createdAtKey] = now
	sess.Save(c.Request(), c.Response())
	return nil
}

func Get(c echo.Context) (uint32, string, error) {
	sess, err := echosession.Get(sessionName, c)
	if err != nil {
		return 0, "", err
	}

	version, ok := sess.Values[versionKey]
	if !ok {
		return 0, "", ErrorSessionNotStored
	}
	versionUint, ok := version.(uint32)
	if !ok {
		return 0, "", ErrorSessionParseError
	}
	if versionUint != CurrentSessionFuncVersion {
		return 0, "", ErrorSessionVersionMissMatch
	}

	createdAt, ok := sess.Values[createdAtKey]
	if !ok {
		return 0, "", ErrorSessionNotStored
	}
	expiration, ok := createdAt.(int64)
	if !ok {
		return 0, "", ErrorSessionParseError
	}
	expiration += 60 * 60 * 24
	if time.Now().Unix() > expiration {
		return 0, "", ErrorSessionExpired
	}

	uuid, ok := sess.Values[uuidKey]
	if !ok {
		return 0, "", ErrorSessionNotStored
	}

	uuidStr, ok := uuid.(string)
	if !ok {
		return 0, "", ErrorSessionParseError
	}

	return versionUint, uuidStr, nil
}

func GetAndVerify(c echo.Context) (SessionData, error) {
	sessVersion, userUuid, err := Get(c)
	if err == ErrorSessionNotStored {

		return SessionData{
			Version: 0,
			Status:  NotStored,
			Uuid:    "",
		}, err
	} else if err == ErrorSessionExpired || err == ErrorSessionParseError {
		return SessionData{
			Version: 0,
			Status:  Rejected,
			Uuid:    "",
		}, err
	} else if err != nil {
		return SessionData{
			Version: 0,
			Status:  Error,
			Uuid:    "",
		}, err
	}

	if _, err = uuid.Parse(userUuid); err != nil {
		return SessionData{
			Version: 0,
			Status:  Rejected,
			Uuid:    "",
		}, err
	}

	return SessionData{
		Version: sessVersion,
		Status:  Ok,
		Uuid:    userUuid,
	}, nil
}

func RequireSession(c echo.Context) (SessionData, error) {
	sess, err := GetAndVerify(c)
	switch sess.Status {
	case Ok:
		return sess, nil
	case Error:
		c.Logger().Error(err)
		return sess, quick.ServiceError()
	case Rejected:
		fallthrough
	case NotStored:
		c.Logger().Warn(err)
		return sess, quick.NotAllowed()
	default:
		c.Logger().Error("not implemented")
		c.Logger().Error(err)
		return sess, quick.ServiceError()
	}
}
