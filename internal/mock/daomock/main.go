package daomock

//go:generate mockgen -destination=mock.go -package=$GOPACKAGE github.com/tsw303005/Dcard-URL-Shortener/internal/dao ShortenerDAO
