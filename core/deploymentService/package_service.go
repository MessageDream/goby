package deploymentService

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	log "qiniupkg.com/x/log.v7"

	"github.com/go-macaron/cache"
	"github.com/go-xorm/xorm"
	"github.com/sergi/go-diff/diffmatchpatch"
	. "gopkg.in/ahmetb/go-linq.v3"

	. "github.com/MessageDream/goby/core"
	"github.com/MessageDream/goby/model"
	"github.com/MessageDream/goby/module/infrastructure"
	"github.com/MessageDream/goby/module/setting"
	"github.com/MessageDream/goby/module/storage"
)

var (
	packageTypeIOS     = errors.New(BUNDLE_IOS)
	packageTypeAndroid = errors.New(BUNDLE_ANDROID)
)

const (
	MANIFEST_FILE_NAME = "manifest.json"
	CONTENTS_NAME      = "contents"

	BUNDLE_IOS     = "main.jsbundle"
	BUNDLE_ANDROID = "android.bundle"
	PATCH_FILE_EXT = ".patch"

	PACKAGE_TYPE_IOS = iota + 1
	PACKAGE_TYPE_ANDROID
)

func checkPackageFileType(packageDir string) (int, error) {
	if err := filepath.Walk(packageDir, func(filename string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fi.IsDir() {
			return nil
		}

		if strings.HasSuffix(filename, BUNDLE_IOS) {
			return packageTypeIOS
		}

		if strings.HasSuffix(filename, BUNDLE_ANDROID) {
			return packageTypeAndroid
		}

		return nil
	}); err != nil {
		switch err {
		case packageTypeIOS:
			return PACKAGE_TYPE_IOS, nil
		case packageTypeAndroid:
			return PACKAGE_TYPE_ANDROID, nil
		default:
			return 0, err
		}
	}

	return 0, ErrDeploymentPackageContentsUnrecognized
}

