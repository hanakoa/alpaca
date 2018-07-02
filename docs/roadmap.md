- [ ] vgo #17
- [ ] set jwt cookie on successful response

### better packaging
right now, we they're ~16mb
```bash
ls -la bin
total 91928
-rwxr-xr-x   1 kevinchen  staff    16M Jul  2 10:18 alpaca-auth*
-rwxr-xr-x   1 kevinchen  staff    13M Jul  2 10:18 alpaca-mfa*
-rwxr-xr-x   1 kevinchen  staff    16M Jul  2 10:18 alpaca-password-reset*
```

if we package smarter, we might be able to reduce this size