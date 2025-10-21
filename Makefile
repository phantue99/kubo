all:
	@gmake $@
.PHONY: all

.DEFAULT:
	@gmake $@

STAGING=157.230.195.37
STAGING_DIR=/mnt/staging_data/w3ipfs_stg

PROD=164.172.4.178
PROD_DIR=/mnt/w3ipfs/pinning-service

build:
	go build -o ./bin ./cmd/ipfs

update-ipfs-production: build
	ssh root@$(PROD)  mv $(PROD_DIR)/ipfs $(PROD_DIR)/ipfs.bk`date +%Y%m%d%H%M%S`
	scp ./bin/ipfs root@$(PROD):$(PROD_DIR)

update-ipfs-staging: build
	ssh root@$(STAGING) mv $(STAGING_DIR)/ipfs $(STAGING_DIR)/ipfs.bk`date +%Y%m%d%H%M%S`
	scp ./bin/ipfs root@$(STAGING):$(STAGING_DIR)