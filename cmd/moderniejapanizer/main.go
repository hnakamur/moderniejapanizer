package main

import (
	"fmt"
	"os"

	"github.com/hnakamur/ieversionlocker"
	"github.com/hnakamur/moderniejapanizer"
	"github.com/hnakamur/w32syscall"
	"github.com/hnakamur/w32timezone"
	"github.com/hnakamur/w32version"
	"github.com/mattn/go-ole"
)

func main() {
	ole.CoInitialize(0)
	defer ole.CoUninitialize()

	fmt.Println("Start modern.IE Japanizer. Please wait until reboot.")

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
	fmt.Println("Locked IE version.")

	tzi, err := w32timezone.BuildDynamicTimeZoneInformation("Tokyo Standard Time")
	if err != nil {
		panic(err)
	}
	err = w32timezone.SetSystemTimeZone(tzi)
	if err != nil {
		panic(err)
	}
	fmt.Println("Set timezone.")

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
	fmt.Println("Set location.")

	version, err := w32version.GetVersion()
	if err != nil {
		panic(err)
	}

	err = moderniejapanizer.SwitchInputMethodJa(version)
	if err != nil {
		fmt.Printf("Error while setting display language: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Switched language to Japanese.")

	err = moderniejapanizer.InstallLangPackJa(version)
	if err != nil {
		fmt.Printf("Error while installing Japanese Language Pack: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Installed Japanese language pack.")

	err = moderniejapanizer.SetDisplayLanguage(moderniejapanizer.JapaneseDisplayLanguageCode)
	if err != nil {
		fmt.Printf("Error while setting display language to Japanese: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Scheduled display change after reboot.")

	err = moderniejapanizer.Reboot(w32syscall.SHTDN_REASON_MINOR_SECURITYFIX)
	if err != nil {
		fmt.Printf("Error while rebooting: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Rebooting...")
}
