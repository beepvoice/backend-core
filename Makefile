# See https://git.makerforce.io/beep/best-practices/wiki/Makefile


#
# Programs
#
GOCMD?=go
GORUN?=$(GOCMD) run
GOBUILD?=$(GOCMD) build
GOCLEAN?=$(GOCMD) clean
GOTEST?=$(GOCMD) test
GOFMT_PROG?=gofmt
DOCKERCOMPOSE?=docker-compose

#
# Local
#
BINARY_NAME=core

#
# Files
#
DOCKERCOMPOSE_INTEGRATION_CONFIG?=docker-compose.integration.yml


#
# Tasks
#


# Let's do a quick unit test and then build backend-core
all: test_fmt test_unit build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v

test: test_fmt test_unit test_integration

test_fmt:
	$(GOFMT_PROG) -l .

test_unit:
	$(GOTEST) -tags=unit -v -cover

test_integration: test_integration_prepare
	$(GOTEST) -tags=integration -v -cover
test_integration_prepare:
	$(GORUN) scripts/testutils.go isrunning || $(DOCKERCOMPOSE) -f $(DOCKERCOMPOSE_INTEGRATION_CONFIG) up -d
	$(GORUN) scripts/testutils.go wait
test_integration_sql_shell:
	$(DOCKERCOMPOSE) -f $(DOCKERCOMPOSE_INTEGRATION_CONFIG) exec pg psql -d core
test_integration_cleanup:
	$(DOCKERCOMPOSE) -f $(DOCKERCOMPOSE_INTEGRATION_CONFIG) down

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
