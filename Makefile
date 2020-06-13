builddir=.build
binary=main

build:
	mkdir -p $(builddir)
	cd $(builddir) \
	  && GOPATH=$(shell pwd) go build ../src/ \
	  && mv src $(binary)

install:
	# 	Binary
	mkdir -p /opt/fanctl
	install $(builddir)/$(binary) \
	  --target-directory=/opt/fanctl/ \
	  --mode=775 --owner=root --group=root 

	# 	Icon and clickable shortcut
	xdg-icon-resource install --size 128 misc/fanctl-fanctl.png
	install misc/fanctl.desktop \
	  --target-directory=/opt/fanctl/ \
	  --mode=664 --owner=root --group=root 

	# 	Systemd service
	install \
	  --target-directory=/usr/lib/systemd/system/ \
	  --mode=664 --owner=root --group=root \
	  misc/fanctl.service

clean: 
	rm -rf $(builddir)

uninstall: 
	rm /usr/lib/systemd/system/fanctl.service
	rm -rf /opt/fanctl
