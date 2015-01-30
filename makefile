BLDDIR = build
BINARIES = oda-toriko

all: $(BINARIES)

$(BLDDIR)/%:
	go get .
	go build -o $(BLDDIR)/oda-toriko

$(BINARIES): %: $(BLDDIR)/%

clean: 
	rm -rf $(BLDDIR)


.PHONY: all
.PHONY: $(BINARIES)
