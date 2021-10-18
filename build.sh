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
echo "[BUILD] building controller"
make generate
go build -o build/controller cmd/controller/main.go

echo "[BUILD] building images"
docker build -f Dockerfile -t controller:latest .
docker build -f Dockerfile_internal_webhook -t webhook:latest .
docker tag webhook:latest tangjiaming1999/casbin-kubesphere-authz:v1

echo "[BUILD] pushing images"
docker push tangjiaming1999/casbin-kubesphere-authz:v1

