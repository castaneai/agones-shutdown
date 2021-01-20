MINIKUBE_PROFILE := agones-shutdown

run:
	skaffold run --minikube-profile $(MINIKUBE_PROFILE) --tail

up:
	minikube start -p $(MINIKUBE_PROFILE) --cpus=3 --memory=2500mb --kubernetes-version=v1.17.13 --vm-driver virtualbox
	minikube profile $(MINIKUBE_PROFILE)
	helm repo add agones https://agones.dev/chart/stable
	helm repo update
	helm upgrade --install my-release --namespace agones-system --create-namespace agones/agones

down:
	minikube delete -p $(MINIKUBE_PROFILE)
