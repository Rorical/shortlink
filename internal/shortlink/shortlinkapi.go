package shortlink

import (
	"shortlink/internal/config"
	"shortlink/internal/levelgo"
	"shortlink/internal/regex"
)

type ShortLinkApi struct {
	db *levelgo.LevelRpcClient
}

func NewApi(conf *config.Setting) *ShortLinkApi {

	db := levelgo.RpcClient(conf.LevelGo.URI)
	db.Connect()
	return &ShortLinkApi{
		db: db,
	}
}

func (self *ShortLinkApi) SetLink(id string, link string) error {
	if regex.IsProperId(id) && regex.IsUrl(link) {
		idbytes := config.StringIn(id)
		linkbytes := config.StringIn(link)
		isexist, err := self.db.Has(idbytes)
		if isexist {
			return ErrAlreadyExists
		}
		err = self.db.Set(idbytes, linkbytes)
		if err != nil {
			return err
		}
		return nil
	}
	return ErrIllegalParameters
}

func (self *ShortLinkApi) GetLink(id string) (string, error) {
	if regex.IsProperId(id) {
		idbytes := config.StringIn(id)
		link, err := self.db.Get(idbytes)
		if err != nil {
			if err == levelgo.ErrNotFound {
				return "", ErrDoesNotExists
			}
			return "", err
		}
		return config.StringOut(link), nil
	}
	return "", ErrIllegalParameters
}

func (self *ShortLinkApi) IsLinkExist(id string) (bool, error) {
	if regex.IsProperId(id) {
		idbytes := config.StringIn(id)
		isexist, err := self.db.Has(idbytes)
		if isexist {
			return true, nil
		} else {
			if err != nil {
				return false, err
			} else {
				return false, nil
			}
		}
	}
	return false, ErrIllegalParameters
}
