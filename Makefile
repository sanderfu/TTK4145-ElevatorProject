GO_CMD = go
GO_BUILD = $(GO_CMD) build

ELEVATOR_FILE_NAME = elevator.go
WATCHDOG_FILE_NAME = watchdog.go
START_ELEVATOR_FILE_NAME = startElevator.go

ELEVATOR_PROJECT_NAME = elevator
WATCHDOG_PROJECT_NAME = watchdog
START_ELEVATOR_PROJECT_NAME = startElevator


ELEVATOR_FILES = $(shell find cmd -name "$(ELEVATOR_FILE_NAME)")
WATCHDOG_FILES = $(shell find cmd -name "$(WATCHDOG_FILE_NAME)")
START_ELEVATOR_FILES = $(shell find cmd -name "$(START_ELEVATOR_FILE_NAME)")

BUILD_DIR = build
ASSET_DIR = assets
MKDIR_BUILD = mkdir -p $(BUILD_DIR)

all: elevator watchdog startelevator

elevator:
	@$(MKDIR_BUILD)
	$(GO_BUILD) -o $(BUILD_DIR)/$(ELEVATOR_PROJECT_NAME) $(ELEVATOR_FILES)

watchdog:
	@$(MKDIR_BUILD)
	$(GO_BUILD) -o $(BUILD_DIR)/$(WATCHDOG_PROJECT_NAME) $(WATCHDOG_FILES)

startelevator:
	@$(MKDIR_BUILD)
	$(GO_BUILD) -o $(BUILD_DIR)/$(START_ELEVATOR_PROJECT_NAME) $(START_ELEVATOR_FILES)

run:
	$(BUILD_DIR)/$(START_ELEVATOR_PROJECT_NAME)

.PHONY: clean

clean:
	rm -rf $(BUILD_DIR)
	rm -rf $(ASSET_DIR)