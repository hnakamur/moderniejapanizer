package moderniejapanizer

import (
	"os/exec"
	"syscall"
	"time"

	"github.com/hnakamur/w32registry"
	wa "github.com/hnakamur/w32uiautomation"
	"github.com/hnakamur/w32version"
)

func SwitchInputMethodJa(version w32version.W32Version) error {
	if version == w32version.Windows7 || version == w32version.WindowsVista {
		return installLangPackJaWindows7()
	} else {
		return switchInputMethodJaWin8()
	}
}

// NOTE: You need to logoff or reboot for the display language to be changed
func SetDisplayLanguage(displayLanguageCode string) error {
	return w32registry.SetKeyValueMultiString(syscall.HKEY_CURRENT_USER, `Control Panel\Desktop`, "PreferredUILanguagesPending", []string{displayLanguageCode})
}

func switchInputMethodJaWin8() error {
	err := exec.Command("control.exe", "/name", "Microsoft.Language").Start()
	if err != nil {
		return err
	}
	time.Sleep(time.Second)

	auto, err := wa.NewUIAutomation()
	if err != nil {
		return err
	}

	root, err := auto.GetRootElement()
	if err != nil {
		return err
	}
	defer root.Release()

	languageWin, err := findElementByName(auto, root, "Language")
	if err != nil {
		return err
	}

	enUsListItem, err := findParentElementByChildName(auto, languageWin, "en-US")
	if err != nil {
		return err
	}
	err = wa.Select(enUsListItem)
	if err != nil {
		return err
	}
	moveDownLink, err := findElementByName(auto, languageWin, "Move down")
	if err != nil {
		return err
	}
	err = wa.Invoke(moveDownLink)
	if err != nil {
		return err
	}
	removeLink, err := findElementByName(auto, languageWin, "Remove")
	if err != nil {
		return err
	}
	err = wa.Invoke(removeLink)
	if err != nil {
		return err
	}

	return nil
}

func findElementByName(auto *wa.IUIAutomation, start *wa.IUIAutomationElement, elementName string) (*wa.IUIAutomationElement, error) {
	return wa.WaitFindFirstWithBreadthFirstSearch(
		auto, start, wa.NewElemMatcherFuncWithName(elementName))
}

func findParentElementByChildName(auto *wa.IUIAutomation, start *wa.IUIAutomationElement, childName string) (*wa.IUIAutomationElement, error) {
	child, err := findElementByName(auto, start, childName)
	if err != nil {
		return nil, err
	}
	return getParentElement(auto, child)
}

func getParentElement(auto *wa.IUIAutomation, element *wa.IUIAutomationElement) (*wa.IUIAutomationElement, error) {
	walker, err := wa.NewTreeWalker(auto)
	if err != nil {
		return nil, err
	}
	defer walker.Release()
	return walker.GetParentElement(element)
}
