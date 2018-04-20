BUILD_DIR := build

default: clox

debug:
	@ $(MAKE) -f util/c.make NAME=cloxd MODE=debug SOURCE_DIR=c

clean:
	@ rm -rf $(BUILD_DIR)
	@ rm -rf gen

clox:
	@ $(MAKE) -f util/c.make NAME=clox MODE=release SOURCE_DIR=c
	@ cp build/clox clox # for convenience, copy the interpreter to the top level.

.PHONY: clean clox debug
