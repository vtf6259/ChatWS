.PHONY: all backend frontend
all: backend frontend
frontend:
	./buildFrontend.sh
backend:
	./buildBackend.sh
clean:
	rm -f backend/ChatWS
	rm -f build/ChatWS
	rm -rf frontend/dist
	rm -rf build/dist