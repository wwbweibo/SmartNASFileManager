build-docker:
	${HOME}/tools/flutter/bin/flutter build web
	docker build -f Dockerfile -t 192.168.163.65:5100/file-server-web . --push