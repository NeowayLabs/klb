#!/usr/bin/env nash

fn vendor() {
        cwdir <= pwd | xargs echo -n
        vendordir = $cwdir + "/tests/vendor"
        rm -rf $vendordir

        bindir = $vendordir + "/bin"
        srcdir = $vendordir + "/src"
        pkgdir = $vendordir + "/pkg"
        mkdir -p $bindir $srcdir $pkgdir

        setenv GOPATH = $vendordir
        setenv GOBIN = $vendordir

        go get -v ./tests/...

        rawpaths <= ls $srcdir
        paths <= split($paths, "\n")
        for path in $paths {
                mv $srcdir + $path $vendor
        }
        rm -rf $bindir $srcdir $pkgdir
}

vendor()
