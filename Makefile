.PHONY: all backend frontend
all: backend frontend
frontend:
	./buildFrontend.sh
backend:
	./buildBackend.sh
clean:
	rm backend/ChatWS
	rm -r build/ChatWS
	rm -r frontend/dist
	rm -r build/dist