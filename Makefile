PREFIX ?= /usr
DESTDIR ?=

.PHONY: build clean deb install uninstall

build:
	go build -o filex .

clean:
	rm -f filex

deb: clean
	dpkg-buildpackage -b -us -uc

install: build
	install -Dm755 filex $(DESTDIR)$(PREFIX)/bin/filex
	install -Dm644 assets/filex.desktop $(DESTDIR)$(PREFIX)/share/applications/filex.desktop
	install -Dm644 assets/filex.svg $(DESTDIR)$(PREFIX)/share/icons/hicolor/scalable/apps/filex.svg

uninstall:
	rm -f $(DESTDIR)$(PREFIX)/bin/filex
	rm -f $(DESTDIR)$(PREFIX)/share/applications/filex.desktop
	rm -f $(DESTDIR)$(PREFIX)/share/icons/hicolor/scalable/apps/filex.svg
