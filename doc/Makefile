output_name=../doc.pdf
name=main
tmp=.build


#
# 	Compilation (Complete or Quick)
#

quick: init
	cd $(tmp) \
		&& pdflatex -interaction=batchmode $(name).tex \
		|| cat $(name).log
	cp $(tmp)/$(name).pdf .
	cp $(tmp)/$(name).pdf $(output_name)

all: init
	cd $(tmp) \
		&& pdflatex -interaction=errorstopmode -halt-on-error $(name).tex \
		&& biber $(name) \
		&& pdflatex -interaction=errorstopmode -halt-on-error $(name).tex \
		&& pdflatex -interaction=errorstopmode -halt-on-error $(name).tex \
		&& pdflatex -interaction=errorstopmode -halt-on-error $(name).tex 
	cp $(tmp)/$(name).pdf .
	cp $(tmp)/$(name).pdf $(output_name)

init: *.tex *.bib *.sty
	mkdir -p $(tmp)
	cp -au *.tex *.bib *.sty figures $(tmp)

clean:
	if [ -d $(tmp) ]; then rm $(tmp) -rf  ; fi
	if [ -f $(name).pdf ]; then rm $(name).pdf; fi


#
# 	Spell checking 
#

spell: 
	for f in ?-*.tex; do aspell --master=en_US --mode=tex --dont-backup check $$f; done;
	resources/spot_dupl $$(ls ?-*.tex) 
	grep -n ' a [aeio]' \
		$$(ls ?-*.tex) --color=auto || [[ $$? == 1 ]]
	grep -n ' can not ' \
		$$(ls ?-*.tex) --color=auto || [[ $$? == 1 ]]
	grep -n ' an [bcdfghjklmnpqrstvwxyzBCDFGHJKLMNPQRSTVWXYZ]' \
		$$(ls ?-*.tex) --color=auto || [[ $$? == 1 ]]
	grep -n "n't" \
		$$(ls ?-*.tex) --color=auto || [[ $$? == 1 ]]
	grep -n "[Rr]unners*" \
		$$(ls ?-*.tex) --color=auto || [[ $$? == 1 ]]
	grep -n " phi| psi" \
		$$(ls ?-*.tex) --color=auto || [[ $$? == 1 ]]
	grep -n " mtlfreeze" \
		$$(ls ?-*.tex) --color=auto || [[ $$? == 1 ]]
	grep -n " mtlprenex" \
		$$(ls ?-*.tex) --color=auto || [[ $$? == 1 ]]
	grep -n " partitions" \
		$$(ls ?-*.tex) --color=auto || [[ $$? == 1 ]]