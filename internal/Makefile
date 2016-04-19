# This Makefile pulls the latest hdrhistogram (with dependencies),
# rewrites the imports to match this location,
# and ensures that all the tests pass.

go: clean clone rewrite

clean:
	rm -rf hdrhistogram*

clone:
	git clone https://github.com/codahale/hdrhistogram.git && rm -rf hdrhistogram/.git

rewrite:
	grep -rl --exclude Makefile 'github.com/codahale' . | xargs sed -i '' 's#github.com/codahale#github.com/smartystreets/metrics/internal#g'

test:
	#go test github.com/smartystreets/metrics/...
	go test github.com/smartystreets/metrics/internal/...
