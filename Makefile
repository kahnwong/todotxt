shell:
	flox activate --start-services
backend-logs:
	flox services logs -f backend

build-frontend:
	cd frontend && yarn install && yarn build

test:
	hurl hurl/today.hurl
