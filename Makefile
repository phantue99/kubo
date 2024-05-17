all:
	@gmake $@
.PHONY: all

.DEFAULT:
	@gmake $@


update-ipfs-production:
	ssh root@167.172.4.178  mv /mnt/w3ipfs/pinning-service/ipfs /mnt/w3ipfs/pinning-service/ipfs.bk`date +%Y%m%d%H%M%S`
	scp /cmd/ipfs/ipfs root@167.172.4.178:/mnt/w3ipfs/pinning-service