func calcAllFileSHA256(dirPath string) (map[string]string, error) {
	data := map[string]string{}
	if err := filepath.Walk(dirPath, func(filename string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fi.IsDir() {
			return nil
		}

		fname := fi.Name()
		if fname == ".DS_Store" || fname == "__MACOSX" {
			return nil
		}

		relativePath, err := filepath.Rel(dirPath, filename)
		relativePath = filepath.ToSlash(relativePath)

		if err != nil {
			return err
		}

		if data[relativePath], err = infrastructure.FileSHA256(filename); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}
	return data, nil
}

func calcPackageHashWithAllFileHashMap(hashOfMap map[string]string) (string, error) {
	result := make([]string, 0, 10)
	From(hashOfMap).Select(func(kv interface{}) interface{} {
		return kv
	}).Sort(func(i, j interface{}) bool {
		first := i.(KeyValue)
		second := j.(KeyValue)
		return strings.Compare(first.Key.(string), second.Key.(string)) > 0
	}).Select(func(kv interface{}) interface{} {
		item := kv.(KeyValue)
		return item.Key.(string) + ":" + item.Value.(string)
	}).ToSlice(&result)
	sort.Sort(sort.StringSlice(result))

	bytes, err := json.Marshal(result)
	if err != nil {
		return "", err
	}

	return infrastructure.EncodeSHA256(string(bytes)), nil
}

func rearrangingPackage(from string) (string, string, string, string, map[string]string, error) {

	manifestMap, err := calcAllFileSHA256(from)
	if err != nil {
		return "", "", "", "", nil, err
	}

	packageHash, err := calcPackageHashWithAllFileHashMap(manifestMap)
	if err != nil {
		return "", "", "", "", nil, err
	}
	packageHashPath := path.Join(setting.AttachmentPath, packageHash)
	manifestFile := path.Join(packageHashPath, MANIFEST_FILE_NAME)
	contentPath := path.Join(packageHashPath, CONTENTS_NAME)

	if err := os.MkdirAll(contentPath, 0755); err != nil {
		return "", "", "", "", nil, err
	}

	if err := infrastructure.CopyDirFiles(from, contentPath); err != nil {
		return "", "", "", "", nil, err
	}

	manifestBytes, err := json.Marshal(manifestMap)
	if err != nil {
		return "", "", "", "", nil, err
	}

	if err = writeTextContentToFile(manifestFile, manifestBytes); err != nil {
		return "", "", "", "", nil, err
	}

	manifestHash, err := infrastructure.FileSHA256(manifestFile)
	if err != nil {
		return "", "", "", "", nil, err
	}

	return packageHash, contentPath, manifestHash, manifestFile, manifestMap, nil

}

func createPackage(deployment *model.Deployment,
	isMandatory,
	isDisabled bool,
	rollout uint8,
	size int64,
	appVersion,
	packageHash,
	manifestURL,
	blobURL,
	releaseMethod,
	releasedBy,
	description,
	originalLabel,
	originalDeployment string) (*model.Package, error) {

	pkgCheck := &model.Package{
		DeployID: deployment.ID,
		Hash:     packageHash,
	}

	if exist, err := pkgCheck.Exist(); err != nil || exist {
		if err != nil {
			return nil, err
		}
		if exist && releaseMethod != "Rollback" {
			return nil, ErrDeploymentPackageAlreadyExist
		}
	}

	result, err := model.Transaction(func(sess *xorm.Session) (interface{}, error) {

		version := &model.DeploymentVersion{
			DeployID:   deployment.ID,
			AppVersion: appVersion,
		}

		exist, err := version.Get()
		if err != nil {
			return nil, err
		}
		if !exist {
			if err := version.Create(sess); err != nil {
				return nil, err
			}
		}

		pkg := &model.Package{
			DeployVersionID:    version.ID,
			DeployID:           deployment.ID,
			Hash:               packageHash,
			BlobURL:            blobURL,
			ManifestBlobURL:    manifestURL,
			Size:               size,
			Label:              "v" + strconv.Itoa(deployment.LabelCursor+1),
			ReleaseMethod:      releaseMethod,
			ReleasedBy:         releasedBy,
			IsMandatory:        isMandatory,
			IsDisabled:         isDisabled,
			Rollout:            rollout,
			OriginalLabel:      originalLabel,
			OriginalDeployment: originalDeployment,
			Description:        description,
		}

		if err := pkg.Create(sess); err != nil {
			return nil, err
		}

		pkgMetrics := &model.PackageMetrics{
			PackageID: pkg.ID,
		}

		if err := pkgMetrics.Create(sess); err != nil {
			return nil, err
		}

		deployHistory := &model.DeploymentHistory{
			DeployID:  deployment.ID,
			PackageID: pkg.ID,
		}

		if err := deployHistory.Create(sess); err != nil {
			return nil, err
		}

		version.PackageID = pkg.ID

		if err := version.Update(sess, "package_id"); err != nil {
			return nil, err
		}

		deployment.LabelCursor++
		deployment.LastVersionID = version.ID

		if err := deployment.Update(sess, "last_version_id", "label_cursor"); err != nil {
			return nil, err
		}

		return pkg, nil

	})

	if err != nil {
		return nil, err
	}

	return result.(*model.Package), nil
}

func createDiffPackage(originalPackage *model.Package) ([]*model.PackageDiff, error) {

	maxDiff := setting.PackageConfig.MAXDiffCount

	lastPackages, err := originalPackage.FindPrePackages(maxDiff, false)
	if err != nil {
		return nil, err
	}

	firstTwoPackages, err := originalPackage.FindPrePackages(2, true)
	if err != nil {
		return nil, err
	}

	destPackages := make([]*model.Package, 0, maxDiff+2)

	From(lastPackages).Union(From(firstTwoPackages)).ToSlice(&destPackages)

	if len(destPackages) == 0 {
		return nil, nil
	}

	diffPkgs := make([]*model.PackageDiff, 0, len(destPackages))
	for _, v := range destPackages {
		diff, err := createOneDiffPackage(originalPackage, v)
		if err != nil {
			return nil, err
		}
		diffPkgs = append(diffPkgs, diff)
	}

	return diffPkgs, nil

}

func createOneDiffPackage(originalPackage, prePackage *model.Package) (*model.PackageDiff, error) {
	pkgDiff := &model.PackageDiff{
		PackageID:              originalPackage.ID,
		DiffAgainstPackageHash: prePackage.Hash,
	}

	if exist, err := pkgDiff.Exist(); err != nil || exist {
		if exist {
			return nil, nil
		}
		return nil, err
	}

	localPkgCacheDir := setting.AttachmentPath

	localOriginalPkgContentDir := path.Join(localPkgCacheDir, originalPackage.Hash, CONTENTS_NAME)
	localOriginalManifestFile := path.Join(localPkgCacheDir, originalPackage.ManifestBlobURL)
	originalManifestMap := map[string]string{}

	if !(infrastructure.FileExist(localOriginalPkgContentDir) && infrastructure.FileExist(localOriginalManifestFile)) {
		_, pkgContentDir, _, _, manifestMap, err := downloadAndUnzipPackage(originalPackage.BlobURL)
		if err != nil {
			return nil, err
		}

		localOriginalPkgContentDir = pkgContentDir
		originalManifestMap = manifestMap
		defer os.RemoveAll(path.Join(pkgContentDir, "../"))

	} else {

		content, err := readTextFileContent(localOriginalManifestFile)

		err = json.Unmarshal(content, originalManifestMap)
		if err != nil {
			return nil, err
		}
	}

	localPrePkgContentDir := path.Join(localPkgCacheDir, prePackage.Hash, CONTENTS_NAME)
	localPreManifestFile := path.Join(localPkgCacheDir, prePackage.ManifestBlobURL)
	preManifestMap := map[string]string{}

	if !(infrastructure.FileExist(localPrePkgContentDir) && infrastructure.FileExist(localPreManifestFile)) {
		_, pkgContentDir, _, _, manifestMap, err := downloadAndUnzipPackage(prePackage.BlobURL)
		if err != nil {
			return nil, err
		}

		localPrePkgContentDir = pkgContentDir
		preManifestMap = manifestMap
		defer os.RemoveAll(path.Join(pkgContentDir, "../"))

	} else {

		content, err := readTextFileContent(localPreManifestFile)

		err = json.Unmarshal(content, preManifestMap)
		if err != nil {
			return nil, err
		}
	}

	diffFiles, originalOnlyFiles, preOnlyFiles := diffManifestMap(originalManifestMap, preManifestMap)

	hotGobyMap := map[string]interface{}{
		"deletedFiles": preOnlyFiles,
	}

	var remainFiles []string
	var hotGobyFileName = "hotcodepush.json"

	if setting.PackageConfig.EnableGoogleDiff {
		hotGobyFileName = "hotgoby.json"
		needRemovedFileIndexies := make([]int, 0, len(diffFiles))
		patchFiles := make([]string, 0, len(diffFiles))
		for idx, fi := range diffFiles {
			ext := path.Ext(fi)
			if ext != ".jsbundle" && ext != ".bundle" {
				continue
			}
			originFilePath := path.Join(localOriginalPkgContentDir, fi)
			originBytes, err := readTextFileContent(originFilePath)
			if err != nil {
				continue
			}

			// if originBytes[0] != 0xef || originBytes[1] != 0xbb || originBytes[2] != 0xbf {
			// 	continue
			// }

			preFilePath := path.Join(localPrePkgContentDir, fi)
			preBytes, err := readTextFileContent(preFilePath)
			if err != nil {
				continue
			}

			newText := string(originBytes)
			oldText := string(preBytes)

			dmp := diffmatchpatch.New()
			diffs := dmp.DiffMain(oldText, newText, true)
			patches := dmp.PatchMake(oldText, diffs)
			patchText := dmp.PatchToText(patches)
			patchFilePath := originFilePath + PATCH_FILE_EXT
			if err := writeTextContentToFile(patchFilePath, []byte(patchText)); err != nil {
				continue
			}
			patchFiles = append(patchFiles, fi+PATCH_FILE_EXT)
			needRemovedFileIndexies = append(needRemovedFileIndexies, idx)
		}

		for _, v := range needRemovedFileIndexies {
			diffFiles = append(diffFiles[:v], diffFiles[v+1:]...)
		}

		From(diffFiles).Concat(From(patchFiles)).Concat(From(originalOnlyFiles)).ToSlice(&remainFiles)

		hotGobyMap["patchedFiles"] = patchFiles
	} else {
		From(diffFiles).Concat(From(originalOnlyFiles)).ToSlice(&remainFiles)
	}
	hotGobyText, err := json.Marshal(hotGobyMap)
	if err != nil {
		return nil, err
	}
	gobyFilePath := path.Join(localOriginalPkgContentDir, hotGobyFileName)

	if err := writeTextContentToFile(gobyFilePath, hotGobyText); err != nil {
		return nil, err
	}

	remainFiles = append(remainFiles, hotGobyFileName)

	zipTo := path.Join(os.TempDir(), infrastructure.GetRandomString(32), infrastructure.GetRandomString(16)+".zip")
	if err := infrastructure.Compress(localOriginalPkgContentDir, remainFiles, zipTo); err != nil {
		return nil, err
	}

	zipHash, err := infrastructure.FileSHA256(zipTo)
	if err != nil {
		return nil, err
	}

	fileInfo, err := os.Stat(zipTo)
	if err != nil {
		return nil, err
	}

	diffSize := fileInfo.Size()
	zipURLPath, err := uploadPackage(zipHash, zipTo)
	if err != nil {
		return nil, err
	}
	defer func() {
		os.Remove(gobyFilePath)
		os.RemoveAll(path.Join(zipTo, "../"))
	}()

	diffPkg := &model.PackageDiff{
		PackageID:              originalPackage.ID,
		DiffAgainstPackageHash: prePackage.Hash,
		DiffBlobURL:            zipURLPath,
		DiffSize:               diffSize,
	}

	if err := diffPkg.Create(nil); err != nil {
		return nil, err
	}

	return diffPkg, nil

}

func uploadPackage(key, packagePath string) (string, error) {
	return storage.Upload(key, packagePath)
}

func readTextFileContent(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return ioutil.ReadAll(file)
}

func writeTextContentToFile(filePath string, content []byte) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err = io.Copy(file, bytes.NewReader(content)); err != nil {
		return err
	}

	return nil
}

