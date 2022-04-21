# Dcard-URL-Shortener
![workflow checking](https://github.com/tsw303005/Dcard-URL-Shortener/actions/workflows/main.yaml/badge.svg)

## About The Project
This is Dcard's homework about implementing url shortener. This repository includes all the modules with the mono-repo architecture. Also, if you want to know more about this homework detail, you can refer to [this spec](https://drive.google.com/file/d/1AreBiHDUYXH6MI5OqWpKP-f6-W0zA8np/view).

<p align="right">(<a href="#top">back to top</a>)</p>

## Built With
The following are the packages used in this Golang project.

- [gin-gonic](https://github.com/gin-gonic/gin): as framework
- [go-pg](https://github.com/go-pg/pg): to connect Postgres DB and store url information.
- [go-redis](https://github.com/go-redis/redis): to enhance the performance of redirecting url request.
- [go-migrate](https://github.com/golang-migrate/migrate): to create db table.
- [google-uuid](https://pkg.go.dev/github.com/google/uuid): to create a new uuid as a shorten url.
- [golang/mock](https://github.com/golang/mock), [ginkgo](https://github.com/onsi/ginkgo) and [gomega](https://github.com/onsi/gomega): to do the unit test
- [golangci](https://github.com/golangci/golangci-lint): to check the style

<p align="right">(<a href="#top">back to top</a>)</p>

## Getting Started
Before starting this program, please make sure that your docker is running.

### Unit Testing
This project has unit tests with [ginkgo](https://github.com/onsi/ginkgo) framework for DAO and toolkit in pkg directory. 

1. Test whole project
    ```
    make dc.test
    ```
2. Test pkg
    ```
    make dc.pkg.test
    ```

### Style Check
This project uses [golangci](https://github.com/golangci/golangci-lint) to check the style.
1. Check whole project's style
    ```
    make dc.lint
    ```
2. Check only modules' style
    ```
    make dc.internal.lint
    ```
3. Check only pkg's style
    ```
    make dc.pkg.lint
    ```

### Build Image
1. Build shorten url api image
    ```
    make dc.image
    ```

<p align="right">(<a href="#top">back to top</a>)</p>

## Future Work
1. This time, I only use [google-uuid](https://pkg.go.dev/github.com/google/uuid) to simply generate shorten url. The next goal is to adopt a suitable shorten url algorithm.
2. This implementation is kind of like image-to-url converter. So I want to extend this project to contain image-to-url converter in the future.

<p align="right">(<a href="#top">back to top</a>)</p>
