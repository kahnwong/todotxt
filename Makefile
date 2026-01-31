start-frontend:
	cd frontend && yarn install && yarn dev
start-backend:
	go run .

build-frontend:
	cd frontend && yarn build
