package yapi

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Client 封装 YAPI 所有 HTTP API 调用
type Client struct {
	baseURL    string
	token      string
	projectID  int
	httpClient *http.Client
}

// NewClient 创建 YAPI 客户端实例
func NewClient(baseURL, token string, projectID int, skipTLSVerify bool) *Client {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	if skipTLSVerify {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true} //nolint:gosec
	}
	return &Client{
		baseURL:   strings.TrimRight(baseURL, "/"),
		token:     token,
		projectID: projectID,
		httpClient: &http.Client{
			Timeout:   30 * time.Second,
			Transport: transport,
		},
	}
}

// GetProject 获取项目基本信息
// YAPI API: GET /api/project/get?token=xxx
func (c *Client) GetProject() (*ProjectInfo, error) {
	params := url.Values{"token": {c.token}}
	var resp Response[ProjectInfo]
	if err := c.doGet("/api/project/get", params, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// GetCategories 获取分类列表
// YAPI API: GET /api/interface/getCatMenu?project_id=xxx&token=xxx
func (c *Client) GetCategories() ([]CategoryInfo, error) {
	params := url.Values{
		"project_id": {strconv.Itoa(c.projectID)},
		"token":      {c.token},
	}
	var resp Response[[]CategoryInfo]
	if err := c.doGet("/api/interface/getCatMenu", params, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// GetMenuWithInterfaces 获取菜单树（分类 + 接口列表）
// YAPI API: GET /api/interface/list_menu?project_id=xxx&token=xxx
func (c *Client) GetMenuWithInterfaces() ([]MenuCategory, error) {
	params := url.Values{
		"project_id": {strconv.Itoa(c.projectID)},
		"token":      {c.token},
	}
	var resp Response[[]MenuCategory]
	if err := c.doGet("/api/interface/list_menu", params, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// GetInterfaceList 获取接口列表（分页）
// YAPI API: GET /api/interface/list?project_id=xxx&token=xxx&page=1&limit=20
func (c *Client) GetInterfaceList(page, limit int) (*InterfaceListResult, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	params := url.Values{
		"project_id": {strconv.Itoa(c.projectID)},
		"token":      {c.token},
		"page":       {strconv.Itoa(page)},
		"limit":      {strconv.Itoa(limit)},
	}
	var resp Response[InterfaceListResult]
	if err := c.doGet("/api/interface/list", params, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// GetInterfaceListByCat 获取分类下的接口列表
// YAPI API: GET /api/interface/list_cat?catid=xxx&token=xxx&page=1&limit=100
func (c *Client) GetInterfaceListByCat(catID, page, limit int) (*InterfaceListResult, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 100
	}
	params := url.Values{
		"catid": {strconv.Itoa(catID)},
		"token": {c.token},
		"page":  {strconv.Itoa(page)},
		"limit": {strconv.Itoa(limit)},
	}
	var resp Response[InterfaceListResult]
	if err := c.doGet("/api/interface/list_cat", params, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// GetInterfaceDetail 获取接口详情
// YAPI API: GET /api/interface/get?id=xxx&token=xxx
func (c *Client) GetInterfaceDetail(interfaceID int) (*InterfaceDetail, error) {
	params := url.Values{
		"id":    {strconv.Itoa(interfaceID)},
		"token": {c.token},
	}
	var resp Response[InterfaceDetail]
	if err := c.doGet("/api/interface/get", params, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// SaveInterface 新增或更新接口
// 有 ID → POST /api/interface/up; 无 ID → POST /api/interface/add
func (c *Client) SaveInterface(params SaveInterfaceParams) (*SaveInterfaceResult, error) {
	params.Token = c.token

	path := "/api/interface/add"
	if params.ID > 0 {
		path = "/api/interface/up"
	}

	var resp Response[SaveInterfaceResult]
	if err := c.doPost(path, params, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// doGet 执行 GET 请求
func (c *Client) doGet(path string, params url.Values, result any) error {
	reqURL := fmt.Sprintf("%s%s?%s", c.baseURL, path, params.Encode())
	resp, err := c.httpClient.Get(reqURL)
	if err != nil {
		return fmt.Errorf("HTTP 请求失败: %w", err)
	}
	defer resp.Body.Close()
	return c.decodeResponse(resp, result)
}

// doPost 执行 POST 请求（JSON body）
func (c *Client) doPost(path string, body any, result any) error {
	jsonData, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("JSON 序列化失败: %w", err)
	}

	reqURL := fmt.Sprintf("%s%s", c.baseURL, path)
	resp, err := c.httpClient.Post(reqURL, "application/json", bytes.NewReader(jsonData))
	if err != nil {
		return fmt.Errorf("HTTP 请求失败: %w", err)
	}
	defer resp.Body.Close()
	return c.decodeResponse(resp, result)
}

// decodeResponse 解析 HTTP 响应并检查 YAPI 错误码
func (c *Client) decodeResponse(resp *http.Response, result any) error {
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %w", err)
	}

	if err := json.Unmarshal(body, result); err != nil {
		return fmt.Errorf("JSON 解析失败: %w", err)
	}

	// 检查 YAPI 业务错误码（通过 raw JSON 读取 errcode）
	var errCheck struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}
	if err := json.Unmarshal(body, &errCheck); err == nil && errCheck.ErrCode != 0 {
		return fmt.Errorf("YAPI 错误 [%d]: %s", errCheck.ErrCode, errCheck.ErrMsg)
	}

	return nil
}
