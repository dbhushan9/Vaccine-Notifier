
#  docker run -d --name vaccine_alerts dbhushan9/vaccine-alerts-go:v0.1
# docker update --restart=always 0576df221c0b
#docker run -v "D:\workspace\test\shared":"/app/shared" --name vaccine_alerts --restart=always dbhushan9/vaccine-alerts-go:v2.0-rc-1

build-local:
	docker build -f Dockerfile.dev .

build:
	docker buildx build --platform linux/amd64,linux/arm64,linux/arm/v7 -t dbhushan9/vaccine-alerts-go:v2.0-rc-1 -f Dockerfile.dev --push .

run
	docker run -v "":"/app/shared" --name vaccine_alerts --restart=always dbhushan9/vaccine-tracker-go:v2.0-rc-1