func downloadAndUnzipPackage(key string) (string, string, string, string, map[string]string, error) {
	pkgPath, err := storage.Download(key)
	if err != nil {
		return "", "", "", "", nil, err
	}
	defer os.Remove(pkgPath)

	unzipTo := pkgPath + "_unzip"

	if err := infrastructure.DeCompress(pkgPath, unzipTo); err != nil {
		return "", "", "", "", nil, err
	}
	return rearrangingPackage(unzipTo)
}

func diffManifestMap(orignal, dest map[string]string) (diff, orignalOnly, destOnly []string) {

	diff = []string{}
	orignalOnly = []string{}
	destOnly = []string{}
	destCopy := map[string]string{}

	for k, v := range dest {
		destCopy[k] = v
	}

	for k, v := range orignal {
		if v2, ok := dest[k]; ok {
			if v != v2 {
				diff = append(diff, k)
			}
			delete(destCopy, k)
		} else {
			orignalOnly = append(orignalOnly, k)
		}
	}

	for k := range destCopy {
		destOnly = append(destOnly, k)
	}

	return
}

func clearCache(cache cache.Cache, deploymentKey, appVersion string) {
	cacheKey := infrastructure.EncodeSHA256(deploymentKey + appVersion)
	if cache.IsExist(cacheKey) {
		if err := cache.Delete(cacheKey); err != nil {
			log.Error(4, "cache Delete %s error:%v", cacheKey, err)
		}
	}
}
