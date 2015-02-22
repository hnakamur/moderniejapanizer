package moderniejapanizer

import (
	"fmt"
	"os/exec"

	wa "github.com/hnakamur/w32uiautomation"
	"github.com/hnakamur/w32version"
)

func SwitchInputMethodJa(version w32version.W32Version) error {
	if version == w32version.Windows7 || version == w32version.WindowsVista {
		return SetKeyboards([]string{JapaneseJapanKeyboardCode})
	} else {
		return switchInputMethodJaWin8()
	}
}

func switchInputMethodJaWin8() error {
	err := exec.Command("control.exe", "/name", "Microsoft.Language").Start()
	if err != nil {
		return err
	}
	fmt.Println("Opened control panel. Waiting the window to populated.")

	auto, err := wa.NewUIAutomation()
	if err != nil {
		return err
	}
	fmt.Println("Created UIAutomation")

	root, err := auto.GetRootElement()
	if err != nil {
		return err
	}
	defer root.Release()
	fmt.Println("Got RootElement")

	languageWin, err := findChildElementByName(auto, root, "Language")
	if err != nil {
		return err
	}
	fmt.Println(`Found "Language" window`)

	languageWinName, err := languageWin.Get_CurrentName()
	if err != nil {
		return err
	}
	fmt.Printf("languageWinName=%s\n", languageWinName)

	addALanguageLink, err := findElementByName(auto, languageWin, "Add a language")
	if err != nil {
		return err
	}
	fmt.Println(`Found "Add a language" link`)
	err = wa.Invoke(addALanguageLink)
	if err != nil {
		return err
	}
	fmt.Println(`Invoked "Add a language" link`)

	addLanguagesWin, err := findChildElementByName(auto, root, "Add languages")
	if err != nil {
		return err
	}
	fmt.Println(`Found "Add languages" window`)
	japaneseListItem, err := findElementByName(auto, addLanguagesWin, "Japanese")
	if err != nil {
		return err
	}
	fmt.Println(`Found "Japanese" listItem`)
	err = wa.Invoke(japaneseListItem)
	if err != nil {
		return err
	}
	fmt.Println(`Invoked "Japanese" listItem`)

	languageWin, err = findChildElementByName(auto, root, "Language")
	if err != nil {
		return err
	}
	fmt.Println(`Found "Language" window`)
	enUsListItem, err := findParentElementByChildName(auto, languageWin, "en-US")
	if err != nil {
		return err
	}
	fmt.Println("Found en-US listItem")
	err = wa.Select(enUsListItem)
	if err != nil {
		return err
	}
	fmt.Println("Selected en-US listItem")
	moveDownLink, err := findElementByName(auto, languageWin, "Move down")
	if err != nil {
		return err
	}
	fmt.Println(`Found "Move down" link`)
	err = wa.Invoke(moveDownLink)
	if err != nil {
		return err
	}
	fmt.Println(`Invoked "Move down" link`)
	removeLink, err := findElementByName(auto, languageWin, "Remove")
	if err != nil {
		return err
	}
	fmt.Println(`Found "Remove" link`)
	err = wa.Invoke(removeLink)
	if err != nil {
		return err
	}
	fmt.Println(`Invoked "Remove" link`)
	closeButton, err := findElementByName(auto, languageWin, "Close")
	if err != nil {
		return err
	}
	fmt.Println(`Found "Close" button`)
	err = wa.Invoke(closeButton)
	if err != nil {
		return err
	}
	fmt.Println(`Invoked "Close" button`)

	return nil
}

func findChildElementByName(auto *wa.IUIAutomation, start *wa.IUIAutomationElement, elementName string) (*wa.IUIAutomationElement, error) {
	condVal := wa.NewVariantString(elementName)
	condition, err := auto.CreatePropertyCondition(wa.UIA_NamePropertyId, condVal)
	if err != nil {
		return nil, err
	}
	return wa.WaitFindFirst(start, wa.TreeScope_Children, condition)
}

func findElementByName(auto *wa.IUIAutomation, start *wa.IUIAutomationElement, elementName string) (*wa.IUIAutomationElement, error) {
	condVal := wa.NewVariantString(elementName)
	condition, err := auto.CreatePropertyCondition(wa.UIA_NamePropertyId, condVal)
	if err != nil {
		return nil, err
	}
	return wa.WaitFindFirst(start, wa.TreeScope_Subtree, condition)
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
