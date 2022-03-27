package service

import (
	"context"
	"net/http"
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
		router       *gin.Engine
		w            *httptest.ResponseRecorder
		svc          *Service
		req          *http.Request
		shortener    *dao.Shortener
		ctx          context.Context
	)

	const shortenURL = "fake shorten url"

	BeforeEach(func() {
		controller = gomock.NewController(GinkgoT())
		w = httptest.NewRecorder()
		svc = NewService(shortenerDAO)
		shortenerDAO = daomock.NewMockShortenerDAO(controller)
		router = gin.Default()
	})

	AfterEach(func() {
		controller.Finish()
	})

	Describe("GetURL", func() {
		BeforeEach(func() {
			req, _ := http.NewRequest()
		})

		JustBeforeEach(func() {
			router.ServeHTTP(w, req)
		})

		When("success", func() {
			BeforeEach(func() {
				shortener = &dao.Shortener{
					URL: "fake url",
				}

				shortenerDAO.EXPECT().Get(ctx, &dao.Shortener{
					ShortenURL: ginContext.Query("shorten_url"),
				}).Return(shortener, nil)
			})

			It("redircts with no error", func() {
				Expect(ginContext.Writer.Status()).To(Equal(message.URLRedirect))
				Expect(ginContext.FullPath()).To(Equal("fake url"))
			})
		})

	})
})
