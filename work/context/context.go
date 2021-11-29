package context

import (
	"github.com/amazing-gao/wechat/v2/credential"
	"github.com/amazing-gao/wechat/v2/work/config"
)

// Context struct
type Context struct {
	*config.Config
	credential.AccessTokenHandle
}
