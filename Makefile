builddir=.build

build:
	mkdir -p $(builddir)
	cd $(builddir) \
	  && GOPATH=$(shell pwd) go build ../src/ \
	  && mv src fanctl

install:
	install -Dm755 $(builddir)/fanctl \
	  $(DESTDIR)/usr/bin/fanctl
	
	install -Dm644 dist/fanctl.service \
	  $(DESTDIR)/usr/lib/systemd/system/fanctl.service

	install -Dm644 dist/fanctl.desktop \
	  $(DESTDIR)/usr/share/applications/fanctl.desktop
	
	install -D -m644 dist/fanctl.png \
	  $(DESTDIR)/usr/share/pixmaps/fanctl.png

clean: 
	rm -rf $(builddir)

uninstall: 
	rm $(DESTDIR)/usr/bin/fanctl
	rm $(DESTDIR)/usr/lib/systemd/system/fanctl.service
	rm $(DESTDIR)/usr/share/applications/fanctl.desktop
	rm $(DESTDIR)/usr/share/pixmaps/fanctl.png
	
