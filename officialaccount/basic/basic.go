package basic

import (
	"fmt"

	"github.com/amazing-gao/wechat/v2/officialaccount/context"
	"github.com/amazing-gao/wechat/v2/util"
)

// Basic struct
type Basic struct {
	*context.Context
}

// NewBasic 实例
func NewBasic(context *context.Context) *Basic {
	basic := new(Basic)
	basic.Context = context
	return basic
}

// IPListRes 获取微信服务器IP地址 返回结果
type IPListRes struct {
	util.CommonError
	IPList []string `json:"ip_list"`
}

// GetCallbackIP 获取微信callback IP地址
func (basic *Basic) GetCallbackIP() ([]string, error) {
	ak, err := basic.GetAccessToken()
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s/cgi-bin/getcallbackip?access_token=%s", basic.Server, ak)
	data, err := util.HTTPGet(url)
	if err != nil {
		return nil, err
	}
	ipListRes := &IPListRes{}
	err = util.DecodeWithError(data, ipListRes, "GetCallbackIP")
	return ipListRes.IPList, err
}

// GetAPIDomainIP 获取微信API接口 IP地址
func (basic *Basic) GetAPIDomainIP() ([]string, error) {
	ak, err := basic.GetAccessToken()
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s/cgi-bin/get_api_domain_ip?access_token=%s", basic.Server, ak)
	data, err := util.HTTPGet(url)
	if err != nil {
		return nil, err
	}
	ipListRes := &IPListRes{}
	err = util.DecodeWithError(data, ipListRes, "GetAPIDomainIP")
	return ipListRes.IPList, err
}

// ClearQuota 清理接口调用次数
func (basic *Basic) ClearQuota() error {
	ak, err := basic.GetAccessToken()
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/cgi-bin/clear_quota?access_token=%s", basic.Server, ak)
	data, err := util.PostJSON(url, map[string]string{
		"appid": basic.AppID,
	})
	if err != nil {
		return err
	}
	return util.DecodeWithCommonError(data, "ClearQuota")
}
