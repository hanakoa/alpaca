## Install / upgrade minikube
```bash
make mk-upgrade
```

## Start cluster
```bash
make mk-start
```

## Stop cluster
```bash
make mk-stop
```

## Build Docker images
```bash
make mk-build
```

## Rebuild Docker images
```bash
make mk-rebuild
```

## Start µServices
```bash
make kb-create
```

## View logs
```bash
kubectl logs $(kubectl get po -l app=alpaca-auth,tier=api -o jsonpath="{.items[0].metadata.name}") -f
```

## Stop µServices
```bash
make kb-delete
```
