### Written in Golang
- The language of choice for containerized environments
([Docker](https://github.com/moby/moby) and [Kubernetes](https://github.com/kubernetes/kubernetes) are both written in Go).
- Lighter memory footprint than our Spring Boot apps (see this [gist](https://gist.github.com/kevinmichaelchen/22ac37452979b05f78e99f775e249659)).
- [More performant](https://benchmarksgame.alioth.debian.org/u64q/go.html).
- Way faster build / startup times than Spring Boot.
- More maintainable / less boilerplate / fewer lines of code.
- DB connection retries (Spring will fail if DB isn't up).
- Able to configure frequency of scheduled tasks, unlike Spring's [`@Scheduled(fixedDelay = 500)`](https://docs.spring.io/spring-framework/docs/current/javadoc-api/org/springframework/scheduling/annotation/Scheduled.html).
- Async functions feel natural. No more creating new files just to invoke 
  [an `@Async`-annotated method](https://stackoverflow.com/questions/24898547/spring-async-method-called-from-another-async-method).

#### Lightweight
- Alpine Postgres DB uses 8MB (no data).
- Golang server (e.g., `auth`) takes 10MB (under no load).
- Alpine RabbitMQ uses 80MB (under no load). Not much I can do about that.
- So I could easily run 20 Golang REST APIs and still only use 440MB (well under a Raspberry Pi's 1GB RAM).

#### Insanely fast build times!
```bash
$ time go build -o ./bin/alpaca-auth ./auth
1.52s user 0.59s system 130% cpu 1.618 total
```

Running feels instantaneous. Doing a `./gradlew bootRun` on a Spring Boot app, on the other hand, takes nearly 15 seconds.

#### Fewer lines of code!
`tokei` reports 3293 LOC.
In Java, the lines of code could easily be 4 times this.

### Microservices, not monolith!
Here's the [original RFC](https://gist.github.com/hanakoa/6614b0799f09144ef348b3cb9a871820) 
that about decomposing the original monolith.

Microservices lead to 
- independent development (easier to onboard people, less surface area for stuff to go wrong)
- independent deployment (no need to bring everything down to update one thing)
- independent scaling (e.g., the CPU-intensive password hashing service need not be tethered to everything else)

### Database Changes
- Postgres instead of MySQL.
  - There is no official Alpine MySQL Docker image as of Jan 2018.
- [Snowflake PKs](https://developer.twitter.com/en/docs/basics/twitter-ids).
  - 8-bytes instead of 16-byte UUID PKs. (see PR #5).
  - Snowflakes are harder to guess, but not unguessable, like Tweet IDs.
  - Where we need unguessability, such as with any of our (reset, 2FA, or confirmation) codes, we use v4 UUIDs.
- [Cursor pagination](https://developer.twitter.com/en/docs/basics/cursoring).
  - [Offset-based pagination doesn't scale for large offsets](http://use-the-index-luke.com/no-offset).
- Better varchar constraints: 50 for names (see Facebook), 25 for usernames (compromise between Github's 39 and Twitter's 15).

### Security Updates
- Salt is stored on Password, not Person, per [OWASP](https://www.owasp.org/index.php/Password_Storage_Cheat_Sheet#Use_a_cryptographically_strong_credential-specific_salt).
  - "Generate a unique salt upon creation of each stored credential (not just per user or system wide)"
- Dropping [LUDS](https://www.usenix.org/conference/usenixsecurity16/technical-sessions/presentation/wheeler)
in favor of password complexity, with Dropbox's [zxcvbn](https://blogs.dropbox.com/tech/2012/04/zxcvbn-realistic-password-strength-estimation/).
- Passwords and people no longer expire.
- Self-calibrating iteration count. App will determine how many password hash 
iterations it must perform such that hashing takes roughly a second, or some other given value.

### Style Updates
- UI relies on [Create-React-App](https://github.com/facebook/create-react-app).
- UI uses [Material UI](http://www.material-ui.com/#/) instead of Bootstrap.
- Outbound emails use [Hermes](https://github.com/matcornic/hermes) for templating.

### Nomenclature Changes
- "Multi-factor" instead of "two-factor". 
- "Claims" instead of "roles".

### FUTURE
- Security -- Backup codes
- Security -- "New device", based on IP address
- Security -- support [YubiKeys](https://www.yubico.com/)
- Security -- support [Authy](https://authy.com/), [Google Authenticator](https://en.wikipedia.org/wiki/Google_Authenticator)
- Database -- look into [CockroachDB](https://github.com/cockroachdb/cockroach)
- Ability to merge accounts