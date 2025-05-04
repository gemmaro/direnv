all: ja/README.md
.PHONY: all

ja/README.md: po/en.pot po/ja.po
	LANG=C po4a po4a.cfg

po/en.pot:
	touch $@
