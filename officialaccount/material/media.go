package material

import (
	"encoding/json"
	"fmt"

	"github.com/amazing-gao/wechat/v2/util"
)

// MediaType 媒体文件类型
type MediaType string

const (
	// MediaTypeImage 媒体文件:图片
	MediaTypeImage MediaType = "image"
	// MediaTypeVoice 媒体文件:声音
	MediaTypeVoice MediaType = "voice"
	// MediaTypeVideo 媒体文件:视频
	MediaTypeVideo MediaType = "video"
	// MediaTypeThumb 媒体文件:缩略图
	MediaTypeThumb MediaType = "thumb"
)

// Media 临时素材上传返回信息
type Media struct {
	util.CommonError

	Type         MediaType `json:"type"`
	MediaID      string    `json:"media_id"`
	ThumbMediaID string    `json:"thumb_media_id"`
	CreatedAt    int64     `json:"created_at"`
}

// MediaUpload 临时素材上传
func (material *Material) MediaUpload(mediaType MediaType, filename string) (media Media, err error) {
	var accessToken string
	accessToken, err = material.GetAccessToken()
	if err != nil {
		return
	}

	uri := fmt.Sprintf("%s/cgi-bin/media/upload?access_token=%s&type=%s", material.Server, accessToken, mediaType)
	var response []byte
	response, err = util.PostFile("media", filename, uri)
	if err != nil {
		return
	}
	err = json.Unmarshal(response, &media)
	if err != nil {
		return
	}
	if media.ErrCode != 0 {
		err = fmt.Errorf("MediaUpload error : errcode=%v , errmsg=%v", media.ErrCode, media.ErrMsg)
		return
	}
	return
}

// GetMediaURL 返回临时素材的下载地址供用户自己处理
// NOTICE: URL 不可公开，因为含access_token 需要立即另存文件
func (material *Material) GetMediaURL(mediaID string) (mediaURL string, err error) {
	var accessToken string
	accessToken, err = material.GetAccessToken()
	if err != nil {
		return
	}
	mediaURL = fmt.Sprintf("%s/cgi-bin/media/get?access_token=%s&media_id=%s", material.Server, accessToken, mediaID)
	return
}

// resMediaImage 图片上传返回结果
type resMediaImage struct {
	util.CommonError

	URL string `json:"url"`
}

// ImageUpload 图片上传
func (material *Material) ImageUpload(filename string) (url string, err error) {
	var accessToken string
	accessToken, err = material.GetAccessToken()
	if err != nil {
		return
	}

	uri := fmt.Sprintf("%s/cgi-bin/media/uploadimg?access_token=%s", material.Server, accessToken)
	var response []byte
	response, err = util.PostFile("media", filename, uri)
	if err != nil {
		return
	}
	var image resMediaImage
	err = json.Unmarshal(response, &image)
	if err != nil {
		return
	}
	if image.ErrCode != 0 {
		err = fmt.Errorf("UploadImage error : errcode=%v , errmsg=%v", image.ErrCode, image.ErrMsg)
		return
	}
	url = image.URL
	return
}
