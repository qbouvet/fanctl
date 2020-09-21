builddir=.build
binary=main

build:
	mkdir -p $(builddir)
	cd $(builddir) \
	  && GOPATH=$(shell pwd) go build ../src/ \
	  && mv src $(binary)

install:
	# 	Binary, icon and .desktop
	mkdir -p $(DESTDIR)/opt/fanctl
	install -Dm755 $(builddir)/$(binary) $(DESTDIR)/opt/fanctl/
	install -Dm644 misc/fanctl.png $(DESTDIR)/opt/fanctl/
	install -Dm644 misc/fanctl.desktop $(DESTDIR)/opt/fanctl/

	# 	Systemd service
	mkdir -p $(DESTDIR)/usr/lib/systemd/system/
	install -Dm644 misc/fanctl.service $(DESTDIR)/usr/lib/systemd/system/

clean: 
	rm -rf $(builddir)

uninstall: 
	rm /usr/lib/systemd/system/fanctl.service
	rm -rf /opt/fanctl


# https://bbs.archlinux.org/viewtopic.php?id=250189
