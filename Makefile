all:
	@echo This makefile is for releasing cat-o-licious for various OSes.

release-macos:
	bash hack/release-macos.sh

release-ubuntu:
	bash hack/release-ubuntu.sh

release-windows:
	bash hack/release-windows.sh
