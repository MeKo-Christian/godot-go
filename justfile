# Default recipe
default: build

# Environment detection
GOOS := env('GOOS', `go env GOOS`)
GOARCH := env('GOARCH', `go env GOARCH`)
GOIMPORTS := `which goimports || true`
CLANG_FORMAT := `which clang-format || which clang-format-10 || which clang-format-11 || which clang-format-12 || true`
GODOT := env('GODOT', `which godot || true`)
CWD := `pwd`

OUTPUT_PATH := "test/demo/lib"
TEST_MAIN := "test/main.go"

# Determine binary extension based on OS
TEST_BINARY_PATH := if GOOS == "windows" {
    OUTPUT_PATH + "/libgodotgo-test-windows-" + GOARCH + ".dll"
} else if GOOS == "darwin" {
    OUTPUT_PATH + "/libgodotgo-test-macos-" + GOARCH + ".framework"
} else if GOOS == "linux" {
    OUTPUT_PATH + "/libgodotgo-test-linux-" + GOARCH + ".so"
} else {
    OUTPUT_PATH + "/libgodotgo-test-" + GOOS + "-" + GOARCH + ".so"
}

# Display Go environment
goenv:
    go env

# Install dependencies
installdeps:
    go install golang.org/x/tools/cmd/goimports@latest

# Generate Go/C bindings from Godot headers
generate: clean
    #!/usr/bin/env bash
    set -euo pipefail
    go generate
    if [ -n "{{CLANG_FORMAT}}" ]; then
        "{{CLANG_FORMAT}}" -i pkg/ffi/ffi_wrapper.gen.h
        "{{CLANG_FORMAT}}" -i pkg/ffi/ffi_wrapper.gen.c
    fi
    find pkg -name '*.gen.go' -exec go fmt {} \;
    if [ -n "{{GOIMPORTS}}" ]; then
        find pkg -name '*.gen.go' -exec {{GOIMPORTS}} -w {} \;
    fi

# Update godot_headers from the godot binary
update_godot_headers_from_binary:
    DISPLAY=:0 "{{GODOT}}" --dump-extension-api --headless
    mv extension_api.json godot_headers/extension_api.json
    DISPLAY=:0 "{{GODOT}}" --dump-gdextension-interface --headless
    mv gdextension_interface.h godot_headers/godot/

# Build the extension library
build: goenv
    CGO_ENABLED=1 \
    GOOS={{GOOS}} \
    GOARCH={{GOARCH}} \
    CGO_CFLAGS='-fPIC -g -ggdb -O0' \
    CGO_LDFLAGS='-g3 -g -O0' \
    go build -gcflags=all="-N -l" -tags tools -buildmode=c-shared -v -x -trimpath -o "{{TEST_BINARY_PATH}}" {{TEST_MAIN}}

# Full debug build with enhanced debugging symbols
build-full: goenv
    CGO_ENABLED=1 \
    GOOS={{GOOS}} \
    GOARCH={{GOARCH}} \
    CGO_CFLAGS='-g3 -g -gdwarf -DX86=1 -fPIC -O0' \
    CGO_LDFLAGS='-g3 -g' \
    go build -gcflags="-N -l" -ldflags=-compressdwarf=0 -tags tools -buildmode=c-shared -v -x -trimpath -o "{{TEST_BINARY_PATH}}" {{TEST_MAIN}}

# Remove generated source files
clean_src:
    find pkg -name '*.gen.go' -delete
    find pkg -name '*.gen.c' -delete
    find pkg -name '*.gen.h' -delete

# Remove generated files and test binaries
clean: clean_src
    rm -f test/demo/lib/libgodotgo-*

# Run test with GDB remote debugging
remote_debug_test:
    CI=1 \
    LOG_LEVEL=info \
    GOTRACEBACK=crash \
    GODEBUG=asyncpreemptoff=1,cgocheck=1,invalidptr=1,clobberfree=1,tracebackancestors=5 \
    gdbserver --once :55555 "{{GODOT}}" --headless --verbose --debug --path test/demo/

# Initialize test project (one-time setup)
ci_gen_test_project_files:
    #!/usr/bin/env bash
    set -euo pipefail
    CI=1 \
    LOG_LEVEL=info \
    GOTRACEBACK=1 \
    GODEBUG=gctrace=1,asyncpreemptoff=1,cgocheck=1,invalidptr=1,clobberfree=1,tracebackancestors=5 \
    "{{GODOT}}" --headless --verbose --path test/demo/ --editor --quit
    # hack until fix lands: https://github.com/godotengine/godot/issues/84460
    if [ ! -f "test/demo/.godot/extension_list.cfg" ]; then
        echo 'res://example.gdextension' >> test/demo/.godot/extension_list.cfg
    fi

# Run tests in headless mode
test:
    CI=1 \
    LOG_LEVEL=info \
    GOTRACEBACK=single \
    GODEBUG=gctrace=1,asyncpreemptoff=1,cgocheck=1,invalidptr=1,clobberfree=1 \
    "{{GODOT}}" --headless --verbose --path test/demo/ --quit

# Run interactive test with debug output
interactive_test:
    LOG_LEVEL=info \
    GOTRACEBACK=1 \
    GODEBUG=gctrace=1,asyncpreemptoff=1,cgocheck=1,invalidptr=1,clobberfree=1,tracebackancestors=5 \
    "{{GODOT}}" --verbose --debug --path test/demo/

# Open demo project in Godot editor
open_demo_in_editor:
    DISPLAY=:0 \
    LOG_LEVEL=info \
    GOTRACEBACK=1 \
    GODEBUG=gctrace=1,asyncpreemptoff=1,cgocheck=1,invalidptr=1,clobberfree=1,tracebackancestors=5 \
    "{{GODOT}}" --verbose --debug --path test/demo/ --editor
