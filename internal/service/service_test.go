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
		svc          *Service
		ctx          context.Context
	)

	const shortenURL = "fake_shorten_url"

	BeforeEach(func() {
		ctx = context.Background()
		controller = gomock.NewController(GinkgoT())
		shortenerDAO = daomock.NewMockShortenerDAO(controller)
		svc = NewService(shortenerDAO)
		router = gin.Default()

		router.GET("/test_get", func(c *gin.Context) {
			svc.GetURL(c)
		})

		router.POST("/test_shorten", func(c *gin.Context) {
			svc.ShortenURL(c)
		})
	})

	AfterEach(func() {
		controller.Finish()
	})

	Describe("GetURL", func() {
		var (
			req       *http.Request
			resp      *httptest.ResponseRecorder
			shortener *dao.Shortener
			err       error
		)

		BeforeEach(func() {
			resp = httptest.NewRecorder()
		})

		JustBeforeEach(func() {
			router.ServeHTTP(resp, req)
		})

		When("success", func() {
			BeforeEach(func() {
				shortener = &dao.Shortener{
					URL: "fake_url",
				}

				req, err = http.NewRequestWithContext(ctx, "GET", "/test_get?shorten_url="+shortenURL, nil)
				Expect(err).NotTo(HaveOccurred())

				shortenerDAO.EXPECT().Get(req.Context(), &dao.Shortener{
					ShortenURL: shortenURL,
				}).Return(shortener, nil)
			})

			It("redircts with no error", func() {
				Expect(resp.Code).To(Equal(message.URLRedirect))
				Expect(resp.Result().Header.Get("Location")).To(Equal("/fake_url"))
				Expect(resp.Result().Body.Close()).NotTo(HaveOccurred())
			})
		})

		When("date expired", func() {
			BeforeEach(func() {
				req, err = http.NewRequestWithContext(ctx, "GET", "/test_get?shorten_url="+shortenURL, nil)
				Expect(err).NotTo(HaveOccurred())

				shortenerDAO.EXPECT().Get(req.Context(), &dao.Shortener{
					ShortenURL: shortenURL,
				}).Return(nil, dao.ErrExpiredat)
			})

			It("returns the error", func() {
				Expect(resp.Code).To((Equal(message.URLExpired)))
				Expect(resp.Result().Body.Close()).NotTo(HaveOccurred())
			})
		})
	})
})
