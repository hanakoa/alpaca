# alpaca
[![forthebadge](http://forthebadge.com/images/badges/built-with-love.svg)](http://forthebadge.com)
[![GoDoc](https://godoc.org/github.com/hanakoa/alpaca?status.svg)](https://godoc.org/github.com/hanakoa/alpaca)
[![Go report](http://goreportcard.com/badge/hanakoa/alpaca)](http://goreportcard.com/report/hanakoa/alpaca)
[![CircleCI](https://circleci.com/gh/hanakoa/alpaca.svg?style=svg)](https://circleci.com/gh/hanakoa/alpaca)

<p align="center">
 <figure>
  <img src="https://image.flaticon.com/icons/svg/371/371645.svg" alt="Alpaca" width="304" height="228">
  <figcaption>
  <div>
  </div>
  </figcaption>
</figure> 
</p>

## Intro
Alpaca is a WIP microservices system that handles authentication and authorization.

It provides:
- login with email, username, or password
- two-factor auth with phone call, SMS, or Yubikey
- password resets
- email confirmation codes
- phone number confirmation codes

## FAQ
- [What stack do we use?](./docs/stack.md)
- [Why does this project exist?](./docs/differences.md)
- [How do I run locally?](./docs/running-locally.md)
- [How do I run in Docker?](./docs/running-with-docker.md)
- [How do I run in minikube?](./docs/running-with-minikube.md)
- [How do I use the REST API?](./docs/using-rest-api.md)

## Notes on Contributing
- `make protoc` will regenerate Go code from Protocol Buffers
- `make lint` [lints](https://github.com/golang/lint) code
- `make fmt` [formats](https://golang.org/cmd/gofmt/) code
- `make vet` [vets](https://golang.org/cmd/vet/) code