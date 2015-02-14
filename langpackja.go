package moderniejapanizer

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/hnakamur/w32version"
	"github.com/hnakamur/windowsupdate"
)

const (
	Windows7JapaneseLanguagePackUpdateID = "00a156d4-3876-4cd5-bd38-517679c6ba59"
)

type langPackUrlAndSum struct {
	Url    string
	Md5Sum string
}

// These  urls are copied from two sites below. Thanks for nice information!
// http://www.froggie.sk/jp/lp32sp2.html
// http://winaero.com/blog/download-official-mui-language-packs-for-windows-8-1-windows-8-and-windows-7/

var langPackUrls = map[w32version.W32Version]langPackUrlAndSum{
	w32version.WindowsVista: langPackUrlAndSum{
		Url:    "http://download.windowsupdate.com/msdownload/update/software/updt/2009/06/lp-ja-jp_30178fcc94adb29cad5a14b535efeb555ab39e0b.exe",
		Md5Sum: "266e170497339e10bf7f52c7aa1a0879",
	},
	w32version.Windows8: langPackUrlAndSum{
		Url:    "http://download.windowsupdate.com/msdownload/update/software/updt/2012/10/windows8-kb2607607-x86-jpn_c0daf9262b007af8d3c968435f81b671a8765519.cab",
		Md5Sum: "2abf931cf582cc24289817c4ff4968a8",
	},
	w32version.Windows8_1: langPackUrlAndSum{
		Url:    "http://fg.v4.download.windowsupdate.com/d/msdownload/update/software/updt/2013/09/lp_1816a277746d350dea3eaf2d88c2b2786e8aa185.cab",
		Md5Sum: "b80008dda049eb4ca7a823855fbac1d1",
	},
}

func InstallLangPackJa(version w32version.W32Version) error {
	if version == w32version.Windows7 {
		return installLangPackJaWindows7()
	}

	path := buildLangPackPath()
	err := downloadLangPackJa(version, path)
	if err != nil {
		return err
	}

	var cmd *exec.Cmd
	if version == w32version.WindowsVista {
		cmd = exec.Command(os.Getenv("ComSpec"), "/c", "start", "/wait", path)
	} else {
		cmd = exec.Command("lpksetup", "/i", "ja-JP", "/r", "/p",
			filepath.Dir(path))
	}
	err = cmd.Run()
	if err != nil {
		return err
	}
	fmt.Printf("install Jpanese language pack done\n")
	return nil
}

// NOTE: You must call ole.CoInitialize(0) in advance.
func installLangPackJaWindows7() error {
	session, err := windowsupdate.NewSession()
	if err != nil {
		return err
	}
	defer session.Release()

	fmt.Printf("Start searching...\n")
	update, err := session.FindByUpdateID(Windows7JapaneseLanguagePackUpdateID)
	if err != nil {
		return err
	}

	if update.Installed {
		fmt.Printf("already installed. exiting\n")
		return nil
	}

	updates := []windowsupdate.Update{update}

	if update.Downloaded {
		fmt.Printf("already downloaded, skip downloading\n")
	} else {
		err = session.Download(updates)
		if err != nil {
			return err
		}
	}

	result, err := session.Install(updates)
	if err != nil {
		return err
	}

	fmt.Printf("ResultCode=%d, RebootRequired=%v\n", result.ResultCode, result.RebootRequired)
	for i, ur := range result.UpdateResults {
		fmt.Printf("UpdateResult[%d] ResultCode=%d, RebootRequired=%v\n", i, ur.ResultCode, ur.RebootRequired)
	}

	return nil
}

func buildLangPackPath() string {
	return filepath.Join(os.Getenv("TEMP"), "lp.mlc")
}

func downloadLangPackJa(version w32version.W32Version, path string) error {
	urlAndSum := langPackUrls[version]
	fmt.Printf("url=%s, sum=%s\n", urlAndSum.Url, urlAndSum.Md5Sum)

	err := downloadFileIfMd5NotMatch(urlAndSum.Md5Sum, urlAndSum.Url, path)
	if err != nil {
		return err
	}

	return nil
}

func downloadFileIfMd5NotMatch(md5, url, path string) error {
	localMd5, err := calcMd5OfFile(path)
	if err != nil {
		return err
	}
	if strings.EqualFold(localMd5, md5) {
		return nil
	}

	fmt.Printf("Start downloading %s from url %s ...\n", path, url)
	localMd5, err = downloadFileAndCalcMd5(url, path)
	if err != nil {
		return err
	}
	if !strings.EqualFold(localMd5, md5) {
		return fmt.Errorf("Md5 unmatched. remote=%s, local=%s",
			md5, localMd5)
	}
	fmt.Printf("Finished downloading %s from url %s\n", path, url)

	return nil
}

func downloadFileAndCalcMd5(url, path string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	writer, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer writer.Close()

	h := md5.New()
	reader := io.TeeReader(resp.Body, h)
	_, err = io.Copy(writer, reader)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func calcMd5OfFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	defer f.Close()
	return calcMd5(f)
}

func calcMd5(rd io.Reader) (string, error) {
	h := md5.New()
	_, err := io.Copy(h, rd)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
