# Backend workspace scripts (like package.json scripts). Run from backend/.
.PHONY: help check dev-ms-auth dev-api-gateway run-ms-auth run-api-gateway

help:
	@echo "WiChat backend — common targets:"
	@echo "  make check              lint + test (wi-shared, ms-auth, api-gateway)"
	@echo "  make dev-ms-auth        ms-auth with Air live reload"
	@echo "  make dev-api-gateway    api-gateway with Air live reload"
	@echo "  make run-ms-auth        ms-auth once (go run, no reload)"
	@echo "  make run-api-gateway    api-gateway once (go run, no reload)"
	@echo ""
	@echo "Per module: cd ms-auth|api-gateway|wi-shared && make dev|run|check"

check:
	$(MAKE) -C wi-shared check
	$(MAKE) -C ms-auth check
	$(MAKE) -C api-gateway check

dev-ms-auth:
	$(MAKE) -C ms-auth dev

dev-api-gateway:
	$(MAKE) -C api-gateway dev

run-ms-auth:
	$(MAKE) -C ms-auth run

run-api-gateway:
	$(MAKE) -C api-gateway run
