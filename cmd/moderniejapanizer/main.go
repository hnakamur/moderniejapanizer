package main

import (
	"fmt"
	"os"

	"github.com/hnakamur/ieversionlocker"
	"github.com/hnakamur/moderniejapanizer"
	"github.com/hnakamur/w32syscall"
	"github.com/hnakamur/w32timezone"
	"github.com/hnakamur/w32version"
	"github.com/hnakamur/windowsupdate"
	"github.com/mattn/go-ole"
)

func main() {
	ole.CoInitialize(0)
	defer ole.CoUninitialize()

	fmt.Println("start modern.IE Japanizer.")
	ieVersion, err := ieversionlocker.CurrentVersion()
	if err != nil {
		fmt.Println("Failed detect IE version: %s", err)
		os.Exit(1)
	}

	err = ieversionlocker.Lock(ieVersion)
	if err != nil {
		fmt.Println("Failed to lock IE version: %s", err)
		os.Exit(1)
	}
	fmt.Println("lock IE version. done")

	tzi, err := w32timezone.BuildDynamicTimeZoneInformation("Tokyo Standard Time")
	if err != nil {
		panic(err)
	}
	err = w32timezone.SetSystemTimeZone(tzi)
	if err != nil {
		panic(err)
	}
	fmt.Println("set timezone. done")

	version, err := w32version.GetVersion()
	if err != nil {
		panic(err)
	}

	err = moderniejapanizer.InstallLangPackJa(version)
	if err != nil {
		fmt.Printf("Error while installing Japanese Language Pack: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("install Japanese language pack. done")

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
	fmt.Println("set location. done")

	err = moderniejapanizer.SetKeyboards([]string{
		moderniejapanizer.JapaneseJapanKeyboardCode,
	})
	if err != nil {
		fmt.Printf("Error while setting keyboards: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("set keyboards. done")

	err = moderniejapanizer.SetDisplayLanguage(moderniejapanizer.JapaneseDisplayLanguageCode)
	if err != nil {
		fmt.Printf("Error while setting display language: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("schedule display change after reboot. done")

	fmt.Println("start installing important windows updates...")
	_, _, err = windowsupdate.InstallImportantUpdates()
	if err != nil {
		fmt.Printf("Error while installing important windows updates: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("install important windows updates. done")

	err = moderniejapanizer.Reboot(w32syscall.SHTDN_REASON_MINOR_SECURITYFIX)
	if err != nil {
		fmt.Printf("Error while rebooting: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("rebooting...")
}
