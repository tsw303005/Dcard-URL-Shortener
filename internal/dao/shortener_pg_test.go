package dao

import (
	"context"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"github.com/onsi/gomega/types"
)

var _ = Describe("PGShortenerDAO", func() {
	var shortenerDAO *pgShortenerDAO
	var ctx context.Context

	BeforeEach(func() {
		shortenerDAO = NewPGShortenerDAO(pgClient)
		ctx = context.Background()
	})

	Describe("Shorten", func() {
		var (
			req        *Shortener
			id         uuid.UUID
			shortenURL string
			err        error
		)

		BeforeEach(func() {
			req = &Shortener{}
			req.ExpiredAt = "2037-04-08T09:20:41Z"
			req.URL = "fake_url"
		})

		JustBeforeEach(func() {
			id, shortenURL, err = shortenerDAO.Shorten(ctx, req)
		})

		When("success", func() {
			AfterEach(func() {
				deleteShortener(id)
			})

			It("returns id and shorten url with no error", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(id).NotTo(Equal(uuid.Nil))
				Expect(shortenURL).To(Equal("http://localhost/" + id.String()))
			})
		})

		When("fail", func() {
			BeforeEach(func() { req.ExpiredAt = expiredAt })

			It("returns the shorten url fail", func() {
				Expect(err).To(MatchError(ErrShortenURLFail))
				Expect(id).To(Equal(uuid.Nil))
				Expect(shortenURL).To(BeEmpty())
			})
		})
	})

	Describe("Get", func() {
		var (
			req       *Shortener
			shortener *Shortener
			err       error
		)

		BeforeEach(func() {
			req = NewFakeShortener("fake_url")
		})

		JustBeforeEach(func() {
			shortener, err = shortenerDAO.Get(ctx, req)
		})

		When("success", func() {
			BeforeEach(func() {
				insertShortener(req)
			})

			AfterEach(func() {
				deleteShortener(req.ID)
			})

			It("returns shortener successfully", func() {
				Expect(err).NotTo(HaveOccurred())
				matchShortener(shortener)
			})
		})

		When("expiredAt expires", func() {
			BeforeEach(func() {
				req.ExpiredAt = expiredAt
				insertShortener(req)
			})

			AfterEach(func() {
				deleteShortener(req.ID)
			})

			It("returns expiredAt error", func() {
				Expect(err).To(MatchError(ErrExpiredat))
				Expect(shortener).To(BeNil())
			})
		})

		When("shorten url not found", func() {
			It("returns ErrShortenURLNotFound error", func() {
				Expect(err).To(MatchError(ErrShortenURLNotFound))
				Expect(shortener).To(BeNil())
			})
		})
	})
})

func insertShortener(shortener *Shortener) {
	query := "INSERT INTO shorteners (id, url, shorten_url, expired_at) VALUES (?, ?, ?, ?);"

	pgExec(query, shortener.ID, shortener.URL, shortener.ShortenURL, shortener.ExpiredAt)
}

func deleteShortener(id uuid.UUID) {
	query := "DELETE FROM shorteners WHERE id = ?;"

	pgExec(query, id)
}

func matchShortener(shortener *Shortener) types.GomegaMatcher {
	return PointTo(MatchFields(IgnoreExtras, Fields{
		"ID":         Equal(shortener.ID),
		"URL":        Equal(shortener.URL),
		"ShortenURL": Equal(shortener.ShortenURL),
		"ExpiredAt":  Equal(shortener.ExpiredAt),
	}))
}
