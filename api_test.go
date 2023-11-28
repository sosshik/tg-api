package api

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"git.foxminded.ua/foxstudent106264/tgapi/mocks"
	"github.com/stretchr/testify/assert"
)

func TestSendTextToTelegramChat(t *testing.T) {
	token := "234khnc82347in29289324nkdsc"

	mockClient := &mocks.HTTPClientInterface{}
	app := &Api{
		SendMessageURL: "https://api.telegram.org/bot" + token + "/sendMessage",
		HTTPClient:     mockClient,
	}

	tests := []struct {
		name           string
		expectedURL    string
		expectedChatID int
		expectedText   string
		mockResponse   *http.Response
		mockError      error
		expectedResult string
		expectedError  bool
	}{
		{
			"Positive Test",
			app.SendMessageURL,
			123,
			"your_message",
			&http.Response{Body: io.NopCloser(strings.NewReader("message your_message successfuly distributed to chat id 123")), StatusCode: http.StatusOK},
			nil,
			"message your_message successfuly distributed to chat id 123",
			false,
		},
		{
			"Negative Test",
			app.SendMessageURL,
			123,
			"your_message",
			nil,
			errors.New("err"),
			"",
			true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			mockClient.On("PostForm", test.expectedURL, url.Values{
				"chat_id": {strconv.Itoa(test.expectedChatID)},
				"text":    {test.expectedText},
			}).Return(test.mockResponse, test.mockError)

			response, err := app.SendTextToTelegramChat(test.expectedChatID, test.expectedText)

			if test.expectedError {
				assert.Empty(t, response, "Response should be empty on error")
			} else {
				assert.NoError(t, err, "Unexpected error")
				assert.Equal(t, test.expectedResult, response, "Unexpected response")
			}

			mockClient.AssertExpectations(t)
		})
	}
}
