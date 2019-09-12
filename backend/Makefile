
.PHONY: clean all

all:
	-$(MAKE) -C $(WVB_BACKEND)
	-$(MAKE) -C $(WVB_CONFIG)
	-$(MAKE) -C $(GET_IMG)
	-$(MAKE) -C $(ADD_IMG)
clean:
	-$(MAKE) -C $(WVB_BACKEND) clean
	-$(MAKE) -C $(WVB_CONFIG) clean
	-$(MAKE) -C $(GET_IMG) clean
	-$(MAKE) -C $(ADD_IMG) clean
