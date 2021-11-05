#pid=$(netstat -nap | grep 8081 | tail -n1 | awk '{printf("%d/n"), $7}' | awk -F/ '{printf("%d\n"), $1}')
echo "[BUILD] build start"
mkdir build  2>/dev/null
echo "[BUILD] building webhook"
rm -rf build/*
go build -o build/webhook cmd/webhook/main.go
mkdir build/config 2>/dev/null
cp -r config/certificate build/config/certificate
cp -r config/config  build/config/config
cp -r config/crd build/config/crd

echo "[BUILD] building images"

docker build . -t webhook:latest 
docker tag webhook:latest tangjiaming1999/casbin-kubesphere-authz:v1

echo "[BUILD] pushing images"
docker push tangjiaming1999/casbin-kubesphere-authz:v1



# #install all crd
# make install
# #load existing models and policies
# cd deployments
# python3 load_crd.py
# #the webhook deployment uses this account
# kubectl create serviceaccount my-sa # if necessary, namespace should be added 
# kubectl apply -f webhook_register_internal_step1.yaml
# kubectl apply -f webhook_register_internal_step2.yaml
# #grant necessary access authority to the webhook
# kubectl create clusterrolebinding my-sa-view \
#   --clusterrole=cluster-admin \
#   --serviceaccount=default:my-sa 
