.PHONY: all
all: build flash

.PHONY: build
build:
	tinygo build -target=elecrow-rp2350 -o main.bin .

.PHONY: dev
dev:
	tinygo flash -target=elecrow-rp2350 -monitor -port=/dev/ttyACM0 example/example.go

.PHONY: flash
flash:
	tinygo flash -target=elecrow-rp2350 -port=/dev/ttyACM0 example/example.go