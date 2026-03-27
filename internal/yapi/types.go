package yapi

import "encoding/json"

// Response YAPI 通用响应包装
type Response[T any] struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	Data    T      `json:"data"`
}

// ProjectInfo 项目基本信息
type ProjectInfo struct {
	ID          int    `json:"_id"`
	Name        string `json:"name"`
	Description string `json:"desc"`
	BasePath    string `json:"basepath"`
	GroupID     int    `json:"group_id"`
	GroupName   string `json:"group_name"`
	ProjectType string `json:"project_type"`
}

// CategoryInfo 接口分类信息
type CategoryInfo struct {
	ID          int    `json:"_id"`
	Name        string `json:"name"`
	Description string `json:"desc"`
	ProjectID   int    `json:"project_id"`
	AddTime     int64  `json:"add_time"`
	UpTime      int64  `json:"up_time"`
}

// InterfaceSummary 接口摘要（列表中使用）
type InterfaceSummary struct {
	ID        int    `json:"_id"`
	Title     string `json:"title"`
	Path      string `json:"path"`
	Method    string `json:"method"`
	ProjectID int    `json:"project_id"`
	CatID     int    `json:"catid"`
	Status    string `json:"status"`
	Tag       []string `json:"tag"`
}

// InterfaceDetail 接口完整详情
type InterfaceDetail struct {
	ID                  int               `json:"_id"`
	ProjectID           int               `json:"project_id"`
	CatID               int               `json:"catid"`
	Title               string            `json:"title"`
	Path                string            `json:"path"`
	Method              string            `json:"method"`
	Status              string            `json:"status"`
	Description         string            `json:"desc"`
	Markdown            string            `json:"markdown"`
	Tag                 []string          `json:"tag"`
	ReqBodyType         string            `json:"req_body_type"`
	ReqBodyIsJSONSchema bool              `json:"req_body_is_json_schema"`
	ReqBodyOther        string            `json:"req_body_other"`
	ReqBodyForm         []ReqBodyFormItem `json:"req_body_form"`
	ReqParams           []ReqParam        `json:"req_params"`
	ReqHeaders          []ReqHeader       `json:"req_headers"`
	ReqQuery            []ReqQuery        `json:"req_query"`
	ResBodyType         string            `json:"res_body_type"`
	ResBodyIsJSONSchema bool              `json:"res_body_is_json_schema"`
	ResBody             string            `json:"res_body"`
	AddTime             int64             `json:"add_time"`
	UpTime              int64             `json:"up_time"`
}

// ReqBodyFormItem 表单请求体字段
type ReqBodyFormItem struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Example  string `json:"example"`
	Desc     string `json:"desc"`
	Required string `json:"required"`
}

// ReqParam 路径参数
type ReqParam struct {
	Name    string `json:"name"`
	Example string `json:"example"`
	Desc    string `json:"desc"`
}

// ReqHeader 请求头
type ReqHeader struct {
	Name     string `json:"name"`
	Value    string `json:"value"`
	Example  string `json:"example"`
	Desc     string `json:"desc"`
	Required string `json:"required"`
}

// ReqQuery 查询参数
type ReqQuery struct {
	Name     string `json:"name"`
	Example  string `json:"example"`
	Desc     string `json:"desc"`
	Required string `json:"required"`
}

// MenuCategory 菜单分类（含分类下的接口列表）
type MenuCategory struct {
	ID          int                `json:"_id"`
	Name        string             `json:"name"`
	Description string             `json:"desc"`
	ProjectID   int                `json:"project_id"`
	List        []InterfaceSummary `json:"list"`
}

// InterfaceListResult 接口列表分页结果
type InterfaceListResult struct {
	Count int                `json:"count"`
	Total int                `json:"total"`
	List  []InterfaceSummary `json:"list"`
}

// SaveInterfaceParams 新增/更新接口参数
type SaveInterfaceParams struct {
	ID                  int               `json:"id,omitempty"`
	ProjectID           int               `json:"project_id,omitempty"`
	Token               string            `json:"token"`
	CatID               int               `json:"catid"`
	Title               string            `json:"title"`
	Path                string            `json:"path"`
	Method              string            `json:"method"`
	Status              string            `json:"status,omitempty"`
	Description         string            `json:"desc,omitempty"`
	Markdown            string            `json:"markdown,omitempty"`
	Tag                 []string          `json:"tag,omitempty"`
	ReqBodyType         string            `json:"req_body_type,omitempty"`
	ReqBodyIsJSONSchema bool              `json:"req_body_is_json_schema,omitempty"`
	ReqBodyOther        string            `json:"req_body_other,omitempty"`
	ReqBodyForm         []ReqBodyFormItem `json:"req_body_form,omitempty"`
	ReqParams           []ReqParam        `json:"req_params,omitempty"`
	ReqHeaders          []ReqHeader       `json:"req_headers,omitempty"`
	ReqQuery            []ReqQuery        `json:"req_query,omitempty"`
	ResBodyType         string            `json:"res_body_type,omitempty"`
	ResBodyIsJSONSchema bool              `json:"res_body_is_json_schema,omitempty"`
	ResBody             string            `json:"res_body,omitempty"`
	SwitchNotice        bool              `json:"switch_notice,omitempty"`
}

// SaveInterfaceResult 保存接口的返回结果
type SaveInterfaceResult struct {
	ID json.Number `json:"_id"`
}

// AddCategoryParams 新增接口分类参数
type AddCategoryParams struct {
	Name        string `json:"name"`
	Description string `json:"desc"`
	ProjectID   int    `json:"project_id"`
	Token       string `json:"token"`
}

// AddCategoryResult 新增分类返回结果
type AddCategoryResult struct {
	ID          int    `json:"_id"`
	Name        string `json:"name"`
	Description string `json:"desc"`
	ProjectID   int    `json:"project_id"`
}

// ImportDataParams 数据导入参数
type ImportDataParams struct {
	Type  string `json:"type"`
	JSON  string `json:"json"`
	Merge string `json:"merge"`
	Token string `json:"token"`
	URL   string `json:"url,omitempty"`
}

// ImportDataResult 数据导入返回结果
type ImportDataResult struct {
	SuccessMessage string `json:"successMessage,omitempty"`
}
