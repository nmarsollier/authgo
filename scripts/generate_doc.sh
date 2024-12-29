set -e

swag fmt
swag init --parseDependency --parseInternal
npx swagger-markdown -i ./docs/swagger.yaml -o README-API.md
