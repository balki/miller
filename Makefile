PREFIX=/usr/local
INSTALLDIR=$(PREFIX)/bin

# ================================================================
# General-use targets

# This must remain the first target in this file, which is what 'make' with no
# arguments will run.
build:
	go build github.com/johnkerl/miller/cmd/mlr
	@echo "Build complete. The Miller executable is ./mlr (or .\mlr.exe on Windows)."
	@echo "You can use 'make check' to run tests".

# For interactive use, 'mlr regtest' offers more options and transparency.
check: unit_test regression_test
	@echo "Tests complete. You can use 'make install' if you like, optionally preceded"
	@echo "by './configure --prefix=/your/install/path' if you wish to install to"
	@echo "somewhere other than /usr/local/bin -- the default prefix is /usr/local."

# DESTDIR is for package installs; nominally blank when this is run interactively.
# See also https://www.gnu.org/prep/standards/html_node/DESTDIR.html
install: build
	cp mlr $(DESTDIR)/$(INSTALLDIR)
	make -C man install

# ================================================================
# Dev targets

# ----------------------------------------------------------------
# Unit tests (small number)
unit-test ut:
	go test github.com/johnkerl/miller/internal/pkg/...

# Keystroke-savers
lib-unbackslash-test:
	go test internal/pkg/lib/unbackslash_test.go internal/pkg/lib/unbackslash.go
lib_regex_test:
	go test internal/pkg/lib/regex_test.go internal/pkg/lib/regex.go
lib-tests:
	go test github.com/johnkerl/miller/internal/pkg/lib/...

mlrval-new-test:
	go test internal/pkg/mlrval/new_test.go \
	  internal/pkg/mlrval/mlrval_type.go \
	  internal/pkg/mlrval/mlrval_constants.go \
	  internal/pkg/mlrval/mlrval_new.go \
	  internal/pkg/mlrval/mlrval_infer.go
mlrval-is-test:
	go test internal/pkg/mlrval/is_test.go \
	  internal/pkg/mlrval/mlrval_type.go \
	  internal/pkg/mlrval/mlrval_constants.go \
	  internal/pkg/mlrval/mlrval_new.go \
	  internal/pkg/mlrval/mlrval_infer.go \
	  internal/pkg/mlrval/mlrval_is.go
mlrval-get-test:
	go test internal/pkg/mlrval/get_test.go \
	  internal/pkg/mlrval/mlrval_type.go \
	  internal/pkg/mlrval/mlrval_constants.go \
	  internal/pkg/mlrval/mlrval_new.go \
	  internal/pkg/mlrval/mlrval_infer.go \
	  internal/pkg/mlrval/mlrval_is.go \
	  internal/pkg/mlrval/mlrval_get.go
mlrval-output-test:
	go test internal/pkg/mlrval/output_test.go \
	  internal/pkg/mlrval/mlrval_type.go \
	  internal/pkg/mlrval/mlrval_constants.go \
	  internal/pkg/mlrval/mlrval_new.go \
	  internal/pkg/mlrval/mlrval_infer.go \
	  internal/pkg/mlrval/mlrval_is.go \
	  internal/pkg/mlrval/mlrval_get.go \
	  internal/pkg/mlrval/mlrval_output.go \
	  internal/pkg/mlrval/mlrval_format.go
mlrval-format-test:
	go test internal/pkg/mlrval/format_test.go \
	  internal/pkg/mlrval/mlrval_type.go \
	  internal/pkg/mlrval/mlrval_constants.go \
	  internal/pkg/mlrval/mlrval_new.go \
	  internal/pkg/mlrval/mlrval_infer.go \
	  internal/pkg/mlrval/mlrval_is.go \
	  internal/pkg/mlrval/mlrval_get.go \
	  internal/pkg/mlrval/mlrval_output.go \
	  internal/pkg/mlrval/mlrval_format.go
mlrval-tests:
	go test github.com/johnkerl/miller/internal/pkg/mlrval/...

mlrmap-new-test:
	go test internal/pkg/mlrval/mlrmap_new_test.go \
	  internal/pkg/mlrval/mlrmap.go \
	  internal/pkg/mlrval/mlrval_type.go \
	  internal/pkg/mlrval/mlrval_constants.go \
	  internal/pkg/mlrval/mlrval_new.go \
	  internal/pkg/mlrval/mlrval_infer.go \
	  internal/pkg/mlrval/mlrval_is.go \
	  internal/pkg/mlrval/mlrval_get.go \
	  internal/pkg/mlrval/mlrval_output.go \
	  internal/pkg/mlrval/mlrval_format.go
