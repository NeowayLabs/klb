#!/usr/bin/env nash

cwdir <= pwd | xargs echo -n
vendordir = $cwdir + "/tests/vendor"

fn vendor() {
        rm -rf $vendordir

        bindir = $vendordir + "/bin"
        srcdir = $vendordir + "/src"
        pkgdir = $vendordir + "/pkg"
        mkdir -p $bindir $srcdir $pkgdir

        setenv GOPATH = $vendordir
        setenv GOBIN = $vendordir

        go get -t -v ./tests/...

        rawpaths <= ls $srcdir
        paths <= split($paths, "\n")
        for path in $paths {
                mv $srcdir + $path $vendor
        }
        rm -rf $bindir $srcdir $pkgdir
}

fn govend() {
    rm -rf $vendordir
    mkdir -p $vendordir

    go get github.com/govend/govend
    govend -v
}

#vendor()
govend()
