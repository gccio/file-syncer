package caiyun

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
)

const defaultCatalog string = "00019700101000000001"

type diskRequest struct {
	CatalogID         string            `json:"catalogID"`
	SortDirection     int               `json:"sortDirection"`
	StartNumber       int               `json:"startNumber"`
	EndNumber         int               `json:"endNumber"`
	FilterType        int               `json:"filterType"`
	CatalogSortType   int               `json:"catalogSortType"`
	ContentSortType   int               `json:"contentSortType"`
	CommonAccountInfo commonAccountInfo `json:"commonAccountInfo"`
}

func (o *diskRequest) ToJSON() []byte {
	bs, _ := json.Marshal(o)
	return bs
}

type commonAccountInfo struct {
	Account     string `json:"account"`
	AccountType int    `json:"accountType"`
}

type diskResponse struct {
	Success bool   `json:"success"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Result struct {
			ResultCode string `json:"resultCode"`
			ResultDesc any    `json:"resultDesc"`
		} `json:"result"`
		GetDiskResult *DiskResult `json:"getDiskResult"`
	} `json:"data"`
}

type DiskResult struct {
	ParentCatalogID string        `json:"parentCatalogID"`
	NodeCount       int           `json:"nodeCount"`
	CatalogList     []CatalogItem `json:"catalogList"`
	ContentList     []ContentItem `json:"contentList"`
	IsCompleted     int           `json:"isCompleted"`
}

type CatalogItem struct {
	CatalogID       string `json:"catalogID"`
	CatalogName     string `json:"catalogName"`
	CatalogType     int    `json:"catalogType"`
	CreateTime      string `json:"createTime"`
	UpdateTime      string `json:"updateTime"`
	IsShared        bool   `json:"isShared"`
	CatalogLevel    int    `json:"catalogLevel"`
	ShareDoneeCount int    `json:"shareDoneeCount"`
	OpenType        int    `json:"openType"`
	ParentCatalogID string `json:"parentCatalogId"`
	DirEtag         int    `json:"dirEtag"`
	Tombstoned      int    `json:"tombstoned"`
	ProxyID         any    `json:"proxyID"`
	Moved           int    `json:"moved"`
	IsFixedDir      int    `json:"isFixedDir"`
	IsSynced        any    `json:"isSynced"`
	Owner           string `json:"owner"`
	Modifier        string `json:"modifier"`
	Path            string `json:"path"`
	ShareType       int    `json:"shareType"`
	SoftLink        any    `json:"softLink"`
	ExtProp1        any    `json:"extProp1"`
	ExtProp2        any    `json:"extProp2"`
	ExtProp3        any    `json:"extProp3"`
	ExtProp4        any    `json:"extProp4"`
	ExtProp5        any    `json:"extProp5"`
	ETagOprType     int    `json:"ETagOprType"`
}

type ContentItem struct {
	ContentID       string `json:"contentID"`
	ContentName     string `json:"contentName"`
	ContentSuffix   string `json:"contentSuffix"`
	ContentSize     int64  `json:"contentSize"`
	ContentDesc     string `json:"contentDesc"`
	ContentType     int    `json:"contentType"`
	ContentOrigin   int    `json:"contentOrigin"`
	UpdateTime      string `json:"updateTime"`
	CommentCount    int    `json:"commentCount"`
	ThumbnailURL    string `json:"thumbnailURL"`
	BigthumbnailURL string `json:"bigthumbnailURL"`
	PresentURL      string `json:"presentURL"`
	PresentLURL     string `json:"presentLURL"`
	PresentHURL     string `json:"presentHURL"`
	ContentTAGList  any    `json:"contentTAGList"`
	ShareDoneeCount int    `json:"shareDoneeCount"`
	Safestate       int    `json:"safestate"`
	Transferstate   int    `json:"transferstate"`
	IsFocusContent  int    `json:"isFocusContent"`
	UpdateShareTime any    `json:"updateShareTime"`
	UploadTime      string `json:"uploadTime"`
	OpenType        int    `json:"openType"`
	AuditResult     int    `json:"auditResult"`
	ParentCatalogID string `json:"parentCatalogId"`
	Channel         string `json:"channel"`
	GeoLocFlag      string `json:"geoLocFlag"`
	Digest          string `json:"digest"`
	Version         string `json:"version"`
	FileEtag        string `json:"fileEtag"`
	FileVersion     string `json:"fileVersion"`
	Tombstoned      int    `json:"tombstoned"`
	ProxyID         string `json:"proxyID"`
	Moved           int    `json:"moved"`
	MidthumbnailURL string `json:"midthumbnailURL"`
	Owner           string `json:"owner"`
	Modifier        string `json:"modifier"`
	ShareType       int    `json:"shareType"`
	ExtInfo         struct {
		Uploader string `json:"uploader"`
	} `json:"extInfo"`
	Exif struct {
		CreateTime    string `json:"createTime"`
		Longitude     any    `json:"longitude"`
		Latitude      any    `json:"latitude"`
		LocalSaveTime any    `json:"localSaveTime"`
	} `json:"exif"`
	CollectionFlag any  `json:"collectionFlag"`
	TreeInfo       any  `json:"treeInfo"`
	IsShared       bool `json:"isShared"`
	ETagOprType    int  `json:"ETagOprType"`
}

func (c *caiYun) GetDisk(catalogID string) (*DiskResult, error) {
	if catalogID == "" {
		catalogID = defaultCatalog
	}

	var (
		catalogList = make([]CatalogItem, 0)
		contentList = make([]ContentItem, 0)
		result      = &DiskResult{}
		err         error
	)

	isCompleted := 0
	for i := 0; isCompleted == 0; i++ {
		result, err = sendRequest(c.cli, catalogID, 1+i*100)
		if err != nil {
			return nil, err
		}

		catalogList = append(catalogList, result.CatalogList...)
		contentList = append(contentList, result.ContentList...)
		isCompleted = result.IsCompleted
	}

	result.CatalogList = catalogList
	result.ContentList = contentList
	return result, nil
}

func sendRequest(cli *http.Client, path string, start int) (*DiskResult, error) {
	var (
		req  *http.Request
		resp *http.Response
		err  error
	)
	reqBody := &diskRequest{
		CatalogID:       path,
		SortDirection:   1,
		StartNumber:     start,
		EndNumber:       start + 100,
		FilterType:      0,
		CatalogSortType: 0,
		ContentSortType: 0,
		CommonAccountInfo: commonAccountInfo{
			Account:     os.Getenv("CAIYUN_TELEPHONE"),
			AccountType: 1,
		},
	}
	req, err = NewRequest(
		"/orchestration/personalCloud/catalog/v1.0/getDisk",
		http.MethodPost,
		reqBody.ToJSON(),
	)
	if err != nil {
		return nil, err
	}

	resp, err = cli.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bs := make([]byte, 0)
	bs, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	obj := &diskResponse{}
	if err := json.Unmarshal(bs, obj); err != nil {
		return nil, err
	}

	return obj.Data.GetDiskResult, nil
}
