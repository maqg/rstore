#!/bin/sh

DISTDIR=package

rm -rf $DISTDIR
mkdir $DISTDIR

TODAY=$(date +"%m%d%Y")

echo "Building..."
go build -o $DISTDIR/rstore main.go
go build -o $DISTDIR/apihelper cmd/apihelper/main.go
go build -o $DISTDIR/rstcli cmd/cli/main.go

cd $DISTDIR

mkdir -p var/logs

if [ ! -x rstore ] || [ ! -x apihelper ] || [ ! -x rstcli ]; then
	echo "Build rstore error"
	exit 1
fi

echo "Stripping binaries..."
strip rstore
strip apihelper
strip rstcli

PACKAGE=RSTORE-$TODAY".tgz"

echo "Packaging..."
cp ../config.yml config-raw.yml
tar -zcf $PACKAGE ../frontend rstore rstcli apihelper var *.yml

echo ""
echo "Build Successfully, got file $PACKAGE"
echo ""
