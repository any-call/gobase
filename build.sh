##version 24052701
##打包成 ios .a
export CFLAGS="-arch x86_64 -miphoneos-version-min=9.0 -isysroot "$(xcrun -sdk iphonesimulator --show-sdk-path)
export CGO_LDFLAGS="-arch x86_64 -miphoneos-version-min=9.0 -isysroot "$(xcrun -sdk iphonesimulator --show-sdk-path)
CGO_ENABLED=1 GOARCH=amd64 GOOS=darwin CC="clang $CFLAGS $CGO_LDFLAGS" go build -tags ios -ldflags=-w -trimpath -v -o "mylib.a" -buildmode c-archive


##打包成 mac .a
export CFLAGS="-mmacosx-version-min=10.9 -isysroot "$(xcrun -sdk macosx --show-sdk-path)
export CGO_LDFLAGS="-mmacosx-version-min=10.9 -isysroot "$(xcrun -sdk macosx --show-sdk-path)
CGO_ENABLED=1 GOARCH=amd64 GOOS=darwin CC="clang $CFLAGS $CGO_LDFLAGS" go build -tags macosx -ldflags=-w -trimpath -v -o "mylib.a" -buildmode c-archive


## 编译成 .so for flutter
go build --buildmode=c-shared -o ../convert_gui/library/json_to_model.so main.go


还有编译为 Android :  https://blog.csdn.net/HeroRazor/article/details/121436261