mlrmap-accessors-test:
	go test internal/pkg/mlrval/mlrmap_accessors_test.go \
	  internal/pkg/mlrval/mlrmap.go \
	  internal/pkg/mlrval/mlrmap_accessors.go \
	  internal/pkg/mlrval/mlrval_type.go \
	  internal/pkg/mlrval/mlrval_constants.go \
	  internal/pkg/mlrval/mlrval_new.go \
	  internal/pkg/mlrval/mlrval_cmp.go \
	  internal/pkg/mlrval/mlrval_copy.go \
	  internal/pkg/mlrval/mlrval_infer.go \
	  internal/pkg/mlrval/mlrval_is.go \
	  internal/pkg/mlrval/mlrval_get.go \
	  internal/pkg/mlrval/mlrval_output.go \
	  internal/pkg/mlrval/mlrval_format.go

mlrmap-tests: mlrmap-new-test mlrmap-accessors-test

input-dkvp-test:
	go test internal/pkg/input/record_reader_dkvp_test.go \
	  internal/pkg/input/record_reader.go \
	  internal/pkg/input/record_reader_dkvp_nidx.go
input-tests: input-dkvp-test

bifs-arithmetic-test:
	go test internal/pkg/bifs/arithmetic_test.go \
	  internal/pkg/bifs/base.go \
	  internal/pkg/bifs/arithmetic.go
bifs-bits-test:
	go test internal/pkg/bifs/bits_test.go \
	  internal/pkg/bifs/base.go \
	  internal/pkg/bifs/arithmetic.go \
	  internal/pkg/bifs/bits.go
bifs-collections-test:
	go test internal/pkg/bifs/collections_test.go \
	  internal/pkg/bifs/base.go \
	  internal/pkg/bifs/arithmetic.go \
	  internal/pkg/bifs/collections.go
bifs-hashing-test:
	go test internal/pkg/bifs/hashing_test.go \
	  internal/pkg/bifs/base.go \
	  internal/pkg/bifs/arithmetic.go \
	  internal/pkg/bifs/hashing.go
bifs-sort-test:
	go test internal/pkg/bifs/sort_test.go \
	  internal/pkg/bifs/base.go \
	  internal/pkg/bifs/arithmetic.go \
	  internal/pkg/bifs/sort.go

bifs-tests: bifs-arithmetic-test bifs-bits-test bifs-collections-test bifs-hashing-test bifs-sort-test

#mlrval_functions_test:
#	go test internal/pkg/mlrval/mlrval_functions_test.go $(ls internal/pkg/types/*.go | grep -v test)
#mlrval_format_test:
#	go test internal/pkg/mlrval/mlrval_format_test.go $(ls internal/pkg/types/*.go|grep -v test)

tests-in-order: mlrval-tests mlrmap-tests input-tests bifs-tests

# ----------------------------------------------------------------
# Regression tests (large number)
#
# See ./regression_test.go for information on how to get more details
# for debugging.  TL;DR is for CI jobs, we have 'go test -v'; for
# interactive use, instead of 'go test -v' simply use 'mlr regtest
# -vvv' or 'mlr regtest -s 20'. See also internal/pkg/auxents/regtest.
regression-test:
	go test -v regression_test.go

# go fmt ./... finds experimental C files which we want to ignore.
fmt:
	-go fmt ./cmd/...
	-go fmt ./internal/pkg/...
	-go fmt ./regression_test.go

# ----------------------------------------------------------------
# For developers before pushing to GitHub.
#
# These steps are done in a particular order:
# go:
# * builds the mlr executable
# man:
# * creates manpage mlr.1 and manpage.txt using mlr from the $PATH
# * copies the latter to docs/src
# docs:
# * turns *.md.in into *.md (live code samples), using mlr from the $PATH
# * note the man/manpage.txt becomes some of the HTML content
# * turns *.md into docs/site HTML and CSS files
dev:
	-make fmt
	make build
	make check
	make -C man build
	make -C docs/src forcebuild
	make -C docs
	@echo DONE

# ----------------------------------------------------------------
# Keystroke-savers
it: build check
so: install
sure: build check
mlr:
	go build github.com/johnkerl/miller/cmd/mlr
mprof:
	go build github.com/johnkerl/miller/cmd/mprof
mprof2:
	go build github.com/johnkerl/miller/cmd/mprof2
mprof3:
	go build github.com/johnkerl/miller/cmd/mprof3
mprof4:
	go build github.com/johnkerl/miller/cmd/mprof4
mprof5:
	go build github.com/johnkerl/miller/cmd/mprof5
mall: mprof5 mprof4 mprof3 mprof2 mprof mlr

# ----------------------------------------------------------------
# Please see comments in ./create-release-tarball as well as
# https://miller.readthedocs.io/en/latest/build/#creating-a-new-release-for-developers
release_tarball: build check
	./create-release-tarball

# ================================================================
# Go does its own dependency management, outside of make.
.PHONY: build mlr mprof mprof2 mprof3 mprof4 mprof5 check unit_test regression_test fmt dev
