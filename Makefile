
#  docker run -d --name vaccine_alerts dbhushan9/vaccine-tracker-go:v0.1

build:
	docker buildx build --platform linux/amd64,linux/arm64,linux/arm/v7 -t dbhushan9/vaccine-alerts-go:v0.2 -f Dockerfile.dev --push .


