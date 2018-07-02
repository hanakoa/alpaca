# Multi-factor Authentication
When a user logs in, if they're account has enabled
two-factor auth, the primary auth µService calls out
to this service and tells it to send a 6-digit 2FA code,
passing along the user's `personId` and `phoneNumber`. 