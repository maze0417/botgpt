package chatgpt

import (
	"botgpt/internal/clients/chatgpt"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"

	http "github.com/bogdanfinn/fhttp"
	"io"
	"strings"
)

type ChatgptController struct {
}

func NewChatgptController() *ChatgptController {
	return &ChatgptController{}
}

//goland:noinspection GoUnhandledErrorResult
func init() {
	//go func() {
	//	ticker := time.NewTicker(time.Minute)
	//	for {
	//		select {
	//		case <-ticker.C:
	//			req, _ := http.NewRequest(http.MethodGet, heartBeatUrl, nil)
	//			req.Header.Set("User-Agent", userAgent)
	//			chatgpt.Client.Do(req)
	//		}
	//	}
	//}()
}

//goland:noinspection GoUnhandledErrorResult
func (h *ChatgptController) GetConversations(c *gin.Context) {
	offset, ok := c.GetQuery("offset")
	if !ok {
		offset = "0"
	}
	limit, ok := c.GetQuery("limit")
	if !ok {
		limit = "20"
	}
	handleGet(c, apiPrefix+"/conversations?offset="+offset+"&limit="+limit, getConversationsErrorMessage)
}

//goland:noinspection GoUnhandledErrorResult
func (h *ChatgptController) CreateConversation(c *gin.Context) {
	var request CreateConversationRequest
	if err := c.BindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, chatgpt.ReturnMessage(parseJsonErrorMessage))
		return
	}

	if request.ConversationID == nil || *request.ConversationID == "" {
		request.ConversationID = nil
	}
	if request.Messages[0].Author.Role == "" {
		request.Messages[0].Author.Role = defaultRole
	}
	if request.VariantPurpose == "" {
		request.VariantPurpose = "none"
	}

	jsonBytes, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", apiPrefix+"/conversation", bytes.NewBuffer(jsonBytes))
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Authorization", chatgpt.GetAccessToken(c.GetHeader(chatgpt.AuthorizationHeader)))
	req.Header.Set("Accept", "text/event-stream")
	resp, err := chatgpt.Client.Do(req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, chatgpt.ReturnMessage(err.Error()))
		return
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK || err != nil {
		switch resp.StatusCode {
		case http.StatusUnauthorized:
			c.AbortWithStatusJSON(resp.StatusCode, chatgpt.ReturnMessage(conversationErrorMessage401))
			return
		case http.StatusForbidden:
			c.AbortWithStatusJSON(resp.StatusCode, chatgpt.ReturnMessage(conversationErrorMessage403))
			return
		case http.StatusNotFound:
			c.AbortWithStatusJSON(resp.StatusCode, chatgpt.ReturnMessage(conversationErrorMessage404))
			return
		case http.StatusRequestEntityTooLarge:
			c.AbortWithStatusJSON(resp.StatusCode, chatgpt.ReturnMessage(conversationErrorMessage413))
			return
		case http.StatusUnprocessableEntity:
			c.AbortWithStatusJSON(resp.StatusCode, chatgpt.ReturnMessage(conversationErrorMessage422))
			return
		case http.StatusTooManyRequests:
			c.AbortWithStatusJSON(resp.StatusCode, chatgpt.ReturnMessage(conversationErrorMessage429))
			return
		case http.StatusInternalServerError:
			c.AbortWithStatusJSON(resp.StatusCode, chatgpt.ReturnMessage(conversationErrorMessage500))
			return
		}
	}

	chatgpt.HandleConversationResponse(c, resp)
}

//goland:noinspection GoUnhandledErrorResult
func (h *ChatgptController) GenerateTitle(c *gin.Context) {
	var request GenerateTitleRequest
	if err := c.BindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, chatgpt.ReturnMessage(parseJsonErrorMessage))
		return
	}

	jsonBytes, _ := json.Marshal(request)
	handlePost(c, apiPrefix+"/conversation/gen_title/"+c.Param("id"), string(jsonBytes), generateTitleErrorMessage)
}

//goland:noinspection GoUnhandledErrorResult
func (h *ChatgptController) GetConversation(c *gin.Context) {
	handleGet(c, apiPrefix+"/conversation/"+c.Param("id"), getContentErrorMessage)
}

//goland:noinspection GoUnhandledErrorResult
func (h *ChatgptController) UpdateConversation(c *gin.Context) {
	var request PatchConversationRequest
	if err := c.BindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, chatgpt.ReturnMessage(parseJsonErrorMessage))
		return
	}

	// bool default to false, then will hide (delete) the conversation
	if request.Title != nil {
		request.IsVisible = true
	}
	jsonBytes, _ := json.Marshal(request)
	handlePatch(c, apiPrefix+"/conversation/"+c.Param("id"), string(jsonBytes), updateConversationErrorMessage)
}

//goland:noinspection GoUnhandledErrorResult
func (h *ChatgptController) FeedbackMessage(c *gin.Context) {
	var request FeedbackMessageRequest
	if err := c.BindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, chatgpt.ReturnMessage(parseJsonErrorMessage))
		return
	}

	jsonBytes, _ := json.Marshal(request)
	handlePost(c, apiPrefix+"/conversation/message_feedback", string(jsonBytes), feedbackMessageErrorMessage)
}

//goland:noinspection GoUnhandledErrorResult
func (h *ChatgptController) ClearConversations(c *gin.Context) {
	jsonBytes, _ := json.Marshal(PatchConversationRequest{
		IsVisible: false,
	})
	handlePatch(c, apiPrefix+"/conversations", string(jsonBytes), clearConversationsErrorMessage)
}

//goland:noinspection GoUnhandledErrorResult
func (h *ChatgptController) GetModels(c *gin.Context) {
	handleGet(c, apiPrefix+"/models", getModelsErrorMessage)
}

