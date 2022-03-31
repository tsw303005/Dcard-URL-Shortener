package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
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
			respValue map[string]string
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

				req, err = http.NewRequestWithContext(ctx, "GET", "/test_get?shorten_url="+shortenURL, http.NoBody)
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
				req, err = http.NewRequestWithContext(ctx, "GET", "/test_get?shorten_url="+shortenURL, http.NoBody)
				Expect(err).NotTo(HaveOccurred())

				shortenerDAO.EXPECT().Get(req.Context(), &dao.Shortener{
					ShortenURL: shortenURL,
				}).Return(nil, dao.ErrExpiredat)
			})

			It("returns the url expired error", func() {
				err = json.Unmarshal(resp.Body.Bytes(), &respValue)

				Expect(err).NotTo(HaveOccurred())
				Expect(respValue["error"]).To(Equal("this shorten url has already expired"))
				Expect(respValue["request"]).To(Equal("redirect url request"))
				Expect(resp.Code).To(Equal(message.URLExpired))
				Expect(resp.Result().Body.Close()).NotTo(HaveOccurred())
			})
		})

		When("shorten url not found", func() {
			BeforeEach(func() {
				req, err = http.NewRequestWithContext(ctx, "GET", "/test_get?shorten_url="+shortenURL, http.NoBody)
				Expect(err).NotTo(HaveOccurred())

				shortenerDAO.EXPECT().Get(req.Context(), &dao.Shortener{
					ShortenURL: shortenURL,
				}).Return(nil, dao.ErrShortenURLNotFound)
			})

			It("returns the url not found error", func() {
				err = json.Unmarshal(resp.Body.Bytes(), &respValue)

				Expect(err).NotTo(HaveOccurred())
				Expect(respValue["error"]).To(Equal("this shorten url not found"))
				Expect(respValue["request"]).To(Equal("redirect url request"))
				Expect(resp.Code).To(Equal(message.BadRequest))
				Expect(resp.Result().Body.Close()).NotTo(HaveOccurred())
			})
		})

		When("internal server error", func() {
			BeforeEach(func() {
				req, err = http.NewRequestWithContext(ctx, "GET", "/test_get?shorten_url="+shortenURL, http.NoBody)
				Expect(err).NotTo(HaveOccurred())

				shortenerDAO.EXPECT().Get(req.Context(), &dao.Shortener{
					ShortenURL: shortenURL,
				}).Return(nil, errors.New("internal server error"))
			})

			It("returns the url not found error", func() {
				err = json.Unmarshal(resp.Body.Bytes(), &respValue)

				Expect(err).NotTo(HaveOccurred())
				Expect(respValue["error"]).To(Equal("internal server error"))
				Expect(respValue["request"]).To(Equal("redirect url request"))
				Expect(resp.Code).To(Equal(message.ServiceError))
				Expect(resp.Result().Body.Close()).NotTo(HaveOccurred())
			})
		})
	})

	Describe("shorten url", func() {
		var (
			req       *http.Request
			resp      *httptest.ResponseRecorder
			id        uuid.UUID
			url       string
			expiredAt string
			jsonData  []byte
			respValue map[string]string
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
				id = uuid.New()
				url = "fake_url"
				expiredAt = "fake_expired_date"
				jsonData = []byte(`{
					"url": "fake_url",
					"expiredAt": "fake_expired_date"
				}`)

				req, err = http.NewRequestWithContext(ctx, "POST", "/test_shorten", bytes.NewBuffer(jsonData))
				Expect(err).NotTo(HaveOccurred())

				req.Header.Set("Content-Type", "application/json")

				shortenerDAO.EXPECT().Shorten(req.Context(), &dao.Shortener{
					URL:       url,
					ExpiredAt: expiredAt,
				}).Return(id, shortenURL, nil)
			})

			It("shorten url with no error", func() {
				err = json.Unmarshal(resp.Body.Bytes(), &respValue)

				Expect(err).ToNot(HaveOccurred())
				Expect(resp.Code).To(Equal(message.SuccessRequest))
				Expect(respValue["id"]).To(Equal(id.String()))
				Expect(respValue["shortenURL"]).To(Equal(shortenURL))
				Expect(resp.Result().Body.Close()).NotTo(HaveOccurred())
			})
		})

		When("argument error", func() {
			BeforeEach(func() {
				jsonData = []byte(`{
					"url": "fake_url",
				}`)

				req, err = http.NewRequestWithContext(ctx, "POST", "/test_shorten", bytes.NewBuffer(jsonData))
				Expect(err).NotTo(HaveOccurred())

				req.Header.Set("Content-Type", "application/json")
			})

			It("returns the argument error", func() {
				err = json.Unmarshal(resp.Body.Bytes(), &respValue)

				Expect(err).ToNot(HaveOccurred())
				Expect(resp.Code).To(Equal(message.BadRequest))
				Expect(respValue["error"]).To(Equal("argument error"))
				Expect(respValue["request"]).To(Equal("shorten url request"))
				Expect(resp.Result().Body.Close()).NotTo(HaveOccurred())
			})
		})

		When("expriedAt argument has already expired", func() {
			BeforeEach(func() {
				url = "fake_url"
				expiredAt = "fake_expired_date"
				jsonData = []byte(`{
					"url": "fake_url",
					"expiredAt": "fake_expired_date"
				}`)

				req, err = http.NewRequestWithContext(ctx, "POST", "/test_shorten", bytes.NewBuffer(jsonData))
				Expect(err).NotTo(HaveOccurred())

				req.Header.Set("Content-Type", "application/json")

				shortenerDAO.EXPECT().Shorten(req.Context(), &dao.Shortener{
					URL:       url,
					ExpiredAt: expiredAt,
				}).Return(uuid.Nil, "", dao.ErrShortenURLFail)
			})

			It("returns the error that the argument expiredAt has already expired error", func() {
				err = json.Unmarshal(resp.Body.Bytes(), &respValue)

				Expect(err).ToNot(HaveOccurred())
				Expect(resp.Code).To(Equal(message.BadRequest))
				Expect(respValue["error"]).To(Equal("expiredAt has already expired"))
				Expect(respValue["request"]).To(Equal("shorten url request"))
				Expect(resp.Result().Body.Close()).NotTo(HaveOccurred())
			})
		})

		When("internal server error", func() {
			BeforeEach(func() {
				url = "fake_url"
				expiredAt = "fake_expired_date"
				jsonData = []byte(`{
					"url": "fake_url",
					"expiredAt": "fake_expired_date"
				}`)

				req, err = http.NewRequestWithContext(ctx, "POST", "/test_shorten", bytes.NewBuffer(jsonData))
				Expect(err).NotTo(HaveOccurred())

				req.Header.Set("Content-Type", "application/json")

				shortenerDAO.EXPECT().Shorten(req.Context(), &dao.Shortener{
					URL:       url,
					ExpiredAt: expiredAt,
				}).Return(uuid.Nil, "", errors.New("internal server error"))
			})

			It("returns the internal server error", func() {
				err = json.Unmarshal(resp.Body.Bytes(), &respValue)

				Expect(err).ToNot(HaveOccurred())
				Expect(resp.Code).To(Equal(message.ServiceError))
				Expect(respValue["error"]).To(Equal("internal server error"))
				Expect(respValue["request"]).To(Equal("shorten url request"))
				Expect(resp.Result().Body.Close()).NotTo(HaveOccurred())
			})
		})
	})
})
