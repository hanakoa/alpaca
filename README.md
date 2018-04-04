# alpaca
[![forthebadge](http://forthebadge.com/images/badges/built-with-love.svg)](http://forthebadge.com)

<p align="center">
 <figure>
  <img src="https://image.flaticon.com/icons/svg/371/371645.svg" alt="Alpaca" width="304" height="228">
  <figcaption>
  <div>
  </div>
  </figcaption>
</figure> 
</p>

Alpaca is a WIP microservices system that handles authentication and authorization.

It provides:
- login with email, username, or password
- two-factor auth with phone call, SMS, or Yubikey
- password resets
- email confirmation codes
- phone number confirmation codes

This project is built on Golang microservices,
which communicate via [gRPC](https://grpc.io/)
and [RabbitMQ](https://www.rabbitmq.com/).

The frontend uses:
- [React](https://reactjs.org/),
- [Redux](https://redux.js.org/),
- [Material UI](https://www.material-ui.com/#/),
- [redux-saga](https://redux-saga.js.org/),
- [redux-form](https://redux-form.com/7.3.0/)

<img width="488" alt="alpaca-login-screen" src="https://user-images.githubusercontent.com/5129994/38286303-b6f8d120-3792-11e8-8ca7-313459e99d90.png">

[Why does this project exist?](./docs/differences.md)

## Getting started
[Running locally](./docs/running-locally.md) is discouraged.
You should instead run with Docker.

To build images and spin them up, run
```bash
make docker
```

To bring it all down, run
```bash
make docker-stop
```
