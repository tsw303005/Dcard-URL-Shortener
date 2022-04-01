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
	const url = "fake_url"
	const expiredAt = "fake_expired_date"

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

		AfterEach(func() {
			Expect(resp.Result().Body.Close()).NotTo(HaveOccurred())
		})

		JustBeforeEach(func() {
			router.ServeHTTP(resp, req)
		})

		Context("success", func() {
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

			When("shorten url found with no error", func() {
				It("redirects url successfully", func() {
					Expect(resp.Code).To(Equal(message.URLRedirect))
					Expect(resp.Result().Header.Get("Location")).To(Equal("/fake_url"))
				})
			})
		})

		Context("fail", func() {
			BeforeEach(func() {
				req, err = http.NewRequestWithContext(ctx, "GET", "/test_get?shorten_url="+shortenURL, http.NoBody)
				Expect(err).NotTo(HaveOccurred())
			})

			JustBeforeEach(func() {
				err = json.Unmarshal(resp.Body.Bytes(), &respValue)
				Expect(err).NotTo(HaveOccurred())
				Expect(respValue["request"]).To(Equal("redirect url request"))
			})

			When("date expired", func() {
				BeforeEach(func() {
					shortenerDAO.EXPECT().Get(req.Context(), &dao.Shortener{
						ShortenURL: shortenURL,
					}).Return(nil, dao.ErrExpiredat)
				})

				It("returns the url expired error", func() {
					Expect(respValue["error"]).To(Equal("this shorten url has already expired"))
					Expect(resp.Code).To(Equal(message.URLExpired))
				})
			})

			When("shorten url not found", func() {
				BeforeEach(func() {
					shortenerDAO.EXPECT().Get(req.Context(), &dao.Shortener{
						ShortenURL: shortenURL,
					}).Return(nil, dao.ErrShortenURLNotFound)
				})

				It("returns the shorten url not found error", func() {
					Expect(respValue["error"]).To(Equal("this shorten url not found"))
					Expect(resp.Code).To(Equal(message.BadRequest))
				})
			})

			When("internal server error", func() {
				BeforeEach(func() {
					shortenerDAO.EXPECT().Get(req.Context(), &dao.Shortener{
						ShortenURL: shortenURL,
					}).Return(nil, errors.New("internal server error"))
				})

				It("returns the internal server error", func() {
					Expect(respValue["error"]).To(Equal("internal server error"))
					Expect(resp.Code).To(Equal(message.ServiceError))
				})
			})
		})
	})

	Describe("shorten url", func() {
		var (
			req       *http.Request
			resp      *httptest.ResponseRecorder
			id        uuid.UUID
			jsonData  []byte
			respValue map[string]string
			err       error
		)

		BeforeEach(func() {
			resp = httptest.NewRecorder()
		})

		AfterEach(func() {
			Expect(resp.Result().Body.Close()).NotTo(HaveOccurred())
		})

		JustBeforeEach(func() {
			router.ServeHTTP(resp, req)
		})

		Context("success", func() {
			BeforeEach(func() {
				id = uuid.New()
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

			JustBeforeEach(func() {
				err = json.Unmarshal(resp.Body.Bytes(), &respValue)
				Expect(err).ToNot(HaveOccurred())
			})

			When("shorten url", func() {
				It("returns shorten url and id with no error", func() {
					Expect(resp.Code).To(Equal(message.SuccessRequest))
					Expect(respValue["id"]).To(Equal(id.String()))
					Expect(respValue["shortenURL"]).To(Equal(shortenURL))
				})
			})
		})

		Context("fail", func() {
			JustBeforeEach(func() {
				err = json.Unmarshal(resp.Body.Bytes(), &respValue)
				Expect(err).ToNot(HaveOccurred())
				Expect(respValue["request"]).To(Equal("shorten url request"))
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
					Expect(resp.Code).To(Equal(message.BadRequest))
					Expect(respValue["error"]).To(Equal("argument error"))
				})
			})

			When("expiredAt argument has already expired", func() {
				BeforeEach(func() {
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

				It("returns the expired error", func() {
					Expect(resp.Code).To(Equal(message.BadRequest))
					Expect(respValue["error"]).To(Equal("expiredAt has already expired"))
				})
			})

			When("internal server error", func() {
				BeforeEach(func() {
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
					Expect(resp.Code).To(Equal(message.ServiceError))
					Expect(respValue["error"]).To(Equal("internal server error"))
				})
			})
		})
	})
})
