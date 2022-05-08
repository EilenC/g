package version

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/voidint/g/pkg/checksum"
	myhttp "github.com/voidint/g/pkg/http"
)

var (
	// ErrVersionNotFound 版本不存在
	ErrVersionNotFound = errors.New("version not found")
	// ErrPackageNotFound 版本包不存在
	ErrPackageNotFound = errors.New("installation package not found")
)

// FindVersion 返回指定名称的版本
func FindVersion(all []*Version, name string) (*Version, error) {
	for i := range all {
		if all[i].Name == name {
			return all[i], nil
		}
	}
	return nil, ErrVersionNotFound
}

// Version go版本
type Version struct {
	Name     string // 版本名，如'1.12.4'
	Packages []*Package
}

// FindPackage 返回指定操作系统和硬件架构的版本包
func (v *Version) FindPackage(kind, goos, goarch string) (*Package, error) {
	prefix := fmt.Sprintf("go%s.%s-%s", v.Name, goos, goarch)
	for i := range v.Packages {
		if v.Packages[i] == nil || !strings.EqualFold(v.Packages[i].Kind, kind) || !strings.HasPrefix(v.Packages[i].FileName, prefix) {
			continue
		}
		return v.Packages[i], nil
	}

	return nil, ErrPackageNotFound
}

// FindPackages 返回指定操作系统和硬件架构的版本包
func (v *Version) FindPackages(kind, goos, goarch string) (pkgs []*Package, err error) {
	prefix := fmt.Sprintf("go%s.%s-%s", v.Name, goos, goarch)
	for i := range v.Packages {
		if v.Packages[i] == nil || !strings.EqualFold(v.Packages[i].Kind, kind) || !strings.HasPrefix(v.Packages[i].FileName, prefix) {
			continue
		}
		pkgs = append(pkgs, v.Packages[i])
	}
	if len(pkgs) == 0 {
		return nil, ErrPackageNotFound
	}
	return pkgs, nil
}

// Package go版本安装包
type Package struct {
	FileName  string
	URL       string
	Kind      string
	OS        string
	Arch      string
	Size      string
	Checksum  string
	Algorithm string // checksum algorithm
}

const (
	// SourceKind go安装包种类-源码
	SourceKind = "Source"
	// ArchiveKind go安装包种类-压缩文件
	ArchiveKind = "Archive"
	// InstallerKind go安装包种类-可安装程序
	InstallerKind = "Installer"
)

// DownloadWithProgress 下载版本另存为指定文件且显示下载进度
func (pkg *Package) DownloadWithProgress(dst string) (size int64, err error) {
	return myhttp.Download(pkg.URL, dst, os.O_CREATE|os.O_WRONLY, 0644, true)
}

const (
	// SHA256 校验和算法-sha256
	SHA256 = "SHA256"
	// SHA1 校验和算法-sha1
	SHA1 = "SHA1"
)

// VerifyChecksum 验证目标文件的校验和与当前安装包的校验和是否一致
func (pkg *Package) VerifyChecksum(filename string) (err error) {
	var algo checksum.Algorithm
	switch pkg.Algorithm {
	case SHA256:
		algo = checksum.SHA256
	case SHA1:
		algo = checksum.SHA1
	default:
		return checksum.ErrUnsupportedChecksumAlgorithm
	}
	return checksum.VerifyFile(algo, pkg.Checksum, filename)
}
