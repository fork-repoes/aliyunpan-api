package openapi

import (
	"encoding/json"
	"fmt"
	"github.com/tickstep/library-go/logger"
	"strings"
)

type (
	PartInfoItem struct {
		// PartNumber 分片序列号，从 1 开始。单个文件分片最大限制5GB，最小限制100KB
		PartNumber int `json:"part_number"`
		// UploadUrl 分片上传URL地址
		UploadUrl int `json:"upload_url"`
		// PartSize 分片大小
		PartSize int64 `json:"part_size"`
	}

	// FileUploadCreateParam 文件创建参数
	FileUploadCreateParam struct {
		// DriveId 网盘ID
		DriveId string `json:"drive_id"`
		// ParentFileId 父目录id，上传到根目录时填写 root
		ParentFileId string `json:"parent_file_id"`
		// Name 文件名称，按照 utf8 编码最长 1024 字节，不能以 / 结尾
		Name string `json:"name"`
		// Type file | folder
		Type string `json:"type"`
		// CheckNameMode auto_rename 自动重命名，存在并发问题 ,refuse 同名不创建 ,ignore 同名文件可创建
		CheckNameMode string `json:"check_name_mode"`
		// Size 文件大小，单位为 byte。秒传必须
		Size int64 `json:"size"`
		// 最大分片数量 10000
		PartInfoList []*PartInfoItem `json:"part_info_list"`
		// PreHash 针对大文件sha1计算非常耗时的情况， 可以先在读取文件的前1k的sha1， 如果前1k的sha1没有匹配的， 那么说明文件无法做秒传， 如果1ksha1有匹配再计算文件sha1进行秒传，这样有效边避免无效的sha1计算。
		PreHash string `json:"pre_hash"`
		// ContentHash 文件内容 hash 值，需要根据 content_hash_name 指定的算法计算，当前都是sha1算法
		ContentHash string `json:"content_hash"`
		// ContentHashName 秒传必须 ,默认都是 sha1
		ContentHashName string `json:"content_hash_name"`
		// ProofCode 防伪码，秒传必须
		ProofCode string `json:"proof_code"`
		// ProofVersion 固定 v1
		ProofVersion string `json:"proof_version"`
		// LocalCreatedAt 本地创建时间，只对文件有效，格式yyyy-MM-dd'T'HH:mm:ss.SSS'Z'
		LocalCreatedAt string `json:"local_created_at"`
		// LocalModifiedAt 本地修改时间，只对文件有效，格式yyyy-MM-dd'T'HH:mm:ss.SSS'Z'
		LocalModifiedAt string `json:"local_modified_at"`
	}

	FileUploadCreateResult struct {
		// DriveId 网盘ID
		DriveId string `json:"drive_id"`
		// ParentFileId 父目录id，上传到根目录时填写 root
		ParentFileId string `json:"parent_file_id"`
		// FileId 文件ID
		FileId string `json:"file_id"`
		// FileName 文件名称
		FileName string `json:"file_name"`
		// Status
		Status string `json:"status"`
		// UploadId 创建文件夹返回空
		UploadId string `json:"upload_id"`
		// Available
		Available bool `json:"available"`
		// Exist 是否存在同名文件
		Exist bool `json:"exist"`
		// RapidUpload 是否能秒传
		RapidUpload bool `json:"rapid_upload"`
		// 最大分片数量 10000
		PartInfoList []*PartInfoItem `json:"part_info_list"`
	}
)

// FileUploadCreate 文件（文件夹）创建
func (a *AliPanClient) FileUploadCreate(param *FileUploadCreateParam) (*FileUploadCreateResult, *AliApiErrResult) {
	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s/adrive/v1.0/openFile/create", OPENAPI_URL)
	logger.Verboseln("do request url: " + fullUrl.String())

	// parameters
	postData := map[string]interface{}{
		"drive_id":        param.DriveId,
		"parent_file_id":  param.ParentFileId,
		"name":            param.Name,
		"type":            param.Type,
		"check_name_mode": param.CheckNameMode,
	}
	if len(param.CheckNameMode) == 0 {
		postData["check_name_mode"] = "auto_rename"
	}
	if strings.ToLower(param.Type) == "folder" {
		// 文件夹
	} else if strings.ToLower(param.Type) == "file" {
		// 文件
	}

	// request
	resp, err := a.httpclient.Req("POST", fullUrl.String(), postData, a.Headers())
	if err != nil {
		logger.Verboseln("file create error ", err)
		return nil, NewAliApiHttpError(err.Error())
	}

	// handler common error
	var body []byte
	var apiErrResult *AliApiErrResult
	if body, apiErrResult = ParseCommonOpenApiError(resp); apiErrResult != nil {
		return nil, apiErrResult
	}

	// parse result
	r := &FileUploadCreateResult{}
	if err2 := json.Unmarshal(body, r); err2 != nil {
		logger.Verboseln("parse file create result json error ", err2)
		return nil, NewAliApiAppError(err2.Error())
	}
	return r, nil
}
