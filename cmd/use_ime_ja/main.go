package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	wa "github.com/hnakamur/w32uiautomation"
	"github.com/mattn/go-ole"
)

func switchInputMethodJa() error {
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

func main() {
	ole.CoInitialize(0)
	defer ole.CoUninitialize()

	err := switchInputMethodJa()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
