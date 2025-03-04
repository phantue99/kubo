all:
	@gmake $@
.PHONY: all

.DEFAULT:
	@gmake $@


update-ipfs-production:
	ssh root@157.230.195.37  mv /mnt/w3ipfs/pinning-service/ipfs /mnt/w3ipfs/pinning-service/ipfs.bk`date +%Y%m%d%H%M%S`
	scp /cmd/ipfs/ipfs root@157.230.195.37:/mnt/w3ipfs/pinning-service
