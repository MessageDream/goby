package tokenService

import (
	"time"

	. "github.com/MessageDream/goby/core"
	"github.com/MessageDream/goby/model"
	"github.com/MessageDream/goby/module/infrastructure"
)

func List(user *model.User) ([]*model.UserToken, error) {
	tokens, err := user.Tokens()
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

func Create(user *model.User, createdBy, friendlyName, description string, ttl int64) (*model.UserToken, error) {

	if ttl == 0 {
		ttl = 2 * 30 * 24 * 60 * 60 * 1000
	}

	token := &model.UserToken{
		UID:         user.ID,
		Name:        friendlyName,
		Token:       infrastructure.GetRandomString(28) + user.Rands,
		CreatedBy:   createdBy,
		Description: description,
		ExpiresAt:   time.Now().Add(time.Duration(ttl) * time.Millisecond),
	}

	if exist, err := token.Exist(); err != nil || exist {
		if exist {
			return nil, ErrTokenAlreadyExist
		}
		return nil, err
	}

	return token, token.Create(nil)
}

func Update(uid uint64, fromName, toName string, ttl int64) (*model.UserToken, error) {
	token := &model.UserToken{
		UID:  uid,
		Name: toName,
	}

	if exist, err := token.Exist(); err != nil || exist {
		if exist {
			return nil, ErrTokenAlreadyExist
		}
		return nil, err
	}

	token.Name = fromName

	if exist, err := token.Get(); err != nil || !exist {
		if !exist {
			return nil, ErrTokenNotExist
		}
		return nil, err
	}

	token.Name = toName
	token.ExpiresAt = token.ExpiresAt.Add(time.Duration(ttl) * time.Millisecond)

	return token, token.Update(nil)
}

func DeleteByUIDAndName(uid uint64, name string) error {
	token := &model.UserToken{UID: uid, Name: name}

	if exist, err := token.Exist(); err != nil || !exist {
		if !exist {
			return ErrTokenNotExist
		}
		return err
	}

	return token.DeleteByNameAndUID(nil)
}

func Delete(uid uint64, creator string) error {
	token := &model.UserToken{UID: uid, CreatedBy: creator}

	if exist, err := token.Exist(); err != nil || !exist {
		if !exist {
			return ErrTokenNotExist
		}
		return err
	}

	return token.DeleteByCreatorAndUID(nil)
}

func AuthWithToken(token string, isSession bool) error {
	_, tok, err := checkToken(token)
	if err != nil {
		return err
	}

	tok.IsSession = isSession

	return tok.Update(nil, "is_session")

}

func CheckTokenSession(token string) (*model.User, error) {
	user, tok, err := checkToken(token)
	if err != nil {
		return nil, err
	}
	if !tok.IsSession {
		return nil, ErrTokenNotExist
	}
	return user, nil
}

func checkToken(token string) (*model.User, *model.UserToken, error) {
	if len(token) < 37 {
		return nil, nil, ErrTokenNotExist
	}
	rs := []rune(token)
	rands := string(rs[28:])

	user := &model.User{Rands: rands}

	if exist, err := user.Get(); err != nil || !exist {
		if !exist {
			return nil, nil, ErrTokenNotExist
		}
		return nil, nil, err
	}

	tok := &model.UserToken{UID: user.ID, Token: token}

	if exist, err := tok.GetUnexpired(); err != nil || !exist {
		if !exist {
			return nil, nil, ErrTokenNotExist
		}
		return nil, nil, err
	}

	return user, tok, nil
}
