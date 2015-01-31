package main

import (
	"fmt"

	"github.com/hnakamur/moderniejapanizer"
	"github.com/hnakamur/w32syscall"
	"github.com/hnakamur/windowsupdate"
	"github.com/mattn/go-ole"
)

func main() {
	ole.CoInitialize(0)
	defer ole.CoUninitialize()

	_, _, err := windowsupdate.InstallLanguagePack(windowsupdate.JapaneseLanguagePackUpdateID)
	if err != nil {
		fmt.Printf("Error while installing Japanese Language Pack: %v\n", err)
		return
	}

	err = moderniejapanizer.SetLanguageAndRegionalFormats(moderniejapanizer.JapaneseLanguageAndRegionalFormats)
	if err != nil {
		fmt.Printf("Error while setting language and regional formats: %v\n", err)
		return
	}

	err = moderniejapanizer.SetLocation(moderniejapanizer.JapaneseLocationCode)
	if err != nil {
		fmt.Printf("Error while setting location: %v\n", err)
		return
	}

	err = moderniejapanizer.SetKeyboards([]string{
		moderniejapanizer.JapaneseJapanKeyboardCode,
		moderniejapanizer.EnglishUnitedStatesKeyboardCode})
	if err != nil {
		fmt.Printf("Error while setting keyboards: %v\n", err)
		return
	}

	err = moderniejapanizer.SetDisplayLanguage(moderniejapanizer.JapaneseDisplayLanguageCode)
	if err != nil {
		fmt.Printf("Error while setting display language: %v\n", err)
		return
	}

	_, _, err = windowsupdate.InstallImportantUpdates()
	if err != nil {
		fmt.Printf("Error while installing important windows updates: %v\n", err)
		return
	}

	err = moderniejapanizer.Reboot(w32syscall.SHTDN_REASON_MINOR_SECURITYFIX)
	if err != nil {
		fmt.Printf("Error while rebooting: %v\n", err)
		return
	}
}
