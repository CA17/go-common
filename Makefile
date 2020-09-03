push:
	@read -p "type commit message: " cimsg; \
	git ci -am "bss-network: $(shell date "+%F %T") $${cimsg}"
	git push origin master


.PHONY: push


