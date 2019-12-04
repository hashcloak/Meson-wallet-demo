flags=.makeFlags
VPATH=$(flags)
$(shell mkdir -p $(flags))

all: build

clean:
	rm -rf ./.makeFlags

build:
	docker build -t hashcloak/meson-wallet -f Dockerfile .

send-txn: build
	docker run \
		--network nonvoting_testnet_nonvoting_test_net \
		hashcloak/meson-wallet \
		/wallet \
		-pk $(shell cat ./pk) \
		-t gor \
		-s gor \
		-chain 5