func (h *ChatgptController) GetAccountCheck(c *gin.Context) {
	handleGet(c, apiPrefix+"/accounts/check", getAccountCheckErrorMessage)
}

//goland:noinspection GoUnhandledErrorResult
func handleGet(c *gin.Context, url string, errorMessage string) {
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Authorization", chatgpt.GetAccessToken(c.GetHeader(chatgpt.AuthorizationHeader)))
	resp, err := chatgpt.Client.Do(req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, chatgpt.ReturnMessage(err.Error()))
		return
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK || err != nil {
		c.AbortWithStatusJSON(resp.StatusCode, chatgpt.ReturnMessage(errorMessage))
		return
	}

	io.Copy(c.Writer, resp.Body)
}

//goland:noinspection GoUnhandledErrorResult
func handlePost(c *gin.Context, url string, requestBody string, errorMessage string) {
	req, _ := http.NewRequest(http.MethodPost, url, strings.NewReader(requestBody))
	handlePostOrPatch(c, req, errorMessage)
}

//goland:noinspection GoUnhandledErrorResult
func handlePatch(c *gin.Context, url string, requestBody string, errorMessage string) {
	req, _ := http.NewRequest(http.MethodPatch, url, strings.NewReader(requestBody))
	handlePostOrPatch(c, req, errorMessage)
}

//goland:noinspection GoUnhandledErrorResult
func handlePostOrPatch(c *gin.Context, req *http.Request, errorMessage string) {
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Authorization", chatgpt.GetAccessToken(c.GetHeader(chatgpt.AuthorizationHeader)))
	resp, err := chatgpt.Client.Do(req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, chatgpt.ReturnMessage(err.Error()))
		return
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK || err != nil {
		c.AbortWithStatusJSON(resp.StatusCode, chatgpt.ReturnMessage(errorMessage))
		return
	}

	io.Copy(c.Writer, resp.Body)
}

//goland:noinspection GoUnhandledErrorResult
func (h *ChatgptController) UserLogin(c *gin.Context) {
	var loginInfo LoginInfo
	if err := c.ShouldBindJSON(&loginInfo); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, chatgpt.ReturnMessage(parseUserInfoErrorMessage))
		return
	}

	// get csrf token
	req, _ := http.NewRequest(http.MethodGet, csrfUrl, nil)
	req.Header.Set("User-Agent", userAgent)
	resp, err := chatgpt.Client.Do(req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, chatgpt.ReturnMessage(err.Error()))
		return
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK || err != nil {
		c.AbortWithStatusJSON(resp.StatusCode, chatgpt.ReturnMessage(getCsrfTokenErrorMessage))
		return
	}

	data, _ := io.ReadAll(resp.Body)
	responseMap := make(map[string]string)
	json.Unmarshal(data, &responseMap)

	// get authorized url
	params := fmt.Sprintf(
		"callbackUrl=/&csrfToken=%s&json=true",
		responseMap["csrfToken"],
	)
	req, err = http.NewRequest(http.MethodPost, promptLoginUrl, strings.NewReader(params))
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("User-Agent", userAgent)
	resp, err = chatgpt.Client.Do(req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, chatgpt.ReturnMessage(err.Error()))
		return
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK || err != nil {
		c.AbortWithStatusJSON(resp.StatusCode, chatgpt.ReturnMessage(getAuthorizedUrlErrorMessage))
		return
	}

	// get state
	data, _ = io.ReadAll(resp.Body)
	json.Unmarshal(data, &responseMap)
	req, err = http.NewRequest(http.MethodGet, responseMap["url"], nil)
	resp, err = chatgpt.Client.Do(req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, chatgpt.ReturnMessage(err.Error()))
		return
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK || err != nil {
		c.AbortWithStatusJSON(resp.StatusCode, chatgpt.ReturnMessage(getStateErrorMessage))
		return
	}

	// check username
	doc, _ := goquery.NewDocumentFromReader(resp.Body)
	state, _ := doc.Find("input[name=state]").Attr("value")
	params = fmt.Sprintf(
		"state=%s&username=%s&js-available=true&webauthn-available=true&is-brave=false&webauthn-platform-available=false&action=default",
		state,
		loginInfo.Username,
	)
	req, err = http.NewRequest(http.MethodPost, loginUsernameUrl+state, strings.NewReader(params))
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("User-Agent", userAgent)
	resp, err = chatgpt.Client.Do(req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, chatgpt.ReturnMessage(err.Error()))
		return
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK || err != nil {
		c.AbortWithStatusJSON(resp.StatusCode, chatgpt.ReturnMessage(emailInvalidErrorMessage))
		return
	}

	// check username and password
	params = fmt.Sprintf(
		"state=%s&username=%s&password=%s&action=default",
		state,
		loginInfo.Username,
		loginInfo.Password,
	)
	req, err = http.NewRequest(http.MethodPost, loginPasswordUrl+state, strings.NewReader(params))
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("User-Agent", userAgent)
	resp, err = chatgpt.Client.Do(req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, chatgpt.ReturnMessage(err.Error()))
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK || err != nil {
		c.AbortWithStatusJSON(resp.StatusCode, chatgpt.ReturnMessage(emailOrPasswordInvalidErrorMessage))
		return
	}

	// get access token
	req, err = http.NewRequest(http.MethodGet, authSessionUrl, nil)
	req.Header.Set("User-Agent", userAgent)
	resp, err = chatgpt.Client.Do(req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, chatgpt.ReturnMessage(err.Error()))
		return
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK || err != nil {
		c.AbortWithStatusJSON(resp.StatusCode, chatgpt.ReturnMessage(getAccessTokenErrorMessage))
		return
	}

	io.Copy(c.Writer, resp.Body)
}
