
#  docker run -d --name vaccine_alerts dbhushan9/vaccine-tracker-go:v0.1
# docker update --restart=always 0576df221c0b

build-local:
	docker build -f Dockerfile.dev .


build:
	docker buildx build --platform linux/amd64,linux/arm64,linux/arm/v7 -t dbhushan9/vaccine-alerts-go:v0.3 -f Dockerfile.dev --push .


