package main

import (
	"fmt"
	"os"

	"github.com/hnakamur/ieversionlocker"
	"github.com/hnakamur/moderniejapanizer"
	"github.com/hnakamur/w32syscall"
	"github.com/hnakamur/w32timezone"
	"github.com/hnakamur/windowsupdate"
	"github.com/mattn/go-ole"
)

func main() {
	ole.CoInitialize(0)
	defer ole.CoUninitialize()

	version, err := ieversionlocker.CurrentVersion()
	if err != nil {
		fmt.Println("Failed detect IE version: %s", err)
		os.Exit(1)
	}

	err = ieversionlocker.Lock(version)
	if err != nil {
		fmt.Println("Failed to lock IE version: %s", err)
		os.Exit(1)
	}

	tzi, err := w32timezone.BuildDynamicTimeZoneInformation("Tokyo Standard Time")
	if err != nil {
		panic(err)
	}
	err = w32timezone.SetSystemTimeZone(tzi)
	if err != nil {
		panic(err)
	}

	_, _, err = windowsupdate.InstallLanguagePack(windowsupdate.JapaneseLanguagePackUpdateID)
	if err != nil {
		fmt.Printf("Error while installing Japanese Language Pack: %v\n", err)
		os.Exit(1)
	}

	err = moderniejapanizer.SetLanguageAndRegionalFormats(moderniejapanizer.JapaneseLanguageAndRegionalFormats)
	if err != nil {
		fmt.Printf("Error while setting language and regional formats: %v\n", err)
		os.Exit(1)
	}

	err = moderniejapanizer.SetLocation(moderniejapanizer.JapaneseLocationCode)
	if err != nil {
		fmt.Printf("Error while setting location: %v\n", err)
		os.Exit(1)
	}

	err = moderniejapanizer.SetKeyboards([]string{
		moderniejapanizer.JapaneseJapanKeyboardCode,
		moderniejapanizer.EnglishUnitedStatesKeyboardCode})
	if err != nil {
		fmt.Printf("Error while setting keyboards: %v\n", err)
		os.Exit(1)
	}

	err = moderniejapanizer.SetDisplayLanguage(moderniejapanizer.JapaneseDisplayLanguageCode)
	if err != nil {
		fmt.Printf("Error while setting display language: %v\n", err)
		os.Exit(1)
	}

	_, _, err = windowsupdate.InstallImportantUpdates()
	if err != nil {
		fmt.Printf("Error while installing important windows updates: %v\n", err)
		os.Exit(1)
	}

	err = moderniejapanizer.Reboot(w32syscall.SHTDN_REASON_MINOR_SECURITYFIX)
	if err != nil {
		fmt.Printf("Error while rebooting: %v\n", err)
		os.Exit(1)
	}
}
