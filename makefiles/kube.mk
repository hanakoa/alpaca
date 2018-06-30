.PHONY: kb-create
kb-create:
	kubectl create -f ./kube/

.PHONY: kb-delete
kb-delete:
	@kubectl delete svc,deploy --selector="app=alpaca-auth"
