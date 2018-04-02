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
