package service

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/tsw303005/Dcard-URL-Shortener/internal/dao"
	"github.com/tsw303005/Dcard-URL-Shortener/internal/message"
	"github.com/tsw303005/Dcard-URL-Shortener/internal/mock/daomock"
)

func TestService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Test Service")
}

var _ = Describe("Service", func() {
	var (
		controller   *gomock.Controller
		shortenerDAO *daomock.MockShortenerDAO
		ginContext   *gin.Context
		svc          *Service
		shortener    *dao.Shortener
	)

	const shortenURL = "fake shorten url"

	BeforeEach(func() {
		controller = gomock.NewController(GinkgoT())
		svc = NewService(shortenerDAO)
		shortenerDAO = daomock.NewMockShortenerDAO(controller)
		ginContext, _ = gin.CreateTestContext(httptest.NewRecorder())

		ginContext.Params = []gin.Param{
			{
				Key:   "shorten_url",
				Value: shortenURL,
			},
		}
	})

	AfterEach(func() {
		controller.Finish()
	})

	Describe("GetURL", func() {
		JustBeforeEach(func() {
			svc.GetURL(ginContext)
		})

		When("success", func() {
			BeforeEach(func() {
				shortener = &dao.Shortener{
					URL: "fake url",
				}

				shortenerDAO.EXPECT().Get(ginContext.Request.Context(), &dao.Shortener{
					ShortenURL: ginContext.Query("shorten_url"),
				}).Return(shortener, nil)
			})

			It("redircts with no error", func() {
				Expect(ginContext.Writer.Status()).To(Equal(message.URLRedirect))
				Expect(ginContext.Request.Header)
			})
		})

	})
})
