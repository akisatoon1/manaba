# Manaba Go Library

## 概要
manabaへのログイン、ファイルアップロード、ファイル提出、ファイル提出取り消し、アップロード済みのファイル削除、を行うGoライブラリです。

## Functions

### func Login
```go
func Login(jar *cookiejar.Jar, username string, password string) error
```
manabaにログインして, `*jar`にCookie情報を保存します. 以降manabaにログインした状態で行う動作は`*cookiejar.Jar`を使います。

### func UploadFile
```go
func UploadFile(jar *cookiejar.Jar, url string, filePath string) error
```
`filePath`で指定されたファイルを, `url`にアップロードします.

### func SubmitReports
```go
func SubmitReports(jar *cookiejar.Jar, url string) error
```
`url`で指定されたmanabaコースレポートで, アップロード済みのファイルを提出します.

### func CancelSubmittion
```go
func CancelSubmittion(jar *cookiejar.Jar, url string) error
```
`url`で指定されたmanabaコースレポートで, 提出を取り消します.

### func DeleteAllFiles
```go
func DeleteAllFiles(jar *cookiejar.Jar, url string) error
```
`url`で指定されたmanabaコースレポートで, アップロード済みの全てのファイルを削除します.
