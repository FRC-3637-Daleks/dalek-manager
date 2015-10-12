# Makefile for deploying to the RoboRIO arm processor

DEPLOY_DIR := dman
WEB_DIR := web

TOPLEVEL_PKG := dalek-manager
TOPLEVEL_FILE := $(TOPLEVEL_PKG).go

OUTPUT := dalek-manager
ARM_OUTPUT := $(OUTPUT)-arm

TEAMNO := 3637
TARGET_HOSTNAME := roboRIO-$(TEAMNO).local     # In the event of failure switch to IP address
TARGET_USER := lvuser
ROOT_USER := Admin
TARGET_ADDRESS := $(TARGET_USER)@$(TARGET_HOSTNAME)

# Makes it so that user can enter password
SSH_FLAGS := -tt

WEB_FILES := $(wildcard $(WEB_DIR)/dynamic/*.html) $(wildcard $(WEB_DIR)/static/css/*.css) $(wildcard $(WEB_DIR)/static/js/*.js)
START_SCRIPT := scripts/start_server.sh
STOP_SCRIPT := scripts/stop_server.sh
INSTALL_SCRIPT := scripts/install.sh
STARTUP_SCRIPT := scripts/dalek-manager-init.sh
SCRIPTS := $(STOP_SCRIPT) $(START_SCRIPT) $(STARTUP_SCRIPT)
DEPLOYED_FILES := $(WEB_FILES) $(SCRIPTS) $(ARM_OUTPUT)
DEPLOYED_TAR := $(OUTPUT).tar.gz

GOARGS := 

all: $(OUTPUT)

arm: $(ARM_OUTPUT)

install: $(OUTPUT)
	go install

$(OUTPUT): $(TOPLEVEL_FILE)
	go build $(GOARGS) $(TOPLEVEL_FILE)

$(ARM_OUTPUT): $(TOPLEVEL_FILE)
	env GOOS=linux GOARCH=arm GOARM=7 go build $(GOARGS) -o $(ARM_OUTPUT) $(TOPLEVEL_FILE)

deploy: $(ARM_OUTPUT) $(SCRIPTS) $(DEPLOYED_TAR)
	ssh  $(TARGET_ADDRESS) "/etc/init.d/$(OUTPUT) stop || true"
	sftp -oBatchMode=no -b sftp_deploy_file $(TARGET_ADDRESS)
	ssh $(SSH_FLAGS) $(TARGET_ADDRESS) "bash /home/$(TARGET_USER)/install.sh"
	ssh  $(TARGET_ADDRESS) "/etc/init.d/$(OUTPUT) start"

$(DEPLOYED_TAR): $(DEPLOYED_FILES)
	tar -czf $(DEPLOYED_TAR) $(WEB_DIR)/ $(SCRIPTS) $(ARM_OUTPUT)


clean:
	go clean
	rm $(DEPLOYED_TAR)

stop:
	ssh $(SSH_FLAGS) $(TARGET_ADDRESS) "/etc/init.d/$(OUTPUT) stop"

start:
	ssh  $(TARGET_ADDRESS) "/etc/init.d/$(OUTPUT) start"

uninstall-init:
	ssh $(SSH_FLAGS) $(ROOT_USER)@$(TARGET_HOSTNAME) "rm /etc/init.d/$(OUTPUT)"
